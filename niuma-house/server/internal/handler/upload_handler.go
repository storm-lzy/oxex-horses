package handler

import (
	"context"
	"time"

	"niuma-house/pkg/response"
	"niuma-house/pkg/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetPresignedURL 获取上传预签名 URL
func GetPresignedURL(c *gin.Context) {
	var req struct {
		Filename string `json:"filename" binding:"required"`
		FileType string `json:"file_type"` // image, video, file
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeInvalidParams, "参数错误")
		return
	}

	// 生成唯一文件名
	ext := ""
	if len(req.Filename) > 0 {
		for i := len(req.Filename) - 1; i >= 0; i-- {
			if req.Filename[i] == '.' {
				ext = req.Filename[i:]
				break
			}
		}
	}
	objectName := uuid.New().String() + ext

	// 生成预签名 URL (有效期 1 小时)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uploadURL, err := storage.GeneratePresignedPutURL(ctx, objectName, time.Hour)
	if err != nil {
		response.Fail(c, response.CodeServerError, "生成上传链接失败")
		return
	}

	// 生成访问 URL (有效期 24 小时)
	accessURL, err := storage.GeneratePresignedGetURL(ctx, objectName, 24*time.Hour)
	if err != nil {
		response.Fail(c, response.CodeServerError, "生成访问链接失败")
		return
	}

	response.Success(c, gin.H{
		"upload_url": uploadURL,
		"access_url": accessURL,
		"object_key": objectName,
	})
}
