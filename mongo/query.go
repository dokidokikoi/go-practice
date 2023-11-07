package mongo

import (
	"context"
	"practice/mongo/model"
)

func FindOne(condition model.User) (model.User, error) {
	var user model.User
	err := db.Collection("users").FindOne(context.TODO(), condition).Decode(&user)
	return user, err
}

func Find(condition model.User) (model.User, error) {
	var user model.User
	err := db.Collection("users").FindOne(context.TODO(), condition).Decode(&user)
	return user, err
}
