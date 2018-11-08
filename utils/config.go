package utils

import (
	"github.com/OhBonsai/yogo/model"
	"github.com/OhBonsai/yogo/mlog"
	"os"
	"io"
	"io/ioutil"
	"encoding/json"
	"github.com/spf13/viper"
	"reflect"
	"fmt"
	"strings"
	"path/filepath"
	"bytes"
)

const (
	LOG_FILENAME  = "yogo.log"
)

func LoadConfig(fileName string) (*model.Config, string, map[string]interface{}, *model.AppError) {
	var configPath string

	if path, err := EnsureConfigFile(fileName); err != nil {
		appErr := model.NewAppError("LoadConfig", "utils.config.load_config.opening.panic", map[string]interface{}{"Filename": fileName, "Error": err.Error()}, "", 0)
		return nil, "", nil, appErr
	}else {
		configPath = path
	}

	config, envConfig, err := ReadConfigFile(configPath, true, false)
	if err != nil {
		appErr := model.NewAppError("LoadConfig", "utils.config.load_config.decoding.panic", map[string]interface{}{"Filename": fileName, "Error": err.Error()}, "", 0)
		return nil, "", nil, appErr
	}

	config.SetDefaults()
	if err := config.IsValid(); err != nil {
		return nil, "", nil, err
	}

	return config, configPath, envConfig, nil
}


func EnsureConfigFile(fileName string) (string, error) {
	return fileName, nil
}


func ReadConfigFile(path string, allowEnvOverrides bool, allowConsulOverrides bool) (*model.Config, map[string]interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	defer f.Close()
	return ReadConfigStream(f, allowEnvOverrides, allowConsulOverrides)
}

func ReadConfigStream(r io.Reader, allowEnvOverrides bool, allowConsulOverrides bool) (*model.Config, map[string]interface{}, error) {
	configData, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, nil, err
	} else {
		var rawConfig interface{}
		if err := json.Unmarshal(configData, &rawConfig); err != nil {
			return nil, nil, HumanizeJsonError(err, configData)
		}
	}

	var config model.Config
	v := newViper()
	if err := v.ReadConfig(bytes.NewReader(configData)); err != nil {
		return nil, nil, err
	}

	if allowEnvOverrides {
		v.SetEnvPrefix("mm")
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		v.AutomaticEnv()
	}

	if allowConsulOverrides {
		v.AddRemoteProvider("consul", "localhost:8500", "MY_CONSUL_KEY")
		v.SetConfigType("json") // Need to explicitly set this to json
		if err := viper.ReadRemoteConfig(); err != nil {
			return nil, nil, err
		}
	}

	unmarshalError := v.Unmarshal(&config)
	return &config, nil, unmarshalError
}

func newViper() *viper.Viper {
	v := viper.New()
	v.SetConfigType("json")

	defaults := getDefaultsFromStruct(model.Config{})

	for key, value := range defaults {
		v.SetDefault(key, value)
	}

	return v
}

func getDefaultsFromStruct(s interface{}) map[string]interface{} {
	return flattenStructToMap(structToMap(reflect.TypeOf(s)))
}

func flattenStructToMap(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})

	for key, value := range in {
		if valueAsMap, ok := value.(map[string]interface{}); ok {
			sub := flattenStructToMap(valueAsMap)

			for subKey, subValue := range sub {
				out[key+"."+subKey] = subValue
			}
		} else {
			out[key] = value
		}
	}

	return out
}



func structToMap(t reflect.Type) (out map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			mlog.Error(fmt.Sprintf("Panicked in structToMap. This should never happen. %v", r))
		}
	}()

	if t.Kind() != reflect.Struct {
		// Should never hit this, but this will prevent a panic if that does happen somehow
		return nil
	}

	out = map[string]interface{}{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		var value interface{}

		switch field.Type.Kind() {
		case reflect.Struct:
			value = structToMap(field.Type)
		case reflect.Ptr:
			indirectType := field.Type.Elem()

			if indirectType.Kind() == reflect.Struct {
				// Follow pointers to structs since we need to define defaults for their fields
				value = structToMap(indirectType)
			} else {
				value = nil
			}
		default:
			value = reflect.Zero(field.Type).Interface()
		}

		out[field.Name] = value
	}

	return
}


func MloggerConfigFromLoggerConfig(s *model.LogSettings) *mlog.LoggerConfiguration {
	return &mlog.LoggerConfiguration{
		EnableConsole: s.EnableConsole,
		ConsoleJson:   *s.ConsoleJson,
		ConsoleLevel:  strings.ToLower(s.ConsoleLevel),
		EnableFile:    s.EnableFile,
		FileJson:      *s.FileJson,
		FileLevel:     strings.ToLower(s.FileLevel),
		FileLocation:  GetLogFileLocation(s.FileLocation),
	}
}

func GetLogFileLocation(fileLocation string) string {
	if fileLocation == "" {
		fileLocation, _ = FindDir("logs")
	}

	return filepath.Join(fileLocation, LOG_FILENAME)
}

func FindDir(dir string) (string, bool) {
	for _, parent := range []string{".", "..", "../..", "../../.."} {
		foundDir, err := filepath.Abs(filepath.Join(parent, dir))
		if err != nil {
			continue
		} else if _, err := os.Stat(foundDir); err == nil {
			return foundDir, true
		}
	}
	return "./", false
}
