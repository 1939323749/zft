package main

import (
	"fmt"
	"github.com/1939323749/zft/ui"
	"github.com/1939323749/zft/utils"
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
			fmt.Println("zft-v0.0.2")
			fmt.Println("File uploaded successfully")
			fmt.Println("File URL copied to clipboard:", utils.FileURL)
		} else {
			fmt.Println("Operation canceled. ")
		}
	}()
}
