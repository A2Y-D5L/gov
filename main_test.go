package gov_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestUse(t *testing.T) {
	t.Run("Switches to a specified version of Go, setting it as the active version", func(t *testing.T) {
		// Create a temporary directory
		tmpDir, err := os.MkdirTemp("", "gov-test")
		if err != nil {
			t.Fatalf("Failed to create temporary directory: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		// Set the GOPATH environment variable
		if err := os.Setenv("GOPATH", tmpDir); err != nil {
			t.Fatalf("Failed to set GOPATH: %v", err)
		}

		// Set the PATH environment variable
		if err := os.Setenv("PATH", fmt.Sprintf("%s:%s", filepath.Join(tmpDir, "go", "bin"), os.Getenv("PATH"))); err != nil {
			t.Fatalf("Failed to set PATH: %v", err)
		}

		// Create a new Go version
		goVersion := "1.15.6"
		goRoot := filepath.Join(tmpDir, "go", goVersion)
		if err := os.MkdirAll(goRoot, 0o755); err != nil {
			t.Fatalf("Failed to create Go root directory: %v", err)
		}

		// Create a temporary Go binary
		goBin := filepath.Join(goRoot, "bin", "go")
		if err := os.WriteFile(goBin, []byte("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, Go!\")\n}\n"), 0o755); err != nil {
			t.Fatalf("Failed to create Go binary: %v", err)
		}

		// Create a symlink to the Go binary
		goBinLink := filepath.Join(tmpDir, "go", "bin", "go")
		if err := os.Symlink(goBin, goBinLink); err != nil {
			t.Fatalf("Failed to create symlink: %v", err)
		}

		// Set the GOBIN environment variable
		if err := os.Setenv("GOBIN", filepath.Join(tmpDir, "go", "bin")); err != nil {
			t.Fatalf("Failed to set GOBIN: %v", err)
		}

		// Create a temporary Go workspace
		goWorkspace := filepath.Join(tmpDir, "workspace")
		if err := os.MkdirAll(goWorkspace, 0o755); err != nil {
			t.Fatalf("Failed to create Go workspace: %v", err)
		}

		// Set the GOROOT environment variable
		if err := os.Setenv("GOROOT", goRoot); err != nil {
			t.Fatalf("Failed to set GOROOT: %v", err)
		}

		// Set the GO111MODULE environment variable
		if err := os.Setenv("GO111MODULE", "on"); err != nil {
			t.Fatalf("Failed to set GO111MODULE: %v", err)
		}

		// Set the GOOS environment variable
		if err := os.Setenv("GOOS", "linux"); err != nil {
			t.Fatalf("Failed to set GOOS: %v", err)
		}

		// Set the GOARCH environment variable
		if err := os.Setenv("GOARCH", "amd64"); err != nil {
			t.Fatalf("Failed to set GOARCH: %v", err)
		}

		// Set the GOCACHE environment variable
		if err := os.Setenv("GOCACHE", filepath.Join(tmpDir, "cache")); err != nil {
			t.Fatalf("Failed to set GOCACHE: %v", err)
		}

		// Set the GOMODCACHE environment variable
		if err := os.Setenv("GOMODCACHE", filepath.Join(tmpDir, "modcache")); err != nil {
			t.Fatalf("Failed to set GOMODCACHE: %v", err)
		}

		// Set the GOPROXY environment variable
		if err := os.Setenv("GOPROXY", "https://proxy.golang.org,direct"); err != nil {
			t.Fatalf("Failed to set GOPROXY: %v", err)
		}

		// Set the GOSUMDB environment variable
		if err := os.Setenv("GOSUMDB", "sum.golang.org"); err != nil {
			t.Fatalf("Failed to set GOSUMDB: %v", err)
		}
	})
}
