package helper

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type UploadHelper struct{}

func (h *UploadHelper) UploadImage(file *multipart.FileHeader, path string) {
	config := aws.NewConfig().WithRegion("us-east-1")
	service := s3.New(session.New(), config)

	size := file.Size
	buffer := make([]byte, size)

	fileOpen, err := file.Open()

	fileOpen.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	params := &s3.PutObjectInput{
		Bucket:        aws.String("udabayarmedia-bucket"),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}

	resp, err := service.PutObject(params)

	if err != nil {
		panic(err)
	}

	fmt.Printf("response %s", awsutil.StringValue(resp))
}

func NewUploadHelper() *UploadHelper {
	return &UploadHelper{}
}
