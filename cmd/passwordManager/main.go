package main

import (
	"awesomeProject3/internal/model"
	"awesomeProject3/internal/storage"
	"flag"
	"fmt"
	"strings"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Please, enter more than zero arguments.")
		return
	}

	switch args[0] {
	case "create":
		newEntry := make(model.Entry)

		if ans, ok := storage.IsYourPasswordCool(args[2]); ok == false {
			fmt.Println(ans)
			return
		} else {
			fmt.Println(ans)
		}

		if len(args) < 3 {
			fmt.Println("Please, recount arguments amount, maybe, you made a mistake")
			return
		} else if len(args) > 3 {
			fullPassword := strings.Join(args[2:], " ")
			newEntry[args[1]] = fullPassword
		} else {
			newEntry[args[1]] = args[2]
		}
		storage.Save(newEntry)
		fmt.Println("We could write entry into the file.")

	case "list":
		storage.PrintList("passwords.json")
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
