package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type GameWindow struct {
	window      *widgets.QMainWindow
	mainStack   *widgets.QStackedWidget
	info        *StartMenuInfo
	nameInput   *widgets.QLineEdit
	ageInput    *widgets.QSpinBox
	moneyLabel  *widgets.QLabel
	jobLabel    *widgets.QLabel
	statusLabel *widgets.QLabel
	jobList     *widgets.QListWidget
}

type StartMenuInfo struct {
	Name     string    `json:"name"`
	Age      int       `json:"age"`
	Money    float64   `json:"money"`
	Job      string    `json:"job"`
	HaveAJob bool      `json:"haveAJob"`
	LastPaid time.Time `json:"lastPaid"`
}

type Job struct {
	Name   string
	Salary float64
}

const (
	dataFile     = "userdata.json"
	salaryPeriod = 5 * time.Minute
)

var jobs = []Job{
	// Existing jobs
	{"Lawyer", 4500},
	{"Pilot", 6000},
	{"Entrepreneur", 7000},
	{"Programmer", 3000},
	{"Doctor", 5500},
	{"Teacher", 2500},
	{"Engineer", 3500},
	{"Chef", 2000},

	// Healthcare sector
	{"Dentist", 5800},
	{"Pharmacist", 4200},
	{"Veterinarian", 4000},
	{"Nurse", 2800},
	{"Physical Therapist", 3200},

	// Technology sector
	{"Data Scientist", 4800},
	{"Cybersecurity Analyst", 4200},
	{"UI/UX Designer", 3300},
	{"Systems Administrator", 3100},
	{"Cloud Architect", 5200},

	// Business and Finance
	{"Investment Banker", 8000},
	{"Financial Analyst", 3800},
	{"Management Consultant", 5500},
	{"Marketing Manager", 3900},
	{"Human Resources Manager", 3400},

	// Creative and Media
	{"Graphic Designer", 2600},
	{"Video Game Developer", 3600},
	{"Film Director", 4800},
	{"Journalist", 2400},
	{"Content Creator", 2800},

	// Scientific
	{"Research Scientist", 3700},
	{"Biomedical Engineer", 4100},
	{"Environmental Scientist", 3200},
}

func createNewGame() *GameWindow {
	rand.Seed(time.Now().UnixNano())

	game := &GameWindow{
		window:    widgets.NewQMainWindow(nil, 0),
		mainStack: widgets.NewQStackedWidget(nil),
		info:      &StartMenuInfo{LastPaid: time.Now()},
	}

	game.window.SetWindowTitle("Dream Career Game")
	game.window.SetMinimumSize2(800, 600)
	game.window.SetStyleSheet(`
		QPushButton {
			background-color: rgb(35, 35, 35);
		}
	`)
	game.setupUI()

	if err := game.loadData(); err != nil {
		game.mainStack.SetCurrentIndex(0)
	} else {
		game.updateDisplay()
		game.mainStack.SetCurrentIndex(1)
	}

	return game
}

func (g *GameWindow) setupUI() {
	newGameWidget := widgets.NewQWidget(nil, 0)
	newGameLayout := widgets.NewQVBoxLayout()

	nameLayout := widgets.NewQHBoxLayout()
	nameLabel := widgets.NewQLabel2("Name:", nil, 0)
	g.nameInput = widgets.NewQLineEdit(nil)
	nameLayout.AddWidget(nameLabel, 0, core.Qt__AlignLeft)
	nameLayout.AddWidget(g.nameInput, 1, 0)

	ageLayout := widgets.NewQHBoxLayout()
	ageLabel := widgets.NewQLabel2("Age:", nil, 0)
	g.ageInput = widgets.NewQSpinBox(nil)
	g.ageInput.SetRange(18, 100)
	ageLayout.AddWidget(ageLabel, 0, core.Qt__AlignLeft)
	ageLayout.AddWidget(g.ageInput, 1, 0)

	startButton := widgets.NewQPushButton2("Start Game", nil)
	startButton.ConnectClicked(func(bool) { g.startNewGame() })

	newGameLayout.AddLayout(nameLayout, 0)
	newGameLayout.AddLayout(ageLayout, 0)
	newGameLayout.AddStretch(1)
	newGameLayout.AddWidget(startButton, 0, core.Qt__AlignCenter)
	newGameWidget.SetLayout(newGameLayout)

	mainGameWidget := widgets.NewQWidget(nil, 0)
	mainGameLayout := widgets.NewQVBoxLayout()

	g.statusLabel = widgets.NewQLabel(nil, 0)
	g.moneyLabel = widgets.NewQLabel(nil, 0)
	g.jobLabel = widgets.NewQLabel(nil, 0)

	g.jobList = widgets.NewQListWidget(nil)
	g.jobList.ConnectItemDoubleClicked(func(item *widgets.QListWidgetItem) {
		g.applyForJob(item.Text())
	})

	for _, job := range jobs {
		jobText := fmt.Sprintf("%s - $%.2f per %d minutes", job.Name, job.Salary, int(salaryPeriod.Minutes()))
		item := widgets.NewQListWidgetItem(nil, 0)
		item.SetText(jobText)
		g.jobList.AddItem2(item)
	}
	g.jobList.Hide()

	buttonLayout := widgets.NewQHBoxLayout()
	findJobButton := widgets.NewQPushButton2("Find Job", nil)
	quitJobButton := widgets.NewQPushButton2("Quit Job", nil)
	saveButton := widgets.NewQPushButton2("Save & Quit", nil)

	findJobButton.ConnectClicked(func(bool) { g.toggleJobList() })
	quitJobButton.ConnectClicked(func(bool) { g.quitJob() })
	saveButton.ConnectClicked(func(bool) { g.saveAndQuit() })

	buttonLayout.AddWidget(findJobButton, 0, 0)
	buttonLayout.AddWidget(quitJobButton, 0, 0)
	buttonLayout.AddStretch(1)
	buttonLayout.AddWidget(saveButton, 0, 0)

	mainGameLayout.AddWidget(g.statusLabel, 0, core.Qt__AlignTop)
	mainGameLayout.AddWidget(g.moneyLabel, 0, core.Qt__AlignTop)
	mainGameLayout.AddWidget(g.jobLabel, 0, core.Qt__AlignTop)
	mainGameLayout.AddWidget(g.jobList, 1, 0)
	mainGameLayout.AddLayout(buttonLayout, 0)
	mainGameWidget.SetLayout(mainGameLayout)

	g.mainStack.AddWidget(newGameWidget)
	g.mainStack.AddWidget(mainGameWidget)
	g.window.SetCentralWidget(g.mainStack)

	timer := core.NewQTimer(nil)
	timer.ConnectTimeout(func() { g.updateSalary() })
	timer.Start(30000)
}

