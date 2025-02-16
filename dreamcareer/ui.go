package dreamcareer

import (
	"fmt"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

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
	startButton.SetMinimumSize2(200, 50)
	startButton.SetStyleSheet(getButtonStyle())
	startButton.ConnectClicked(func(bool) { g.startNewGame() })

	newGameLayout.AddLayout(nameLayout, 0)
	newGameLayout.AddLayout(ageLayout, 0)
	newGameLayout.AddStretch(1)
	newGameLayout.AddWidget(startButton, 0, core.Qt__AlignCenter)
	newGameWidget.SetLayout(newGameLayout)

	mainGameWidget := widgets.NewQWidget(nil, 0)
	mainGameLayout := widgets.NewQVBoxLayout()

	g.setupMainGameWidgets(mainGameWidget, mainGameLayout)

	g.mainStack.AddWidget(newGameWidget)
	g.mainStack.AddWidget(mainGameWidget)
	g.window.SetCentralWidget(g.mainStack)

	timer := core.NewQTimer(nil)
	timer.ConnectTimeout(func() { g.updateSalary() })
	timer.Start(30000)
}

func (g *GameWindow) setupMainGameWidgets(mainGameWidget *widgets.QWidget, mainGameLayout *widgets.QVBoxLayout) {
	g.statusLabel = widgets.NewQLabel(nil, 0)
	g.moneyLabel = widgets.NewQLabel(nil, 0)
	g.jobLabel = widgets.NewQLabel(nil, 0)

	g.setupJobList()
	g.setupButtons(mainGameLayout)
	g.setupNotificationLabel()

	mainGameLayout.AddWidget(g.statusLabel, 0, core.Qt__AlignTop)
	mainGameLayout.AddWidget(g.moneyLabel, 0, core.Qt__AlignTop)
	mainGameLayout.AddWidget(g.jobLabel, 0, core.Qt__AlignTop)
	mainGameLayout.AddWidget(g.notificationLabel, 0, core.Qt__AlignTop)
	mainGameLayout.AddWidget(g.jobList, 1, 0)
	mainGameWidget.SetLayout(mainGameLayout)
}

func (g *GameWindow) setupJobList() {
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
}

func (g *GameWindow) setupButtons(mainGameLayout *widgets.QVBoxLayout) {
	buttonLayout := widgets.NewQHBoxLayout()
	findJobButton := widgets.NewQPushButton2("Find Job", nil)
	findJobButton.SetStyleSheet(getButtonStyle())
	quitJobButton := widgets.NewQPushButton2("Quit Job", nil)
	quitJobButton.SetStyleSheet(getButtonStyle())
	saveButton := widgets.NewQPushButton2("Save & Quit", nil)
	saveButton.SetStyleSheet(getButtonStyle())

	findJobButton.ConnectClicked(func(bool) { g.toggleJobList() })
	quitJobButton.ConnectClicked(func(bool) { g.quitJob() })
	saveButton.ConnectClicked(func(bool) { g.saveAndQuit() })

	buttonLayout.AddWidget(findJobButton, 0, 0)
	buttonLayout.AddWidget(quitJobButton, 0, 0)
	buttonLayout.AddStretch(1)
	buttonLayout.AddWidget(saveButton, 0, 0)
	mainGameLayout.AddLayout(buttonLayout, 0)
}

func (g *GameWindow) setupNotificationLabel() {
	g.notificationLabel = widgets.NewQLabel(nil, 0)
	g.notificationLabel.SetStyleSheet(`
		QLabel {
			color: white;
			background-color: rgba(0, 0, 0, 0.7);
			padding: 10px;
			border-radius: 5px;
			font-size: 14px;
		}
	`)
	g.notificationLabel.Hide()
}

func getButtonStyle() string {
	return `
		QPushButton {
			background-color: rgb(85, 85, 85);
			color: white;
			font-size: 15px;
			font-weight: bold;
			border-radius: 10px;
			padding: 10px;
		}
		QPushButton:hover {
			background-color: rgb(123, 123, 123);
		}
		QPushButton:pressed {
			background-color: rgb(150, 150, 150);
		}
	`
}
