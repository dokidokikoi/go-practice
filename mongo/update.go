package mongo

import (
	"context"
	"fmt"
	"practice/mongo/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdatById(id string, user model.User) error {
	mId, _ := primitive.ObjectIDFromHex(id)
	// filter := bson.D{{"_id", mId}}
	update := bson.D{{"$set", user}}

	res, err := db.Collection("users").UpdateByID(context.TODO(), mId, update)
	if err != nil {
		return err
	}
	fmt.Printf("update result: %#v", res)
	return nil
}

func Update(condition model.User, user model.User) error {
	filter := condition
	update := bson.D{{"$set", user}}

	res, err := db.Collection("users").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	fmt.Printf("update result: %#v", res)
	return nil
}

func UpdateMany(condition model.User, user model.User) error {
	filter := condition
	update := bson.D{{"$set", user}}

	res, err := db.Collection("users").UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	fmt.Printf("update result: %#v", res)
	return nil
}
