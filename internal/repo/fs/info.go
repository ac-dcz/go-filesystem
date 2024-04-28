package fs

import (
	"crypto/sha256"
	"fmt"
	"strconv"
)

type FileInfo struct {
	FID          string
	Name         string //文件名
	UploadTS     int64
	LastModifyTS int64
	Path         string
}

func NewFileInfo(name, path string, upts, lmts int64) *FileInfo {
	info := &FileInfo{
		Name:         name,
		UploadTS:     upts,
		LastModifyTS: lmts,
		Path:         path,
	}
	key := sha256.Sum256([]byte(strconv.Itoa(int(upts))))
	info.FID = fmt.Sprintf("%x", key)
	return info
}

func (info FileInfo) Hash() string {
	return info.FID
}
