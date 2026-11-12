package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"time"
)

type GameResult struct {
	Date     string `json:"date"`
	Result   string `json:"result"`
	Attempts int    `json:"attempts"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var maxNumber int
	var maxAttempts int
	fmt.Println("Игра 'Угадай число'")
	for {
		fmt.Println("выбери сложность: Easy/Medium/Hard")
		if scanner.Scan() {
			answer := scanner.Text()
			switch answer {
			case "Easy":
				maxNumber = 50
				maxAttempts = 15
				fmt.Println("Угадай число от 1 до 50 за 15 попыток")
			case "Medium":
				maxNumber = 100
				maxAttempts = 10
				fmt.Println("Угадай число от 1 до 100 за 10 попыток")
			case "Hard":
				maxNumber = 200
				maxAttempts = 5
				fmt.Println("Угадай число от 1 до 200 за 5 попыток")
			default:
				fmt.Println("Выбери уровень сложности Easy/Medium/Hard")
				continue
			}
			fmt.Println("Игра началась!")
			result, attempts := findNumber(scanner, maxNumber, maxAttempts)
			err := saveResult(result, attempts)
			if err != nil {
				fmt.Println("Ошибка сохранения", err)
			}
		} else {
			if err := scanner.Err(); err != nil {
				fmt.Println("Ошибка чтения:", err)
			}
			break
		}
		for {
			fmt.Println("Хочешь сыграть еще раз?(да/нет)")
			if scanner.Scan() {
				answer := scanner.Text()
				if answer == "да" {
					break
				} else if answer == "нет" {
					return
				} else {
					fmt.Println("Введите только `да` либо `нет`")
					continue
				}
			} else {
				if err := scanner.Err(); err != nil {
					fmt.Println("Ошибка чтения:", err)
				}
				return
			}
		}

	}
}
func getNumber(maxNumber int) int {
	n := rand.IntN(maxNumber) + 1
	return n
}

func findNumber(scanner *bufio.Scanner, maxNumber int, maxAttempts int) (string, int) {

	secretNumber := getNumber(maxNumber)
	previosNumbers := []int{}
	attempts := 0
	for attempts < maxAttempts {
		if scanner.Scan() {
			input := scanner.Text()
			number, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Введите числовое значение")
				continue
			}
			previosNumbers = append(previosNumbers, number)
			fmt.Println("предыдущие попытки: ", previosNumbers)
			fmt.Println("Ты ввёл:", input)
			attempts++

			if number == secretNumber {
				fmt.Println("ты угадал!")
				return "победа", attempts
			}
			difference := math.Abs(float64(number - secretNumber))
			if difference <= 5 {
				fmt.Println("горячо")
			} else if difference <= 15 {
				fmt.Println("тепло")
			} else {
				fmt.Println("холодно")
			}

			if attempts == maxAttempts {
				fmt.Println("ты проиграл , попробуй еще раз")
				fmt.Println("правильный ответ ", secretNumber)

			}
		} else {
			if err := scanner.Err(); err != nil {
				fmt.Println("Ошибка чтения:", err)
			}
		}
	}
	return "проигрыш", attempts
}
func saveResult(result string, attempts int) error {
	currenttime := time.Now().Format("2006-01-02 15:04:05")
	save := GameResult{
		Date:     currenttime,
		Result:   result,
		Attempts: attempts,
	}
	results := []GameResult{}
	data, err := os.ReadFile("results.json")
	if os.IsNotExist(err) {
	} else if err != nil {
		return err
	} else {
		err = json.Unmarshal(data, &results)
		if err != nil {
			return err
		}
	}
	results = append(results, save)
	data, err = json.MarshalIndent(results, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile("results.json", data, 0666)
	if err != nil {
		return err
	}
	return nil
}
