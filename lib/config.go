package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

type ConfigLib struct {
}

var Config ConfigLib

func (c *ConfigLib) ParseJSONToBytes(filepath string, t any) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("failed to open config %s: %v", filepath, err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("unable to parse bytevalue of %s: %v", filepath, err)
	}
	json.Unmarshal(byteValue, t)
}
