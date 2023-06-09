package api

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Generate a random string for image name
func generateRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func (server Server) HandleUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide a valid file.",
		})
		return
	}

	// Make sure the images directory exists
	err = os.MkdirAll("images", os.ModePerm)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	// Generate a random filename for the new image file
	rand.Seed(time.Now().UnixNano())
	filename := generateRandomString(10) + filepath.Ext(file.Filename)

	if err := c.SaveUploadedFile(file, "images/"+filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to save the file.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": "/images/" + filename,
	})
}
