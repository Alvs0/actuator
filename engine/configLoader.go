package engine

import "github.com/spf13/viper"

func LoadConfig(configName, pathToConfigFile string, destConfigMap interface{}) (err error) {
	viper.SetConfigName(configName)
	viper.AddConfigPath(pathToConfigFile)

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&destConfigMap)

	return
}