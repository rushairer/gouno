package generator

import (
	"fmt"
)

// builtinTemplates 是内置的代码生成模板
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

// loadTemplate 加载内置的代码生成模板
func loadTemplate(typeName string) (string, error) {
	if tmpl, ok := builtinTemplates[typeName]; ok {
		return tmpl, nil
	}
	return "", fmt.Errorf("unknown template type: %s", typeName)
}

