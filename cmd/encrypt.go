package cmd

import (
	"fmt"
	"io"
	"main/crypt"
	"os"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:     "encrypt",
	Short:   "encrypt files",
	Example: "  rcs encrypt -m off -p <password> -s <password2> potato",
	Run: func(cmd *cobra.Command, args []string) {
		cipher, err := newCipherForConfig(cmd)
		if err != nil {
			fmt.Printf("%+v\n", err)
			return
		}

		// encrypt data
		if len(args) == 0 {
			cmd.Help()
			return
		}

		progress, err := cmd.Flags().GetBool("no-progress")
		if err != nil {
			fmt.Printf("%+v\n", err)
			return
		}

		output, _ := cmd.Flags().GetString("output")

		for _, arg := range args {
			err := encrypt(cipher, arg, !progress, output)
			if err != nil {
				fmt.Printf("encrypt %s failed: %+v\n", arg, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	initCipherFlag(encryptCmd)
	encryptCmd.Flags().BoolP("no-progress", "", false, "disable progress bar")
	encryptCmd.Flags().StringP("output", "o", "", "output file location")

}

func encrypt(cipher *crypt.Cipher, f string, progress bool, dest string) error {
	filename := cipher.EncryptFileName(f)
	dest = destParser(dest, filename)

	fmt.Printf("Save to %s\n", dest)

	file, err := os.OpenFile(f, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	encrypter, err := cipher.EncryptData(file)
	if err != nil {
		return err
	}

	newFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer newFile.Close()

	var bar *pb.ProgressBar
	var writer io.Writer = newFile

	if progress {
		// get file size
		stat, err := os.Stat(f)
		if err != nil {
			return err
		}
		size := cipher.EncryptedSize(stat.Size())

		bar = pb.Full.Start64(0)
		bar.Set(pb.Bytes, true)
		bar.SetTotal(size)
		defer bar.Finish()

		writer = &writeCounter{bar: bar, offset: 0, writer: newFile}
	}

	_, err = io.Copy(writer, encrypter)
	if err != nil {
		return err
	}

	return nil
}
