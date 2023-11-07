package main

import (
	"log"

	"github.com/gocolly/colly/v2"
)

func main() {
	// create a new collector
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	// authenticate
	err := c.Post("https://bangumi.tv/login", map[string]string{"email": "h18673628658@outlook.com", "password": "DJSW20070831miku"})
	if err != nil {
		log.Fatal(err)
	}

	// attach callbacks after login
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
	})

	c.OnHTML("html", func(r *colly.HTMLElement) {
		log.Println("response received", r.Text)
	})

	// start scraping
	c.Visit("https://bangumi.tv/game/tag/R18")
}
