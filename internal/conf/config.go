package conf

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

type DataBaseConfig struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	User   string `yaml:"user"`
	PassWD string `yaml:"password"`
	DBName string `yaml:"dbname"`
}
