// ScreamLock Config is a small GUI to choose the microphone and sensitivity for ScreamLock.
// It writes to the same config.json used by the main ScreamLock app.
//
//go:build windows
// +build windows

package main

import (
	"log"

	"github.com/go-ole/go-ole"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/screamlock/screamlock/config"
	"github.com/screamlock/screamlock/internal/audio"
)

type deviceItem struct {
	ID   string
	Name string
}

func main() {
	if err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED); err != nil {
		log.Fatal(err)
	}
	defer ole.CoUninitialize()

	cfg, _, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	devices, err := audio.ListCaptureDevices()
	if err != nil {
		log.Fatal(err)
	}

	// Build model: "(Default)" first, then listed devices
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

	selectedID := cfg.DeviceID
	threshold := cfg.ThresholdDB
	interval := cfg.CheckIntervalSeconds
	if interval < 1 {
		interval = 1
	}

	_, err = declarative.MainWindow{
		AssignTo: &mw,
		Title:    "ScreamLock Config",
		MinSize:  declarative.Size{Width: 420, Height: 220},
		Layout:   declarative.VBox{MarginsZero: true},
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
						MinSize:       declarative.Size{Width: 320},
					},
					declarative.Label{Text: "Sensitivity (dB):"},
					declarative.NumberEdit{
						AssignTo: &thresholdEdit,
						Value:    float64(threshold),
						MinValue: -80,
						MaxValue: -20,
						Decimals: 0,
						Suffix:   " dB (more negative = less sensitive)",
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
						Text: "Cancel",
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
