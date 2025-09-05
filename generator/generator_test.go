package generator_test

import (
	"bytes"
	"os"
	"path/filepath"
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

	return c, buf.String(), err
}

func TestGeneratorService(t *testing.T) {
	cmd := generator.GeneratorCmd

	executeCommandC(cmd, "service", "foo_bar")
	// 判断是否在 ./internal/service 目录下创建了foo文件，并删除
	serviceFilePath := filepath.Join("./internal/service", "foo_bar.go")
	if _, err := os.Stat(serviceFilePath); os.IsNotExist(err) {
		t.Errorf("Service file not created: %s", serviceFilePath)
	}
	executeCommandC(cmd, "service", "foo_bar", "--force")
	if err := os.Remove(serviceFilePath); err != nil {
		t.Errorf("Failed to remove service file: %v", err)
	}
	if err := os.RemoveAll("./internal"); err != nil {
		t.Errorf("Failed to remove internal directory: %v", err)
	}
	// 判断是否在 ./internal/service 目录下删除了foo文件
	if _, err := os.Stat(serviceFilePath); err == nil {
		t.Errorf("Service file not deleted: %s", serviceFilePath)
	}
	// 判断是否在当前目录下删除了internal目录
	if _, err := os.Stat("./internal"); err == nil {
		t.Errorf("Internal directory not deleted: %s", "./internal")
	}

	executeCommandC(cmd, "service", "foo_bar", "--path", "./custom/service")
	// 判断是否在 ./custom/service 目录下创建了foo文件，并删除
	serviceFilePath = filepath.Join("./custom/service", "foo_bar.go")
	if _, err := os.Stat(serviceFilePath); os.IsNotExist(err) {
		t.Errorf("Service file not created: %s", serviceFilePath)
	}
	if err := os.Remove(serviceFilePath); err != nil {
		t.Errorf("Failed to remove service file: %v", err)
	}
	if err := os.RemoveAll("./custom"); err != nil {
		t.Errorf("Failed to remove custom directory: %v", err)
	}
	// 判断是否在 ./custom/service 目录下删除了foo文件
	if _, err := os.Stat(serviceFilePath); err == nil {
		t.Errorf("Service file not deleted: %s", serviceFilePath)
	}
	// 判断是否在当前目录下删除了custom目录
	if _, err := os.Stat("./custom"); err == nil {
		t.Errorf("Custom directory not deleted: %s", "./custom")
	}

}

func TestGeneratorDomain(t *testing.T) {
	cmd := generator.GeneratorCmd

	executeCommandC(cmd, "domain", "foo")
	// 判断是否在 ./internal/domain 目录下创建了foo文件，并删除
	domainFilePath := filepath.Join("./internal/domain", "foo.go")
	if _, err := os.Stat(domainFilePath); os.IsNotExist(err) {
		t.Errorf("Domain file not created: %s", domainFilePath)
	}
	executeCommandC(cmd, "domain", "foo", "--force")
	if err := os.Remove(domainFilePath); err != nil {
		t.Errorf("Failed to remove domain file: %v", err)
	}
	if err := os.RemoveAll("./internal"); err != nil {
		t.Errorf("Failed to remove internal directory: %v", err)
	}
	// 判断是否在 ./internal/domain 目录下删除了foo文件
	if _, err := os.Stat(domainFilePath); err == nil {
		t.Errorf("Domain file not deleted: %s", domainFilePath)
	}
	// 判断是否在当前目录下删除了internal目录
	if _, err := os.Stat("./internal"); err == nil {
		t.Errorf("Internal directory not deleted: %s", "./internal")
	}

	executeCommandC(cmd, "domain", "foo", "--path", "./custom/domain")
	// 判断是否在 ./custom/domain 目录下创建了foo文件，并删除
	domainFilePath = filepath.Join("./custom/domain", "foo.go")
	if _, err := os.Stat(domainFilePath); os.IsNotExist(err) {
		t.Errorf("Domain file not created: %s", domainFilePath)
	}
	if err := os.Remove(domainFilePath); err != nil {
		t.Errorf("Failed to remove domain file: %v", err)
	}
	if err := os.RemoveAll("./custom"); err != nil {
		t.Errorf("Failed to remove custom directory: %v", err)
	}
	// 判断是否在 ./custom/domain 目录下删除了foo文件
	if _, err := os.Stat(domainFilePath); err == nil {
		t.Errorf("Domain file not deleted: %s", domainFilePath)
	}
	// 判断是否在当前目录下删除了custom目录
	if _, err := os.Stat("./custom"); err == nil {
		t.Errorf("Custom directory not deleted: %s", "./custom")
	}
}

