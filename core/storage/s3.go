package storage

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	Client  *s3.Client
	Bucket  string
	BaseURL string
}

func (s *S3Storage) Name() string {
	return "s3"
}

func (s *S3Storage) Upload(file *multipart.FileHeader, path string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	key := fmt.Sprintf("%s/%s", path, file.Filename)

	_, err = s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   src,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", s.BaseURL, key), nil
}

func (s *S3Storage) Delete(path string) error {
	_, err := s.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(path),
	})
	return err
}

func NewS3(region, accessKey, secretKey, bucket, baseURL string) *S3Storage {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)
	if err != nil {
		log.Fatalf("cannot load aws config: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	return &S3Storage{
		Client:  client,
		Bucket:  bucket,
		BaseURL: baseURL,
	}
}

func init() {
	// Hardcode dulu, nanti bisa dibaca dari ENV
	Register(NewS3(
		os.Getenv("AWS_REGION"),
		os.Getenv("AWS_ACCESS_KEY"),
		os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_BUCKET"),
		os.Getenv("AWS_URL"),
	))
}
