package main

import (
	"flag"
	"fmt"

	"github.com/Krev3tka/ShrimPG/internal/storage"
	"github.com/Krev3tka/ShrimPG/internal/utils"
)

const (
	KoolPasswordLength = 4
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
			fmt.Println("Bro, try to learn to count")
			return
		}

		vault, err := storage.Load(masterKey)
		if err != nil {
			fmt.Printf("Something went wrong: %v", err)
			return
		}

		if args[2] == "random" {
			vault[args[1]] = storage.GenerateKoolPassword(KoolPasswordLength)
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
	case "catch":
		if len(args) == 1 {
			fmt.Println("Dude, you forgot to ask password. Please, enter name of service, for which you want to get the password")
			return
		}

		vault, err := storage.Load(masterKey)
		if err != nil {
			fmt.Printf("Something went wrong: %v", err)
			return
		}

		for i := 1; i < len(args); i++ {
			if _, ok := vault[args[i]]; !ok {
				fmt.Printf("Password for service '%s' isn't found\n", args[i])
				continue
			}
			fmt.Printf("Password number %d: %s\n", i, vault[args[i]])
		}

	case "turf":
		if len(args) != 2 {
			fmt.Println("Please, recount arguments amount, maybe, you made a mistake")
			return
		}
		err := storage.Delete(args[1], masterKey)
		if err != nil {
			fmt.Printf("Error while turfing: %v", err)
		}
	case "help":
		if len(args) > 1 {
			fmt.Println("Even for helping you need help...\nRecount arguments or rewrite command.")
			return
		}

		fmt.Println("\nCommand <usage> <...>")
		fmt.Println("create <name> <password> - for creating password / create <name> random - for random passphrase")
		fmt.Println("catch <name> - for getting password by service name")
		fmt.Println("turf <name> - for deleting password for service")
		fmt.Print("list - for printing all your passwords\n\n")

	default:
		fmt.Println("Unknown command. Please, re-read what you entered.\n1) if you sealed up, just write command correctly\n2) otherwise, you can read manual with \"help\" command.")
	}
}
