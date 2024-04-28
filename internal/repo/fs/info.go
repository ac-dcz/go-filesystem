package fs

import "time"

type FileInfo struct {
	FHash        string
	FName        string //文件名
	Fsize        int64
	UploadTS     int64
	LastModifyTS int64
	LocalPath    string
	Status       int //文件当前状态(0可用/1禁用/2删除)
}

func NewFileInfo(fhash, fname, fpath string, fsize int64) *FileInfo {
	info := &FileInfo{
		FHash:        fhash,
		FName:        fname,
		Fsize:        fsize,
		LocalPath:    fpath,
		UploadTS:     time.Now().Unix(),
		LastModifyTS: time.Now().Unix(),
		Status:       0,
	}
	return info
}
