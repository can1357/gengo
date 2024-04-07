package gengo

import (
	"fmt"
	"go/token"
	"strconv"
	"strings"

	"github.com/can1357/gengo/clang"
	"github.com/dave/dst"
)

type deferred []func()

func (d *deferred) Go() {
	list := *d
	*d = nil
	for _, f := range list {
		f()
	}
}
func (d *deferred) Do(f func()) {
	*d = append(*d, f)
}

type initTypeValidation struct {
	ty     TypeRef
	size   int
	align  int
	fields []dst.Expr
}

func ConvertEnumNode(n *clang.EnumDecl, mod Module) {
	// Add the typedef.
	var ty dst.Expr
	for _, c := range n.Inner {
		if c, ok := c.(*clang.EnumConstantDecl); ok {
			ty = mod.Parent.ConvertQualType(c.Type.QualType)
			break
		}
	}
	if ty == nil {
		ty = BuiltinInt.Ref()
	}
	tydef := mod.AddType(TcEnum, n.Name, ty)
	tydef.Rename(mod.Parent.ConvertTypeName(n.Name))
	tydef.Decl.Decs.Start.Append(n.Comments()...)

	// Add the enum constants.
	spec := []dst.Spec{}
	val := -1
	for _, c := range n.Inner {
		if c, ok := c.(*clang.EnumConstantDecl); ok {
			val += 1
			if ival := clang.First[clang.ValueNode](c); ival != nil {
				if cval, ok := ival.(clang.ConstValueNode); ok {
					val, _ = strconv.Atoi(cval.Value())
				} else {
					panic(fmt.Sprintf("unhandled value node: %T", ival))
				}
			}

			spec = append(spec, &dst.ValueSpec{
				Names: []*dst.Ident{dst.NewIdent(mod.Parent.ConvertValueName(c.Name))},
				Type:  tydef.Ref(),
				Values: []dst.Expr{
					&dst.BasicLit{
						Kind:  token.INT,
						Value: strconv.Itoa(val),
					},
				},
				Decs: dst.ValueSpecDecorations{
					Assign: c.Comments(),
				},
			})
		}
	}
	mod.Decls = append(mod.Decls, &dst.GenDecl{
		Tok:   token.CONST,
		Specs: spec,
	})
}

func DefineNaturalType(n *clang.RecordDecl, layout *clang.RecordLayout, mod Module, ty TypeRef, fl *dst.FieldList, vld *[]initTypeValidation) bool {
	// Can't define a natural type if it is a union.
	if n.TagUsed == "union" {
		return false
	}
	if mod.Parent.ForceSyntethic(n.Name) {
		return false
	}

	v := initTypeValidation{
		ty:     ty,
		size:   layout.Size,
		align:  layout.Align,
		fields: []dst.Expr{},
	}

	// Add the fields.
	for _, c := range layout.Fields {
		if c.Name == "" {
			fl.List = append(fl.List, &dst.Field{
				Names: []*dst.Ident{},
				Type:  mod.Parent.ConvertQualType(c.Type),
			})
		} else {
			name := mod.Parent.ConvertFieldName(c.Name)
			fl.List = append(fl.List, &dst.Field{
				Names: []*dst.Ident{dst.NewIdent(name)},
				Type:  mod.Parent.ConvertQualType(c.Type),
			})
			v.fields = append(v.fields,
				&dst.BasicLit{
					Kind:  token.STRING,
					Value: strconv.Quote(name),
				},
				&dst.BasicLit{
					Kind:  token.INT,
					Value: fmt.Sprintf("0x%x", c.Offset),
				},
			)
		}
	}

	*vld = append(*vld, v)
	return true
}

