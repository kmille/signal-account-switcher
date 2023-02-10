package main

import (
	"errors"
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func get_signal_executable() (error, string) {
	switch runtime.GOOS {
	case "linux":
		return nil, "signal-desktop"
	case "darwin":
		return nil, "/Applications/Signal.app/Contents/MacOS/Signal"
	case "windows":
		dir, err := os.UserCacheDir()
		if err != nil {
			return err, ""
		}
		return nil, filepath.Join(dir, "Programs", "signal-desktop", "Signal.exe")
	default:
		fmt.Println("OS is not supported")
		os.Exit(1)
	}
	return errors.New("This is will never be reached"), ""

}

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

func change_text(debug *canvas.Text, text string) {
	debug.Text = text
	debug.Refresh()
}

func run_signal(account_id int, debug *canvas.Text) {
	err, signal_bin := get_signal_executable()
	if err != nil {
		change_text(debug, fmt.Sprintf("Error getting signal binary: %s", err))
		return
	}
	err, data_dir := get_data_dir(account_id)
	if err != nil {
		change_text(debug, fmt.Sprintf("Error creating data_dir: %s", err))
		return
	}
	cmd := exec.Command(signal_bin, "--user-data-dir="+data_dir)
	change_text(debug, fmt.Sprintf("Starting Signal Account %d with data_dir %q. Plase wait... ", account_id, data_dir))
	err = cmd.Run()
	if err != nil {
		change_text(debug, fmt.Sprint("Error executing Signal: %s\n", err))
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("Signal account switcher")
	w.Resize(fyne.NewSize(400, 500))
	w.RequestFocus()

	//debug := canvas.NewText("Hi, please choose account", color.Black)
	debug := canvas.NewText("Hi, please choose account", color.RGBA{R: 0x32, G: 0xCD, B: 0x32, A: 100})
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
