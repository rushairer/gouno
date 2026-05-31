package generator

import "github.com/spf13/cobra"

var suiteCmd = &cobra.Command{
	Use:                   "suite [name]",
	Short:                 "Generate suite (domain, repository, service)",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := generateFile(cmd, args, "domain", defaultDomainPath, domainTemplate); err != nil {
			return err
		}
		if err := generateFile(cmd, args, "repository", defaultRepositoryPath, repositoryTemplate); err != nil {
			return err
		}
		return generateFile(cmd, args, "service", defaultServicePath, serviceTemplate)
	},
}

func init() {
	suiteCmd.Flags().BoolP("force", "f", false, "force overwrite")
}
