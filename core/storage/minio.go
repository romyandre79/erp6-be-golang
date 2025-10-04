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

type MinioStorage struct {
	Client  *s3.Client
	Bucket  string
	BaseURL string
}

func (m *MinioStorage) Name() string {
	return "minio"
}

func (m *MinioStorage) Upload(file *multipart.FileHeader, path string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	key := fmt.Sprintf("%s/%s", path, file.Filename)

	_, err = m.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(m.Bucket),
		Key:    aws.String(key),
		Body:   src,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", m.BaseURL, key), nil
}

func (m *MinioStorage) Delete(path string) error {
	_, err := m.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(m.Bucket),
		Key:    aws.String(path),
	})
	return err
}

func NewMinio(endpoint, accessKey, secretKey, bucket, baseURL string) *MinioStorage {
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: "us-east-1",
			}, nil
		},
	)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		log.Fatalf("cannot load minio config: %v", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &MinioStorage{
		Client:  client,
		Bucket:  bucket,
		BaseURL: baseURL,
	}
}

func init() {
	// sementara hardcode, nanti bisa di-load dari env/config
	Register(NewMinio(
		os.Getenv("MINIO_ENDPOINT"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("MINIO_BUCKET"),
		os.Getenv("MINIO_BASE_URL"),
	))
}
