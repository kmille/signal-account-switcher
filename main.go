package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var signal_bin string

const version = "0.1.0"

func chooseSignalBinary(w fyne.Window) {
	callback := func(uri fyne.URIReadCloser, err error) {
		if uri == nil {
			dialog.ShowInformation("Error", "No file selected", w)
		} else {
			signal_bin = uri.URI().Path()
		}
	}
	dialogOpen := dialog.NewFileOpen(callback, w)
	dialogOpen.Show()
	// call it later, then it will be showed on the top
	dialog.ShowInformation("Error", "Signal application not found. Please specify", w)
}

func find_signal_executable(w fyne.Window) {
	var location string
	var err error

	switch runtime.GOOS {
	case "linux":
		location, err = exec.LookPath("signal-desktop")
		if err != nil {
			chooseSignalBinary(w)
			return
		}
	case "darwin":
		location = "/Applications/Signal.app/Contents/MacOS/Signal"
	case "windows":
		cache_dir, err := os.UserCacheDir()
		if err != nil {
			chooseSignalBinary(w)
			return
		}
		location = filepath.Join(cache_dir, "Programs", "signal-desktop", "Signal.exe")
	}

	if _, err := os.Stat(location); errors.Is(err, fs.ErrNotExist) {
		chooseSignalBinary(w)
		return
	}
	signal_bin = location
}

func get_data_dir(account_id int) (string, error) {
	config_dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	data_dir := filepath.Join(config_dir, fmt.Sprintf("Signal-Account-%d", account_id))
	if _, err := os.Stat(data_dir); errors.Is(err, fs.ErrNotExist) {
		if err := os.Mkdir(data_dir, os.ModePerm); err != nil {
			return "", err
		}
	}
	return data_dir, nil
}

func run_signal(account_id int, w fyne.Window) {
	data_dir, err := get_data_dir(account_id)
	if err != nil {
		dialog.ShowInformation("Error", fmt.Sprintf("Error creating data_dir: %s", err), w)
		return
	}
	dialog.ShowInformation("Info", fmt.Sprintf("Starting Signal Account #r%d with data_dir %q. Plase wait... ", account_id, data_dir), w)

	cmd := exec.Command(signal_bin, "--user-data-dir="+data_dir)
	err = cmd.Run()
	if err != nil {
		dialog.ShowInformation("Error", fmt.Sprintf("Error executing Signal: %s\n", err), w)
		return
	}
}

func main() {
	app := app.New()
	app.Settings().SetTheme(theme.DarkTheme())

	w := app.NewWindow(fmt.Sprintf("Signal account switcher v%s", version))
	w.Resize(fyne.NewSize(600, 800))
	find_signal_executable(w)

	account1 := widget.NewButton("Start signal account #1", func() {
		go run_signal(1, w)
	})
	account2 := widget.NewButton("Start signal account #2", func() {
		go run_signal(2, w)
	})
	account3 := widget.NewButton("Start signal account #3", func() {
		go run_signal(3, w)
	})
	account4 := widget.NewButton("Start signal account #4", func() {
		go run_signal(4, w)
	})

	grid := container.New(layout.NewGridLayout(1), account1, account2, account3, account4)
	w.SetContent(grid)
	w.Show()

	w.ShowAndRun()
	fmt.Println("Exiting ...")
}
