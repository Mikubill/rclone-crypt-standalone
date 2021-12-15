package cmd

import (
	"fmt"
	"main/obscure"

	"github.com/spf13/cobra"
)

// obscureCmd represents the obscure command
var obscureCmd = &cobra.Command{
	Use:     "obscure",
	Short:   "Obscure password for use in somewhere.",
	Example: "  rcs obscure potato",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		for _, arg := range args {
			got, err := obscure.Obscure(arg)
			if err != nil {
				fmt.Printf("%+v\n", err)
			} else {
				fmt.Printf("%s", got)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(obscureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// obscureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// obscureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
