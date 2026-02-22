package main

import (
	"flag"
	"fmt"
	"github.com/Krev3tka/ShrimPG/internal/storage"
	"github.com/Krev3tka/ShrimPG/internal/utils"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Please, enter more than zero arguments.")
		return
	}

	masterKey := utils.GetMasterPassword()

	switch args[0] {
	case "list":
		vault, err := storage.Load(masterKey)
		if err != nil {
			fmt.Printf("Access Denied: %v\n", err)
			return
		}

		storage.PrintVault(vault)

	case "create":
		if len(args) < 3 {
			fmt.Println("Yo, Shrimp! I need a name and a password. Usage: create <name> <password>")
			return
		}

		vault, err := storage.Load(masterKey)

		if args[2] == "random" {
			vault[args[1]] = storage.GenerateKoolPassword(7)
			storage.Save(vault, masterKey)
			fmt.Println("Saved and encrypted")
			fmt.Print("Shrimp is saved. Password is too kool now. Press [s]ee, to see it: ")

			var choice string

			fmt.Scan(&choice)

			if choice == "s" {
				fmt.Printf("Password: %s", vault[args[1]])
			}
			return
		}

		if err != nil {
			fmt.Printf("Access Denied: %v\n", err)
			return
		}

		vault[args[1]] = args[2]

		storage.Save(vault, masterKey)
		fmt.Println("Saved and encrypted")

	case "delete":
		if len(args) != 2 {
			fmt.Println("Please, recount arguments amount, maybe, you made a mistake")
			return
		}
		storage.Delete(args[1])
	default:
		fmt.Println("Unknown command. Please, re-read what you entered.\n1) if you sealed up, just write command correctly\n2) otherwise, you can read manual with \"help\" command.")
	}
}
