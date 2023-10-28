package entity

type Book struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	ReleaseYear string `json:"releaseYear"`
	Pages       int    `json:"pages"`
}

type BookPost struct {
	Id          int    `json:"id"`
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author" binding:"required"`
	ReleaseYear string `json:"releaseYear" binding:"required,min=4,max=4"`
	Pages       int    `json:"pages" binding:"required"`
}

type BookUpdate struct {
	Id          int    `json:"id" binding:"required"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	ReleaseYear string `json:"releaseYear"`
	Pages       int    `json:"pages"`
}

type BookDelete struct {
	Id          int    `json:"id" binding:"required"`
}