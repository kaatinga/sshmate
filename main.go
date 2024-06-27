package main

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kaatinga/sshmate/internal/command"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sshmate",
	Short: "CLI tool for managing SSH key pairs and server connections",
	Long: `sshmate is a CLI application for managing SSH key pairs and server connections.
	It allows adding and deleting SSH key pairs, as well as checking connection to added servers.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to sshmate CLI!")
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new SSH key pair",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TBI")
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an existing SSH key pair",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TBI")
	},
}

var listKeysCmd = &cobra.Command{
	Use:   "list",
	Short: "List all added SSH key pairs",
	Run: func(cmd *cobra.Command, args []string) {
		keyPairs, err := command.GetKeyPairs()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		// print key pairs
		t := table.NewWriter()
		t.SetStyle(table.StyleColoredBright)
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Host", "Available", "Private Key", "Keep Public Key", "Type"})
		for _, keyPair := range keyPairs {
			t.AppendRows([]table.Row{
				{keyPair.Host, keyPair.Available, keyPair.PrivateFile, keyPair.PublicFile != "", keyPair.Type},
			})
		}
		t.Render()
	},
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check connection to every added server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TBI")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(listKeysCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
