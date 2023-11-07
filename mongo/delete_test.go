package mongo

import (
	"practice/mongo/model"
	"testing"
)

func TestDelete(t *testing.T) {
	condition := model.User{
		Name: "jone",
	}
	if err := Delete(condition); err != nil {
		t.Errorf("delete error: %s", err)
	}
}

func TestDeleteMany(t *testing.T) {
	condition := model.User{
		Name: "jane",
	}
	if err := DeleteMany(condition); err != nil {
		t.Errorf("delete error: %s", err)
	}
}
