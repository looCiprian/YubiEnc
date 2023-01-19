package cmd

import (
	internalCommand "yubienc/internal/command"

	"github.com/spf13/cobra"
)

var (
	keySource    string
	verifykeyCmd = &cobra.Command{
		Use:   "verifykey",
		Short: "Check if the encrypted keyset located in <source> can be decrypted using Yubikey",
		Long:  `Check if the encrypted keyset located in <source> can be decrypted using Yubikey`,
		RunE: func(cmd *cobra.Command, args []string) error {

			verifyKeyCommand := internalCommand.NewKeyVerify(keySource)
			if err := verifyKeyCommand.ExecuteKeyVerify(); err != nil {
				return err
			}

			return nil
		},
	}
)
