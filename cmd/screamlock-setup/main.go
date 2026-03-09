// ScreamLock Setup is a "next, next, next" installer: asks about autostart,
// adds Task Scheduler if requested, then opens the microphone/config dialog
// with a live level meter for testing (no lock).
//
//go:build windows
// +build windows

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-ole/go-ole"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/screamlock/screamlock/config"
	"github.com/screamlock/screamlock/internal/audio"
)

const taskName = "ScreamLock"

type deviceItem struct {
	ID   string
	Name string
}

func main() {
	if err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED); err != nil {
		log.Fatal(err)
	}
	defer ole.CoUninitialize()

	runWizard()
}

func runWizard() {
	var mw *walk.MainWindow
	var page1, page2, page3 *walk.Composite
	var nextBtn, backBtn, finishBtn *walk.PushButton
	var autostartCheck *walk.CheckBox
	currentPage := 1

	updateButtons := func() {
		if backBtn != nil {
			backBtn.SetEnabled(currentPage > 1)
		}
		if nextBtn != nil {
			nextBtn.SetVisible(currentPage < 3)
		}
		if finishBtn != nil {
			finishBtn.SetVisible(currentPage == 3)
		}
	}

	_, err := declarative.MainWindow{
		AssignTo: &mw,
		Title:    "ScreamLock Setup",
		MinSize:  declarative.Size{Width: 460, Height: 300},
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.Composite{
				AssignTo: &page1,
				Layout:   declarative.VBox{},
				Children: []declarative.Widget{
					declarative.Label{Text: "Welcome to ScreamLock Setup", Font: declarative.Font{Bold: true, PointSize: 12}},
					declarative.Label{Text: "This wizard will help you:\n\n• Optionally run ScreamLock when you log on to Windows\n• Choose your microphone and set the sensitivity\n• Test the input level without locking the computer\n\nClick Next to continue.", MaxSize: declarative.Size{Width: 400}},
				},
			},
			declarative.Composite{
				AssignTo: &page2,
				Visible:  false,
				Layout:   declarative.VBox{},
				Children: []declarative.Widget{
					declarative.Label{Text: "Run at Windows startup?", Font: declarative.Font{Bold: true, PointSize: 12}},
					declarative.Label{Text: "Should ScreamLock start automatically when you log on?\n\nYes = It will run in the background after every restart.\nNo = You can start it manually when needed.", MaxSize: declarative.Size{Width: 400}},
					declarative.CheckBox{
						AssignTo: &autostartCheck,
						Text:     "Yes, run ScreamLock when I log on to Windows",
						Checked:  true,
					},
				},
			},
			declarative.Composite{
				AssignTo: &page3,
				Visible:  false,
				Layout:   declarative.VBox{},
				Children: []declarative.Widget{
					declarative.Label{Text: "Ready to finish", Font: declarative.Font{Bold: true, PointSize: 12}},
					declarative.Label{Text: "Click Finish to open microphone settings.\n\nYou can choose the microphone, set the lock threshold, and use the live level meter to test your voice without locking the computer.", MaxSize: declarative.Size{Width: 400}},
				},
			},
			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.HSpacer{},
					declarative.PushButton{
						AssignTo: &backBtn,
						Text:     "Back",
						Enabled:  false,
						OnClicked: func() {
							if currentPage == 2 {
								page2.SetVisible(false)
								page1.SetVisible(true)
								currentPage = 1
							} else if currentPage == 3 {
								page3.SetVisible(false)
								page2.SetVisible(true)
								currentPage = 2
							}
							updateButtons()
						},
					},
					declarative.PushButton{
						AssignTo: &nextBtn,
						Text:     "Next",
						OnClicked: func() {
							if currentPage == 1 {
								page1.SetVisible(false)
								page2.SetVisible(true)
								currentPage = 2
							} else if currentPage == 2 {
								page2.SetVisible(false)
								page3.SetVisible(true)
								currentPage = 3
							}
							updateButtons()
						},
					},
					declarative.PushButton{
						AssignTo:  &finishBtn,
						Text:     "Finish",
						Visible:  false,
						OnClicked: func() {
							if autostartCheck != nil && autostartCheck.Checked() {
								exePath, _ := os.Executable()
								screamlockExe := filepath.Join(filepath.Dir(exePath), "screamlock.exe")
								if _, err := os.Stat(screamlockExe); err == nil {
									_ = exec.Command("schtasks", "/Create", "/TN", taskName, "/TR", screamlockExe, "/SC", "ONLOGON", "/F").Run()
								}
							}
							mw.Close()
							runConfigWithMeter()
						},
					},
				},
			},
		},
	}.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// runConfigWithMeter shows the config dialog with a live input level meter (test only, no lock).
