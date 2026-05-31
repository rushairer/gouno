package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// GounoConfig 项目级 .gouno.yaml 配置
type GounoConfig struct {
	TemplateSet string `yaml:"template-set"`
}

const configFileName = ".gouno.yaml"
const templateDirName = ".gouno"
const templatesDirName = "templates"

// builtinTemplates 是内置的默认模板集，作为兜底
var builtinTemplates = map[string]string{
	"domain":     domainTemplate,
	"repository": repositoryTemplate,
	"service":    serviceTemplate,
	"controller": controllerTemplate,
	"task":       taskTemplate,
}

const domainTemplate = `package domain

import "context"

type %s struct {
}

func New%s() *%s {
	return &%s{}
}

func (d *%s) Foo(ctx context.Context) (bar string, err error) {
	return
}`

const repositoryTemplate = `package repository

import "context"

type %sRepository struct {
}

func New%sRepository() *%sRepository {
	return &%sRepository{}
}

func (r *%sRepository) Foo(ctx context.Context) (bar string, err error) {
	return
}`

const serviceTemplate = `package service

import "context"

type %sService struct {
}

func New%sService() *%sService {
	return &%sService{}
}

func (s *%sService) Foo(ctx context.Context) (bar string, err error) {
	return
}`

const controllerTemplate = `package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/rushairer/gouno"
)

type %sController struct {
}

func New%sController() *%sController {
	return &%sController{}
}

func (c *%sController) Foo(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gouno.NewSuccessResponse("bar"))
}`

const taskTemplate = `package task

import "context"

type %sTask struct {
}

func New%sTask() *%sTask {
	return &%sTask{}
}

func (t *%sTask) Run(ctx context.Context) error {
	return nil
}`

// resolveTemplateSet 确定要使用的模板集名称
// 优先级：--template-set flag > .gouno.yaml > "default"
func resolveTemplateSet(cmd *cobra.Command) string {
	// 1. 命令行 flag
	if flag := cmd.Flag("template-set"); flag != nil {
		if v := flag.Value.String(); v != "" {
			return v
		}
	}
	// 2. .gouno.yaml
	if cfg := loadProjectConfig(); cfg != nil && cfg.TemplateSet != "" {
		return cfg.TemplateSet
	}
	// 3. 默认
	return "default"
}

// loadProjectConfig 从当前目录加载 .gouno.yaml
func loadProjectConfig() *GounoConfig {
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	data, err := os.ReadFile(filepath.Join(cwd, configFileName))
	if err != nil {
		return nil
	}
	var cfg GounoConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil
	}
	return &cfg
}

// loadTemplate 加载指定模板集中的模板
// 搜索路径：
// 1. ~/.gouno/templates/<templateSet>/<typeName>.tmpl
// 2. 内置模板
func loadTemplate(cmd *cobra.Command, templateSet, typeName string) (string, error) {
	// 1. 用户模板目录
	homeDir, err := os.UserHomeDir()
	if err == nil {
		localPath := filepath.Join(homeDir, templateDirName, templatesDirName, templateSet, typeName+".tmpl")
		if content, err := os.ReadFile(localPath); err == nil {
			cmd.Printf("Using template: %s\n", localPath)
			return string(content), nil
		}
	}

	// 2. 内置模板（仅 default 模板集）
	if tmpl, ok := builtinTemplates[typeName]; ok {
		if templateSet == "default" || templateSet == "" {
			return tmpl, nil
		}
		return "", fmt.Errorf("template set %q not found (run: gouno-cli template install %s <url>)", templateSet, templateSet)
	}

	return "", fmt.Errorf("unknown template type: %s", typeName)
}

// templateSetDir 返回用户模板集的根目录 ~/.gouno/templates/
func templateSetDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, templateDirName, templatesDirName), nil
}
