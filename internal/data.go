package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Page struct {
	Posts []Post
}

type Post struct {
	Title   string `selector:".td_title a" db:"title"`
	Date    string `selector:".td_last_post > div:nth-child(1)" db:"date"`
	Author  string `selector:".td_last_post > div:nth-child(2) a" db:"author"`
	Link    string `selector:".td_title a" attr:"href" db:"link"`
	Replies string `selector:".td_replies > div:nth-child(1)" db:"replies"`
	Views   string `selector:".td_replies > div:nth-child(2)" db:"views"`
}

func (s *Post) ConvertTime() {
	dateString := s.Date
	if dateString == "" {
		return
	}

	if strings.HasPrefix(dateString, "Idag") {
		today := time.Now().Format("2006-01-02")
		timePart := strings.TrimPrefix(dateString, "Idag ")
		dateString = fmt.Sprintf("%s %s", today, timePart)
	}

	if strings.HasPrefix(dateString, "Igår") {
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

func (s *Post) StripReplies() {
	r := strings.Replace(s.Replies, "svar", "", -1)
	r = strings.TrimSpace(r)
	r = cleanNumericString(r)
	replies, err := strconv.Atoi(r)
	if err != nil {
		fmt.Errorf("Failed to convert %s", err)
	}

	v := strings.Replace(s.Views, "visningar", "", -1)
	v = strings.TrimSpace(v)
	v = cleanNumericString(v)
	views, err := strconv.Atoi(v)

	if err != nil {
		fmt.Errorf("Failed to convert %s", err)
	}

	s.Replies = strconv.Itoa(replies) //replies
	s.Views = strconv.Itoa(views)     //iews
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

type Gpt4AllResponse struct {
}
