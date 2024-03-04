package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

type book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func main() {
	var err error
	connStr := "user=postgres password=1403 dbname=library sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", addBook)
	router.DELETE("/books/:id", delBook)
	router.Run("localhost:8080")
}

func getBooks(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, title, author FROM books")
	if err != nil {
		log.Fatal(err)
	}

	var books []book
	for rows.Next() {
		var a book
		err := rows.Scan(&a.ID, &a.Title, &a.Author)
		if err != nil {
			log.Fatal(err)
		}
		books = append(books, a)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	c.IndentedJSON(http.StatusOK, books)
}

func addBook(c *gin.Context) {
	var awesomeBook book
	if err := c.BindJSON(&awesomeBook); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	stmt, err := db.Prepare("INSERT INTO books (id, title, author) VALUES ($1, $2, $3)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(awesomeBook.ID, awesomeBook.Title, awesomeBook.Author); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusCreated, awesomeBook)
}

func delBook(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM books WHERE id=$1", id)
	if err != nil {
		log.Fatal(err)
	}

}
