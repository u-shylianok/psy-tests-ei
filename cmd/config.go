package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Tests []Test `json:"Тесты"`
}

type Test struct {
	Name            string           `json:"Название теста"`
	Path            string           `json:"Путь к файлу"`
	AnswerOptions   []AnswerOption   `json:"Варианты ответов"`
	Interpretations []Interpretation `json:"Интерпретация"`
}

type AnswerOption struct {
	Option       string `json:"Вариант ответа"`
	DirectValue  int    `json:"Прямое значение"`
	ReverseValue int    `json:"Обратное значение"`
}

type Interpretation struct {
	Scale             string `json:"Шкала"`
	DirectStatements  []int  `json:"Прямые утверждения"`
	ReverseStatements []int  `json:"Обратные утверждения"`
}

func loadConfig(fileName string) *Config {

	// Open config.json
	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	// Read config file
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	config := &Config{}
	json.Unmarshal(byteValue, config)

	return config
}
