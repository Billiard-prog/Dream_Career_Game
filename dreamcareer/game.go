package dreamcareer

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func createNewGame() *GameWindow {
	rand.Seed(time.Now().UnixNano())

	game := &GameWindow{
		window:            widgets.NewQMainWindow(nil, 0),
		mainStack:         widgets.NewQStackedWidget(nil),
		info:              &StartMenuInfo{LastPaid: time.Now()},
		notificationLabel: widgets.NewQLabel(nil, 0),
	}

	game.window.SetStyleSheet(`
		QWidget {
			background-color: rgb(66, 66, 66);
			color: white;
		}
		QLabel {
			color: white;
		}
		QSpinBox, QLineEdit {
			background-color: rgb(100, 100, 100);
			color: white;
			border: 1px solid rgb(85, 85, 85);
			padding: 5px;
			border-radius: 5px;
		}
		QListWidget {
			background-color: rgb(100, 100, 100);
			color: white;
			border: 1px solid rgb(85, 85, 85);
		}
	`)

	game.window.SetWindowTitle("Dream Career Game")
	game.window.SetMinimumSize2(800, 600)
	game.setupUI()

	if err := game.loadData(); err != nil {
		game.mainStack.SetCurrentIndex(0)
	} else {
		game.updateDisplay()
		game.mainStack.SetCurrentIndex(1)
	}

	return game
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
				g.showNotification(fmt.Sprintf("You earned $%.2f from your job!", salary))
				break
			}
		}
	}
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
			g.showNotification(fmt.Sprintf("You got the job as a %s!", job.Name))
			return
		}
	}
}

func (g *GameWindow) showNotification(message string) {
	g.notificationLabel.SetText(message)
	g.notificationLabel.Show()

	timer := core.NewQTimer(nil)
	timer.ConnectTimeout(func() {
		g.notificationLabel.Hide()
		timer.DeleteLater()
	})
	timer.Start(3000)
}
