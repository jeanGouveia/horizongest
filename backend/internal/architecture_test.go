//go:build ignore_test

// Architecture Tests - Garantem que a arquitetura de camadas é respeitada
// Este arquivo deve ser executado com: go test -tags=ignore_test ./internal/...

package internal

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDomainDoesNotImportInfra(t *testing.T) {
	domainDir := "internal/domain"
	infraDir := "internal/infra"

	err := filepath.Walk(domainDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}

		for _, imp := range node.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			if strings.HasPrefix(importPath, infraDir) {
				t.Errorf("Arquivo %s no domain importa infra: %s", path, importPath)
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Erro ao percorrer diretório domain: %v", err)
	}
}

func TestDomainDoesNotImportRepository(t *testing.T) {
	domainDir := "internal/domain"
	repoDir := "internal/repository"

	err := filepath.Walk(domainDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}

		for _, imp := range node.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			if strings.HasPrefix(importPath, repoDir) {
				t.Errorf("Arquivo %s no domain importa repository: %s", path, importPath)
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Erro ao percorrer diretório domain: %v", err)
	}
}

func TestServiceDoesNotImportHandler(t *testing.T) {
	serviceDir := "internal/service"
	handlerDir := "internal/handler"

	err := filepath.Walk(serviceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}

		for _, imp := range node.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			if strings.HasPrefix(importPath, handlerDir) {
				t.Errorf("Arquivo %s no service importa handler: %s", path, importPath)
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Erro ao percorrer diretório service: %v", err)
	}
}

func TestRepositoryDoesNotImportHandler(t *testing.T) {
	repoDir := "internal/repository"
	handlerDir := "internal/handler"

	err := filepath.Walk(repoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}

		for _, imp := range node.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			if strings.HasPrefix(importPath, handlerDir) {
				t.Errorf("Arquivo %s no repository importa handler: %s", path, importPath)
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Erro ao percorrer diretório repository: %v", err)
	}
}

func TestRepositoryDoesNotImportService(t *testing.T) {
	repoDir := "internal/repository"
	serviceDir := "internal/service"

	err := filepath.Walk(repoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}

		for _, imp := range node.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			if strings.HasPrefix(importPath, serviceDir) {
				t.Errorf("Arquivo %s no repository importa service: %s", path, importPath)
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Erro ao percorrer diretório repository: %v", err)
	}
}
