package generator

import "github.com/spf13/cobra"

var suiteCmd = &cobra.Command{
	Use:                   "suite [name]",
	Short:                 "Generate suite (domain, repository, service)",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runDomain(cmd, args)
		runRepository(cmd, args)
		runService(cmd, args)
	},
}

func init() {
	suiteCmd.Flags().BoolP("force", "f", false, "force overwrite")
}
