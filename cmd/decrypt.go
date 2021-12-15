package cmd

import (
	"fmt"
	"io"
	"main/crypt"
	"os"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:     "decrypt",
	Short:   "decrypt files",
	Example: "  rcs decrypt -m off -p <password> -s <password2> potato",
	Run: func(cmd *cobra.Command, args []string) {
		cipher, err := newCipherForConfig(cmd)
		if err != nil {
			fmt.Printf("%+v\n", err)
			return
		}

		// decrypt data
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
			err := decrypt(cipher, arg, !progress, output)
			if err != nil {
				fmt.Printf("decrypt %s failed: %+v\n", arg, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)
	initCipherFlag(decryptCmd)
	decryptCmd.Flags().BoolP("no-progress", "", false, "disable progress bar")
	decryptCmd.Flags().StringP("output", "o", "", "output file location")

}

func decrypt(cipher *crypt.Cipher, f string, progress bool, dest string) error {
	filename, err := cipher.DecryptFileName(f)
	if err != nil {
		fmt.Printf("Warning: error decrypting filename: %+v\n", err)
		// generate random string
		filename = RandStringBytesRmndr(16)
	}
	dest = destParser(dest, filename)

	fmt.Printf("Save to %s\n", dest)

	file, err := os.OpenFile(f, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	decrypter, err := cipher.DecryptData(file)
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
		size, err := cipher.DecryptedSize(stat.Size())
		if err != nil {
			return err
		}

		bar = pb.Full.Start64(0)
		bar.Set(pb.Bytes, true)
		bar.SetTotal(size)
		defer bar.Finish()

		writer = &writeCounter{bar: bar, offset: 0, writer: newFile}
	}

	_, err = io.Copy(writer, decrypter)
	if err != nil {
		return err
	}

	return nil
}
