package main

import "fmt"

type TestInfo struct {
	Test      Test
	Responses []Response
}

type Response struct {
	DateTime string
	Sex      string
	Answers  []string
}

func main() {
	config := *loadConfig("config.json")

	var testInfos = []TestInfo{}
	for _, test := range config.Tests {
		responses := *loadAnswers(test.Path)
		testInfos = append(testInfos, TestInfo{test, responses})
	}

	var errorIndexes = []int{}

	// testInfo - отдельно каждый тест и ответы на него
	for _, testInfo := range testInfos {
		options := testInfo.Test.AnswerOptions
		interpretations := testInfo.Test.Interpretations

		// .CSV
		csvResult := [][]string{}
		// Формируем шапку .CSV
		csvShapka := []string{"Время прохождения", "Пол"}
		for _, interpretation := range interpretations {
			max := ((len(interpretation.DirectStatements) * findDirectMax(options)) + (len(interpretation.ReverseStatements) * findReverseMax(options)))
			csvShapka = append(csvShapka, interpretation.Scale+" (макс - "+fmt.Sprint(max)+")")
		}
		// Закончили формировать шапку. Время прохождения индекс 0, пол 1, дальше по номеру интерпретации. 2+j
		csvResult = append(csvResult, csvShapka)

		// response - отдельный ответ отдельного человека на тест
		for _, response := range testInfo.Responses {
			fmt.Println("Время прохождения:", response.DateTime, "Пол:", response.Sex)

			// .CSV
			csvResponse := []string{response.DateTime, response.Sex}
			// .CSV

			// interpretation - у каждого теста в интерпретации есть разные шкалы, тут происходит деление, i - номер интерпретации
			for _, interpretation := range interpretations {
				fmt.Println("\tШкала:", interpretation.Scale)
				var sum int

				// questionNumber - номер вопроса в интерпретации прямых утверждений этой шкалы
				for _, questionNumber := range interpretation.DirectStatements {
					isExist, directValue, _ := getValues(options, response.Answers[questionNumber-1])
					if isExist {
						sum += directValue
					} else {
						errorIndexes = append(errorIndexes, questionNumber)
					}
				}
				// questionNumber - номер вопроса в интерпретации обратных утверждений этой шкалы
				for _, questionNumber := range interpretation.ReverseStatements {
					isExist, _, reverseValue := getValues(options, response.Answers[questionNumber-1])
					if isExist {
						sum += reverseValue
					} else {
						errorIndexes = append(errorIndexes, questionNumber)
					}
				}
				fmt.Println("\t\tРезультат:", sum)
				csvResponse = append(csvResponse, fmt.Sprint(sum))
			}
			csvResult = append(csvResult, csvResponse)
		}
		writeResults("Результат: "+testInfo.Test.Name+".result.csv", csvResult)
	}

	fmt.Println("Error count:", len(errorIndexes), "\nError indexes:", errorIndexes)
}

func getValues(options []AnswerOption, answer string) (bool, int, int) {
	var isExist bool
	var directValue int
	var reverseValue int

	for _, option := range options {
		if option.Option == answer {
			directValue = option.DirectValue
			reverseValue = option.ReverseValue
			isExist = true
			break
		}
	}

	return isExist, directValue, reverseValue
}

func findDirectMax(slice []AnswerOption) int {
	if len(slice) == 0 {
		return 0
	}

	max := slice[0].DirectValue
	for _, value := range slice {
		if value.DirectValue > max {
			max = value.DirectValue
		}
	}
	return max
}

func findReverseMax(slice []AnswerOption) int {
	if len(slice) == 0 {
		return 0
	}

	max := slice[0].ReverseValue
	for _, value := range slice {
		if value.ReverseValue > max {
			max = value.ReverseValue
		}
	}
	return max
}
