package mongo

import (
	"practice/mongo/model"
	"testing"
)

func TestInsert(t *testing.T) {
	newUser := model.User{
		Name: "jane",
		Pwd:  "123",
	}
	if err := Insert(newUser); err != nil {
		t.Errorf("insert error: %s", err)
	}
}

func TestInsertMany(t *testing.T) {
	newUsers := []interface{}{
		model.User{Name: "jone", Pwd: "123456"},
		model.User{Name: "will", Pwd: "123456"},
		model.User{Name: "alan", Pwd: "123456"},
	}
	if err := InsertMany(newUsers...); err != nil {
		t.Errorf("insert many error: %s", err)
	}
}
