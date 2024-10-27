package storage

import (
	"bytes"
	"context"
	"cowlibrary/constants"
	"cowlibrary/utils"
	"cowlibrary/vo"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	cowStorage "appconnections/minio"

	"github.com/minio/minio-go/v7"
)

func UploadFileToMinio(minioFileVO *vo.MinioFileVO) (*vo.MinioFileInfoVO, error) {
	minoEndpoint := utils.Getenv(constants.EnvMinioLoginURL, "cowstorage:9000")
	// Supress log
	log.SetOutput(io.Discard)

	minioClient, err := cowStorage.RegisterMinio(minoEndpoint, utils.Getenv(constants.EnvMinioRootUser, ""), utils.Getenv(constants.EnvMinioRootPassword, ""), minioFileVO.BucketName)

	if err != nil {
		return nil, err
	}

	if minioClient == nil {
		return nil, errors.New("cannot create minio client")
	}

	folderPath := path.Join(minioFileVO.Path, minioFileVO.FileName)

	fileExtension := filepath.Ext(minioFileVO.FileName)
	contentType := "text/csv/jpg/jpeg"
	if utils.IsNotEmpty(fileExtension) {
		contentType = fileExtension
	}

	bucketName, prefix := cowStorage.GetBucketAndPrefix(minioFileVO.BucketName)

	folderPath = prefix + folderPath

	_, err = minioClient.PutObject(context.Background(), bucketName, folderPath, bytes.NewBuffer(minioFileVO.FileContent), int64(len(minioFileVO.FileContent)), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return nil, err
	}
	var urlValues url.Values

	url, err := minioClient.PresignedGetObject(context.Background(), bucketName, folderPath, 7*time.Hour*24, urlValues)

	url.RawQuery = ""

	if err != nil {
		return nil, err
	}
	return &vo.MinioFileInfoVO{FileURL: url.String()}, nil

}

func DownloadFile(minioFileInfoVO *vo.MinioFileInfoVO, additionalInfoVO *vo.AdditionalInfo) (*vo.MinioFileVO, *vo.ErrorResponseVO) {

	fileURL, err := url.Parse(minioFileInfoVO.FileURL)
	if err != nil {
		return nil, &vo.ErrorResponseVO{StatusCode: http.StatusBadRequest, Error: &vo.ErrorVO{
			Message: "Cannot parse the file rule", Description: "Cannot parse the file rule",
			ErrorDetails: utils.GetValidationError(fmt.Errorf("invalid file URL: %w", err))}}
	}
	// Default bucketName
	bucketName := "demo"
	if fileURL.Scheme == "http" {
		splitPath := strings.Split(fileURL.Path, "/")
		if len(splitPath) > 2 {
			bucketName = splitPath[1]
		}
	}

	minoEndpoint := utils.Getenv(constants.EnvMinioLoginURL, "cowstorage:9000")
	// Supress log
	log.SetOutput(io.Discard)
	fileURL.Host = minoEndpoint

	objectPath := strings.TrimPrefix(fileURL.Path, fmt.Sprintf("/%v", bucketName))

	bucketName, prefix := cowStorage.GetBucketAndPrefix(bucketName)
	objectPath = prefix + objectPath

	if cowStorage.IsAmazonS3Host(minioFileInfoVO.FileURL) {
		parts := strings.Split(fileURL.Path, "/")
		if len(parts) < 4 {
			return nil, &vo.ErrorResponseVO{StatusCode: http.StatusBadRequest, Error: &vo.ErrorVO{
				Message: "invalid URL structure", Description: "invalid URL structure, cannot extract bucket and object",
				ErrorDetails: utils.GetValidationError(fmt.Errorf("invalid URL structure, cannot extract bucket and object: %w", err))}}
		}
		bucketName = parts[3]
		objectPath = strings.Join(parts[4:], "/")
	}

	if prefix := fileURL.Query().Get("prefix"); prefix != "" {
		objectPath = prefix
	}

	minioClient, err := cowStorage.RegisterMinio(minoEndpoint, utils.Getenv(constants.EnvMinioRootUser, ""), utils.Getenv(constants.EnvMinioRootPassword, ""), bucketName)

	if err != nil {
		return nil, &vo.ErrorResponseVO{StatusCode: http.StatusBadRequest, Error: &vo.ErrorVO{
			Message: "Cannot connect to the minio system", Description: "Cannot connect to the minio system",
			ErrorDetails: utils.GetValidationError(fmt.Errorf("cannot connect to the minio system: %w", err))}}
	}

	if minioClient == nil {
		return nil, &vo.ErrorResponseVO{StatusCode: http.StatusBadRequest, Error: &vo.ErrorVO{
			Message: "cannot create minio client", Description: "cannot create minio client",
			ErrorDetails: utils.GetValidationError(errors.New("cannot create minio client"))}}
	}

	fileObject, err := minioClient.GetObject(context.Background(), bucketName, objectPath, minio.GetObjectOptions{})
	if err != nil {
		return nil, &vo.ErrorResponseVO{StatusCode: http.StatusBadRequest, Error: &vo.ErrorVO{
			Message: "cannot find the file", Description: "cannot find the file",
			ErrorDetails: utils.GetValidationError(fmt.Errorf("cannot find the file: %w", err))}}
	}
	fileContent, err := io.ReadAll(fileObject)
	if err != nil {
		return nil, &vo.ErrorResponseVO{StatusCode: http.StatusBadRequest, Error: &vo.ErrorVO{
			Message: "cannot read the file", Description: "cannot read the file",
			ErrorDetails: utils.GetValidationError(fmt.Errorf("cannot read the file: %w", err))}}
	}

	fmt.Println("fileContent :", string(fileContent))

	return &vo.MinioFileVO{FileContent: fileContent, FileName: filepath.Base(objectPath)}, nil

}