func DefineSyntheticType(layout *clang.RecordLayout, mod Module, ty TypeRef, fl *dst.FieldList) {
	// Add byte[size] field.
	primitive := BuiltinByte
	n := layout.Size
	if layout.Align >= 8 {
		primitive = BuiltinInt64
		n /= 8
	} else if layout.Align >= 4 {
		primitive = BuiltinInt32
		n /= 4
	} else if layout.Align >= 2 {
		primitive = BuiltinInt16
		n /= 2
	}
	fl.List = append(fl.List, &dst.Field{
		Names: []*dst.Ident{dst.NewIdent("Raw")},
		Type: &dst.ArrayType{
			Len: &dst.BasicLit{
				Kind:  token.INT,
				Value: strconv.Itoa(n),
			},
			Elt: primitive.Ref(),
		},
	})

	sliceOffset := func(offset int) dst.Expr {
		return &dst.CallExpr{
			Fun: &dst.SelectorExpr{
				X:   &dst.Ident{Name: "unsafe"},
				Sel: &dst.Ident{Name: "Add"},
			},
			Args: []dst.Expr{
				&dst.CallExpr{
					Fun: &dst.SelectorExpr{
						X:   &dst.Ident{Name: "unsafe"},
						Sel: &dst.Ident{Name: "Pointer"},
					},
					Args: []dst.Expr{
						&dst.CallExpr{
							Fun: &dst.SelectorExpr{
								X:   &dst.Ident{Name: "unsafe"},
								Sel: &dst.Ident{Name: "SliceData"},
							},
							Args: []dst.Expr{
								&dst.SliceExpr{
									X: &dst.SelectorExpr{
										X:   dst.NewIdent("s"),
										Sel: dst.NewIdent("Raw"),
									},
								},
							},
						},
					},
				},
				&dst.BasicLit{
					Kind:  token.INT,
					Value: strconv.Itoa(offset),
				},
			},
		}
	}

	// Add accessors for each field.
	for _, c := range layout.Fields {
		name := mod.Parent.ConvertFieldName(c.Name)
		mod.Decls = append(mod.Decls, &dst.FuncDecl{
			Name: dst.NewIdent(name),
			Recv: &dst.FieldList{
				List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent("s")},
						Type:  ty.Ref(),
					},
				},
			},
			Type: &dst.FuncType{
				Params: &dst.FieldList{},
				Results: &dst.FieldList{
					List: []*dst.Field{
						{
							Names: []*dst.Ident{},
							Type:  mod.Parent.ConvertQualType(c.Type),
						},
					},
				},
			},
			Body: &dst.BlockStmt{
				List: []dst.Stmt{
					&dst.ReturnStmt{
						Results: []dst.Expr{
							&dst.CallExpr{
								Fun: &dst.IndexExpr{
									X: &dst.SelectorExpr{
										X:   dst.NewIdent("gengort"),
										Sel: dst.NewIdent("ReadBitcast"),
									},
									Index: mod.Parent.ConvertQualType(c.Type),
								},
								Args: []dst.Expr{
									sliceOffset(c.Offset),
								},
							},
						},
					},
				},
			},
		})
		mod.Decls = append(mod.Decls, &dst.FuncDecl{
			Name: dst.NewIdent("Set" + name),
			Recv: &dst.FieldList{
				List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent("s")},
						Type:  &dst.StarExpr{X: ty.Ref()},
					},
				},
			},
			Type: &dst.FuncType{
				Params: &dst.FieldList{
					List: []*dst.Field{
						{
							Names: []*dst.Ident{dst.NewIdent("v")},
							Type:  mod.Parent.ConvertQualType(c.Type),
						},
					},
				},
				Results: &dst.FieldList{},
			},
			Body: &dst.BlockStmt{
				List: []dst.Stmt{
					&dst.ExprStmt{
						X: &dst.CallExpr{
							Fun: &dst.SelectorExpr{
								X:   dst.NewIdent("gengort"),
								Sel: dst.NewIdent("WriteBitcast"),
							},
							Args: []dst.Expr{
								sliceOffset(c.Offset),
								dst.NewIdent("v"),
							},
						},
					},
				},
			},
		})
	}
}

