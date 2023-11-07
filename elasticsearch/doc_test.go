package main

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Book struct {
	ID      string     `json:"id,omitempty"`
	Author  string     `json:"author,omitempty"`
	Name    string     `json:"name,omitempty"`
	Pages   int        `json:"pages,omitempty"`
	Price   float64    `json:"price,omitempty"`
	PubDate *time.Time `json:"pubDate,omitempty"`
	Summary string     `json:"summary,omitempty"`
}

func TestCreateDocument(t *testing.T) {
	a := assert.New(t)
	body := &bytes.Buffer{}
	pubDate := time.Now()
	err := json.NewEncoder(body).Encode(&Book{
		Author:  "金庸",
		Price:   96.0,
		Name:    "天龙八部",
		Pages:   1978,
		PubDate: &pubDate,
		Summary: "...",
	})
	a.Nil(err)
	response, err := es.Create("book", "10001", body)
	a.Nil(err)
	t.Log(response)
}

func TestUpdateDocument(t *testing.T) {
	a := assert.New(t)
	body := &bytes.Buffer{}
	pubDate := time.Now()
	err := json.NewEncoder(body).Encode(&Book{
		Author:  "金庸",
		Price:   96.0,
		Name:    "天龙八部",
		Pages:   1978,
		PubDate: &pubDate,
		Summary: "......",
	})
	a.Nil(err)
	response, err := es.Index("book", body, es.Index.WithDocumentID("10001"))
	a.Nil(err)
	t.Log(response)
}

func TestDeleteDocument(t *testing.T) {
	a := assert.New(t)
	response, err := es.Delete("book", "10001")
	a.Nil(err)
	t.Log(response)
}

type Hit struct {
	Source map[string]interface{} `json:"_source"`
}

func TestGetDocument(t *testing.T) {
	a := assert.New(t)
	resp, err := es.Get("book", "10001")
	a.Nil(err)
	result, _ := io.ReadAll(resp.Body)
	t.Log(string(result))
	meta := Hit{}
	json.Unmarshal(result, &meta)
	bs, _ := json.Marshal(meta.Source)
	book := Book{}
	json.Unmarshal(bs, &book)
	t.Log(book)
}
