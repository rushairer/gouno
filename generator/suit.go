package generator

import "github.com/spf13/cobra"

var suitCmd = &cobra.Command{
	Use:                   "suit [name]",
	Short:                 "Generate suit (domain, repository, service)",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runDomain(cmd, args)
		runRepository(cmd, args)
		runService(cmd, args)
	},
}
