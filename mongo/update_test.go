package mongo

import (
	"practice/mongo/model"
	"testing"
)

func TestUpdateById(t *testing.T) {
	newUser := model.User{
		Pwd: "abcdef",
	}
	if err := UpdatById("639ad311c8b0cd95b1639de8", newUser); err != nil {
		t.Errorf("update error: %s", err)
	}
}

func TestUpdate(t *testing.T) {
	condition := model.User{
		Name: "jone",
	}
	newUser := model.User{
		Pwd: "abcde",
	}
	if err := Update(condition, newUser); err != nil {
		t.Errorf("update error: %s", err)
	}
}

func TestUpdateMany(t *testing.T) {
	condition := model.User{
		Name: "jone",
	}
	newUser := model.User{
		Pwd: "abc",
	}
	if err := UpdateMany(condition, newUser); err != nil {
		t.Errorf("update error: %s", err)
	}
}
