package main

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/upload", uploadImage)

	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}

func uploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve the file"})
		return
	}

	err = saveImageToS3(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload the image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully"})
}

func saveImageToS3(file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1"),
		Credentials: credentials.NewStaticCredentials(
			"YOUR_AWS_ACCESS_KEY_ID",
			"YOUR_AWS_SECRET_ACCESS_KEY",
			""),
	})
	if err != nil {
		return err
	}

	s3Svc := s3.New(sess)

	fileName := filepath.Base(file.Filename)
	s3Key := fmt.Sprintf("images/%s", fileName)

	_, err = s3Svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("YOUR_S3_BUCKET_NAME"),
		Key:    aws.String(s3Key),
		Body:   src,
	})
	if err != nil {
		return err
	}

	return nil
}
