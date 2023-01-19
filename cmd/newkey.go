package cmd

import (
	internalCommand "yubienc/internal/command"

	"github.com/spf13/cobra"
)

var (
	keyDestination string
	newkeyCmd      = &cobra.Command{
		Use:   "newkey",
		Short: "Create a new key and save it encrypted in <destination>",
		Long:  `Create a new key and save it encrypted in <destination>`,
		RunE: func(cmd *cobra.Command, args []string) error {

			newKeyCommand := internalCommand.NewNewkey(keyDestination)
			if err := newKeyCommand.ExecuteNewkey(); err != nil {
				return err
			}

			return nil
		},
	}
)
