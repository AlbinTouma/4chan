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

package main

import (
  "fmt"
  "time"
  "github.com/gocolly/colly"
  "database/sql"
_ "github.com/mattn/go-sqlite3"
)

/*I want to scrape the table rows for information about each post.
In the Scrape Posts function, I pass in the current page as e and I return rows from each page as an array maps
*/

type Page struct {
  Posts []Post
}

type Post struct {
  Title string  `selector:".td_title a"`
  Date  string  `selector:".td_last_post > div:nth-child(1)"`
  Author  string  `selector:".td_last_post > div:nth-child(2) a"`
  Link  string `selector:".td_title a" attr:"href"`
  Replies  string `selector:".td_replies > div:nth-child(1)"`
  Views string  `selector:".td_replies > div:nth-child(2)"`
}


func CreateDB(db *sql.DB){
  createTable := `CREATE TABLE IF NOT EXISTS posts (
      title STRING,
      date  STRING,
      author STRING,
      link  STRING,
      replies STRING,
      views  STRING
  )`
  _, err := db.Exec(createTable)
  if err != nil {
    fmt.Println("Error in creating table")
    panic(err)
  }
    
}

func InsertRow(db *sql.DB, Post Post){
  query := `INSERT INTO posts values (?,?,?,?,?,?);`
  _, err := db.Exec(query, Post.Title, Post.Date, Post.Author, Post.Link, Post.Replies, Post.Views)
  if err != nil {
    fmt.Println("Error inserting rows")
    panic(err)
  }
}

func ScrapePosts(e *colly.HTMLElement)  Page {
  var Page Page
  e.ForEach("tr", func(_ int, item*colly.HTMLElement){
      s := &Post{} 
      err:= e.Unmarshal(s)
      if err != nil{
        panic(err)
      }
    fmt.Printf("s is %", s.Title)
      Page.Posts = append(Page.Posts, *s)
  })
  return Page
}

func main(){
  var db *sql.DB
  db, err := sql.Open("sqlite3", "./data/forum.db")
  if err != nil {

    fmt.Println("Error in creating db")
    panic(err)
  }
  
  CreateDB(db)


  c := colly.NewCollector()
  c.Limit(&colly.LimitRule{
    DomainGlob: "*flashback.*",
    Parallelism: 2,
    Delay: 1 *time.Second,
  })

  //Handle pagination
  c.OnHTML("ul.pagination li.next a", func(e *colly.HTMLElement){
    nextPage := e.Attr("href")
    if nextPage != "" {
      url := e.Request.AbsoluteURL(nextPage)
      e.Request.Visit(url)
    }
  })

  c.OnHTML("tbody", func(e *colly.HTMLElement){
    page := ScrapePosts(e)
    for _, row := range page.Posts {
      //fmt.Printf("Row title %s", row.Title)
      InsertRow(db, row)
    }
	})

  c.OnRequest(func(r *colly.Request){
    fmt.Println("Visiting", r.URL.String())
  })

  c.Visit("https://www.flashback.org/f358-forsvunna-personer-50121")//"https://www.flashback.org/f249-aktuella-brott-och-kriminalfall-50121")
}
