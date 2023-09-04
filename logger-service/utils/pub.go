package utils

import (
	"context"
	"errors"
	"log"
	"os/exec"

	"github.com/spf13/viper"
)

func ViperEnvVariable(key string) (string, error) {

	viper.SetConfigFile("../.././.env")
	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
		return "", err
	}

	value, ok := viper.Get(key).(string)

	if !ok {
		return "", errors.New("Invalid type assertion")
	}

	return value, nil
}

func RunCommand(ctx context.Context, str string) error {
	cmd := exec.Command(str)

	// err := cmd.Run()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	result, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(string(result))

	return nil

}
