package cmd

import (
	"github.com/Laica-Lunasys/awsctx/service"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List aws accounts",
	Long:  `List available AWS accounts & account id.`,
	Run: func(cmd *cobra.Command, args []string) {
		p, err := service.GetSAML2AWS()
		if err != nil {
			panic(err)
		}
		token, err := cmd.Flags().GetString("mfa")
		if profiles, err := p.GetProfiles(); err != nil {
			panic(err)
		} else if len(profiles) >= 1 {
			if err := p.ListRoles(profiles[0],
				func() *service.MFA {
					if token != "" {
						return &service.MFA{Token: token}
					}
					return nil
				}(),
			); err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
