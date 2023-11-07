package mongo

import (
	"context"
	"fmt"
	"practice/mongo/model"
)

func Insert(user model.User) error {
	res, err := db.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	fmt.Printf("insert result id: %#v", res.InsertedID)
	return nil
}

func InsertMany(users ...interface{}) error {
	res, err := db.Collection("users").InsertMany(context.TODO(), users)
	if err != nil {
		return err
	}

	fmt.Printf("many insert result ids: %#v", res.InsertedIDs...)
	return nil
}
