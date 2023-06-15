package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"zft/cmd"
	"zft/ui"
	"zft/utils"
)

var iscmd bool
var rootCmd = &cobra.Command{
	Use:   "zft",
	Short: "Zft is a command line tool for file operations",
	Run: func(cmd *cobra.Command, args []string) {
		iscmd = true
		ui.Run()
	},
}

func main() {
	iscmd = false
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
		if ui.Errors != nil {
			fmt.Println(ui.Errors)
		}
		if utils.FileURL != "" {
			fmt.Println("File uploaded successfully")
			fmt.Println("File URL copied to clipboard:", utils.FileURL)
		} else {
			if iscmd {
				fmt.Println("Operation canceled.")
			}
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(cmd.Confcmd)
}
