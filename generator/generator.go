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
	GeneratorCmd.AddCommand(
		controllerCmd,
		serviceCmd,
		repositoryCmd,
		domainCmd,
		suiteCmd,
		taskCmd,
	)
}
