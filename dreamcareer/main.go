package dreamcareer

import (
	"os"

	"github.com/therecipe/qt/widgets"
)

func Run() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	game := createNewGame()
	game.window.Show()
	app.Exec()
}
