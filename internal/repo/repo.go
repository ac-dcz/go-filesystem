package repo

import (
	"go-fs/common/store"
	"go-fs/internal/repo/fs"
)

var (
	FileMetaDB store.Store[fs.FileInfo]
)

func init() {
	FileMetaDB = store.NewMemoryDB[fs.FileInfo](store.DefaultMetaPath)
}
