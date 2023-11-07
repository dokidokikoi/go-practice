package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type BaseModel = any

type searchResult[T BaseModel] struct {
	Hits struct {
		Total int `json:"total"`
		Hits  []struct {
			Source T `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

// 下面代码演示了go-elasticsearch提供的搜索查询功能，实际中查询请求体建议使用go template动态生成。
func TestSearch(t *testing.T) {
	a := assert.New(t)
	body := &bytes.Buffer{}
	content := map[string]interface{}{
		"_source": map[string]interface{}{
			"excludes": []string{"author"},
		},
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"sort": []interface{}{
			map[string]interface{}{
				"pages": map[string]interface{}{
					"order": "desc",
				},
			},
		},
		"from": 0,
		"size": 5,
	}
	if err := json.NewEncoder(body).Encode(content); err != nil {
		panic(err)
	}
	resp, err := es.Search(es.Search.WithIndex("book"), es.Search.WithBody(body))
	a.Nil(err)
	result, _ := io.ReadAll(resp.Body)
	t.Log(string(result))
	meta := searchResult[Book]{}
	json.Unmarshal(result, &meta)
	// bs, _ := json.Marshal(meta.Hits.Hits)
	// book := Book{}
	// json.Unmarshal(bs, &book)
	t.Log(meta)
}

func TestSearchString(t *testing.T) {
	a := assert.New(t)
	// body := &bytes.Buffer{}
	// content := map[string]interface{}{
	// 	"_source": map[string]interface{}{
	// 		"excludes": []string{"author"},
	// 	},
	// 	"query": map[string]interface{}{
	// 		"match_all": map[string]interface{}{},
	// 	},
	// 	"sort": []interface{}{
	// 		map[string]interface{}{
	// 			"pages": map[string]interface{}{
	// 				"order": "desc",
	// 			},
	// 		},
	// 	},
	// 	"from": 0,
	// 	"size": 5,
	// }
	query := `
		{
			"query": {
				"match_all": {
				}
			},
			"_source": {
				"excludes": ["author"]
			},
			"sort": [
				{
					"pages": {
						"order": "desc"
					}
				}
			],
			"from": 0,
			"size": 5
		}
	`
	query = fmt.Sprintf(`{"query":{"match":{"name":"%s"}},"sort": [{"version":{"order":"desc"}}],"from":0,"size":1}`, "test4_1")
	// if err := json.NewEncoder(body).Encode(content); err != nil {
	// 	panic(err)
	// }
	resp, err := es.Search(es.Search.WithIndex("metadata"), es.Search.WithBody(strings.NewReader(query)))
	a.Nil(err)
	result, _ := io.ReadAll(resp.Body)
	t.Log(string(result))
	// meta := searchResult[Book]{}
	// json.Unmarshal(result, &meta)
	// // bs, _ := json.Marshal(meta.Hits.Hits)
	// // book := Book{}
	// // json.Unmarshal(bs, &book)
	// t.Log(meta)
}
