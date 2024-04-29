package repo

import (
	"go-fs/common/geeorm"
	"go-fs/internal/conf"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const driver = "mysql"

var (
	defaultEngine   *geeorm.Engine
	fileinfoSession *geeorm.Session
	userinfoSession *geeorm.Session
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
	userinfoSession = defaultEngine.NewSession().SetTableName("userinfo")
	return nil
}

func bindVars(num int) string {
	temp := []string{}
	for i := 0; i < num; i++ {
		temp = append(temp, "?")
	}
	return strings.Join(temp, ",")
}
