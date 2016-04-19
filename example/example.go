package main

import (
	"fmt"

	"github.com/richgieg/lmaodb/example/models"
)

func main() {
	u := models.User{
		FirstName: "Fred",
		LastName:  "Savage",
		UserName:  "fsavage",
		Age:       40,
	}
	u.Put()

	p := models.Post{
		UserID: u.ID,
		Text:   "Hello, world!",
	}
	p.Put()

	fmt.Println(models.GetUsers())
	fmt.Println(models.GetPosts())
	fmt.Println(models.QueryPosts("Text", "Hello, world!"))

	p.Text = "Goodbye, world!"
	p.Put()
	fmt.Println(models.GetPosts())

	p.Delete()
	fmt.Println(models.GetPosts())

	u = models.User{
		FirstName: "Sarah",
		LastName:  "Smith",
		UserName:  "ssmith",
		Age:       27,
	}
	u.Put()

	fmt.Println(models.QueryUsers("FirstName", "Fred"))
	fmt.Println(models.QueryUsers("Age", 40))
	fmt.Println(models.QueryUsers("ID", int64(1)))

	users, _ := models.GetUsers()
	fmt.Println(users)
	models.SortUsers(users, "LastName", true)
	fmt.Println(users)
	models.SortUsers(users, "FirstName", false)
	fmt.Println(users)
	models.SortUsers(users, "Age", false)
	fmt.Println(users)
}