func TestGeneratorRepository(t *testing.T) {
	cmd := generator.GeneratorCmd

	executeCommandC(cmd, "repository", "foo")
	// 判断是否在 ./internal/repository 目录下创建了foo文件，并删除
	repositoryFilePath := filepath.Join("./internal/repository", "foo.go")
	if _, err := os.Stat(repositoryFilePath); os.IsNotExist(err) {
		t.Errorf("Repository file not created: %s", repositoryFilePath)
	}
	executeCommandC(cmd, "repository", "foo", "--force")
	if err := os.Remove(repositoryFilePath); err != nil {
		t.Errorf("Failed to remove repository file: %v", err)
	}
	if err := os.RemoveAll("./internal"); err != nil {
		t.Errorf("Failed to remove internal directory: %v", err)
	}
	// 判断是否在 ./internal/repository 目录下删除了foo文件
	if _, err := os.Stat(repositoryFilePath); err == nil {
		t.Errorf("Repository file not deleted: %s", repositoryFilePath)
	}
	// 判断是否在当前目录下删除了internal目录
	if _, err := os.Stat("./internal"); err == nil {
		t.Errorf("Internal directory not deleted: %s", "./internal")
	}

	executeCommandC(cmd, "repository", "foo", "--path", "./custom/repository")
	// 判断是否在 ./custom/repository 目录下创建了foo文件，并删除
	repositoryFilePath = filepath.Join("./custom/repository", "foo.go")
	if _, err := os.Stat(repositoryFilePath); os.IsNotExist(err) {
		t.Errorf("Repository file not created: %s", repositoryFilePath)
	}
	if err := os.Remove(repositoryFilePath); err != nil {
		t.Errorf("Failed to remove repository file: %v", err)
	}
	if err := os.RemoveAll("./custom"); err != nil {
		t.Errorf("Failed to remove custom directory: %v", err)
	}
	// 判断是否在 ./custom/repository 目录下删除了foo文件
	if _, err := os.Stat(repositoryFilePath); err == nil {
		t.Errorf("Repository file not deleted: %s", repositoryFilePath)
	}
	// 判断是否在当前目录下删除了custom目录
	if _, err := os.Stat("./custom"); err == nil {
		t.Errorf("Custom directory not deleted: %s", "./custom")
	}
}

func TestGeneratorTask(t *testing.T) {
	cmd := generator.GeneratorCmd

	executeCommandC(cmd, "task", "foo")
	// 判断是否在 ./internal/task 目录下创建了foo文件，并删除
	taskFilePath := filepath.Join("./internal/task", "foo.go")
	if _, err := os.Stat(taskFilePath); os.IsNotExist(err) {
		t.Errorf("Task file not created: %s", taskFilePath)
	}
	executeCommandC(cmd, "task", "foo", "--force")
	if err := os.Remove(taskFilePath); err != nil {
		t.Errorf("Failed to remove task file: %v", err)
	}
	if err := os.RemoveAll("./internal"); err != nil {
		t.Errorf("Failed to remove internal directory: %v", err)
	}
	// 判断是否在 ./internal/task 目录下删除了foo文件
	if _, err := os.Stat(taskFilePath); err == nil {
		t.Errorf("Task file not deleted: %s", taskFilePath)
	}
	// 判断是否在当前目录下删除了internal目录
	if _, err := os.Stat("./internal"); err == nil {
		t.Errorf("Internal directory not deleted: %s", "./internal")
	}

	executeCommandC(cmd, "task", "foo", "--path", "./custom/task")
	// 判断是否在 ./custom/task 目录下创建了foo文件，并删除
	taskFilePath = filepath.Join("./custom/task", "foo.go")
	if _, err := os.Stat(taskFilePath); os.IsNotExist(err) {
		t.Errorf("Task file not created: %s", taskFilePath)
	}
	if err := os.Remove(taskFilePath); err != nil {
		t.Errorf("Failed to remove task file: %v", err)
	}
	if err := os.RemoveAll("./custom"); err != nil {
		t.Errorf("Failed to remove custom directory: %v", err)
	}
	// 判断是否在 ./custom/task 目录下删除了foo文件
	if _, err := os.Stat(taskFilePath); err == nil {
		t.Errorf("Task file not deleted: %s", taskFilePath)
	}
	// 判断是否在当前目录下删除了custom目录
	if _, err := os.Stat("./custom"); err == nil {
		t.Errorf("Custom directory not deleted: %s", "./custom")
	}
}

