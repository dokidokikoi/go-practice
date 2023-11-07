package mongo

import (
	"context"
	"fmt"
	"practice/mongo/model"
)

func Delete(condition model.User) error {
	res, err := db.Collection("users").DeleteOne(context.TODO(), condition)
	if err != nil {
		return err
	}
	fmt.Printf("delete result: %#v", res)
	return nil
}

func DeleteMany(condition model.User) error {
	res, err := db.Collection("users").DeleteMany(context.TODO(), condition)
	if err != nil {
		return err
	}
	fmt.Printf("delete result: %#v", res)
	return nil
}
