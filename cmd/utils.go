package cmd

import (
	"errors"
	"fmt"
	"main/crypt"
	"main/obscure"
	"math/rand"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

type writeCounter struct {
	bar    *pb.ProgressBar
	offset int64
	writer *os.File
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n, err := wc.writer.WriteAt(p, wc.offset)
	if err != nil {
		return 0, err
	}
	wc.offset += int64(n)
	if wc.bar != nil {
		wc.bar.Add(n)
	}
	return n, nil
}

func initCipherFlag(c *cobra.Command) {
	c.Flags().StringP("filename-encoding", "e", "base32", "How to encode the encrypted filename to text string.")
	c.Flags().StringP("filename-encrypt-mode", "m", "standard", "How to encrypt the filenames, standard, off or obfuscated")
	c.Flags().StringP("password", "p", "", "Password or pass phrase for encryption.")
	c.Flags().StringP("salt", "s", "", "Password or pass phrase for salt, sometimes named password2.")
	c.Flags().BoolP("dirname-encrypt", "d", true, "Option to either encrypt directory names or leave them intact.")
}

// newCipherForConfig constructs a Cipher for the given config name
func newCipherForConfig(cmd *cobra.Command) (*crypt.Cipher, error) {
	filename_enc := cmd.Flag("filename-encoding").Value.String()
	password_raw := cmd.Flag("password").Value.String()
	salt_raw := cmd.Flag("salt").Value.String()
	filename_encrypt_mode := cmd.Flag("filename-encrypt-mode").Value.String()
	dirname_encrypt, _ := cmd.Flags().GetBool("dirname-encrypt")

	if filename_enc == "off" {
		dirname_encrypt = false
	}

	mode, err := crypt.NewNameEncryptionMode(filename_encrypt_mode)
	if err != nil {
		return nil, err
	}
	if password_raw == "" {
		return nil, errors.New("password not set in config file")
	}
	password, err := obscure.Reveal(password_raw)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password: %w", err)
	}
	var salt string
	if salt_raw != "" {
		salt, err = obscure.Reveal(salt_raw)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt password2(salt): %w", err)
		}
	}
	enc, err := crypt.NewNameEncoding(filename_enc)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("using password: %s\n", password)
	// fmt.Printf("using salt: %s\n", salt)
	cipher, err := crypt.NewCipher(mode, password, salt, dirname_encrypt, enc)
	if err != nil {
		return nil, fmt.Errorf("failed to make cipher: %w", err)
	}

	return cipher, nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytesRmndr(n int) string {
	rand.Seed(time.Hour.Microseconds())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func destParser(dest, filename string) string {
	if dest != "" {
		if IsDir(dest) {
			if dest[len(dest)-1] != '/' {
				dest = dest + "/"
			}
			dest = dest + filename
		}

		return dest
	}
	return filename
}
