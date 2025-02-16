package dreamcareer

import (
	"time"

	"github.com/therecipe/qt/widgets"
)

type GameWindow struct {
	window            *widgets.QMainWindow
	mainStack         *widgets.QStackedWidget
	info              *StartMenuInfo
	nameInput         *widgets.QLineEdit
	ageInput          *widgets.QSpinBox
	moneyLabel        *widgets.QLabel
	jobLabel          *widgets.QLabel
	statusLabel       *widgets.QLabel
	jobList           *widgets.QListWidget
	notificationLabel *widgets.QLabel
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
	{"Lawyer", 4500},
	{"Pilot", 6000},
	{"Entrepreneur", 7000},
	{"Programmer", 3000},
	{"Doctor", 5500},
	{"Teacher", 2500},
	{"Engineer", 3500},
	{"Chef", 2000},
	{"Dentist", 5800},
	{"Pharmacist", 4200},
	{"Veterinarian", 4000},
	{"Nurse", 2800},
	{"Physical Therapist", 3200},
	{"Data Scientist", 4800},
	{"Cybersecurity Analyst", 4200},
	{"UI/UX Designer", 3300},
	{"Systems Administrator", 3100},
	{"Cloud Architect", 5200},
	{"Investment Banker", 8000},
	{"Financial Analyst", 3800},
	{"Management Consultant", 5500},
	{"Marketing Manager", 3900},
	{"Human Resources Manager", 3400},
	{"Graphic Designer", 2600},
	{"Video Game Developer", 3600},
	{"Film Director", 4800},
	{"Journalist", 2400},
	{"Content Creator", 2800},
	{"Research Scientist", 3700},
	{"Biomedical Engineer", 4100},
	{"Environmental Scientist", 3200},
}
