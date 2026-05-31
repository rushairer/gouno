package generator

import (
	"github.com/spf13/cobra"
)

// GeneratorCmd is the root Cobra command for the code generator.
// It provides subcommands to scaffold DDD layers: domain, repository, service,
// controller, task, and suite (all three domain layers at once).
// Aliases: "gen".
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
