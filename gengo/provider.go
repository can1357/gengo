package gengo

import (
	"fmt"
	"go/token"
	"strconv"
	"strings"

	"github.com/can1357/gengo/clang"
	"github.com/dave/dst"
	"github.com/iancoleman/strcase"
)

type TypeRef struct {
	Identifier
	Decl *dst.GenDecl
}

type TypeClass int

const (
	TcEnum TypeClass = iota
	TcStruct
	TcUnion
	TcTypedef
	TcCount
)

// Builtin types.
const (
	BuiltinVoid       = RelaxedIdentifier("void")
	BuiltinInt8       = RelaxedIdentifier("int8")
	BuiltinInt16      = RelaxedIdentifier("int16")
	BuiltinInt32      = RelaxedIdentifier("int32")
	BuiltinInt64      = RelaxedIdentifier("int64")
	BuiltinUint8      = RelaxedIdentifier("uint8")
	BuiltinUint16     = RelaxedIdentifier("uint16")
	BuiltinUint32     = RelaxedIdentifier("uint32")
	BuiltinUint64     = RelaxedIdentifier("uint64")
	BuiltinUint       = RelaxedIdentifier("uint")
	BuiltinInt        = RelaxedIdentifier("int")
	BuiltinFloat32    = RelaxedIdentifier("float32")
	BuiltinFloat64    = RelaxedIdentifier("float64")
	BuiltinComplex64  = RelaxedIdentifier("complex64")
	BuiltinComplex128 = RelaxedIdentifier("complex128")
	BuiltinBool       = RelaxedIdentifier("bool")
	BuiltinByte       = RelaxedIdentifier("byte")
	BuiltinRune       = RelaxedIdentifier("rune")
	BuiltinUintptr    = RelaxedIdentifier("uintptr")
	BuiltinAny        = RelaxedIdentifier("any")
)

// Reserved identifiers.
var reservedIdentifiers = map[string]struct{}{
	"break":       {},
	"default":     {},
	"func":        {},
	"interface":   {},
	"select":      {},
	"case":        {},
	"defer":       {},
	"go":          {},
	"map":         {},
	"struct":      {},
	"chan":        {},
	"else":        {},
	"goto":        {},
	"package":     {},
	"switch":      {},
	"const":       {},
	"fallthrough": {},
	"if":          {},
	"range":       {},
	"type":        {},
	"continue":    {},
	"for":         {},
	"import":      {},
	"return":      {},
	"var":         {},
}

// Default mappings of builtin C types.
var defaultBuiltinMap = map[string]Identifier{
	"void":                 BuiltinVoid,
	"char":                 BuiltinByte,
	"short":                BuiltinInt16,
	"int":                  BuiltinInt32,
	"long":                 BuiltinInt64,
	"long long":            BuiltinInt64,
	"signed char":          BuiltinInt8,
	"unsigned char":        BuiltinUint8,
	"signed short":         BuiltinInt16,
	"unsigned short":       BuiltinUint16,
	"signed int":           BuiltinInt32,
	"unsigned int":         BuiltinUint32,
	"signed long":          BuiltinInt64,
	"unsigned long":        BuiltinUint64,
	"signed long long":     BuiltinInt64,
	"unsigned long long":   BuiltinUint64,
	"signed":               BuiltinInt32,
	"unsigned":             BuiltinUint32,
	"float":                BuiltinFloat32,
	"double":               BuiltinFloat64,
	"long double":          BuiltinFloat64,
	"size_t":               BuiltinUint,
	"uint8_t":              BuiltinUint8,
	"uint16_t":             BuiltinUint16,
	"uint32_t":             BuiltinUint32,
	"uint64_t":             BuiltinUint64,
	"int8_t":               BuiltinInt8,
	"int16_t":              BuiltinInt16,
	"int32_t":              BuiltinInt32,
	"int64_t":              BuiltinInt64,
	"uintptr_t":            BuiltinUintptr,
	"intmax_t":             BuiltinInt64,
	"uintmax_t":            BuiltinUint64,
	"intptr_t":             BuiltinInt,
	"bool":                 BuiltinBool,
	"_Bool":                BuiltinBool,
	"char16_t":             BuiltinInt16,
	"char32_t":             BuiltinInt32,
	"float _Complex":       BuiltinComplex64,
	"double _Complex":      BuiltinComplex128,
	"long double _Complex": BuiltinComplex128,
}

