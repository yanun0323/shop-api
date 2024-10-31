package config

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	conf Config
	once sync.Once
)

type Config struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`

	Http struct {
		Port string `json:"port"`
	} `json:"http"`

	Token struct {
		Secret     string `json:"secret"`
		Expiration struct {
			Refresh time.Duration `json:"refresh"`
			Access  time.Duration `json:"access"`
		} `json:"expiration"`
	} `json:"token"`

	OTP struct {
		Expiration time.Duration `json:"expiration"`
	} `json:"otp"`

	Product struct {
		Recommendation struct {
			Expiration time.Duration `json:"expiration"`
		} `json:"recommendation"`
	} `json:"product"`

	MySQL struct {
		Host         string `json:"host"`
		Port         int    `json:"port"`
		User         string `json:"user"`
		Password     string `json:"password"`
		Database     string `json:"database"`
		MaxIdleConns int    `json:"maxIdleConns"`
		MaxOpenConns int    `json:"maxOpenConns"`
		MaxLifetime  int    `json:"maxLifetime"`
	} `json:"mysql"`

	Redis struct {
		Addr         string `json:"addr"`
		Username     string `json:"username"`
		Password     string `json:"password"`
		Database     int    `json:"database"`
		MaxIdleConns int    `json:"maxIdleConns"`
		MaxOpenConns int    `json:"maxOpenConns"`
	} `json:"redis"`
}

func Load() Config {
	return conf
}

func Init(cfgName string, relativePaths ...string) error {
	var (
		dir     string
		execDir string
		err     error
	)

	once.Do(func() {
		if len(cfgName) == 0 {
			err = errors.Errorf("empty config name")
			return
		}

		dir, err = os.Getwd()
		if err != nil {
			err = errors.Errorf("get wd: %+v", err)
			return
		}

		execDir, err = os.Executable()
		if err != nil {
			err = errors.Errorf("get executable path: %+v", err)
			return
		}

		for _, p := range relativePaths {
			path := filepath.Join(dir, p)
			excPath := filepath.Join(execDir, p)

			viper.AddConfigPath(path)
			viper.AddConfigPath(excPath)
		}

		viper.AddConfigPath(".")
		viper.SetConfigName(cfgName)
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.SetConfigType("yaml")

		err = viper.ReadInConfig()
		if err != nil {
			err = errors.Errorf("read in config: %+v", err)
			return
		}

		if err = viper.Unmarshal(&conf); err != nil {
			err = errors.Errorf("unmarshal config: %+v", err)
			return
		}

	})

	return err
}
