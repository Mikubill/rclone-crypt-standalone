/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rcs",
	Short: "A brief description of your application",
	Long:  "An encryption and decryption tool based on rclone crypt and nacl/secretbox. It can be used to encrypt/decrypt files encrypted by rclone, or you can generate your own encrypted file and decrypt it on rclone.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
