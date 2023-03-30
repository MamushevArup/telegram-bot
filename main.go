package main

import "fmt"

type User struct {
	ID    string `bson:"_id"`
	Name  string `bson:"telegram_tag"`
	Email string `bson:"email"`
}

func main() {
	main := "Медеуский район"
	var res string
	for _, r := range main {
		if r == ' ' {
			break
		}
		res += string(r)
	}
	fmt.Println(res)
}
