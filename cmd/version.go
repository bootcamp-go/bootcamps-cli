package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCMD = &cobra.Command{
	Use:   "version",
	Short: "Get tool version",
	Long:  `Get the tool version`,
	Run: func(cmd *cobra.Command, args []string) {
		// get the version installed in the system and print it
		fmt.Println("Version: v1.4.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCMD)
}
