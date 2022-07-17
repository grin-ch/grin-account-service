package cfg

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

const (
	SERVER_NAME = "GRIN_ACCOUNT_SERVICE"
)

// 获取服务配置
func GetConfig() *ServerConfig {
	return &ServerConfig{
		Name: viper.GetString("server.name"),
		Host: viper.GetString("server.host"),
		Port: viper.GetInt("server.port"),

		RegProvider: viper.GetString("registry.provider"),
		RegPath:     viper.GetString("registry.path"),
		RegEndpoint: viper.GetStringSlice("registry.endpoint"),
		RegTimeout:  viper.GetInt("registry.timeout"),

		LogPath:   viper.GetString("log.path"),
		LogLevel:  viper.GetInt("log.level"),
		LogColor:  viper.GetBool("log.color"),
		LogCaller: viper.GetBool("log.caller"),

		RedisAddr: viper.GetString("redis.addr"),
		RedisPass: viper.GetString("redis.pass"),
		RedisDB:   viper.GetInt("redis.db"),

		DbPort: viper.GetInt("database.port"),
		DbHost: viper.GetString("database.host"),
		DbName: viper.GetString("database.name"),
		DbUser: viper.GetString("database.user"),
		DbPass: viper.GetString("database.pass"),
	}
}

type ServerConfig struct {
	Name string
	Host string
	Port int

	RegProvider string
	RegPath     string
	RegEndpoint []string
	RegTimeout  int

	LogPath   string
	LogLevel  int
	LogColor  bool
	LogCaller bool

	RedisAddr string
	RedisPass string
	RedisDB   int

	DbPort int
	DbHost string
	DbName string
	DbUser string
	DbPass string
}

func (c *ServerConfig) Dsn() string {
	dsnFormat := "%s:%s@tcp(%s:%d)/%s"
	params := "?charset=utf8mb4&parseTime=True&loc=Local&timeout=30s"
	return fmt.Sprintf(dsnFormat, c.DbUser, c.DbPass, c.DbHost, c.DbPort, c.DbName) + params
}

// 注册中心地址
func (c *ServerConfig) RegistryHost() string {
	return viper.GetString("registry.endpoint")
}

// SetServerConfig 初始化服务配置
func SetServerConfig(file string, paths ...string) {
	for _, v := range paths {
		viper.AddConfigPath(v)
	}
	loadRemoteConfig(
		viper.GetString("registry.provider"),
		viper.GetString("registry.endpoint"),
		viper.GetString("registry.path"),
	)
	setFileConfig(file)
	setEnvConfig()
}

// LoadRemoteConfig 加载远程配置
func loadRemoteConfig(provider, endpoint, path string) {
	if provider != "" && endpoint != "" && path != "" {
		// 支持 json, toml, yaml, yml, properties, props, prop, env, dotenv
		viper.AddRemoteProvider(provider, endpoint, path)
		err := viper.ReadRemoteConfig()
		if err != nil {
			panic(err)
		}
	}
}

func setEnvConfig() {
	viper.SetEnvPrefix(SERVER_NAME)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

// file 需要携带扩展名
func setFileConfig(file string) {
	viper.SetConfigName(file)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
