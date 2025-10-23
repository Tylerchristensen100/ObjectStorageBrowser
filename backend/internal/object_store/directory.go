package object_store

import (
	"fmt"
	"strings"

	"github.com/minio/minio-go/v7"
)

func (o *ObjectStore) ListBuckets(app AppContext) ([]Bucket, error) {
	buckets, err := o.client.ListBuckets(app.Context())
	if err != nil {
		app.Logger().Error("object_storage/directory: error listing buckets", "error", err)
		return nil, err
	}

	var listedBuckets []Bucket
	for _, bucket := range buckets {
		listedBuckets = append(listedBuckets, Bucket{Name: bucket.Name, info: bucket})
	}

	o.Buckets = listedBuckets

	return listedBuckets, nil
}

func (o *ObjectStore) ListDirectory(app AppContext, bucket string, prefix string, recursive bool) ([]Listing, error) {
	prefix = strings.TrimPrefix(prefix, "/")
	opts := minio.ListObjectsOptions{
		// The path to list
		Prefix: prefix,
		// Set to true for a simple listing
		Recursive: recursive,
		// Get a manageable number of results
		MaxKeys: 1000,
	}
	var listings []Listing = make([]Listing, 0)
	objects := o.client.ListObjects(app.Context(), bucket, opts)
	for object := range objects {
		if object.Err != nil {
			app.Logger().Error("object_storage/directory: error listing directory", "error", object.Err)
			return nil, object.Err
		}
		stats, err := o.client.StatObject(app.Context(), bucket, strings.TrimPrefix(object.Key, "/"), minio.StatObjectOptions{})
		if err != nil {
			app.Logger().Error("object_storage/directory: error getting stats", "error", err)
			// return nil, err
		}

		if object.Key == prefix {
			// Skip the directory itself
			continue
		}
		if !strings.HasPrefix(object.Key, "/") {
			object.Key = fmt.Sprintf("/%s", object.Key)
		}
		if strings.HasSuffix(object.Key, "/") && object.Key != prefix {
			listings = append(listings, Listing{
				Path:  object.Key,
				IsDir: true,
				Size:  object.Size,
				Stats: stats,
			})
		} else {
			listings = append(listings, Listing{
				Path:  object.Key,
				IsDir: false,
				Size:  object.Size,
				Stats: stats,
			})
		}
	}
	return listings, nil
}

func (o *ObjectStore) Directory(app AppContext, bucket string, prefix string) (*DirectoryItem, error) {
	// Lists all the objects and directories at the given prefix and constructs a tree structure.
	root := prefix

	if root != "" {
		if !strings.HasPrefix(root, "/") {
			root = fmt.Sprintf("/%s", root)
		}
		if !strings.HasSuffix(root, "/") {
			root = fmt.Sprintf("%s/", root)
		}
	}

	dir := &DirectoryItem{
		Path:     root,
		IsDir:    true,
		Children: []DirectoryItem{},
		Depth:    0,
	}

	listing, err := o.ListDirectory(app, bucket, prefix, true)
	if err != nil {
		return nil, err
	}
	// current := dir
	for _, item := range sortListingByPath(listing) {
		item.Path = strings.TrimPrefix(item.Path, root)
		dir.AddItemFromRoot(item)
	}
	return dir, nil
}

func (o *ObjectStore) DeleteDirectory(app AppContext, bucket string, prefix string) (int, error) {
	prefix = strings.TrimPrefix(prefix, "/")
	opts := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}
	deletedFiles := 0
	objects := o.client.ListObjects(app.Context(), bucket, opts)
	for object := range objects {
		if object.Err != nil {
			app.Logger().Error("object_storage/directory: error listing directory", "error", object.Err)
			return 0, object.Err
		}

		if object.Key == prefix {
			// Skip the directory itself
			continue
		}

		err := o.client.RemoveObject(app.Context(), bucket, object.Key, minio.RemoveObjectOptions{})
		if err != nil {
			app.Logger().Error(fmt.Sprintf("object_storage/directory: error deleting `%s`", object.Key), "error", err)
			return 0, err
		}

		deletedFiles++
	}
	app.Logger().Info("object_storage/directory: deleted files", "count", deletedFiles)
	return deletedFiles, nil
}
