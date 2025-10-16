package object_store

import (
	"context"
	"log/slog"
	"path"
	"sort"
	"strings"

	"github.com/minio/minio-go/v7"
)

type ObjectStore struct {
	client  *minio.Client
	Buckets []Bucket
}

type Listing struct {
	Path  string           `json:"path"`
	IsDir bool             `json:"dir"`
	Size  int64            `json:"size"`
	Stats minio.ObjectInfo `json:"stats"`
}

type Bucket struct {
	Name string `json:"name"`
	info minio.BucketInfo
}

type AppContext interface {
	Logger() *slog.Logger
	Context() context.Context
}

type DirectoryItem struct {
	Path     string            `json:"path"`
	IsDir    bool              `json:"dir"`
	IsFile   bool              `json:"file"`
	Depth    int               `json:"depth"`
	Stats    *minio.ObjectInfo `json:"stats,omitempty"`
	Children []DirectoryItem   `json:"children"`
}

func sortListingByPath(listing []Listing) []Listing {
	sort.Slice(listing, func(i, j int) bool {
		return listing[i].Path < listing[j].Path
	})
	return listing
}

func (d *DirectoryItem) AddItemFromRoot(item Listing) {
	//From the root, splits the path, creates the missing nodes, and adds the item to the bottom of the tree

	parts := strings.Split(item.Path, "/")
	if len(parts) == 0 {
		return
	}

	current := d
	for i, part := range parts {
		if part == "" {
			continue
		}

		isLast := i == len(parts)-1
		if isLast {
			// Add file at the final node
			current.AddFile(item)
			return
		}

		// Traverse or create directory
		current = current.ExistsOrAddDirectory(part)
	}

}

func (d *DirectoryItem) ExistsOrAddDirectory(path string) *DirectoryItem {

	for i := range d.Children {
		if d.Children[i].IsDir && d.Children[i].Path == path {
			return &d.Children[i]
		}
	}

	newDir := DirectoryItem{
		Path:  path,
		IsDir: true,
		Depth: d.Depth + 1,
	}
	d.Children = append(d.Children, newDir)
	return &d.Children[len(d.Children)-1]

}
func (d *DirectoryItem) AddFile(item Listing) *DirectoryItem {

	fileName := path.Base(item.Path)

	// Check if file already exists
	for _, child := range d.Children {
		if !child.IsDir && child.Path == fileName {
			return &child // File already exists
		}
	}

	newFile := DirectoryItem{
		Path:  fileName,
		IsDir: false,
		Depth: d.Depth + 1,
		Stats: &item.Stats,
	}
	d.Children = append(d.Children, newFile)
	return &newFile
}
