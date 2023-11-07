package mongo

import (
	"fmt"
	"practice/mongo/model"
	"testing"
)

func TestFindOne(t *testing.T) {
	condition := model.User{
		Pwd: "abc",
	}
	user, err := FindOne(condition)
	if err != nil {
		t.Errorf("FindOne error: %#v", err)
	}
	fmt.Printf("user: %#v", user)
}
