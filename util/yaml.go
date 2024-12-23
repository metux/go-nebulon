package util

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadYaml(fn string, target interface{}) error {
	yamlFile, err := os.ReadFile(fn)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, target)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return err
	}

	return nil
}
