package cache

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type R2Config struct {
	Endpoint        string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

type R2 struct {
	client *s3.Client
	bucket string
}

func NewR2(cfg R2Config) *R2 {
	client := s3.New(s3.Options{
		Region:       cfg.Region,
		BaseEndpoint: aws.String(cfg.Endpoint),
		UsePathStyle: true,
		Credentials:  credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	})

	return &R2{client: client, bucket: cfg.Bucket}
}

func (r *R2) Get(ctx context.Context, key string) ([]byte, error) {
	object, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var missing *types.NoSuchKey
		if errors.As(err, &missing) {
			return nil, ErrMiss
		}
		return nil, fmt.Errorf("get object %q: %w", key, err)
	}
	defer object.Body.Close()

	audio, err := io.ReadAll(object.Body)
	if err != nil {
		return nil, fmt.Errorf("read object %q: %w", key, err)
	}
	return audio, nil
}

func (r *R2) Put(ctx context.Context, key string, audio []byte) error {
	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       aws.String(r.bucket),
		Key:          aws.String(key),
		Body:         bytes.NewReader(audio),
		ContentType:  aws.String("audio/wav"),
		CacheControl: aws.String("public, max-age=31536000, immutable"),
	})
	if err != nil {
		return fmt.Errorf("put object %q: %w", key, err)
	}
	return nil
}
