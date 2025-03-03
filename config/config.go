package config

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
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
	consulURL := os.Getenv("CONSUL_HTTP_URL")

	err := BindFromFile(&Cold, "config.cold.json", ".")
	if err != nil {
		fmt.Printf("failed load cold config from file: %s", viper.ConfigFileUsed())
		err = BindFromConsul(
			&Cold,
			consulURL,
			fmt.Sprintf("%s/%s", os.Getenv("CONSUL_HTTP_KEY"), "cold"),
		)
		if err != nil {
			panic(err)
		}
	}

	err = BindFromFile(&Hot, "config.hot.json", ".")
	if err != nil {
		fmt.Printf("failed load hot config from file: %s", viper.ConfigFileUsed())

		interval, err := LoadConsulIntervalFromEnv()
		if err != nil {
			panic(err)
		}

		err = BindAndWatchFromConsul(
			&Hot,
			consulURL,
			fmt.Sprintf("%s/%s", os.Getenv("CONSUL_HTTP_KEY"), "hot"),
			interval,
		)
		if err != nil {
			panic(err)
		}
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

// BindFromConsul load config from remote consul then assign to destination
func BindFromConsul(dest any, endPoint, path string) error {
	v := viper.New()

	v.SetConfigType("json")
	err := v.AddRemoteProvider("consul", endPoint, path)
	if err != nil {
		return err
	}

	err = v.ReadRemoteConfig()
	if err != nil {
		return err
	}

	fmt.Printf("using config from consul: %s/%s.\n", endPoint, path)

	err = v.Unmarshal(dest)
	if err != nil {
		fmt.Printf("failed to unmarshal config dest: %+v\n", err)
		return err
	}

	err = SetEnvFromConsulKV(v)
	if err != nil {
		fmt.Printf("failed to set env from consul: %+v\n", err)
		return err
	}

	return nil
}

// BindAndWatchFromConsul load and watch config from remote consul then assign to destination
func BindAndWatchFromConsul(dest any, endPoint, path string, interval int) error {
	location, err := time.LoadLocation(LoadTimeZoneFromEnv())
	if err != nil {
		fmt.Printf("failed load location: %s\n", location)
	}

	err = BindFromConsul(dest, endPoint, path)
	if err != nil {
		fmt.Printf("failed reloading consul: %+v\n", err)
		return err
	}

	scheduler := gocron.NewScheduler(location)
	_, err = scheduler.Every(interval).Seconds().Do(func() {
		er := BindFromConsul(dest, endPoint, path)
		if er != nil {
			fmt.Printf("failed reloading consul: %+v\n", er)
		}
	})

	if err != nil {
		fmt.Printf("failed run scheduler specify jobFunc: %+v\n", err)
		return err
	}

	scheduler.StartAsync()

	return nil
}

// LoadConsulIntervalFromEnv get interval value for loading config from consul
func LoadConsulIntervalFromEnv() (int, error) {
	fromEnv := os.Getenv(common.ConsulWatchInterval)
	if len(fromEnv) <= 0 {
		return common.DefaultLoadConsulInterval, nil
	}

	interval, err := strconv.Atoi(fromEnv)
	if err != nil {
		fmt.Printf("failed convert %s: %+v\n", common.ConsulWatchInterval, err)
		return 0, err
	}

	return interval, nil
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
