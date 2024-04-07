package gengo

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"path/filepath"

	"github.com/can1357/gengo/clang"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

type Package struct {
	dst.Package
	Restorer *decorator.Restorer
	Provider
}

func NewPackage(name string, prov Provider) *Package {
	return &Package{
		Package: dst.Package{
			Name:    name,
			Scope:   dst.NewScope(nil),
			Imports: make(map[string]*dst.Object),
			Files:   make(map[string]*dst.File),
		},
		Restorer: decorator.NewRestorer(),
		Provider: prov,
	}
}

func (p *Package) Transform(module string, opt *clang.Options) error {
	ast, layouts, err := clang.Transform(opt)
	if err != nil {
		return fmt.Errorf("failed to transform %+v: %w", opt.Sources, err)
	}
	main := p.Upsert(module)
	Convert(ast, layouts, main)
	return nil
}

func (p *Package) Fprint(fn func(path string) (io.WriteCloser, error)) error {
	for k, f := range p.Files {
		file, err := fn(k + ".go")
		if err != nil {
			return err
		}
		p.Restorer.Fprint(file, f)
		if err := file.Close(); err != nil {
			return err
		}
	}
	return nil
}
func (p *Package) Upsert(module string) Module {
	f, ok := p.Files[module]
	if !ok {
		f = &dst.File{
			Name:  dst.NewIdent(p.Name),
			Scope: dst.NewScope(nil),
		}
		f.Decs.Start.Append("// Code generated by gengo. DO NOT EDIT.\n")
		p.Files[module] = f

		// Import from unsafe & gengort.
		f.Decls = append(f.Decls, &dst.GenDecl{
			Tok: token.IMPORT,
			Specs: []dst.Spec{
				&dst.ImportSpec{
					Path: &dst.BasicLit{
						Kind:  token.STRING,
						Value: `"unsafe"`,
					},
				},
				&dst.ImportSpec{
					Path: &dst.BasicLit{
						Kind:  token.STRING,
						Value: `"github.com/can1357/gengo/gengort"`,
					},
				},
			},
		})
	}
	return Module{p, f}
}

type stdoutCloser struct{ io.Writer }

func (stdoutCloser) Close() error {
	fmt.Println()
	return nil
}
func (p *Package) Print() {
	p.Fprint(func(path string) (io.WriteCloser, error) {
		fmt.Println("//// ", path)
		return stdoutCloser{os.Stdout}, nil
	})
}
func (p *Package) WriteToDir(dir string) error {
	os.Mkdir(dir, 0755)
	return p.Fprint(func(path string) (io.WriteCloser, error) {
		return os.Create(filepath.Join(dir, path))
	})
}

type Module struct {
	Parent *Package
	*dst.File
}

func (m Module) AddType(tc TypeClass, name string, decl dst.Expr) TypeRef {
	ref := m.Parent.AddType(tc, name, decl)
	m.Decls = append(m.Decls, ref.Decl)
	return ref
}

func (m Module) Go() (*ast.File, error) {
	return m.Parent.Restorer.RestoreFile(m.File)
}
func (m Module) Fprint(w io.Writer) {
	m.Parent.Restorer.Fprint(w, m.File)
}
func (m Module) String() string {
	buf := &bytes.Buffer{}
	m.Fprint(buf)
	return buf.String()
}
func (m Module) Print() {
	m.Fprint(os.Stdout)
}
