package main

// package doc https://pkg.go.dev/github.com/elastic/go-elasticsearch/esapi

import (
	"io/ioutil"

	"github.com/elastic/go-elasticsearch/v8"
)

var es *elasticsearch.Client

func init() {
	cert, _ := ioutil.ReadFile("http_ca.crt")
	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200",
		},
		Username: "elastic",
		Password: "ADR*piFezssmbUhhN8*S",
		CACert:   cert,
	}
	es, _ = elasticsearch.NewClient(cfg)
}

func main() {

}
