package generator

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:                   "service [name]",
	Short:                 "Generate service",
	Aliases:               []string{"s"},
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		generateFile(cmd, args, "service", defaultServicePath, serviceTemplate)
	},
}

var defaultServicePath = filepath.Join("internal", "service")

func init() {
	serviceCmd.Flags().StringP("path", "p", defaultServicePath, "path to service")
	serviceCmd.Flags().BoolP("force", "f", false, "force overwrite")
}

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
