package storage

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type CloudflareR2Storage struct {
	bucketName string
	s3Client   *s3.Client
	publicURL  string
}

func NewCloudflareR2Storage(client *s3.Client) *CloudflareR2Storage {
	return &CloudflareR2Storage{
		s3Client:   client,
		bucketName: os.Getenv("CLOUDFLARE_R2_BUCKET_NAME"),
		publicURL:  os.Getenv("PUBLIC_R2_PUBLIC_URL"),
	}
}

func (s *CloudflareR2Storage) UploadPicture(ctx context.Context, fileBytes []byte, fileName string) (string, error) {
	// failsafe to ensure the env is loaded
	if s.bucketName == "" || s.publicURL == "" {
		return "", fmt.Errorf("R2 bucket name or public URL is not configured")
	}

	// Sniff the MIME type (e.g., "image/jpeg", "image/png")
	// This ensures browsers render the image instead of forcing a download
	contentType := http.DetectContentType(fileBytes)

	// Prepare the AWS S3 Upload Object
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(fileName),       // e.g., "profiles/uuid/avatar.jpg"
		Body:        bytes.NewReader(fileBytes), // Convert []byte to io.Reader stream
		ContentType: aws.String(contentType),
	}

	// Execute the Upload to Cloudflare
	_, err := s.s3Client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("[r2_image_storage.go] UploadProfilePicture: failed to put object in bucket: %w", err)
	}

	// Construct the final public URL to save in the database
	// Example: https://cdn.the-daily-sunshine.com/profiles/123/avatar.jpg
	finalURL := fmt.Sprintf("%s/%s", s.publicURL, fileName)

	return finalURL, nil
}

func (s *CloudflareR2Storage) DeletePicture(ctx context.Context, fullImageURL string) error {
	// failsafe to ensure the env is loaded
	if s.bucketName == "" || s.publicURL == "" {
		return fmt.Errorf("R2 bucket name or public URL is not configured")
	}

	baseURL := s.publicURL
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	s3Key := strings.Replace(fullImageURL, baseURL, "", 1)

	// Prepare the AWS S3 Delete Object
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(s3Key), // e.g., "profiles/uuid/avatar.jpg"
	}

	// Execute the delete request to Cloudflare
	_, err := s.s3Client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("[r2_image_storage.go] DeleteProfilePicture: failed to delete object from bucket: %w", err)
	}

	return nil
}
