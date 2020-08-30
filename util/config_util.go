package util

import "github.com/spf13/viper"

func GetDBConfig() (interface{}, interface{}, interface{}, interface{}, interface{}, interface{}) {
	viper.SetConfigName("mysql_config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic("配置文件错误：" + err.Error())
	}

	user := viper.Get("db.user")
	pass := viper.Get("db.pass")
	host := viper.Get("db.host")
	port := viper.Get("db.port")
	dbname := viper.Get("db.dbname")
	charset := viper.Get("db.charset")

	return user, pass, host, port, dbname, charset
}
