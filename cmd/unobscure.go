package cmd

import (
	"fmt"
	"main/obscure"

	"github.com/spf13/cobra"
)

// unobscureCmd represents the unobscure command
var unobscureCmd = &cobra.Command{
	Use:     "unobscure",
	Short:   "Unobscure password for use in somewhere.",
	Example: "  rcs unobscure E9RDEKki5MKohzOGuga7wyGsPy0",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		for _, arg := range args {
			got, err := obscure.Reveal(arg)
			if err != nil {
				fmt.Printf("%+v\n", err)
			} else {
				fmt.Printf("%s", got)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(unobscureCmd)
}
