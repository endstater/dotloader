package main

import (
	"os"
	"fmt"
)

/*
TODO:
rsync options in config
git sync
improve project repo
*/

func main() {
	dotloader, err := NewDotloader()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	
	args := os.Args
	if len(args) < 2 {	
		dotloader.Help()
		return
	}

	switch args[1] {
	case "sync", "s":
		err = dotloader.Sync()
	case "load", "l":
		err = dotloader.Load()
	case "help", "--help", "-h", "man", "h":
		dotloader.Help()
	case "vesion", "v":
		fmt.Printf("\nDotloader\nSimple dotfiles manager working with rsync and git written in go.\nVersion: 0.1\nAuthor: endstater\n")
	default:
		fmt.Println("Unknown option:", args[1])
	}
	if err != nil {
		fmt.Printf("%v", err)
	}
}
