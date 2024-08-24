package ui

import (
	"fmt"
	"sort"
)

func (app *App) showNetworks() {
	app.updateNetworks()
	app.SetFocus(app.Table)
}

func (app *App) updateNetworks() {
	networks, err := app.DockerClient.ListNetworks()
	if err != nil {
		app.showError("Error listing networks: " + err.Error())
		return
	}
	sort.Slice(networks, func(i, j int) bool {
		return networks[i][1] < networks[j][1]
	})
	app.updateTable("Networks", []string{"NETWORK ID", "NAME", "DRIVER", "SCOPE"}, networks)
}

func (app *App) handleNetworkAction(action string) {
	if app.Table.GetTitle() != "Networks" {
		return
	}

	row, _ := app.Table.GetSelection()
	if row == 0 {
		return
	}

	networkID := app.Table.GetCell(row, 0).Text
	var err error

	switch action {
	case "delete":
		app.showConfirmationDialog("Delete Network", "Are you sure you want to delete this network?", func() {
			err = app.DockerClient.DeleteNetwork(networkID)
			if err != nil {
				app.showError(fmt.Sprintf("Error deleting network: %v", err))
			} else {
				app.updateNetworks()
			}
		})
	default:
		app.showError("Unknown action: " + action)
		return
	}
}
