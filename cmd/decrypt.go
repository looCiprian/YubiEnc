package cmd

import (
	internalCommand "yubienc/internal/command"

	"github.com/spf13/cobra"
)

var (
	decryptSource      string
	decryptDestination string
	keyPathDecrypt     string
	threatDecrypt      int

	decryptCmd = &cobra.Command{
		Use:   "decrypt",
		Short: "Decrypt command will decrypt the <source> file/directory to the <destination> directory",
		Long:  `Decrypt command will decrypt the <source> file/directory to the <destination> directory`,
		RunE: func(cmd *cobra.Command, args []string) error {

			decryptCommand := internalCommand.NewDecryption(decryptSource, decryptDestination, keyPathDecrypt, threatDecrypt)
			if err := decryptCommand.ExecuteDecrypt(); err != nil {
				return err
			}

			return nil
		},
	}
)
