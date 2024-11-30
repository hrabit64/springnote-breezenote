package config

import (
	"github.com/gin-gonic/gin"
	"github.com/hrabit64/springnote-breezenote/pkg/utils"
	"github.com/spf13/viper"
	"io"
	"os"
	"path"
)

type Config struct {
	Profile             string `mapstructure:"PROFILE"`
	SavePath            string `mapstructure:"SAVE_PATH"`
	ConvertTargetFormat string `mapstructure:"CONVERT_TARGET_FORMAT"`
	MaxImageLen         int    `mapstructure:"MAX_IMAGE_LEN"`
	AllowImageUploadUid string `mapstructure:"ALLOW_IMAGE_UPLOAD_UID"`
	DBConnURL           string `mapstructure:"DB_CONN_URL"`
	FireBaseConfig      string `mapstructure:"FIRE_BASE_CONFIG"`
}

var RootConfig *Config

func loadConfig() (*Config, error) {
	env := Config{}

	profile := os.Getenv("BREEZENOTE_PROFILE")
	viper.AddConfigPath(utils.GetRootPath())
	viper.SetConfigName(".env." + profile)
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}

// SetupConfig 설정파일을 로드하고, router 설정을 합니다.
func SetupConfig() error {
	env, err := loadConfig()
	if err != nil {
		return err
	}

	RootConfig = env

	if RootConfig.Profile == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.DisableConsoleColor()

	f, _ := os.Create("./logs/gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	return nil
}

// GetSavePath 이미지 저장 경로를 반환합니다.
func GetSavePath() string {
	return path.Join(utils.GetRootPath(), RootConfig.SavePath)
}
