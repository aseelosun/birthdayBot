package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Configuration struct {
	TelegramBotToken string `json:"tgBotToken"`
	DbConfig         Db     `json:"db"`
}

type Db struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

func LoadConfiguration(filePath string) (Configuration, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var configs Configuration
	json.Unmarshal(byteValue, &configs)
	return configs, err

}
