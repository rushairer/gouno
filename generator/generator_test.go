package generator_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rushairer/gouno/generator"
	"github.com/spf13/cobra"
)

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	c, err = root.ExecuteC()

	// 重置所有子命令的标志到默认值，防止跨测试状态污染
	for _, subCmd := range root.Commands() {
		if f := subCmd.Flag("path"); f != nil {
			f.Value.Set(f.DefValue)
		}
		if f := subCmd.Flag("force"); f != nil {
			f.Value.Set(f.DefValue)
		}
	}

	return c, buf.String(), err
}

// chdir 切换到临时目录作为工作目录，测试结束后自动还原并清理
func chdir(t *testing.T) string {
	t.Helper()
	tmpDir, err := filepath.Abs(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.Chdir(cwd)
		os.RemoveAll(tmpDir)
	})
	return tmpDir
}

func TestGeneratorService(t *testing.T) {
	tmpDir := chdir(t)

	t.Run("default path", func(t *testing.T) {
		_, _, err := executeCommandC(generator.GeneratorCmd, "service", "foo_bar")
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		filePath := filepath.Join(tmpDir, "internal", "service", "foo_bar.go")
		assertFileExists(t, filePath)
		assertFileContains(t, filePath, "package service")
		assertFileContains(t, filePath, "FooBarService")
	})

	t.Run("force overwrite", func(t *testing.T) {
		_, _, err := executeCommandC(generator.GeneratorCmd, "service", "foo_bar", "--force")
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		assertFileExists(t, filepath.Join(tmpDir, "internal", "service", "foo_bar.go"))
	})

	t.Run("custom path", func(t *testing.T) {
		_, _, err := executeCommandC(generator.GeneratorCmd, "service", "foo_bar", "--path", "./custom/service")
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		filePath := filepath.Join(tmpDir, "custom", "service", "foo_bar.go")
		assertFileExists(t, filePath)
		assertFileContains(t, filePath, "package service")
	})
}

func TestGeneratorDomain(t *testing.T) {
	tmpDir := chdir(t)

	t.Run("default path", func(t *testing.T) {
		_, _, err := executeCommandC(generator.GeneratorCmd, "domain", "foo")
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		filePath := filepath.Join(tmpDir, "internal", "domain", "foo.go")
		assertFileExists(t, filePath)
		assertFileContains(t, filePath, "package domain")
	})

	t.Run("force overwrite", func(t *testing.T) {
		_, _, err := executeCommandC(generator.GeneratorCmd, "domain", "foo", "--force")
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		assertFileExists(t, filepath.Join(tmpDir, "internal", "domain", "foo.go"))
	})

	t.Run("custom path", func(t *testing.T) {
		_, _, err := executeCommandC(generator.GeneratorCmd, "domain", "foo", "--path", "./custom/domain")
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		filePath := filepath.Join(tmpDir, "custom", "domain", "foo.go")
		assertFileExists(t, filePath)
		assertFileContains(t, filePath, "package domain")
	})
}

func TestGeneratorRepository(t *testing.T) {
	tmpDir := chdir(t)

	t.Run("default path", func(t *testing.T) {
		_, _, err := executeCommandC(generator.GeneratorCmd, "repository", "foo")
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		filePath := filepath.Join(tmpDir, "internal", "repository", "foo.go")
		assertFileExists(t, filePath)
		assertFileContains(t, filePath, "package repository")
		assertFileContains(t, filePath, "FooRepository")
	})

	t.Run("force overwrite", func(t *testing.T) {
		_, _, err := executeCommandC(generator.GeneratorCmd, "repository", "foo", "--force")
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		assertFileExists(t, filepath.Join(tmpDir, "internal", "repository", "foo.go"))
	})

	t.Run("custom path", func(t *testing.T) {
		_, _, err := executeCommandC(generator.GeneratorCmd, "repository", "foo", "--path", "./custom/repository")
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		filePath := filepath.Join(tmpDir, "custom", "repository", "foo.go")
		assertFileExists(t, filePath)
		assertFileContains(t, filePath, "package repository")
	})
}

func TestGeneratorTask(t *testing.T) {
	tmpDir := chdir(t)

	t.Run("default path", func(t *testing.T) {
		_, _, err := executeCommandC(generator.GeneratorCmd, "task", "foo")
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		filePath := filepath.Join(tmpDir, "internal", "task", "foo.go")
		assertFileExists(t, filePath)
		assertFileContains(t, filePath, "package task")
		assertFileContains(t, filePath, "FooTask")
	})

	t.Run("force overwrite", func(t *testing.T) {
		_, _, err := executeCommandC(generator.GeneratorCmd, "task", "foo", "--force")
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		assertFileExists(t, filepath.Join(tmpDir, "internal", "task", "foo.go"))
	})

	t.Run("custom path", func(t *testing.T) {
		_, _, err := executeCommandC(generator.GeneratorCmd, "task", "foo", "--path", "./custom/task")
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		filePath := filepath.Join(tmpDir, "custom", "task", "foo.go")
		assertFileExists(t, filePath)
		assertFileContains(t, filePath, "package task")
	})
}

