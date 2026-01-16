package storage

import (
	"context"
	"log"
	"sync"
	"time"

	"niuma-house/pkg/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	minioClient *minio.Client
	bucketName  string
	once        sync.Once
)

// InitMinIO 初始化 MinIO 客户端
func InitMinIO(cfg *config.MinIOConfig) *minio.Client {
	once.Do(func() {
		var err error
		minioClient, err = minio.New(cfg.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure: cfg.UseSSL,
		})
		if err != nil {
			log.Fatalf("Failed to connect to MinIO: %v", err)
		}

		bucketName = cfg.Bucket

		// 确保 bucket 存在
		ctx := context.Background()
		exists, err := minioClient.BucketExists(ctx, bucketName)
		if err != nil {
			log.Fatalf("Failed to check bucket existence: %v", err)
		}
		if !exists {
			err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
			if err != nil {
				log.Fatalf("Failed to create bucket: %v", err)
			}
			log.Printf("Bucket '%s' created successfully", bucketName)
		}

		log.Println("MinIO connected successfully")
	})
	return minioClient
}

// GetMinIO 获取 MinIO 客户端单例
func GetMinIO() *minio.Client {
	if minioClient == nil {
		log.Fatal("MinIO not initialized. Call InitMinIO first.")
	}
	return minioClient
}

// GetBucketName 获取 bucket 名称
func GetBucketName() string {
	return bucketName
}

// GeneratePresignedPutURL 生成预签名上传 URL
func GeneratePresignedPutURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := minioClient.PresignedPutObject(ctx, bucketName, objectName, expiry)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

// GeneratePresignedGetURL 生成预签名下载 URL
func GeneratePresignedGetURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := minioClient.PresignedGetObject(ctx, bucketName, objectName, expiry, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
