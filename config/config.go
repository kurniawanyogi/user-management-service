package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strconv"
	"strings"
	"user-management-service/common"
)

var (
	Cold ColdFlat
	Hot  HotFlat
)

type ColdFlat struct {
	AppEnv      string `json:"appEnv" yaml:"appEnv"`
	AppName     string `json:"appName" yaml:"appName"`
	AppPort     uint   `json:"appPort" yaml:"appPort"`
	AppTimezone string `json:"appTimezone" yaml:"appTimezone"`
	AppApiKey   string `json:"appApiKey" yaml:"appApiKey"`
	SecretKey   string `json:"secretKey" yaml:"secretKey"`

	// DB Mysql
	DBMysqlDriver                string `json:"dbMysqlDriver" yaml:"dbMysqlDriver"`
	DBMysqlHost                  string `json:"dbMysqlHost" yaml:"dbMysqlHost"`
	DBMysqlPort                  int    `json:"dbMysqlPort" yaml:"dbMysqlPort"`
	DBMysqlDBName                string `json:"dbMysqlDbName" yaml:"dbMysqlDbName"`
	DBMysqlUser                  string `json:"dbMysqlUser" yaml:"dbMysqlUser"`
	DBMysqlPassword              string `json:"dbMysqlPassword" yaml:"dbMysqlPassword"`
	DBMysqlSSLMode               string `json:"dbMysqlSslMode" yaml:"dbMysqlrSslMode"`
	DBMysqlMaxOpenConnections    int    `json:"dbMysqlMaxOpenConnections" yaml:"dbMysqlMaxOpenConnections"`
	DBMysqlMaxLifeTimeConnection int    `json:"dbMysqlMaxLifeTimeConnection" yaml:"dbMysqlMaxLifeTimeConnection"`
	DBMysqlMaxIdleConnections    int    `json:"dbMysqlMaxIdleConnections" yaml:"dbMysqlMaxIdleConnections"`
	DBMysqlMaxIdleTimeConnection int    `json:"dbMysqlMaxIdleTimeConnection" yaml:"dbMysqlMaxIdleTimeConnection"`
}

type HotFlat struct {
	AppDebug              bool   `json:"appDebug" yaml:"appDebug"`
	LoggerDebug           bool   `json:"loggerDebug" yaml:"loggerDebug"`
	ShutDownDelayInSecond uint64 `json:"shutDownDelayInSecond" yaml:"shutDownDelayInSecond"`
}

func Init() {
	err := BindFromFile(&Cold, "config.cold.json", ".")
	if err != nil {
		fmt.Printf("failed load cold config from file: %s", viper.ConfigFileUsed())
		panic(err)
	}

	err = BindFromFile(&Hot, "config.hot.json", ".")
	if err != nil {
		fmt.Printf("failed load hot config from file: %s", viper.ConfigFileUsed())
		panic(err)
	}
}

// BindFromFile load config from filename then assign to destination
func BindFromFile(dest any, filename string, paths ...string) error {
	v := viper.New()

	v.SetConfigType("json")
	v.SetConfigName(filename)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	for _, path := range paths {
		v.AddConfigPath(path)
	}

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	fmt.Printf("using config file: %s.\n", filename)

	err = v.Unmarshal(dest)
	if err != nil {
		return err
	}

	err = SetEnvFromConsulKV(v)
	if err != nil {
		fmt.Printf("failed to set env from file: %+v\n", err)
		return err
	}

	return nil
}

func SetEnvFromConsulKV(v *viper.Viper) error {
	env := make(map[string]any)

	err := v.Unmarshal(&env)
	if err != nil {
		fmt.Printf("failed to unmarshal config env: %+v\n", err)
		return err
	}

	for k, v := range env {
		var (
			valOf = reflect.ValueOf(v)
			val   string
		)

		switch valOf.Kind() {
		case reflect.String:
			val = valOf.String()
		case reflect.Int:
			val = strconv.Itoa(int(valOf.Int()))
		case reflect.Uint:
			val = strconv.Itoa(int(valOf.Uint()))
		case reflect.Float64:
			val = strconv.Itoa(int(valOf.Float()))
		case reflect.Float32:
			val = strconv.Itoa(int(valOf.Float()))
		case reflect.Bool:
			val = strconv.FormatBool(valOf.Bool())
		}

		os.Setenv(k, val)
	}

	return nil
}

func LoadTimeZoneFromEnv() string {
	tz := os.Getenv(common.Timezone)
	if len(tz) <= 0 {
		return common.DefaultTimeZone
	}
	return tz
}