func TestGeneratorController(t *testing.T) {
	tmpDir := chdir(t)

	t.Run("default path", func(t *testing.T) {
		_, _, err := executeCommandC(generator.GeneratorCmd, "controller", "foo")
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		filePath := filepath.Join(tmpDir, "controller", "foo.go")
		assertFileExists(t, filePath)
		assertFileContains(t, filePath, "package controller")
		assertFileContains(t, filePath, "FooController")
	})

	t.Run("force overwrite", func(t *testing.T) {
		_, _, err := executeCommandC(generator.GeneratorCmd, "controller", "foo", "--force")
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		assertFileExists(t, filepath.Join(tmpDir, "controller", "foo.go"))
	})

	t.Run("custom path", func(t *testing.T) {
		_, _, err := executeCommandC(generator.GeneratorCmd, "controller", "foo", "--path", "./custom/controller")
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		filePath := filepath.Join(tmpDir, "custom", "controller", "foo.go")
		assertFileExists(t, filePath)
		assertFileContains(t, filePath, "package controller")
	})
}

func TestGeneratorSuite(t *testing.T) {
	tmpDir := chdir(t)

	t.Run("generates domain, repository and service", func(t *testing.T) {
		_, _, err := executeCommandC(generator.GeneratorCmd, "suite", "foo")
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}

		domainPath := filepath.Join(tmpDir, "internal", "domain", "foo.go")
		assertFileExists(t, domainPath)
		assertFileContains(t, domainPath, "package domain")

		repositoryPath := filepath.Join(tmpDir, "internal", "repository", "foo.go")
		assertFileExists(t, repositoryPath)
		assertFileContains(t, repositoryPath, "package repository")

		servicePath := filepath.Join(tmpDir, "internal", "service", "foo.go")
		assertFileExists(t, servicePath)
		assertFileContains(t, servicePath, "package service")
	})

	t.Run("force overwrite", func(t *testing.T) {
		_, _, err := executeCommandC(generator.GeneratorCmd, "suite", "foo", "--force")
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		assertFileExists(t, filepath.Join(tmpDir, "internal", "domain", "foo.go"))
		assertFileExists(t, filepath.Join(tmpDir, "internal", "repository", "foo.go"))
		assertFileExists(t, filepath.Join(tmpDir, "internal", "service", "foo.go"))
	})
}

func TestGeneratorAliases(t *testing.T) {
	tmpDir := chdir(t)

	// 每个别名对应的实际路径（与 defaultPath 一致）
	paths := map[string]string{
		"c": filepath.Join("controller", "foo.go"),
		"d": filepath.Join("internal", "domain", "foo.go"),
		"r": filepath.Join("internal", "repository", "foo.go"),
		"s": filepath.Join("internal", "service", "foo.go"),
		"t": filepath.Join("internal", "task", "foo.go"),
	}
	pkgNames := map[string]string{
		"c": "package controller",
		"d": "package domain",
		"r": "package repository",
		"s": "package service",
		"t": "package task",
	}

	for alias, relPath := range paths {
		t.Run(alias, func(t *testing.T) {
			_, _, err := executeCommandC(generator.GeneratorCmd, alias, "foo")
			if err != nil {
				t.Fatalf("command %q failed: %v", alias, err)
			}
			assertFileContains(t, filepath.Join(tmpDir, relPath), pkgNames[alias])
		})
	}
}

func TestGeneratorNoArgs(t *testing.T) {
	commands := []string{"service", "domain", "repository", "task", "controller"}
	for _, name := range commands {
		t.Run(name, func(t *testing.T) {
			_, _, err := executeCommandC(generator.GeneratorCmd, name)
			if err == nil {
				t.Errorf("expected error for %q with no args, got nil", name)
			}
		})
	}
}

func TestGeneratorCamelCaseConversion(t *testing.T) {
	tmpDir := chdir(t)

	tests := []struct {
		input    string
		expected string
	}{
		{"foo_bar", "FooBar"},
		{"my_service", "MyService"},
		{"simple", "Simple"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, _, err := executeCommandC(generator.GeneratorCmd, "service", tt.input)
			if err != nil {
				t.Fatalf("command failed: %v", err)
			}
			filePath := filepath.Join(tmpDir, "internal", "service", tt.input+".go")
			assertFileContains(t, filePath, tt.expected)
		})
	}
}

func TestGeneratorSkipExisting(t *testing.T) {
	tmpDir := chdir(t)
	filePath := filepath.Join(tmpDir, "internal", "service", "foo.go")

	// 第一次创建文件
	_, _, err := executeCommandC(generator.GeneratorCmd, "service", "foo")
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	assertFileExists(t, filePath)

	// 记录文件修改时间
	info1, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	// 再次执行相同命令（不带 --force），文件应被跳过且内容不变
	_, _, err = executeCommandC(generator.GeneratorCmd, "service", "foo")
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	info2, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}
	if info1.ModTime() != info2.ModTime() {
		t.Errorf("file should not be modified when skipped (modtime changed)")
	}

	// 带 --force 应覆盖
	_, _, err = executeCommandC(generator.GeneratorCmd, "service", "foo", "--force")
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	info3, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}
	if info2.ModTime() == info3.ModTime() {
		t.Errorf("file should be modified with --force (modtime unchanged)")
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist: %s", path)
	}
}

func assertFileContains(t *testing.T, path, substr string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}
	if !strings.Contains(string(content), substr) {
		t.Errorf("file %s does not contain %q, got:\n%s", path, substr, string(content))
	}
}