// Provider is a type provider allowing hooks for library specific type conversions.
type Provider interface {
	// ConvertFieldName converts a field name to a Go compatible name.
	ConvertFieldName(name string) string
	// ConvertTypeName converts a type name to a Go compatible name.
	ConvertTypeName(name string) string
	// ConvertValueName converts a value name to a Go compatible name.
	ConvertValueName(name string) string
	// ConvertFuncName converts a function name to a Go compatible name.
	ConvertFuncName(name string) string
	// ConvertArgName converts an argument name to a Go compatible name.
	ConvertArgName(name string) string

	// ForceSyntethic returns true if the type should be syntethic (accessors instead of direct fields).
	ForceSyntethic(name string) bool
	// ConvertQualType converts a qualified type to a Go expression.
	ConvertQualType(q string) dst.Expr
	// ConvertTypeExpr converts a type node to a Go expression.
	ConvertTypeExpr(n clang.TypeNode) dst.Expr

	// AddType adds a type to the provider.
	AddType(tc TypeClass, name string, decl dst.Expr) TypeRef
	// RemapType remaps a given qualified typename to a type reference.
	RemapType(name string, tr TypeRef)

	// InferMethod returns the receiver type given the method name.
	InferMethod(name string) (rcv string, newName string)
}

// Normalize anonymous type names.
func normalizeAnonName(name string) string {
	if strings.IndexByte(name, ')') != -1 {
		// struct ZydisFormatter_::(anonymous at ../rsrc/Zydis.h:11222:5)
		tag, rest, _ := strings.Cut(name, " ")
		_, pos, _ := strings.Cut(rest, " at ")
		return tag + "@" + pos
	}
	return name
}

type BasicProviderOption func(*BasicProvider)

func WithRemovePrefix(prefixes ...string) BasicProviderOption {
	return func(p *BasicProvider) {
		p.RemovedPrefixes = append(p.RemovedPrefixes, prefixes...)
	}
}
func WithInferredMethods(rules []MethodInferenceRule) BasicProviderOption {
	return func(p *BasicProvider) {
		p.InferredMethods = append(p.InferredMethods, rules...)
	}
}

type MethodInferenceRule struct {
	Name     string
	Receiver string
}

type BasicProvider struct {
	Types           map[string]TypeRef
	Builtins        map[string]Identifier
	RemovedPrefixes []string
	InferredMethods []MethodInferenceRule
}

func NewBasicProvider(opt ...BasicProviderOption) *BasicProvider {
	dc := &BasicProvider{
		Types:           map[string]TypeRef{},
		Builtins:        map[string]Identifier{},
		InferredMethods: []MethodInferenceRule{},
	}
	for k, v := range defaultBuiltinMap {
		dc.Builtins[k] = v
	}
	for _, o := range opt {
		o(dc)
	}
	return dc
}
func (p *BasicProvider) removePrefixes(name string) string {
	for _, prefix := range p.RemovedPrefixes {
		if len(name) > len(prefix) && strings.EqualFold(name[:len(prefix)], prefix) {
			return name[len(prefix):]
		}
	}
	return name
}

