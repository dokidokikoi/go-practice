package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndex(t *testing.T) {
	a := assert.New(t)
	response, err := es.Indices.Create("book", es.Indices.Create.WithBody(strings.NewReader(`
	{
		"aliases": {
			"book":{}
		},
		"settings": {
			"analysis": {
				"normalizer": {
					"lowercase": {
						"type": "custom",
						"char_filter": [],
						"filter": ["lowercase"]
					}
				}
			}
		},
		"mappings": {
			"properties": {
				"name": {
					"type": "keyword",
					"normalizer": "lowercase"
				},
				"price": {
					"type": "double"
				},
				"summary": {
					"type": "text"
				},
				"author": {
					"type": "keyword"
				},
				"pubDate": {
					"type": "date"
				},
				"pages": {
					"type": "integer"
				}
			}
		}
	}
	`)))
	a.Nil(err)
	t.Log(response)
}

func TestGetIndex(t *testing.T) {
	a := assert.New(t)
	response, err := es.Indices.Get([]string{"book"})
	a.Nil(err)
	t.Log(response)
}

func TestDeleteIndex(t *testing.T) {
	a := assert.New(t)
	response, err := es.Indices.Delete([]string{"book-0.1.0"})
	a.Nil(err)
	t.Log(response)
}
