package dreamcareer

import (
	"encoding/json"
	"os"

	"github.com/therecipe/qt/widgets"
)

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
