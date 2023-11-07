package model

type User struct {
	Name string `bson:"name,omitempty"`
	Pwd  string `bson:"pwd,omitempty"`
}