func runConfigWithMeter() {
	cfg, _, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	devices, err := audio.ListCaptureDevices()
	if err != nil {
		log.Fatal(err)
	}

	model := make([]deviceItem, 0, len(devices)+1)
	model = append(model, deviceItem{ID: "", Name: "(Default microphone)"})
	for _, d := range devices {
		display := d.Name
		if display == "" {
			display = d.ID
		}
		model = append(model, deviceItem{ID: d.ID, Name: display})
	}

	var mw *walk.MainWindow
	var cb *walk.ComboBox
	var thresholdEdit *walk.NumberEdit
	var intervalEdit *walk.NumberEdit
	var levelBar *walk.ProgressBar
	var levelLabel *walk.Label

	selectedID := cfg.DeviceID
	threshold := cfg.ThresholdDB
	interval := cfg.CheckIntervalSeconds
	if interval < 1 {
		interval = 1
	}

	var reader *audio.PeakReader
	var meterDone chan struct{}
	var meterStartedOnce sync.Once

	startMeter := func(deviceID string) {
		if meterDone != nil {
			close(meterDone)
			meterDone = nil
		}
		if reader != nil {
			reader.Close()
			reader = nil
		}
		r, err := audio.OpenPeakReader(deviceID)
		if err != nil {
			return
		}
		reader = r
		meterDone = make(chan struct{})
		done := meterDone
		go func() {
			tick := time.NewTicker(50 * time.Millisecond)
			defer tick.Stop()
			for {
				select {
				case <-done:
					return
				case <-tick.C:
					peak, err := reader.Peak()
					if err != nil {
						continue
					}
					pct := int(peak * 100)
					if pct > 100 {
						pct = 100
					}
					db := audio.DBFromLinear(peak)
					if mw != nil {
						mw.Synchronize(func() {
							if levelBar != nil {
								levelBar.SetValue(pct)
							}
							if levelLabel != nil {
								levelLabel.SetText(formatLevel(peak, db))
							}
						})
					}
				}
			}
		}()
	}

	onDeviceChange := func() {
		if cb == nil {
			return
		}
		idx := cb.CurrentIndex()
		if idx >= 0 && idx < len(model) {
			startMeter(model[idx].ID)
		}
	}

	_, err = declarative.MainWindow{
		AssignTo: &mw,
		Title:    "ScreamLock — Microphone & Test",
		MinSize:  declarative.Size{Width: 440, Height: 380},
		// Start the level meter once shortly after the window is shown
		OnBoundsChanged: func() {
			meterStartedOnce.Do(func() {
				go func() {
					time.Sleep(300 * time.Millisecond)
					if mw != nil {
						mw.Synchronize(onDeviceChange)
					}
				}()
			})
		},
		Layout: declarative.VBox{MarginsZero: true},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.Grid{Columns: 2},
				Children: []declarative.Widget{
					declarative.Label{Text: "Microphone:"},
					declarative.ComboBox{
						AssignTo:      &cb,
						Model:         model,
						BindingMember: "ID",
						DisplayMember: "Name",
						Value:         selectedID,
						MinSize:       declarative.Size{Width: 300},
						OnCurrentIndexChanged: func() {
							onDeviceChange()
						},
					},
					declarative.Label{Text: "Sensitivity (dB):"},
					declarative.NumberEdit{
						AssignTo: &thresholdEdit,
						Value:    float64(threshold),
						MinValue: -80,
						MaxValue: -20,
						Decimals: 0,
						Suffix:   " dB (lock when above this)",
					},
					declarative.Label{Text: "Check every (seconds):"},
					declarative.NumberEdit{
						AssignTo: &intervalEdit,
						Value:    float64(interval),
						MinValue: 1,
						MaxValue: 60,
						Decimals: 0,
					},
				},
			},
			declarative.Composite{
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					declarative.Label{Text: "Test level (no lock): speak into the mic — bar shows input level. Computer will not lock here.", TextColor: walk.RGB(80, 80, 80), MaxSize: declarative.Size{Width: 380}},
					declarative.Composite{
						Layout: declarative.HBox{},
						Children: []declarative.Widget{
							declarative.ProgressBar{
								AssignTo: &levelBar,
								MinValue: 0,
								MaxValue: 100,
								Value:    0,
								MinSize:  declarative.Size{Width: 280, Height: 24},
							},
							declarative.Label{AssignTo: &levelLabel, Text: "— dB", MinSize: declarative.Size{Width: 56}},
						},
					},
				},
			},
			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.HSpacer{},
					declarative.PushButton{
						Text: "Save",
						OnClicked: func() {
							idx := cb.CurrentIndex()
							if idx >= 0 && idx < len(model) {
								cfg.DeviceID = model[idx].ID
							}
							cfg.ThresholdDB = thresholdEdit.Value()
							v := int(intervalEdit.Value())
							if v < 1 {
								v = 1
							}
							cfg.CheckIntervalSeconds = v
							if err := config.Save(cfg); err != nil {
								walk.MsgBox(mw, "Error", "Could not save config: "+err.Error(), walk.MsgBoxIconError)
								return
							}
							walk.MsgBox(mw, "Saved", "Configuration saved. ScreamLock will use it on next run.", walk.MsgBoxIconInformation)
							mw.Close()
						},
					},
					declarative.PushButton{
						Text: "Done",
						OnClicked: func() {
							mw.Close()
						},
					},
				},
			},
		},
	}.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func formatLevel(linear float32, dB float64) string {
	if linear <= 0 {
		return "— dB"
	}
	return fmt.Sprintf("%.0f dB", dB)
}