func (p *BasicProvider) ConvertFieldName(name string) string {
	return p.ConvertTypeName(name)
}
func (p *BasicProvider) ConvertFuncName(name string) string {
	return p.ConvertTypeName(name)
}
func (p *BasicProvider) ConvertTypeName(name string) string {
	if strings.HasSuffix(name, "_") {
		return p.ConvertTypeName(name[:len(name)-1]) + "_"
	}
	if strings.HasPrefix(name, "_") {
		return "_" + p.ConvertTypeName(name[1:])
	}
	name = p.removePrefixes(name)
	return strcase.ToCamel(name)
}
func (p *BasicProvider) ConvertValueName(name string) string {
	if strings.HasSuffix(name, "_") {
		return p.ConvertTypeName(name[:len(name)-1]) + "_"
	}
	if strings.HasPrefix(name, "_") {
		return "_" + p.ConvertTypeName(name[1:])
	}
	name = p.removePrefixes(name)
	return strcase.ToScreamingSnake(name)
}
func (p *BasicProvider) ConvertArgName(name string) string {
	if _, ok := reservedIdentifiers[name]; ok {
		return "_" + name
	}
	return name
}
func (p *BasicProvider) ForceSyntethic(name string) bool {
	return false
}
func (p *BasicProvider) ConvertQualType(q string) dst.Expr {
	// Dumb qualifiers.
	q = strings.TrimSpace(q)
	q = strings.ReplaceAll(q, "const ", "")
	q = strings.ReplaceAll(q, "volatile ", "")
	q = strings.ReplaceAll(q, "restrict ", "")
	q = strings.TrimSpace(q)

	// Pointer type.
	if q, ok := strings.CutSuffix(q, "*"); ok {
		res := p.ConvertQualType(q)
		if ident, ok := res.(*dst.Ident); ok && (ident.Name == "any" || ident.Name == "void") {
			return &dst.SelectorExpr{
				X:   dst.NewIdent("unsafe"),
				Sel: dst.NewIdent("Pointer"),
			}
		}
		return &dst.StarExpr{
			X: res,
		}
	}

	// Array type.
	if q, ok := strings.CutSuffix(q, "]"); ok {
		lastBracker := strings.LastIndex(q, "[")
		if lastBracker == -1 {
			fmt.Printf("[WARN] Invalid array type: %s\n", q)
			return BuiltinAny.Ref()
		}
		before := q[:lastBracker]
		after := q[lastBracker+1:]
		n, err := strconv.Atoi(after)
		if err != nil {
			fmt.Printf("[WARN] Invalid array size: %s (%s)\n", q, err)
			return BuiltinAny.Ref()
		}
		return &dst.ArrayType{
			Len: dst.NewIdent(strconv.Itoa(n)),
			Elt: p.ConvertQualType(before),
		}
	}

	// Builtin types.
	if b, ok := p.Builtins[q]; ok {
		return b.Ref()
	}

	// Normalize anonymous type names.
	q = normalizeAnonName(q)

	// Named types.
	if td, ok := p.Types[q]; ok {
		return td.Ref()
	}

	// Unknown type.
	fmt.Printf("[WARN] Unknown type: %s\n", q)
	return BuiltinAny.Ref()
}
func (p *BasicProvider) ConvertTypeExpr(n clang.TypeNode) dst.Expr {
	switch n := n.(type) {
	case *clang.BuiltinType:
		return p.ConvertQualType(n.Type.QualType)
	case *clang.PointerType:
		var innerExpr dst.Expr
		inner := clang.First[clang.TypeNode](n)
		if inner != nil {
			innerExpr = p.ConvertTypeExpr(inner)
			if innerExpr != nil {
				if ident, ok := innerExpr.(*dst.Ident); ok {
					if ident.Name == "void" || ident.Name == "any" {
						innerExpr = nil
					}
				}
			}
		}

		if innerExpr == nil {
			return &dst.SelectorExpr{
				X:   dst.NewIdent("unsafe"),
				Sel: dst.NewIdent("Pointer"),
			}
		}
		return &dst.StarExpr{
			X: innerExpr,
		}
	case *clang.QualType:
		var innerExpr dst.Expr
		inner := clang.First[clang.TypeNode](n)
		if inner != nil {
			innerExpr = p.ConvertTypeExpr(inner)
		}
		if innerExpr == nil {
			return p.ConvertQualType(n.Type.QualType)
		}
		return innerExpr
	case *clang.RecordType:
		return p.ConvertQualType(n.Type.QualType)
	case *clang.EnumType:
		return p.ConvertQualType(n.Type.QualType)
	case *clang.TypedefType:
		return p.ConvertQualType(n.Type.QualType)
	case *clang.ElaboratedType:
		return p.ConvertQualType(n.Type.QualType)
	default:
		//case *clang.FunctionProtoType:
		//case *clang.ParenType:
		return nil
	}
}
func (p *BasicProvider) AddType(tc TypeClass, name string, decl dst.Expr) TypeRef {
	ident := &TrackedIdentifier{Name: name}
	gen := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name:   ident.Ref(),
				Assign: tc == TcTypedef,
				Type:   decl,
			},
		},
	}

	tr := TypeRef{Identifier: ident, Decl: gen}
	switch tc {
	case TcEnum:
		p.Types["enum "+name] = tr
	case TcStruct:
		p.Types["struct "+name] = tr
	case TcUnion:
		p.Types["union "+name] = tr
	case TcTypedef:
		p.Types[name] = tr
	default:
		panic("invalid type class")
	}
	return tr
}
func (p *BasicProvider) RemapType(name string, tr TypeRef) {
	name = normalizeAnonName(name)
	p.Types[name] = tr
}
func (p *BasicProvider) InferMethod(name string) (rcv string, newName string) {
	for _, rule := range p.InferredMethods {
		if strings.HasPrefix(name, rule.Name) {
			return rule.Receiver, strings.TrimPrefix(name, rule.Name)
		}
	}
	return "", name
}
