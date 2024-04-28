package handler

import (
	"encoding/hex"
	"go-fs/common/geeweb"
	"go-fs/internal/repo"
	"go-fs/internal/repo/fs"
	"go-fs/internal/util"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"
)

func createFileInfo(f multipart.File, h *multipart.FileHeader) (*fs.FileInfo, error) {
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		log.Println(err)
		return nil, err
	}
	data, err := io.ReadAll(f)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	sha1 := util.NewHasher().Add(data).Sum()
	fsha1 := hex.EncodeToString(sha1)
	info := fs.NewFileInfo(fsha1, h.Filename, StaticRoot, h.Size)
	if err := repo.InsertFileInfo(info); err != nil {
		log.Println(err)
		return nil, err
	}
	return info, nil
}

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
			if info, err := createFileInfo(f, h); err != nil {
				http.Error(c.W, err.Error(), http.StatusInternalServerError)
			} else {
				c.String(http.StatusOK, "Successfuly upload")
				log.Printf("upload file[%s] successfully\n", info.FHash)
			}
		}
	}
}

// UpdateHandle: POST /file/update
func UpdateHandle(c *geeweb.Context) {
	hash, newname := c.PostForm("hash"), c.PostForm("newname")
	if infos, err := repo.SelectFileInfo(hash); err == nil {
		for _, info := range infos {
			if err := os.Rename(info.LocalPath+info.FName, info.LocalPath+newname); err != nil {
				http.Error(c.W, err.Error(), http.StatusInternalServerError)
			} else {
				info.LastModifyTS = time.Now().Unix()
				info.FName = newname
				repo.UpdateFileInfo(info)
				c.String(http.StatusOK, "Successfully update")
			}
		}
	} else {
		http.Error(c.W, err.Error(), http.StatusInternalServerError)
	}
}

// QueryHandle: GET /file/query?hash=...
func QueryHandle(c *geeweb.Context) {
	hash := c.Query("hash")
	if infos, err := repo.SelectFileInfo(hash); err == nil {
		for _, info := range infos {
			c.JSON(http.StatusOK, info)
		}
	} else {
		http.Error(c.W, err.Error(), http.StatusInternalServerError)
	}
}

// DownloadHandle: GET /file/download?filename=...
func DownloadHandle(c *geeweb.Context) {
	hash := c.Query("hash")
	if infos, err := repo.SelectFileInfo(hash); err == nil {
		for _, info := range infos {
			f, err := os.Open(info.LocalPath + info.FName)
			if os.IsNotExist(err) {
				c.String(http.StatusOK, "%s not exists", hash)
			} else if err != nil {
				http.Error(c.W, err.Error(), http.StatusInternalServerError)
			} else {
				defer f.Close()
				c.AddHeader("Content-Disposition", "attachment; filename="+url.QueryEscape(info.FName))
				c.AddHeader("Content-Type", "application/octet-stream")
				io.Copy(c.W, f)
			}
		}
	} else {
		http.Error(c.W, err.Error(), http.StatusInternalServerError)
	}
}

// DeleteHandle: DELETE /file/delete?filename=...
func DeleteHandle(c *geeweb.Context) {
	hash := c.Query("hash")
	if infos, err := repo.SelectFileInfo(hash); err == nil {
		for _, info := range infos {
			info.LastModifyTS = time.Now().Unix()
			info.Status = 2
			if err := repo.UpdateFileInfo(info); err != nil {
				http.Error(c.W, err.Error(), http.StatusInternalServerError)
			} else {
				c.String(http.StatusOK, "Successfully Delete File %s\n", hash)
			}
		}
	} else {
		http.Error(c.W, err.Error(), http.StatusInternalServerError)
	}
}

// FileListHandle: Get /file/list
func FileListHandle(c *geeweb.Context) {
	h := make(geeweb.H)
	if infos, err := repo.SelectAllFileInfo(); err != nil {
		http.Error(c.W, err.Error(), http.StatusInternalServerError)
	} else {
		for _, info := range infos {
			h[info.FName] = info
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
