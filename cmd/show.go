package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show profiles from ~/.saml2aws file.",
	Long:  `Show profiles list from ~/.saml2aws file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if p := os.Getenv("AWS_PROFILE"); p != "" {
			fmt.Println(p)
		} else {
			fmt.Println("default")
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
