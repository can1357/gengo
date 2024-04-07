package clang

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Options struct {
	ToolkitPath      string
	AdditionalParams []string
	Sources          []string
}

func (o *Options) GetClangPath() string {
	if o.ToolkitPath != "" {
		if stat, err := os.Stat(o.ToolkitPath); err == nil && stat.IsDir() {
			return filepath.Join(o.ToolkitPath, "clang")
		} else {
			return o.ToolkitPath
		}
	}
	return "clang"
}

func GenerateAST(opt *Options) ([]byte, error) {
	opts := append([]string{
		"-fsyntax-only",
		"-nobuiltininc",
		"-Xclang",
		"-ast-dump=json",
	}, opt.Sources...)
	opts = append(opts, opt.AdditionalParams...)
	cmd := exec.Command(opt.GetClangPath(), opts...)
	buf := &bytes.Buffer{}
	cmd.Stdout = buf
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to generate AST: %w", err)
	}
	return buf.Bytes(), nil
}

func GenerateLayout(opt *Options) ([]byte, error) {
	opts := append([]string{
		"-fsyntax-only",
		"-nobuiltininc",
		"-emit-llvm",
		"-Xclang",
		"-fdump-record-layouts",
		"-Xclang",
		"-fdump-record-layouts-complete",
	}, opt.Sources...)
	opts = append(opts, opt.AdditionalParams...)
	cmd := exec.Command(opt.GetClangPath(), opts...)
	buf := &bytes.Buffer{}
	cmd.Stdout = buf
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to generate layout: %w", err)
	}
	return buf.Bytes(), nil
}

func Transform(opt *Options) (Node, *Layouts, error) {
	ast, err := GenerateAST(opt)
	if err != nil {
		return nil, nil, err
	}
	astl, err := GenerateLayout(opt)
	if err != nil {
		return nil, nil, err
	}
	layouts, err := ParseLayout(astl)
	if err != nil {
		return nil, nil, err
	}
	node, err := ParseAST(ast)
	if err != nil {
		return nil, nil, err
	}
	return node, layouts, nil
}
