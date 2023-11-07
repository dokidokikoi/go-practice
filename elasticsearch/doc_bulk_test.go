package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// POST /<target>/_bulk
// { "index" : { "_id" : "1" } }
// { "field1" : "value1" }
// { "delete" : { "_id" : "2" } }
// { "create" : { "_id" : "3" } }
// { "field1" : "value3" }
// { "update" : {"_id" : "1" } }
// { "doc" : {"field2" : "value2"} }

func TestBulkCreateAndDelete(t *testing.T) {
	createBooks := []*Book{
		{
			ID:     "10002",
			Name:   "神雕侠侣",
			Author: "金庸",
		},
		{
			ID:     "10003",
			Name:   "射雕英雄传",
			Author: "金庸",
		},
	}

	updateBooks := []*Book{
		{
			ID:     "10002",
			Name:   "三少爷的剑",
			Author: "古龙",
		},
	}
	deleteBookIds := []string{"10001"}

	a := assert.New(t)
	body := &bytes.Buffer{}
	for _, book := range createBooks {
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%s" } }%s`, book.ID, "\n"))
		data, err := json.Marshal(book)
		a.Nil(err)
		data = append(data, "\n"...)
		body.Grow(len(meta) + len(data))
		body.Write(meta)
		body.Write(data)
	}
	for _, id := range deleteBookIds {
		meta := []byte(fmt.Sprintf(`{ "delete" : { "_id" : "%s" } }%s`, id, "\n"))
		body.Grow(len(meta))
		body.Write(meta)
	}
	for _, book := range updateBooks {
		meta := []byte(fmt.Sprintf(`{ "update" : { "_id" : "%s" } }%s`, book.ID, "\n"))
		var v = &struct {
			Doc Book `json:"doc"`
		}{Doc: *book}
		data, err := json.Marshal(v)
		a.Nil(err)
		data = append(data, "\n"...)
		body.Grow(len(meta) + len(data))
		body.Write(meta)
		body.Write(data)
	}
	t.Log(body.String())

	response, err := es.Bulk(body, es.Bulk.WithIndex("book"))
	a.Nil(err)
	t.Log(response)
}
