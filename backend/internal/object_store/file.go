package object_store

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/Tylerchristensen100/object_browser/internal/constants"
	"github.com/minio/minio-go/v7"
)

func (o *ObjectStore) GetFile(app AppContext, bucket string, path string) ([]byte, error) {
	object, err := o.client.GetObject(app.Context(), bucket, path, minio.GetObjectOptions{})
	if err != nil {
		app.Logger().Error("object_storage/file: error getting object from storage", "error", err)
		return nil, err
	}

	// Check if the object exists and is not empty
	if object == nil {
		return nil, fmt.Errorf("object not found or empty")
	}
	defer object.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(object); err != nil {
		if err.Error() == "The specified key does not exist." {
			return nil, constants.ErrObjectNotFound
		}
		app.Logger().Error("object_storage/file: error reading object data", "error", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

func (o *ObjectStore) SaveFile(app AppContext, bucket string, file multipart.File, header *multipart.FileHeader, filename string) (minio.UploadInfo, error) {

	uploadInfo, err := o.client.PutObject(app.Context(), bucket, filename, file, header.Size, minio.PutObjectOptions{
		ContentType:        header.Header.Get("Content-Type"),
		ContentDisposition: header.Header.Get("Content-Disposition"),
	})
	if err != nil {
		return minio.UploadInfo{}, err
	}

	if uploadInfo.Key[0] != '/' {
		uploadInfo.Key = fmt.Sprintf("/%s", uploadInfo.Key)
	}

	app.Logger().Info("object_storage/file: successfully uploaded file", "bucket", bucket, "filename", filename, "size", uploadInfo.Size)
	return uploadInfo, nil
}

func (o *ObjectStore) FileExists(app AppContext, bucket string, path string) (bool, error) {
	if strings.HasSuffix(path, "/") {
		// Directories are not objects, so we consider them to always exist
		return true, nil
	}
	_, err := o.client.StatObject(app.Context(), bucket, path, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" || minio.ToErrorResponse(err).Code == "AccessDenied" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (o *ObjectStore) DeleteFile(app AppContext, bucket string, path string) error {
	err := o.client.RemoveObject(app.Context(), bucket, path, minio.RemoveObjectOptions{})
	if err != nil {
		app.Logger().Error("object_storage/file: error deleting object from storage", "error", err)
		return err
	}

	return nil
}
