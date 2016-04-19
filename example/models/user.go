package models

import (
	"github.com/richgieg/lmaodb"
)

type User struct {
	ID        int64
	FirstName string
	LastName  string
	UserName  string
	Age       int
}

func init() {
	lmaodb.InitModel(User{})
}

func GetUser(id int64) (*User, error) {
	u := User{}
	err := lmaodb.GetRecord(id, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUsers() ([]User, error) {
	users := []User{}
	err := lmaodb.GetRecords(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func QueryUsers(field string, value interface{}) ([]User, error) {
	users := []User{}
	err := lmaodb.QueryRecords(&users, field, value)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func PutUsers(users []User) error {
	return lmaodb.PutRecords(users)
}

func SortUsers(users []User, field string, desc bool) {
	lmaodb.Sort(users, field, desc)
}

func (u *User) Delete() error {
	return lmaodb.DeleteRecord(u)
}

func (u *User) Put() error {
	return lmaodb.PutRecord(u)
}
