package cmd

import (
	"fmt"
	"os"

	"github.com/Laica-Lunasys/awsctx/service"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Logging into AWS Account",
	Long:  `Logging into AWS Account`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		p, err := service.GetSAML2AWS()
		if err != nil {
			panic(err)
		}
		profiles, err := p.GetProfiles()
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveError
		}
		return profiles, cobra.ShellCompDirectiveDefault
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Invalid args")
			os.Exit(1)
		}
		p, err := service.GetSAML2AWS()
		if err != nil {
			panic(err)
		}

		// Console
		console, err := cmd.Flags().GetBool("console")
		if err != nil {
			panic(err)
		}

		// Link
		link, err := cmd.Flags().GetBool("link")
		if err != nil {
			panic(err)
		}

		// Firefox
		firefox, err := cmd.Flags().GetBool("firefox")
		if err != nil {
			panic(err)
		}

		// MFA
		token, err := cmd.Flags().GetString("mfa")
		if err != nil {
			panic(err)
		}
		if err := p.Login(args[0],
			&service.LoginOption{
				Console:  console,
				LinkOnly: link,
				Firefox:  firefox,
			},
			func() *service.MFA {
				if token != "" {
					return &service.MFA{Token: token}
				}
				return nil
			}(),
		); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().BoolP("console", "c", false, "Open AWS Management Console (browser)")
	loginCmd.Flags().BoolP("firefox", "F", false, "Open as Firefox")
	loginCmd.Flags().BoolP("link", "l", false, "Present link to AWS console instead of opening browser")
}
