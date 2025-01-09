package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// StartMenuInfo содержит всю информацию о пользователе
type StartMenuInfo struct {
	Name     string    `json:"name"`
	Age      int       `json:"age"`
	Money    float64   `json:"money"`
	Job      string    `json:"job"`
	HaveAJob bool      `json:"haveajob"`
	LastPaid time.Time `json:"lastPaid"`
}

// Job определяет структуру работы
type Job struct {
	Name   string
	Salary float64
}

const (
	dataFile     = "userdata.json"
	salaryPeriod = 5 * time.Minute
)

// Список всех доступных работ
var jobs = []Job{
	{"Lawyer", 4500},
	{"Pilot", 6000},
	{"Entrepreneur", 7000},
	{"Programmer", 3000},
	{"Writer", 2000},
	{"Scientist", 4800},
	{"Mechanic", 2600},
	{"Nurse", 3000},
	{"Driver", 2500},
	{"Musician", 1800},
	{"Chef", 2300},
	{"Engineer", 4000},
	{"Doctor", 5000},
	{"Artist", 2200},
	{"Police Officer", 3500},
	{"Firefighter", 3400},
	{"Farmer", 2000},
	{"Designer", 2500},
	{"Teacher", 1500},
	{"Salesperson", 2700},
}

// generateStartingMoney генерирует случайную стартовую сумму денег
func generateStartingMoney() float64 {
	rand.Seed(time.Now().UnixNano())
	money := float64(rand.Intn(50001)) // Случайная сумма от 0 до 50000
	fmt.Printf("The game randomly generated $%.2f as your starting money!\n", money)

	if money < 1000 {
		money += 1000
		fmt.Println("Since it's less than $1000, here's a $1000 bonus! Have a good game!")
	}

	return money
}

// readInput читает строку ввода от пользователя
func readInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return ""
	}
	return strings.TrimSpace(input)
}

// readNumber читает и проверяет числовой ввод
func readNumber(prompt string, min, max float64) (float64, error) {
	input := readInput(prompt)
	num, err := strconv.ParseFloat(input, 64)
	if err != nil || num < min || num > max {
		return 0, fmt.Errorf("please enter a number between %.2f and %.2f", min, max)
	}
	return num, nil
}

// startMenu инициализирует нового игрока
func startMenu() (StartMenuInfo, error) {
	info := StartMenuInfo{
		LastPaid: time.Now(),
	}

	info.Name = readInput("Enter your name: ")
	if info.Name == "" {
		return info, fmt.Errorf("name cannot be empty")
	}

	age, err := readNumber("\nEnter your age: ", 18, 100)
	if err != nil {
		return info, err
	}
	info.Age = int(age)

	// Генерируем случайную стартовую сумму денег
	info.Money = generateStartingMoney()

	return info, nil
}

// saveData сохраняет данные игрока
func saveData(info *StartMenuInfo) error {
	file, err := os.Create(dataFile)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(info); err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}
	return nil
}

// loadData загружает сохраненные данные игрока
func loadData() (StartMenuInfo, error) {
	file, err := os.Open(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return StartMenuInfo{}, fmt.Errorf("no saved game found")
		}
		return StartMenuInfo{}, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var info StartMenuInfo
	if err := json.NewDecoder(file).Decode(&info); err != nil {
		return StartMenuInfo{}, fmt.Errorf("error decoding JSON: %w", err)
	}
	return info, nil
}

// updateSalary обновляет зарплату игрока
func updateSalary(info *StartMenuInfo) {
	if !info.HaveAJob {
		return
	}

	now := time.Now()
	elapsed := now.Sub(info.LastPaid)
	periods := int(elapsed / salaryPeriod)

	if periods > 0 {
		job := jobs[getJobID(info.Job)]
		salary := float64(periods) * job.Salary
		info.Money += salary
		info.LastPaid = info.LastPaid.Add(time.Duration(periods) * salaryPeriod)
		fmt.Printf("You earned $%.2f from your job!\n", salary)
	}
}

