package generator

import (
	"github.com/spf13/cobra"
)

var GeneratorCmd = &cobra.Command{
	Use:     "generator",
	Short:   "Generate go code",
	Aliases: []string{"gen"},
}

func init() {
	GeneratorCmd.AddCommand(controllerCmd)
	GeneratorCmd.AddCommand(serviceCmd)
	GeneratorCmd.AddCommand(repositoryCmd)
	GeneratorCmd.AddCommand(domainCmd)
}
