package main

import (
	"encoding/csv"
	_ "encoding/csv"
	"fmt"
	"math/rand"
	"os"
	_ "os"
	"time"
)

func GeneratePassword(length int, characters string) string {
	var password string
	for i := 0; i < length; i++ {
		n := rand.Intn(len(characters))
		password += string(characters[n])
	}
	return password
}

func PrintPasswords(password_keeper map[string]string) {
	fmt.Println("Все ваши пароли:")

	for key, value := range password_keeper {
		fmt.Printf("%s:\t%s\n", key, value)
	}

}

func WriteData(file *os.File, password_keeper map[string]string, service string, isFirst bool) {
	writer := csv.NewWriter(file)

	if isFirst {
		if err := writer.Write([]string{"Service", "Password"}); err != nil {
			fmt.Println(err)
		}
	}

	err := writer.Write([]string{service, password_keeper[service]})

	if err != nil {
		fmt.Println(err)
	}

	writer.Flush()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var password_keeper = make(map[string]string)

	file, err := os.Create("passwrods.csv")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	isFirst := true

	for {
		var length int
		var service string
		characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

		fmt.Print("Введите название сервиса, для которого вы хотите получить пароль: ")
		fmt.Scan(&service)

		fmt.Print("Введите длину генериуремого пароля: ")
		fmt.Scan(&length)

		if length <= 0 {
			fmt.Println("Длина пароля должна быть больше нуля")
			continue
		}

		var digits string
		var symbols string

		fmt.Print("Хотите добавить цифры? (0123456789) (y/n): ")
		fmt.Scan(&digits)

		if digits != "n" {
			characters += "0123456789"
		}

		fmt.Print("Хотите добавить спец.символы? (!@#$%^&*()_+-=[]{}|;:',.<>?/) (y/n): ")
		fmt.Scan(&symbols)

		if symbols != "n" {
			characters += "!@#$%^&*()_+-=[]{}|;:',.<>?/"
		}

		password_keeper[service] = GeneratePassword(length, characters)

		fmt.Printf("Пароль к сервису %s: %s\n", service, password_keeper[service])

		WriteData(file, password_keeper, service, isFirst)

		var fin string

		fmt.Print("Хотите продолжить? (y/n): ")
		fmt.Scan(&fin)

		if fin == "n" {
			break
		}

		isFirst = false
	}

	PrintPasswords(password_keeper)

}
