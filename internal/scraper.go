/*

I want to scrape Flashback crime section for information about newly committed crimes.
Each section has its own board. I think for now, what I want to do is:
-> Scrape the title, date, metadata of each article.
-> Save the data to a SQL table.
-> Clean the data and try to work out if crime, data + region yields any results ie a heatmap of crimes/type
-> Write a blog post about findings
-> Some point it would be nice to plot the data with dates to a map and maybe geolocate any images.

For now, it's probably safest to avoid scraping the actual boards. Some boards are thousands of pages long and I don't know exactly what I want to do with this information other than save potential links to news paper articles.
Saving articles/links might be useful as a way to triangulate data points and create an instance of "adverse media" that can be cross-referenced against any potential court documents.
I think this would be too much work, for now.

*/

package internal

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gocolly/colly"
	_ "github.com/mattn/go-sqlite3"
)

func ScrapeFlashBack() {
	var db *sql.DB
	db = OpenDB()
	CreateDB(db)

	//Instantiate default collector
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*flashback.*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

	// On every element with href for next button call callback
	c.OnHTML("ul.pagination li.next a", func(e *colly.HTMLElement) {
		next := e.Attr("href")
		if next != "" {
			fmt.Println("Visiting:", next)
			e.Request.Visit(e.Request.AbsoluteURL(next))
		}
	})

	// Scrape feed
	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		page := ScrapePosts(e)
		for _, row := range page.Posts {
			row.StripReplies()
			row.ConvertTime()
			InsertRow(db, row)
		}
	})

	c.Visit("https://www.flashback.org/f249-aktuella-brott-och-kriminalfall-50121")
}

func ScrapePosts(e *colly.HTMLElement) Page {
	var Page Page
	e.ForEach("tr", func(_ int, e *colly.HTMLElement) {
		s := &Post{}
		err := e.Unmarshal(s)
		if err != nil {
			panic(err)
		}
		Page.Posts = append(Page.Posts, *s)
	})
	return Page
}
