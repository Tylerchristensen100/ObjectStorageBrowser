package object_store

import (
	"errors"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func LoadConfig(app AppContext, endpoint, key, password string) (*ObjectStore, error) {
	// endpoint := secrets.ObjectStore.Endpoint
	// key := secrets.ObjectStore.AccessKey
	// password := secrets.ObjectStore.SecretKey

	if key == "" || password == "" {
		app.Logger().Error("object_storage/config: missing object storage credentials in environment variables", "key", key, "password", password)
		return nil, errors.New("missing object storage credentials in environment variables")
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(key, password, ""),
		Secure: true,
	})
	if err != nil {
		return nil, errors.New("object_store/LoadConfig: failed to create S3 client: " + err.Error())
	}

	store := &ObjectStore{
		client: client,
	}

	store.ListBuckets(app)

	return store, nil
}

// func (o *ObjectStore) bucketExists(app AppContext, bucket string) (bool, error) {
// 	exists, err := o.client.BucketExists(app.Context(), bucket)
// 	if err != nil {
// 		app.Logger().Error("object_storage/bucketExists: error checking if bucket exists", "error", err)
// 		return false, errors.New("error checking if bucket exists")
// 	}
// 	if !exists {
// 		app.Logger().Error("object_storage/bucketExists: bucket does not exist", "bucket", bucket)
// 		return false, errors.New("object storage bucket does not exist")
// 	}
// 	return exists, nil
// }