func ConvertRecordNode(n *clang.RecordDecl, layouts *clang.Layouts, mod Module, deferred *deferred, vld *[]initTypeValidation) {
	var layout *clang.RecordLayout
	var defName string
	var qualDefName string
	if n.Name == "" {
		line := fmt.Sprintf(":%d:%d", n.Loc.Line, n.Loc.Col)
		for ty, rec := range layouts.Map {
			if strings.Contains(ty, line) {
				layout = rec
				break
			}
		}
		if layout == nil {
			line = fmt.Sprintf(":%d:", n.Loc.Line)
			for ty, rec := range layouts.Map {
				if strings.Contains(ty, line) {
					layout = rec
					break
				}
			}
		}
		defName = fmt.Sprintf("Anon%d_%d", n.Loc.Line, n.Loc.Col)
		qualDefName = n.TagUsed + " " + defName
	} else {
		defName = n.Name
		qualDefName = n.TagUsed + " " + defName
		layout = layouts.Map[qualDefName]
	}

	if layout == nil {
		fmt.Println("[WARN] No layout for record:", n.Name)
		return
	}

	// Add the struct.
	tc := TcStruct
	if n.TagUsed == "union" {
		tc = TcUnion
	}
	fieldList := new(dst.FieldList)
	str := mod.AddType(tc, defName, &dst.StructType{
		Fields: fieldList,
	})
	str.Rename(mod.Parent.ConvertTypeName(defName))
	if qualDefName != layout.Type {
		mod.Parent.RemapType(layout.Type, str)
	}

	deferred.Do(func() {
		if !DefineNaturalType(n, layout, mod, str, fieldList, vld) {
			DefineSyntheticType(layout, mod, str, fieldList)
		}
	})
}

func ConvertFunctionNode(n *clang.FunctionDecl, mod Module) {
	// Add the import.
	// - var __imp_func = Import("func")
	//
	name := mod.Parent.ConvertFuncName(n.Name)
	importName := "__imp_" + name
	mod.Decls = append(mod.Decls, &dst.GenDecl{
		Tok: token.VAR,
		Specs: []dst.Spec{
			&dst.ValueSpec{
				Names: []*dst.Ident{dst.NewIdent(importName)},
				Values: []dst.Expr{
					&dst.CallExpr{
						Fun: &dst.SelectorExpr{
							X:   dst.NewIdent("gengort"),
							Sel: dst.NewIdent("Import"),
						},
						Args: []dst.Expr{
							&dst.BasicLit{
								Kind:  token.STRING,
								Value: strconv.Quote(mod.Parent.Name),
							},
							&dst.BasicLit{
								Kind:  token.STRING,
								Value: strconv.Quote(n.MangledName),
							},
						},
					},
				},
			},
		},
	})

	// Get the result type.
	var result, result2 dst.Expr
	{
		before, _, ok := strings.Cut(n.Type.QualType, "(")
		if !ok {
			panic("unhandled function type")
		}
		before = strings.TrimSpace(before)
		result = mod.Parent.ConvertQualType(before)
		if ident, ok := result.(*dst.Ident); ok {
			if ident.Name == "void" {
				result = nil
			}
		}
		if result != nil {
			result2 = mod.Parent.ConvertQualType(before)
		}
	}

	// Add the function.
	// - func Func(arg *T) {
	// - 	__imp_func.Call1(
	// - 		marshalsyscall(arg),
	// - 	)
	// - }
	typ := &dst.FuncType{
		Params:  &dst.FieldList{},
		Results: &dst.FieldList{},
	}
	decl := &dst.FuncDecl{
		Name: dst.NewIdent(name),
		Type: typ,
	}
	mod.Decls = append(mod.Decls, decl)
	callExpr := &dst.CallExpr{
		Fun: &dst.SelectorExpr{
			X: dst.NewIdent(importName),
		},
	}

	// Add the parameters.
	paramNodes := clang.All[*clang.ParmVarDecl](n)
	for _, p := range paramNodes {
		arg := mod.Parent.ConvertArgName(p.Name)
		typ.Params.List = append(typ.Params.List, &dst.Field{
			Names: []*dst.Ident{dst.NewIdent(arg)},
			Type:  mod.Parent.ConvertQualType(p.Type.QualType),
		})
		callExpr.Args = append(callExpr.Args, &dst.CallExpr{
			Fun: &dst.SelectorExpr{
				X:   dst.NewIdent("gengort"),
				Sel: dst.NewIdent("MarshallSyscall"),
			},
			Args: []dst.Expr{
				&dst.Ident{Name: arg},
			},
		})
	}
	callExpr.Fun.(*dst.SelectorExpr).Sel = dst.NewIdent("Call" + strconv.Itoa(len(callExpr.Args)))

	// Infer the method receiver.
	rcv, mnameNew := mod.Parent.InferMethod(n.Name)
	idx := -1
	if rcv != "" {
		for i, p := range paramNodes {
			if strings.Contains(p.Type.QualType, rcv) {
				idx = i
				break
			}
		}
		if idx != -1 {
			name = mod.Parent.ConvertFuncName(mnameNew)
			decl.Name = dst.NewIdent(name)
			decl.Recv = &dst.FieldList{
				List: []*dst.Field{
					typ.Params.List[idx],
				},
			}
			typ.Params.List = append(typ.Params.List[:idx], typ.Params.List[idx+1:]...)
		}
	}

	// Add the result.
	if result != nil {
		typ.Results.List = append(typ.Results.List, &dst.Field{
			Names: []*dst.Ident{},
			Type:  result,
		})
		decl.Body = &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.AssignStmt{
					Lhs: []dst.Expr{
						dst.NewIdent("__res"),
					},
					Tok: token.DEFINE,
					Rhs: []dst.Expr{
						callExpr,
					},
				},
				&dst.ReturnStmt{
					Results: []dst.Expr{
						&dst.CallExpr{
							Fun: &dst.IndexExpr{
								X: &dst.SelectorExpr{
									X:   dst.NewIdent("gengort"),
									Sel: dst.NewIdent("UnmarshallSyscall"),
								},
								Index: result2,
							},
							Args: []dst.Expr{
								dst.NewIdent("__res"),
							},
						},
					},
				},
			},
		}
		// func Func(arg *T) T {
		// __res := __imp_func.Call(
		//	...
		// )
		// return unmarshalsyscall[T](__res)

	} else {
		decl.Body = &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.ExprStmt{X: callExpr},
			},
		}
		// func Func(arg *T) {
		// __imp_func.Call(
		//	...
		// )
	}
}

