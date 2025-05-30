package configz

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"reflect"

	awsssm "github.com/PaddleHQ/go-aws-ssm"
	"github.com/go-playground/assert/v2"
	"github.com/spf13/viper"
)

func LoadFromDotenv(filename string, out any) error {
	if filename == "" {
		return errors.New("filename cannot be empty")
	}

	val := reflect.ValueOf(out)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("output parameter must be a pointer")
	}

	v := viper.New()
	v.SetConfigFile(filename)
	v.SetConfigType("env")
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	if err := v.Unmarshal(out); err != nil {
		return err
	}

	slog.Info("Configuration loaded from", "filename", filename)

	return nil
}

// LoadFromAWSParameterStore is a
func LoadFromAWSParameterStore(path, prefix string, out any) error {
	val := reflect.ValueOf(out)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("output parameter must be a pointer")
	}

	pmstore, err := awsssm.NewParameterStore()
	if err != nil {
		return err
	}

	params, err := pmstore.GetAllParametersByPath(path, true)
	if err != nil {
		return err
	}

	if assert.IsEqual(len(params.GetAllValues()), 0) {
		return errors.New("no parameters found")
	}

	paramJson, err := io.ReadAll(params)
	if err != nil {
		return err
	}

	var configMap = make(map[string]any)
	if err := json.Unmarshal(paramJson, &configMap); err != nil {
		return err
	}

	if err := decode(prefix, configMap, out); err != nil {
		return err
	}

	slog.Info("Configuration loaded from", "path", path)

	return nil
}
