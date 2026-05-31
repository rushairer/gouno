package generator

import "github.com/spf13/cobra"

var suiteCmd = &cobra.Command{
	Use:                   "suite [name]",
	Short:                 "Generate suite (domain, repository, service)",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		generateFile(cmd, args, "domain", defaultDomainPath, domainTemplate)
		generateFile(cmd, args, "repository", defaultRepositoryPath, repositoryTemplate)
		generateFile(cmd, args, "service", defaultServicePath, serviceTemplate)
	},
}

func init() {
	suiteCmd.Flags().BoolP("force", "f", false, "force overwrite")
}
