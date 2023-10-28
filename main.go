package main

import (
	"book-api/config"
	"book-api/entity"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

/* global variable */
var db = config.ConnectDb()

/* main function */
func main() {
	fmt.Println("Hello, Welcome!")

	// intial router
	router := gin.Default()

	// get books
	router.GET("/books", getBooks)

	// create book
	router.POST("/books", createBook)

	// update book
	router.PUT("/books", updateBook)

	// delete book
	router.DELETE("/books", deleteBook)

	router.Run(":8080")
}


/*
########################################################################################################
Function
########################################################################################################
*/

/* delete */
func deleteBook(c *gin.Context) {
	book := entity.BookDelete{}
	err := c.ShouldBind(&book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "DELETE FROM mst_book WHERE id=$1;"
	result, err := db.Exec(query, book.Id)
	fmt.Println(result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book successfully deleted!"})
}

/* update book */
func updateBook(c *gin.Context)  {
	book := entity.BookUpdate{}
	err := c.ShouldBind(&book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
		return
	}

	if book.Title != "" {
		query := "UPDATE mst_book SET title=$1 WHERE id=$2;"
		_, err = tx.Exec(query, book.Title, book.Id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
			return
		}
	}

	if book.Author != "" {
		query := "UPDATE mst_book SET author=$1 WHERE id=$2;"
		_, err = tx.Exec(query, book.Author, book.Id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
			return
		}
	}

	if (book.ReleaseYear != ""){
		if len(book.ReleaseYear) != 4 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "release year must be 4 digit!"})
			return
		}
		query := "UPDATE mst_book SET release_year=$1 WHERE id=$2;"
		_, err = tx.Exec(query, book.ReleaseYear, book.Id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
			return
		}
	}

	if book.Pages != 0  {
		if book.Pages < 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Book pages must be greater than 0!"})
			return
		}
		query := "UPDATE mst_book SET pages=$1 WHERE id=$2;"
		_, err = tx.Exec(query, book.Pages, book.Id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
		return
	}

	c.JSON(http.StatusCreated, book)
}

/* insert book */
func createBook(c *gin.Context)  {
	book := entity.BookPost{}
	err := c.ShouldBind(&book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO mst_book (title, author, release_year, pages) VALUES($1, $2, $3, $4) RETURNING id;"

	var bookId int
	err = db.QueryRow(query, book.Title, book.Author, book.ReleaseYear, book.Pages).Scan(&bookId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
		return
	}
	
	book.Id = bookId
	c.JSON(http.StatusCreated, book)
}

/* select book */
func getBooks(c *gin.Context) {
	searchedTitle := c.Query("title")

	query := "SELECT id, title, author, release_year, pages FROM mst_book"
	var rows *sql.Rows
	var err error
	if searchedTitle != "" {
		query += " WHERE title ILIKE '%' || $1 || '%';"
		rows, err = db.Query(query, searchedTitle)
	} else {
		query += ";"
		rows, err = db.Query(query)
	}
	
	defer rows.Close()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
		return
	}

	books := []entity.Book{}
	for rows.Next() {
		aBook := entity.Book{}
		err := rows.Scan(&aBook.Id, &aBook.Title, &aBook.Author, &aBook.ReleaseYear, &aBook.Pages)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
			return
		}
		books = append(books, aBook)
	}

	if len(books) > 0 {
		c.JSON(http.StatusOK, books)
		return
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
}