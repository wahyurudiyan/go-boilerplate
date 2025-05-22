package configz

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/go-playground/assert/v2"
	"github.com/spf13/viper"
)

func decode(prefix string, configMap map[string]any, out any) error {
	if !assert.IsEqual(prefix, "") {
		for k, val := range configMap {
			oldKey := strings.ToUpper(k)
			prefix := strings.ToUpper(prefix + "_")
			if strings.HasPrefix(oldKey, prefix) {
				newKey := strings.TrimPrefix(oldKey, prefix)
				configMap[newKey] = val
				delete(configMap, k)
			}
		}
	}

	viper.SetConfigType("json")
	configByte, err := json.Marshal(configMap)
	if err != nil {
		return err
	}

	bufConfig := bytes.NewReader(configByte)
	err = viper.ReadConfig(bufConfig)
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&out)
	if err != nil {
		return err
	}

	return nil
}
