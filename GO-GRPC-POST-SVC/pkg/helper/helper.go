package helper

import (
	"bytes"
	"fmt"

	"github.com/akhi9550/post-svc/pkg/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func AddImageToAwsS3(file []byte, filename string) (string, error) {
	cfg, _ := config.LoadConfig()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.AWS_REGION),
		Credentials: credentials.NewStaticCredentials(
			cfg.AWS_ACCESS_KEY_ID,
			cfg.AWS_SECRET_ACCESS_KEY,
			"",
		),
	})
	if err != nil {
		return "", err
	}
	uploader := s3manager.NewUploader(sess)
	bucketName := "zhooze"

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(file),
	})

	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, filename)
	return url, nil
}
