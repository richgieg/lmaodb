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

	p.Text = "Goodbye, world!"
	p.Put()
	fmt.Println(models.GetPosts())

	p.Delete()
	fmt.Println(models.GetPosts())
}
