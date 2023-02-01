package helpers

import (
	"context"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

func UploadFile(ctx context.Context, image *multipart.FileHeader, filePath string) error {
	return ctx.(*gin.Context).SaveUploadedFile(image, filePath)
}