// displayJobs показывает список доступных работ
func displayJobs() {
	fmt.Println("\nAvailable jobs:")
	for i, job := range jobs {
		fmt.Printf("%d. %s - $%.2f per %d minutes\n", i+1, job.Name, job.Salary, int(salaryPeriod.Minutes()))
	}
	fmt.Println("0. Back to menu")
}

// findJob позволяет игроку найти работу
func findJob(info *StartMenuInfo) error {
	displayJobs()
	choice, err := readNumber("Choose a job by typing a number: ", 0, float64(len(jobs)))
	if err != nil {
		return err
	}

	if choice == 0 {
		return nil
	}

	selectedJob := jobs[int(choice)-1]
	info.Job = selectedJob.Name
	info.HaveAJob = true
	info.LastPaid = time.Now()

	if err := saveData(info); err != nil {
		return fmt.Errorf("error saving job selection: %w", err)
	}

	fmt.Printf("Congratulations! You're now a %s earning $%.2f every %d minutes!\n",
		selectedJob.Name, selectedJob.Salary, int(salaryPeriod.Minutes()))
	return nil
}

// displayMenu показывает главное меню игры
func displayMenu(info *StartMenuInfo) {
	fmt.Println("\n=== Career Game Menu ===")
	if info.HaveAJob {
		fmt.Println("1. View current job")
		fmt.Println("2. Leave job")
	} else {
		fmt.Println("1. Find a job")
	}
	fmt.Println("3. View account balance")
	fmt.Println("4. Save and quit")
}

// getJobID получает ID работы по её названию
func getJobID(jobName string) int {
	for id, job := range jobs {
		if job.Name == jobName {
			return id
		}
	}
	return -1
}

// displayWelcomeBack показывает приветственное сообщение для вернувшегося игрока
func displayWelcomeBack(info *StartMenuInfo) {
	fmt.Printf("\nWelcome back, %s!\n", info.Name)

	if info.HaveAJob {
		job := jobs[getJobID(info.Job)]
		fmt.Printf("You are currently working as a %s with a salary of $%.2f per %d minutes.\n",
			info.Job, job.Salary, int(salaryPeriod.Minutes()))
		fmt.Printf("Your current balance is: $%.2f\n", info.Money)
	} else {
		fmt.Printf("You currently don't have a job. Your balance is: $%.2f\n", info.Money)
		fmt.Println("Don't forget to check the job market!")
	}
}

func main() {
	var info StartMenuInfo
	var err error

	// Пытаемся загрузить сохраненную игру
	info, err = loadData()
	if err != nil {
		// Если сохранение не найдено - создаем нового игрока
		fmt.Println("Welcome new player! Let's set up your character.")
		info, err = startMenu()
		if err != nil {
			fmt.Printf("Error starting game: %v\n", err)
			return
		}
		if err := saveData(&info); err != nil {
			fmt.Printf("Error saving initial data: %v\n", err)
			return
		}
	} else {
		// Приветствуем вернувшегося игрока с подробной информацией
		displayWelcomeBack(&info)
	}

	// Основной игровой цикл
	for {
		updateSalary(&info)
		displayMenu(&info)

		choice, err := readNumber("Enter your choice: ", 1, 4)
		if err != nil {
			fmt.Println("Invalid input, please try again")
			continue
		}

		switch int(choice) {
		case 1:
			if info.HaveAJob {
				job := jobs[getJobID(info.Job)]
				fmt.Printf("\nCurrent job: %s\nSalary: $%.2f per %d minutes\n",
					info.Job, job.Salary, int(salaryPeriod.Minutes()))
			} else {
				if err := findJob(&info); err != nil {
					fmt.Printf("Error finding job: %v\n", err)
				}
			}
		case 2:
			if info.HaveAJob {
				fmt.Printf("You left your job as %s\n", info.Job)
				info.Job = ""
				info.HaveAJob = false
				if err := saveData(&info); err != nil {
					fmt.Printf("Error saving after leaving job: %v\n", err)
				}
			}
		case 3:
			fmt.Printf("\nCurrent balance: $%.2f\n", info.Money)
		case 4:
			if err := saveData(&info); err != nil {
				fmt.Printf("Error saving game: %v\n", err)
			}
			fmt.Println("Thanks for playing! Goodbye!")
			return
		}
	}
}
