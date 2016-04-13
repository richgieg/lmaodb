package models

import (
	"github.com/richgieg/lmaodb"
)

type Post struct {
	ID     int64
	UserID int64
	Text   string
}

func init() {
	lmaodb.InitModel(Post{})
}

func GetPost(id int64) (*Post, error) {
	p := Post{}
	err := lmaodb.GetRecord(id, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetPosts() ([]Post, error) {
	posts := []Post{}
	err := lmaodb.GetRecords(&posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func PutPosts(posts []Post) error {
	return lmaodb.PutRecords(posts)
}

func (p *Post) Delete() error {
	return lmaodb.DeleteRecord(p)
}

func (p *Post) Put() error {
	return lmaodb.PutRecord(p)
}
