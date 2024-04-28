package conf

import "fmt"

type Config struct {
	ServerCfg  *ServerConfig   `yaml:"server"`
	ReadDBCfg  *DataBaseConfig `yaml:"readDB"`
	WriteDBCfg *DataBaseConfig `yaml:"writeDB"`
}

type ServerConfig struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	Network string `yaml:"network"`
}

//mysql dsn格式
//涉及参数:
//username   数据库账号
//password   数据库密码
//host       数据库连接地址，可以是Ip或者域名
//port       数据库端口
//Dbname     数据库名

//username:password@tcp(host:port)/Dbname?charset=utf8&parseTime=True&loc=Local

// 填上参数后的例子
// username = root
// password = 123456
// host     = localhost
// port     = 3306
// Dbname   = tizi365
// 后面K/V键值对参数含义为：
//
//	charset=utf8 客户端字符集为utf8
//	parseTime=true 支持把数据库datetime和date类型转换为golang的time.Time类型
//	loc=Local 使用系统本地时区
type DataBaseConfig struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	User   string `yaml:"user"`
	PassWD string `yaml:"password"`
	DBName string `yaml:"dbname"`
}

func (cfg *DataBaseConfig) DataSourceName() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.User,
		cfg.PassWD,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
}
