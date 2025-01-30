package internal

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func ParseCrimes() {
	db, err := sqlx.Open("sqlite3", "data/forum.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var page Page
	err = db.Select(&page.Posts, "SELECT title FROM posts LIMIT 100;")
	if err != nil {
		panic(err)
	}

	for _, post := range page.Posts {
		postTitle := post.Title
		fmt.Println(postTitle)

	}

}
