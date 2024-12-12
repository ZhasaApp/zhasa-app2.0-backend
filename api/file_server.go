package api

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

func (server *Server) HandleAvatarUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide a valid file.",
		})
		return
	}

	// Make sure the images directory exists
	err = os.MkdirAll("/images/avatar", os.ModePerm)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	// Generate a random filename for the new image file
	rand.Seed(time.Now().UnixNano())
	filename := generateRandomString(10) + filepath.Ext(file.Filename)

	if err := c.SaveUploadedFile(file, "/images/avatar/"+filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to save the file.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": "https://doschamp.doscar.kz/images/avatar/" + filename,
	})
}

func (server *Server) HandleNewsUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide a valid file.",
		})
		return
	}

	// Make sure the images directory exists
	err = os.MkdirAll("/images/news", os.ModePerm)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	// Generate a random filename for the new image file
	rand.Seed(time.Now().UnixNano())
	filename := generateRandomString(10) + filepath.Ext(file.Filename)

	if err := c.SaveUploadedFile(file, "/images/news/"+filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to save the file.",
		})
		return
	}

	if server.environment == "prod" {
		c.JSON(http.StatusOK, gin.H{
			"url": "https://doschamp.doscar.kz/images/news/" + filename,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"url": "http://185.182.219.90/images/news/" + filename,
		})
	}
}

func (server *Server) HandleManagersUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide a valid file.",
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer src.Close()

	// Create a new reader.
	r := csv.NewReader(src)

	// Read the first record (header)
	_, err = r.Read()
	if err != nil {
		log.Fatalf("error reading header: %s", err)
	}

	// Read remaining records
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("error reading CSV: %s", err)
	}

	for _, rec := range records {
		//	phone, err := NewPhone(transformPhoneNumber(rec[0]))
		fmt.Println(rec[0])
		if err != nil {
			fmt.Println(err)
			continue
		}
		//id, err := server.userService.CreateUser(CreateUserRequest{
		//	Phone:     *phone,
		//	FirstName: Name(rec[1]),
		//	LastName:  Name(rec[2]),
		//})
		if err != nil {
			fmt.Println(err)
			continue
		}
		//	branchId, err := strconv.Atoi(rec[3])
		//err = server.salesManagerService.CreateSalesManager(id, int32(branchId))

		if err != nil {
			fmt.Println(err)
		}
	}
}

func (server *Server) HandleDirectorsUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide a valid file.",
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer src.Close()

	// Create a new reader.
	r := csv.NewReader(src)

	// Read the first record (header)
	_, err = r.Read()
	if err != nil {
		log.Fatalf("error reading header: %s", err)
	}

	// Read remaining records
	//records, err := r.ReadAll()
	//if err != nil {
	//	log.Fatalf("error reading CSV: %s", err)
	//}
	//
	//for _, rec := range records {
	//	//phone, err := NewPhone(transformPhoneNumber(rec[0]))
	//
	//	if err != nil {
	//		fmt.Println(err)
	//		continue
	//	}
	//	//id, err := server.userService.CreateUser(CreateUserRequest{
	//	//	Phone:     *phone,
	//	//	FirstName: Name(rec[1]),
	//	//	LastName:  Name(rec[2]),
	//	//})
	//	if err != nil {
	//		fmt.Println(err)
	//		continue
	//	}
	////	branchId, err := strconv.Atoi(rec[3])
	////	_, err = server.directorService.CreateBranchDirector(id, int32(branchId))
	//
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}
}

func transformPhoneNumber(phone string) string {
	if strings.HasPrefix(phone, "8") {
		return "+7" + phone[1:]
	}
	return phone
}
