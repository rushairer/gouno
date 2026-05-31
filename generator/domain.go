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
		return generateFile(cmd, args, "domain", defaultDomainPath)
	},
}

var defaultDomainPath = filepath.Join("internal", "domain")

func init() {
	domainCmd.Flags().StringP("path", "p", defaultDomainPath, "path to domain")
	domainCmd.Flags().BoolP("force", "f", false, "force overwrite")
	domainCmd.Flags().String("template-set", "", "template set name")
}
