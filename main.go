package main

import (
	"fmt"
	"zft/ui"
	"zft/utils"
)

func main() {
	ui.Run()
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
			fmt.Println("Operation canceled.")
		}
	}()
}
