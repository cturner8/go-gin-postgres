package routes

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	log.Println("getAlbums")

	db := c.MustGet("db").(*sql.DB)

	// An albums slice to hold data from returned rows.
	var albums []album

	rows, err := db.Query("SELECT * FROM recordings.album;")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	log.Println("postAlbums")

	db := c.MustGet("db").(*sql.DB)

	var newAlbum album
	c.BindJSON(&newAlbum)

	log.Println("Body", newAlbum.Title, newAlbum.Artist, newAlbum.Price)

	_, err := db.Exec("INSERT INTO recordings.album (title, artist, price) VALUES ($1, $2, $3)", newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	log.Println("getAlbumByID")

	db := c.MustGet("db").(*sql.DB)

	id := c.Param("id")

	log.Println("ID", id)

	var alb album
	row := db.QueryRow("SELECT * FROM recordings.album WHERE id = $1", id)

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, alb)
}

func RegisterAlbums(api *gin.RouterGroup) *gin.RouterGroup {
	rg := api.Group("/albums")
	{
		rg.GET("", getAlbums)
		rg.POST("", postAlbums)

		rg.GET(":id", getAlbumByID)
	}

	return rg
}
