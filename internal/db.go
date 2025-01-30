package internal

import (
	"database/sql"
	"fmt"
)

func CreateDB(db *sql.DB) {
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

func OpenDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		panic(err)
	}
	return db
}

func InsertRow(db *sql.DB, Post Post) {
	query := `INSERT INTO posts values (?,?,?,?,?,?);`
	_, err := db.Exec(query, Post.Title, Post.Date, Post.Author, Post.Link, Post.Replies, Post.Views)
	if err != nil {
		fmt.Println("Error inserting rows")
		panic(err)
	}
}
