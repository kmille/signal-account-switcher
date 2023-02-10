package main

import (
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func get_data_dir(account_id int) (error, string) {
	config_dir, err := os.UserConfigDir()
	if err != nil {
		return err, ""
	}
	data_dir := filepath.Join(config_dir, fmt.Sprintf("Signal-Account-%d", account_id))
	if _, err := os.Stat(data_dir); os.IsNotExist(err) {
		if err := os.Mkdir(data_dir, os.ModePerm); err != nil {
			return err, ""
		}
	}
	return nil, data_dir
}

func run_signal(account_id int, debug *canvas.Text) {
	err, data_dir := get_data_dir(account_id)
	if err != nil {
		debug.Text = fmt.Sprintf("Error creating data_dir: %s", err)
		return
	}
	cmd := exec.Command("signal-desktop", "--user-data-dir="+data_dir)
	debug.Text = fmt.Sprintf("Staring Signal Account %d with data_dir %q", account_id, data_dir)
	err = cmd.Run()
	if err != nil {
		debug.Text = fmt.Sprint("Error executing Signal: %s\n", err)
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("Signal account switcher")
	w.Resize(fyne.NewSize(200, 400))

	debug := canvas.NewText("Hi, please choose account", color.White)
	account1 := widget.NewButton("Account #1", func() {
		run_signal(1, debug)
	})
	account2 := widget.NewButton("Account #2", func() {
		run_signal(2, debug)
	})
	account3 := widget.NewButton("Account #3", func() {
		run_signal(3, debug)
	})
	account4 := widget.NewButton("Account #4", func() {
		run_signal(4, debug)
	})

	grid := container.New(layout.NewGridLayout(1), debug, account1, account2, account3, account4)
	w.SetContent(grid)

	w.ShowAndRun()
	fmt.Println("Exiting ...")
}
