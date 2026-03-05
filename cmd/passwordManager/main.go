package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/Krev3tka/ShrimPG/internal/storage"
	"github.com/Krev3tka/ShrimPG/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
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

	connStr := "postgres://shrimp_user:shrimp_password@localhost:5433/shrimp_vault"

	dbPool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	err = dbPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Database ping failed: %v\n", err)
	}

	fmt.Println("ShrimPG is swimming in Postgres Sea!")

	myStore := storage.DBStorage{Pool: dbPool}

	err = myStore.InitSchema()
	if err != nil {
		log.Fatalf("Error while initializing database: %v", err)
	}

	err = myStore.SeedCanary(masterKey)
	if err != nil {
		log.Fatalf("Error while initializing database: %v", err)
	}

	_, err = myStore.GetPassword("__MASTER_CHECK__", masterKey)
	if err != nil {
		log.Fatal("Master Password isn't right.\n")
	}

	switch args[0] {
	case "list":
		vault, err := myStore.GetAllPasswords(masterKey)
		if err != nil {
			log.Fatalf("Error while reading passwords from database: %v", err)
		}

		fmt.Print("Passwords list:\n\n")

		for service, passwd := range vault {
			fmt.Printf("%s: %s\n", service, passwd)
		}

		fmt.Println("\nDone.")

	case "create":
		if len(args) < 3 {
			fmt.Println("Bro, try to learn to count")
			return
		}

		if args[2] == "random" {
			err := myStore.SavePassword(args[1], storage.GenerateKoolPassword(KoolPasswordLength), masterKey)
			if err != nil {
				fmt.Printf("Something went wrong: %v", err)
				return
			}
			fmt.Println("Saved and encrypted")
			fmt.Print("Shrimp is saved. Password is too kool now. Press [s]ee, to see it: ")

			var choice string

			fmt.Scan(&choice)

			if choice == "s" {
				passwd, err := myStore.GetPassword(args[1], masterKey)
				if err != nil {
					fmt.Printf("Something went wrong: %v", err)
					return
				}
				fmt.Printf("Password: %s", passwd)
			}
			return
		}

		err := myStore.SavePassword(args[1], args[2], masterKey)
		if err != nil {
			fmt.Printf("Something went wrong: %v", err)
			return
		}
		fmt.Println("Saved and encrypted")
	case "catch":
		if len(args) == 1 {
			fmt.Println("Dude, you forgot to ask password. Please, enter name of service, for which you want to get the password")
			return
		}

		for i := 1; i < len(args); i++ {
			passwd, err := myStore.GetPassword(args[i], masterKey)
			if err != nil {
				fmt.Printf("Something went wrong: %v", err)
				continue
			}

			fmt.Printf("Password number %d: %s\n", i, passwd)
		}

	case "turf":
		if len(args) != 2 {
			fmt.Println("Please, recount arguments amount, maybe, you made a mistake")
			return
		}
		err := myStore.DeletePassword(args[1])
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
