package mongo

import (
	"practice/mongo/model"
	"testing"
)

func TestReplace(t *testing.T) {
	condition := model.User{
		Name: "jack",
	}
	newUser := model.User{
		Name: "jack o",
	}
	if err := Replace(condition, newUser); err != nil {
		t.Errorf("replace error: %s", err)
	}
}
