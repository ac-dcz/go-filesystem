package handler

import (
	"go-fs/common/geeweb"
	"go-fs/internal/repo"
	"go-fs/internal/repo/fs"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const StaticRoot = "./static/"

// UploadHandle: POST /file/upload filedname = upload
func UploadHandle(c *geeweb.Context) {
	if err := c.R.ParseMultipartForm(1 << 20); err != nil {
		http.Error(c.W, err.Error(), http.StatusBadRequest)
	} else {
		if f, h, err := c.R.FormFile("upload"); err != nil {
			http.Error(c.W, err.Error(), http.StatusInternalServerError)
		} else {
			save, err := os.OpenFile(StaticRoot+h.Filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
			if err != nil {
				http.Error(c.W, err.Error(), http.StatusInternalServerError)
				return
			}
			defer save.Close()
			io.Copy(save, f)
			now := time.Now()
			info := fs.NewFileInfo(h.Filename, StaticRoot, now.Unix(), now.Unix())
			repo.FileMetaDB.Put(info.Hash(), *info) //存储文件元数据
			c.String(http.StatusOK, "Successfuly upload")
			log.Printf("upload file[%s] successfully\n", info.Hash())
		}
	}
}

// UpdateHandle: POST /file/update
func UpdateHandle(c *geeweb.Context) {
	hash, newname := c.PostForm("hash"), c.PostForm("newname")
	if info, ok := repo.FileMetaDB.Get(hash); ok {
		if err := os.Rename(info.Path+info.Name, info.Path+newname); err != nil {
			http.Error(c.W, err.Error(), http.StatusInternalServerError)
		} else {
			info.LastModifyTS = time.Now().Unix()
			info.Name = newname
			repo.FileMetaDB.Put(hash, info)
			c.String(http.StatusOK, "Successfully update")
		}
	} else {
		c.String(http.StatusOK, "not found file hash %s\n", hash)
	}
}

// QueryHandle: GET /file/query?hash=...
func QueryHandle(c *geeweb.Context) {
	hash := c.Query("hash")
	if info, ok := repo.FileMetaDB.Get(hash); ok {
		c.JSON(http.StatusOK, info)
	} else {
		c.String(http.StatusOK, "not found file hash %s\n", hash)
	}
}

// DownloadHandle: GET /file/download?filename=...
func DownloadHandle(c *geeweb.Context) {
	hash := c.Query("hash")
	if info, ok := repo.FileMetaDB.Get(hash); ok {
		f, err := os.Open(info.Path + info.Name)
		if os.IsNotExist(err) {
			c.String(http.StatusOK, "%s not exists", hash)
		} else if err != nil {
			http.Error(c.W, err.Error(), http.StatusInternalServerError)
		} else {
			defer f.Close()
			c.AddHeader("Content-Disposition", "attachment; filename="+url.QueryEscape(info.Name))
			c.AddHeader("Content-Type", "application/octet-stream")
			io.Copy(c.W, f)
		}
	} else {
		c.String(http.StatusOK, "not found file hash %s\n", hash)
	}
}

// DeleteHandle: DELETE /file/delete?filename=...
func DeleteHandle(c *geeweb.Context) {
	hash := c.Query("hash")
	if info, ok := repo.FileMetaDB.Get(hash); ok {
		if err := os.Remove(info.Path + info.Name); err != nil {
			http.Error(c.W, err.Error(), http.StatusInternalServerError)
		} else {
			c.String(http.StatusOK, "Successfully delete")
		}
	} else {
		c.String(http.StatusOK, "not found file hash %s\n", hash)
	}
}

// FileListHandle: Get /file/list
func FileListHandle(c *geeweb.Context) {
	keys := repo.FileMetaDB.Keys()
	h := make(geeweb.H)
	for _, key := range keys {
		if info, ok := repo.FileMetaDB.Get(key); ok {
			h[key] = info
		}
	}
	c.JSON(http.StatusOK, h)
}

// RegistryHandleFunc 注册路由
func RegistryHandleFunc(group *geeweb.RouterGroup) {
	group.AddRoute(http.MethodPost, "/file/upload", UploadHandle)
	group.AddRoute(http.MethodPost, "/file/update", UpdateHandle)
	group.AddRoute(http.MethodGet, "/file/query", QueryHandle)
	group.AddRoute(http.MethodGet, "/file/download", DownloadHandle)
	group.AddRoute(http.MethodDelete, "/file/delete", DeleteHandle)
	group.AddRoute(http.MethodGet, "/file/list", FileListHandle)
}