func (g *GameWindow) applyForJob(jobText string) {
	if g.info.HaveAJob {
		widgets.QMessageBox_Warning(nil, "Error", "You already have a job! Quit your current job first.", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}

	for _, job := range jobs {
		if jobText == fmt.Sprintf("%s - $%.2f per %d minutes", job.Name, job.Salary, int(salaryPeriod.Minutes())) {
			g.info.Job = job.Name
			g.info.HaveAJob = true
			g.info.LastPaid = time.Now()
			g.saveData()
			g.updateDisplay()
			g.jobList.Hide()
			widgets.QMessageBox_Information(nil, "Congratulations", fmt.Sprintf("You got the job as a %s!", job.Name), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
			return
		}
	}
}

func (g *GameWindow) startNewGame() {
	if g.nameInput.Text() == "" {
		widgets.QMessageBox_Warning(nil, "Error", "Please enter your name", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}

	g.info.Name = g.nameInput.Text()
	g.info.Age = g.ageInput.Value()
	g.info.Money = float64(rand.Intn(49001) + 1000)
	g.info.LastPaid = time.Now()
	g.info.HaveAJob = false
	g.info.Job = ""

	if err := g.saveData(); err != nil {
		widgets.QMessageBox_Warning(nil, "Error", "Failed to save game data", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}

	g.updateDisplay()
	g.mainStack.SetCurrentIndex(1)
}

func (g *GameWindow) updateDisplay() {
	g.statusLabel.SetText(fmt.Sprintf("Welcome, %s! (Age: %d)", g.info.Name, g.info.Age))
	g.moneyLabel.SetText(fmt.Sprintf("Current Balance: $%.2f", g.info.Money))

	if g.info.HaveAJob {
		g.jobLabel.SetText(fmt.Sprintf("Current Job: %s", g.info.Job))
	} else {
		g.jobLabel.SetText("Currently Unemployed")
	}
}

func (g *GameWindow) toggleJobList() {
	if g.jobList.IsVisible() {
		g.jobList.Hide()
	} else {
		g.jobList.Show()
	}
}

func (g *GameWindow) quitJob() {
	if !g.info.HaveAJob {
		widgets.QMessageBox_Information(nil, "Info", "You don't have a job to quit!", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}

	g.info.Job = ""
	g.info.HaveAJob = false
	g.saveData()
	g.updateDisplay()
	widgets.QMessageBox_Information(nil, "Job Status", "You have quit your job.", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
}

func (g *GameWindow) updateSalary() {
	if !g.info.HaveAJob {
		return
	}

	now := time.Now()
	elapsed := now.Sub(g.info.LastPaid)
	periods := int(elapsed / salaryPeriod)

	if periods > 0 {
		for _, job := range jobs {
			if job.Name == g.info.Job {
				salary := float64(periods) * job.Salary
				g.info.Money += salary
				g.info.LastPaid = g.info.LastPaid.Add(time.Duration(periods) * salaryPeriod)
				g.saveData()
				g.updateDisplay()
				widgets.QMessageBox_Information(nil, "Salary Received",
					fmt.Sprintf("You earned $%.2f from your job!", salary),
					widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
				break
			}
		}
	}
}

func (g *GameWindow) saveData() error {
	file, err := os.Create(dataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(g.info)
}

func (g *GameWindow) loadData() error {
	file, err := os.Open(dataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(g.info)
}

func (g *GameWindow) saveAndQuit() {
	if err := g.saveData(); err != nil {
		widgets.QMessageBox_Warning(nil, "Error", "Failed to save game data", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}
	g.window.Close()
}

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	game := createNewGame()
	game.window.Show()
	app.Exec()
}
