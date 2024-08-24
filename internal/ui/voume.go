package ui

import (
	"fmt"
	"sort"
)

func (app *App) showVolumes() {
	app.updateVolumes()
	app.SetFocus(app.Table)
}

func (app *App) updateVolumes() {
	volumes, err := app.DockerClient.ListVolumes()
	if err != nil {
		app.showError("Error listing volumes: " + err.Error())
		return
	}
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i][1] < volumes[j][1]
	})
	app.updateTable("Volumes", []string{"DRIVER", "VOLUME NAME", "MOUNT POINT"}, volumes)
}

func (app *App) handleVolumeAction(action string) {
	if app.Table.GetTitle() != "Volumes" {
		return
	}

	row, _ := app.Table.GetSelection()
	if row == 0 {
		return
	}

	volumeName := app.Table.GetCell(row, 1).Text
	var err error

	switch action {
	case "delete":
		app.showConfirmationDialog("Delete Volume", "Are you sure you want to delete this volume?", func() {
			err = app.DockerClient.DeleteVolume(volumeName)
			if err != nil {
				app.showError(fmt.Sprintf("Error deleting volume: %v", err))
			} else {
				app.updateVolumes() // Refresh the volume list
			}
		})
	default:
		app.showError("Unknown action: " + action)
		return
	}
}