func TestGeneratorController(t *testing.T) {
	cmd := generator.GeneratorCmd

	executeCommandC(cmd, "controller", "foo")
	// 判断是否在 ./internal/controller 目录下创建了foo文件，并删除
	controllerFilePath := filepath.Join("./controller", "foo.go")
	if _, err := os.Stat(controllerFilePath); os.IsNotExist(err) {
		t.Errorf("Controller file not created: %s", controllerFilePath)
	}
	executeCommandC(cmd, "controller", "foo", "--force")
	if err := os.Remove(controllerFilePath); err != nil {
		t.Errorf("Failed to remove controller file: %v", err)
	}
	if err := os.RemoveAll("./controller"); err != nil {
		t.Errorf("Failed to remove controller directory: %v", err)
	}
	// 判断是否在 ./internal/controller 目录下删除了foo文件
	if _, err := os.Stat(controllerFilePath); err == nil {
		t.Errorf("Controller file not deleted: %s", controllerFilePath)
	}
	// 判断是否在 ./controller 目录下删除了controller目录
	if _, err := os.Stat("./controller"); err == nil {
		t.Errorf("Controller directory not deleted: %s", "./controller")
	}

	executeCommandC(cmd, "controller", "foo", "--path", "./custom/controller")
	// 判断是否在 ./custom/controller 目录下创建了foo文件，并删除
	controllerFilePath = filepath.Join("./custom/controller", "foo.go")
	if _, err := os.Stat(controllerFilePath); os.IsNotExist(err) {
		t.Errorf("Controller file not created: %s", controllerFilePath)
	}
	if err := os.Remove(controllerFilePath); err != nil {
		t.Errorf("Failed to remove controller file: %v", err)
	}
	if err := os.RemoveAll("./custom"); err != nil {
		t.Errorf("Failed to remove controller directory: %v", err)
	}
	// 判断是否在 ./custom/controller 目录下删除了foo文件
	if _, err := os.Stat(controllerFilePath); err == nil {
		t.Errorf("Controller file not deleted: %s", controllerFilePath)
	}
	// 判断是否在当前目录下删除了custom目录
	if _, err := os.Stat("./controller"); err == nil {
		t.Errorf("Controller directory not deleted: %s", "./controller")
	}
}

func TestGeneratorSuite(t *testing.T) {
	cmd := generator.GeneratorCmd

	executeCommandC(cmd, "suite", "foo")
	// 判断是否在 ./internal/repository 目录下创建了foo文件，并删除
	repositoryFilePath := filepath.Join("./internal/repository", "foo.go")
	if _, err := os.Stat(repositoryFilePath); os.IsNotExist(err) {
		t.Errorf("Repository file not created: %s", repositoryFilePath)
	}
	serviceFilePath := filepath.Join("./internal/service", "foo.go")
	if _, err := os.Stat(serviceFilePath); os.IsNotExist(err) {
		t.Errorf("Service file not created: %s", serviceFilePath)
	}
	domainFilePath := filepath.Join("./internal/domain", "foo.go")
	if _, err := os.Stat(domainFilePath); os.IsNotExist(err) {
		t.Errorf("Domain file not created: %s", domainFilePath)
	}

	if err := os.RemoveAll("./internal"); err != nil {
		t.Errorf("Failed to remove repository directory: %v", err)
	}
	// 判断是否在 ./internal/repository 目录下删除了foo文件
	if _, err := os.Stat(repositoryFilePath); err == nil {
		t.Errorf("Repository file not deleted: %s", repositoryFilePath)
	}
	// 判断是否在 ./internal/service 目录下删除了foo文件
	if _, err := os.Stat(serviceFilePath); err == nil {
		t.Errorf("Service file not deleted: %s", serviceFilePath)
	}
	// 判断是否在 ./internal/domain 目录下删除了foo文件
	if _, err := os.Stat(domainFilePath); err == nil {
		t.Errorf("Domain file not deleted: %s", domainFilePath)
	}
	// 判断是否在当前目录下删除了internal目录
	if _, err := os.Stat("./internal"); err == nil {
		t.Errorf("Repository directory not deleted: %s", "./internal")
	}
}
