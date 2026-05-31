package generator

import "github.com/spf13/cobra"

var suiteCmd = &cobra.Command{
	Use:                   "suite [name]",
	Short:                 "Generate suite (domain, repository, service)",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := generateFile(cmd, args, "domain", defaultDomainPath); err != nil {
			return err
		}
		if err := generateFile(cmd, args, "repository", defaultRepositoryPath); err != nil {
			return err
		}
		return generateFile(cmd, args, "service", defaultServicePath)
	},
}

func init() {
	suiteCmd.Flags().BoolP("force", "f", false, "force overwrite")
	suiteCmd.Flags().String("template-set", "", "template set name")
}
