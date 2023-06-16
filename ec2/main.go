package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Static("/images", "/home/ubuntu/images/")
	r.GET("/files", getFilesHandler)
	r.POST("/image/ec2", uploadFileHandler)

	err := r.Run("0.0.0.0:9000")
	if err != nil {
		log.Fatal(err)
	}
}

// Handler for retrieving the list of files
func getFilesHandler(c *gin.Context) {
	dirPath := "/home/ubuntu/images" // Specify the directory path relative to the `/images` route

	// Retrieve the list of files in the directory
	files, err := getFilesInDirectory(dirPath)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving files: %s", err.Error()))
		return
	}

	// Generate a slice of file links
	links := make([]string, 0)
	for _, file := range files {
		fileLink := fmt.Sprintf("http://35.154.84.78:9000/images/%s", file.Name())
		links = append(links, fileLink)
	}

	c.JSON(http.StatusOK, gin.H{
		"files": links,
	})
}

// Handler for uploading a file
func uploadFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// Save the file to the desired location on the EC2 instance
	err = c.SaveUploadedFile(file, "/home/ubuntu/images/"+file.Filename)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H {
		 "File uploaded successfully": file.Filename,
})
}

// Retrieve the list of files in a directory
func getFilesInDirectory(dirPath string) ([]os.FileInfo, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	return files, nil
}
