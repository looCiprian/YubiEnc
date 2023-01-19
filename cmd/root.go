package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "YubiEnc",
		Short: "YubiEnc allows you to encrypt/decrypt a directory using your Yubikey",
	}
)

func init() {

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(versionCmd)

	rootCmd.AddCommand(newkeyCmd)
	newkeyCmd.Flags().StringVarP(&keyDestination, "destination", "d", "", "Directory to store the new keyset")
	newkeyCmd.MarkFlagRequired("destination")

	rootCmd.AddCommand(verifykeyCmd)
	verifykeyCmd.Flags().StringVarP(&keySource, "source", "s", "", "Encrypted keyset file")
	verifykeyCmd.MarkFlagRequired("source")

	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().StringVarP(&encryptSource, "source", "s", "", "Directory to encrypt")
	encryptCmd.Flags().StringVarP(&encryptDestination, "destination", "d", "", "Directory to store encrypted files")
	encryptCmd.Flags().StringVarP(&keyPathEncrypt, "key", "k", "", "Encrypted keyset")
	encryptCmd.Flags().IntVarP(&threatEncrypt, "thread", "t", 1, "Threat")
	encryptCmd.MarkFlagRequired("source")
	encryptCmd.MarkFlagRequired("destination")
	encryptCmd.MarkFlagRequired("key")

	rootCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().StringVarP(&decryptSource, "source", "s", "", "Directory containing encrypted files")
	decryptCmd.Flags().StringVarP(&decryptDestination, "destination", "d", "", "Directory to store decrypted files")
	decryptCmd.Flags().StringVarP(&keyPathDecrypt, "key", "k", "", "Encrypted keyset")
	decryptCmd.Flags().IntVarP(&threatDecrypt, "thread", "t", 1, "Threat")
	decryptCmd.MarkFlagRequired("source")
	decryptCmd.MarkFlagRequired("destination")
	decryptCmd.MarkFlagRequired("key")

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		//fmt.Fprintln(os.Stderr, err)
		//os.Exit(1)
	}
}
