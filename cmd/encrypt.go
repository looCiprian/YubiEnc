package cmd

import (
	internalCommand "yubienc/internal/command"

	"github.com/spf13/cobra"
)

var (
	encryptSource      string
	encryptDestination string
	keyPathEncrypt     string
	threatEncrypt      int

	encryptCmd = &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt command will encrypt the <source> file/directory to the <destination> directory, will create <destination> if not exists",
		Long:  `Encrypt command will encrypt the <source> file/directory to the <destination> directory, will create <destination> if not exists`,
		RunE: func(cmd *cobra.Command, args []string) error {

			encryption := internalCommand.NewEncryption(encryptSource, encryptDestination, keyPathEncrypt, threatEncrypt)
			if err := encryption.ExecuteEncrypt(); err != nil {
				return err
			}

			return nil
		},
	}
)
