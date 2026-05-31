package generator

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:                   "domain [name]",
	Short:                 "Generate domain",
	Aliases:               []string{"d"},
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateFile(cmd, args, "domain", defaultDomainPath, domainTemplate)
	},
}

var defaultDomainPath = filepath.Join("internal", "domain")

func init() {
	domainCmd.Flags().StringP("path", "p", defaultDomainPath, "path to domain")
	domainCmd.Flags().BoolP("force", "f", false, "force overwrite")
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
