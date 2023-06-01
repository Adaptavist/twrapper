package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// TwrapperConfig model

// environmentDecoder detected environment variable references in teh config and substributes them
// Mostly stolen from: https://stackoverflow.com/a/68763349
func environmentDecoder(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	prefix, suffix := "${", "}"
	if f.Kind() == reflect.String {
		stringData := data.(string)
		for {
			// Find prefix position, and bail if its not found
			prePos := strings.Index(stringData, prefix)
			if prePos < 0 {
				break
			}

			// Get suffix position, and error if its not found
			sufPos := strings.Index(stringData, suffix)
			if sufPos < 0 {
				return nil, fmt.Errorf("failed to find closing substribution bracked in \"%s\"", stringData)
			}

			strKey := stringData[prePos : sufPos+1]
			envKey := strKey[2 : len(strKey)-1]
			value, ok := os.LookupEnv(envKey)
			if !ok {
				return nil, fmt.Errorf("failed to find environment variable \"%s\"", envKey)
			}

			stringData = strings.ReplaceAll(stringData, strKey, value)
		}
		return stringData, nil
		// if strings.HasPrefix(stringData, prefix) && strings.HasSuffix(stringData, suffix) {
		// 	envVarName := strings.TrimPrefix(strings.TrimSuffix(stringData, suffix), prefix)
		// 	envVarValue, ok := os.LookupEnv(envVarName)
		// 	if !ok {
		// 		return nil, fmt.Errorf("%s is not found in environment", envVarName)
		// 	}
		// 	return envVarValue, nil
		// }
	}
	return data, nil
}

// bind to environment
func bind() {
	viper.SetEnvPrefix("TW")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// aws
	viper.BindEnv("aws.account_id")
	viper.BindEnv("aws.role_arn")
	viper.BindEnv("aws.role_tf_var")
	viper.BindEnv("aws.try_role_names")

	// backend general
	viper.BindEnv("backend.type")

	// backend s3
	viper.BindEnv("backend.props.bucket")
	viper.BindEnv("backend.props.dynamodb_table")
	viper.BindEnv("backend.props.encrypt")
	viper.BindEnv("backend.props.key")
	viper.BindEnv("backend.props.region")

	// backend path
	viper.BindEnv("backend.props.path")

	// cloud
	viper.BindEnv("cloud.organization")
	viper.BindEnv("cloud.project")
	viper.BindEnv("cloud.workspace")

	// misc
	viper.BindEnv("required_vars")
}

func read(logger *zap.Logger) (conf *Config, err error) {
	config := Config{
		Logger: logger,
	}

	err = viper.ReadInConfig()

	if err != nil {
		logger.Debug("failed to read config files", zap.Error(err))
	}

	err = viper.Unmarshal(&config, viper.DecodeHook(environmentDecoder))

	if err != nil {
		return &config, err
	}

	logger.Debug("loaded config", zap.Any("config", config))

	return &config, nil
}

// NewTest configuration object
func NewTest(logger *zap.Logger) (*Config, error) {
	viper.AddConfigPath("$HOME") // User home directory
	viper.SetConfigName("test.twrapper.yml")
	viper.SetConfigType("yaml")
	bind()
	return read(logger)
}

// New configuration object
func New(logger *zap.Logger) (*Config, error) {
	viper.AddConfigPath(".")     // Current working directory
	viper.AddConfigPath("$HOME") // User home directory
	viper.SetConfigName("terraform.twrapper.yml")
	viper.SetConfigType("yaml")
	bind()
	return read(logger)
}
