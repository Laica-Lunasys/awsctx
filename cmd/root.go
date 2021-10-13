package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "awsctx",
	Short: "Switch AWS_PROFILE with saml2aws",
	Long:  `This application can be switch AWS_PROFILE as one-line.`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringP("mfa", "m", "", "MFA Token (Optional)")
}
