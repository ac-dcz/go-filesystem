package repo

import (
	"fmt"
	"go-fs/common/geeorm"
	"go-fs/internal/conf"
	"go-fs/internal/repo/fs"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const driver = "mysql"

var (
	defaultEngine   *geeorm.Engine
	fileinfoSession *geeorm.Session
)

func SetDataBaseConfig(cfg *conf.DataBaseConfig) error {
	engine, err := geeorm.NewEngine(driver, cfg.DataSourceName())
	if err != nil {
		log.Println(err)
		return err
	} else {
		log.Printf("Successfully Connect DataBase [%v]\n", cfg.DataSourceName())
	}
	defaultEngine = engine
	fileinfoSession = defaultEngine.NewSession().SetTableName("filemeta")
	return nil
}

func bindVars(num int) string {
	temp := []string{}
	for i := 0; i < num; i++ {
		temp = append(temp, "?")
	}
	return strings.Join(temp, ",")
}

func InsertFileInfo(info *fs.FileInfo) error {
	sql := fmt.Sprintf(
		"insert into %s (file_sha1,file_name,file_size,file_local_path,upload_timestamp,last_modify_timestamp,status) values (%s)",
		fileinfoSession.TableName(),
		bindVars(7),
	)
	argv := []any{info.FHash, info.FName, info.Fsize, info.LocalPath, info.UploadTS, info.LastModifyTS, info.Status}
	_, err := fileinfoSession.Raw(sql, argv...).Exex()
	if err != nil {
		log.Println(err)
	}
	return err
}

func DeleteFileInfo(fhash ...string) error {
	n := len(fhash)
	if n == 0 {
		return nil
	}
	sql := fmt.Sprintf("update  %s set status = 2 where file_sha1 in (%s)", fileinfoSession.TableName(), bindVars(n))
	vars := make([]any, 0)
	for _, hash := range fhash {
		vars = append(vars, hash)
	}
	_, err := fileinfoSession.Raw(sql, vars...).Exex()
	if err != nil {
		log.Println(err)
	}
	return err
}

func UpdateFileInfo(info *fs.FileInfo) error {
	sql := fmt.Sprintf("update %s set file_name = ?,last_modify_timestamp = ?,status = ? where file_sha1 = ?", fileinfoSession.TableName())
	vars := []any{info.FName, info.LastModifyTS, info.Status, info.FHash}
	_, err := fileinfoSession.Raw(sql, vars...).Exex()
	if err != nil {
		log.Println(err)
	}
	return err
}

func SelectFileInfo(fhash ...string) (infos []*fs.FileInfo, err error) {
	sql := fmt.Sprintf(
		"select file_sha1,file_name,file_size,file_local_path,upload_timestamp,last_modify_timestamp,status from %s where file_sha1 = ?",
		fileinfoSession.TableName(),
	)
	for _, hash := range fhash {
		row := fileinfoSession.Raw(sql, hash).QueryRow()
		if err = row.Err(); err != nil {
			log.Println(err)
			return infos, err
		}
		info := &fs.FileInfo{}
		err = row.Scan(&info.FHash, &info.FName, &info.LocalPath, &info.UploadTS, &info.LastModifyTS, &info.Status)
		if err != nil {
			return infos, err
		}
		infos = append(infos, info)
	}
	return infos, err
}

func SelectAllFileInfo() (infos []*fs.FileInfo, err error) {
	sql := fmt.Sprintf(
		"select file_sha1,file_name,file_size,file_local_path,upload_timestamp,last_modify_timestamp,status from %s",
		fileinfoSession.TableName(),
	)
	if rows, err := fileinfoSession.Raw(sql).QueryRows(); err != nil {
		log.Println(err)
		return nil, err
	} else {
		for rows.Next() {
			info := &fs.FileInfo{}
			err = rows.Scan(&info.FHash, &info.FName, &info.LocalPath, &info.UploadTS, &info.LastModifyTS, &info.Status)
			if err != nil {
				return infos, err
			}
			infos = append(infos, info)
		}
		return infos, err
	}
}
