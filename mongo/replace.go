package mongo

import (
	"context"
	"fmt"
	"practice/mongo/model"
)

func Replace(condition model.User, user model.User) error {
	res, err := db.Collection("users").ReplaceOne(context.TODO(), condition, user)
	if err != nil {
		return err
	}
	fmt.Printf("replace result: %#v", res)
	return nil
}