func ConvertTypedefNode(n *clang.TypedefDecl, mod Module) {
	// Add the typedef.
	var ty dst.Expr
	if inner := clang.First[clang.TypeNode](n); inner != nil {
		ty = mod.Parent.ConvertTypeExpr(inner)
	} else {
		ty = mod.Parent.ConvertQualType(n.Type.QualType)
	}
	tydef := mod.AddType(TcTypedef, n.Name, ty)
	tydef.Rename(mod.Parent.ConvertTypeName(n.Name))
	tydef.Decl.Decs.Start.Append(n.Comments()...)
}

func Convert(ast clang.Node, layouts *clang.Layouts, mod Module) {
	deferred := deferred{}
	validators := []initTypeValidation{}

	// Define enums.
	clang.Visit(ast, func(ed *clang.EnumDecl) bool {
		ConvertEnumNode(ed, mod)
		return true
	})

	// Define structs and unions.
	clang.Visit(ast, func(rd *clang.RecordDecl) bool {
		if rd.CompleteDefinition {
			ConvertRecordNode(rd, layouts, mod, &deferred, &validators)
		}
		return true
	})

	// Define typedefs.
	clang.Visit(ast, func(td *clang.TypedefDecl) bool {
		ConvertTypedefNode(td, mod)
		return true
	})

	// Define functions.
	clang.Visit(ast, func(fd *clang.FunctionDecl) bool {
		ConvertFunctionNode(fd, mod)
		return true
	})

	// Run deferred functions.
	deferred.Go()

	// Validate types.
	validatorBody := []dst.Stmt{}
	for _, v := range validators {
		args := []dst.Expr{
			&dst.BasicLit{
				Kind:  token.INT,
				Value: fmt.Sprintf("0x%x", v.size),
			},
			&dst.BasicLit{
				Kind:  token.INT,
				Value: fmt.Sprintf("0x%x", v.align),
			},
		}
		args = append(args, v.fields...)

		validatorBody = append(validatorBody, &dst.ExprStmt{
			X: &dst.CallExpr{
				Fun: &dst.IndexExpr{
					X: &dst.SelectorExpr{
						X:   dst.NewIdent("gengort"),
						Sel: dst.NewIdent("Validate"),
					},
					Index: v.ty.Ref(),
				},
				Args: args,
			},
		})
	}

	// Find or create init
	var initFunc *dst.FuncDecl
	for _, decl := range mod.Decls {
		if f, ok := decl.(*dst.FuncDecl); ok && f.Name.Name == "init" {
			initFunc = f
			break
		}
	}
	if initFunc == nil {
		initFunc = &dst.FuncDecl{
			Name: dst.NewIdent("init"),
			Type: &dst.FuncType{},
			Body: &dst.BlockStmt{},
		}
		mod.Decls = append(mod.Decls, initFunc)
	}

	// Add the validation calls to the init function.
	initFunc.Body.List = append(initFunc.Body.List, validatorBody...)
}
