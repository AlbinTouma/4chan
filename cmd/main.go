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
  "strconv"
	"database/sql"
	"fmt"
	"time"
  "strings"
	"github.com/gocolly/colly"
	_ "github.com/mattn/go-sqlite3"
  "unicode"
)

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

//2024-12-06 01:48


func(s *Post)ConvertTime(){
//  dateString := strings.Split(s," ")[0]
  dateString := s.Date
  if dateString == ""{
    return
  }

  if strings.HasPrefix(dateString, "Idag"){
    	today := time.Now().Format("2006-01-02")
    	timePart := strings.TrimPrefix(dateString, "Idag ")
      dateString = fmt.Sprintf("%s %s", today, timePart)
    }

  if strings.HasPrefix(dateString, "Igår"){
    	today := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
    	timePart := strings.TrimPrefix(dateString, "Igår")
      dateString = fmt.Sprintf("%s %s", today, timePart)
    }
 

  	layout := "2006-01-02 15:04"
  	parsedTime, err := time.Parse(layout, dateString)
	  if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}
  s.Date = parsedTime.Format(layout)

}



func(s *Post) StripReplies(){
	r := strings.Replace(s.Replies, "svar", "", -1)
  r = strings.TrimSpace(r)
	r = cleanNumericString(r)
	replies, err := strconv.Atoi(r)
  if err != nil {
    fmt.Errorf("Failed to convert %s", err)
  }

  v := strings.Replace(s.Views, "visningar","", -1)
  v = strings.TrimSpace(v)
	v = cleanNumericString(v)
  views, err := strconv.Atoi(v)

   if err != nil {
    fmt.Errorf("Failed to convert %s", err)
  }

  s.Replies = strconv.Itoa(replies) //replies
  s.Views = strconv.Itoa(views) //iews
}

func cleanNumericString(input string) string {
	// Replace non-breaking spaces and other whitespace characters
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1 // Remove the character
		}
		return r
	}, input)
}

func CreateDB(db *sql.DB){
  createTable := `CREATE TABLE IF NOT EXISTS posts (
      title STRING,
      date  DATETIME,
      author STRING,
      link  STRING,
      replies INTEGER,
      views  INTEGER
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
  e.ForEach("tr", func(_ int, e *colly.HTMLElement){
    s := &Post{} 
    err:= e.Unmarshal(s)
    if err != nil{
      panic(err)
     }
    Page.Posts = append(Page.Posts, *s)
  })
    return Page
}

func OpenDB() *sql.DB {
    db, err := sql.Open("sqlite3", "./data/forum.db")
    if err != nil {
    panic(err)
    }
    return db
}

func main(){
  var db *sql.DB
  db = OpenDB()
  CreateDB(db)

  //Instantiate default collector
  c := colly.NewCollector()
  c.Limit(&colly.LimitRule{
    DomainGlob: "*flashback.*",
    Parallelism: 2,
    Delay: 1 *time.Second,
  })

  // On every element with href for next button call callback
  c.OnHTML("ul.pagination li.next a", func(e *colly.HTMLElement){
    next := e.Attr("href")
    if next != "" {
      fmt.Println("Visiting:", next)
        e.Request.Visit(e.Request.AbsoluteURL(next))
      }
  })

  // Scrape feed
  c.OnHTML("tbody", func(e *colly.HTMLElement){
    page := ScrapePosts(e)
    fmt.Println(page.Posts)
    for _, row := range page.Posts {
      row.StripReplies()
      row.ConvertTime()
      InsertRow(db, row)
    }
	})


  c.Visit("https://www.flashback.org/f249-aktuella-brott-och-kriminalfall-50121")
}
