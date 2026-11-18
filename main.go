package main

import (
	"bufio"
	"encoding/json"
	"errors"
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

const (
	Green  = "\033[32m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Reset  = "\033[0m"
)

var errInputClosed = errors.New("ввод завершен")
var errInvalidNumber = errors.New("Введите только числовое значение")

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
		fmt.Println("Введите число")
		fmt.Println(Yellow+"попытка", attempts+1, Reset)
		for {
			number, err := enterNumber(scanner)
			if errors.Is(err, errInputClosed) {
				fmt.Println("Ввод закончен")
				return "проигрыш", attempts
			} else if errors.Is(err, errInvalidNumber) {
				fmt.Println("Введите только числовое значение")
				continue
			} else if err != nil {
				fmt.Println("Ошибка ввода")
				return "проигрыш", attempts
			}

			attempts++
			previosNumbers = append(previosNumbers, number)
			fmt.Println("предыдущие попытки: ", previosNumbers)
			fmt.Println("Ты ввёл:", number)

			if number == secretNumber {
				fmt.Println(Green + "ты угадал!" + Reset)
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
			if number < secretNumber {
				fmt.Println("загаданное число больше")
			} else {
				fmt.Println("загаданное число меньше")
			}

			if attempts == maxAttempts {
				fmt.Println(Red + "ты проиграл , попробуй еще раз" + Reset)
				fmt.Println("правильный ответ ", secretNumber)

			}
			break
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
func enterNumber(scanner *bufio.Scanner) (int, error) {
	if scanner.Scan() {
		input := scanner.Text()
		number, err := strconv.Atoi(input)
		if err != nil {
			return 0, errInvalidNumber
		} else {
			return number, nil
		}
	} else {
		if err := scanner.Err(); err != nil {
			return 0, errors.New("errorEnter")
		} else {
			return 0, errInputClosed
		}
	}
}
