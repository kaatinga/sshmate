package main

import (
	"fmt"
	"os"

	"github.com/buger/goterm"
	"github.com/kaatinga/sshmate/internal/command"
	"github.com/spf13/cobra"

	"github.com/kaatinga/gocliselect"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an existing SSH key pair or a record about them",
	Run: func(cmd *cobra.Command, args []string) {
		keyPairs, err := command.GetKeyPairs()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		menu := gocliselect.NewMenu[int]("Select the key pair/record", goterm.BLUE)

		for index, keyPair := range keyPairs {
			menu.AddItem(keyPair.Host+" ("+keyPair.PrivateFile+")", index)
		}

		choice := menu.Display()

		// delete the record
		if err = command.DeleteKeyPair(keyPairs[choice]); err != nil {
			fmt.Println("unable to update the config file:", err)
			os.Exit(1)
		}

		fmt.Printf("The key pair '%s' has been deleted\n", keyPairs[choice].Host)
	},
}
