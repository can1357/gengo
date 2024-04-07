package clang

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

type Options struct {
	ToolkitPath      string
	AdditionalParams []string
	Sources          []string
}

func (o *Options) ClangPath() string {
	if o.ToolkitPath != "" {
		if stat, err := os.Stat(o.ToolkitPath); err == nil && stat.IsDir() {
			return filepath.Join(o.ToolkitPath, "clang")
		} else {
			return o.ToolkitPath
		}
	}
	return "clang"
}
func (o *Options) ClangCommand(opt ...string) ([]byte, error) {
	cmd := exec.Command(o.ClangPath(), opt...)
	cmd.Args = append(cmd.Args, o.AdditionalParams...)
	cmd.Args = append(cmd.Args, o.Sources...)
	buf := &bytes.Buffer{}
	cmd.Stdout = buf
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run clang: %w", err)
	}
	return buf.Bytes(), nil
}

func CreateAST(opt *Options) ([]byte, error) {
	return opt.ClangCommand(
		"-fsyntax-only",
		"-nobuiltininc",
		"-Xclang",
		"-ast-dump=json",
	)
}

func CreateLayoutMap(opt *Options) ([]byte, error) {
	return opt.ClangCommand(
		"-fsyntax-only",
		"-nobuiltininc",
		"-emit-llvm",
		"-Xclang",
		"-fdump-record-layouts",
		"-Xclang",
		"-fdump-record-layouts-complete",
	)
}

func Parse(opt *Options) (ast Node, layout *LayoutMap, err error) {
	errg := &errgroup.Group{}
	errg.Go(func() error {
		res, e := CreateAST(opt)
		if e != nil {
			return e
		}
		ast, e = ParseAST(res)
		return e
	})
	errg.Go(func() error {
		res, e := CreateLayoutMap(opt)
		if e != nil {
			return e
		}
		layout, e = ParseLayoutMap(res)
		return e
	})
	if err := errg.Wait(); err != nil {
		return nil, nil, err
	}
	return ast, layout, nil
}
