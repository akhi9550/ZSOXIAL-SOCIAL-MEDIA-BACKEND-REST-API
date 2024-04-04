package helper

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/akhi9550/auth-svc/pkg/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func GetImageMimeType(filename string) string {
	extension := strings.ToLower(strings.Split(filename, ".")[len(strings.Split(filename, "."))-1])

	imageMimeTypes := map[string]string{
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
		"gif":  "image/gif",
		"bmp":  "image/bmp",
		"webp": "image/webp",
	}

	if mimeType, ok := imageMimeTypes[extension]; ok {
		return mimeType
	}

	return "application/octet-stream"
}

func AddImageToAwsS3(file []byte, filename string) (string, error) {
	cfg, _ := config.LoadConfig()
	mimeType := GetImageMimeType(filename)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.AWS_REGION),
		Credentials: credentials.NewStaticCredentials(
			cfg.AWS_ACCESS_KEY_ID,
			cfg.AWS_SECRET_ACCESS_KEY,
			"",
		),
	})
	if err != nil {
		fmt.Println("error in session config", err)
		return "", err
	}
	// Create an S3 uploader with the session and default options
	uploader := s3manager.NewUploader(sess)
	BucketName := "zsoxial"
	//upload data(video or image)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(BucketName),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(file),
		ContentType: aws.String(mimeType),
	})
	if err != nil {
		fmt.Println("error 2", err)
		return "", err
	}
	url := fmt.Sprintf("https://d2jkb5eqmpty2t.cloudfront.net/%s/%s", BucketName, filename)
	return url, nil
}
