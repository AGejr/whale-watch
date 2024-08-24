package ui

import (
	"fmt"
	"regexp"
)

func (app *App) showContainers() {
	app.updateContainers()
	app.SetFocus(app.Table)
}

func (app *App) updateContainers() {
	containers, err := app.DockerClient.ListContainers()
	if err != nil {
		app.showError("Error listing containers: " + err.Error())
		return
	}
	app.updateTable("Containers", []string{"CONTAINER ID", "IMAGE", "COMMAND", "CREATED", "STATUS", "PORTS", "NAMES"}, containers)
}

func (app *App) showContainerLogs() {
	row, _ := app.Table.GetSelection()
	if row == 0 {
		return // Header row selected
	}
	containerID := app.Table.GetCell(row, 0).Text
	containerName := app.Table.GetCell(row, 6).Text
	logs, err := app.DockerClient.GetContainerLogs(containerID)
	tagRegex := regexp.MustCompile("\\[[a-zA-Z0-9_,;: \\-\\.]+\\]")
	logs = tagRegex.ReplaceAllString(logs, "")
	if err != nil {
		app.showError("Error getting container logs: " + err.Error())
		return
	}
	app.TextLayout.SetTitle(fmt.Sprintf("Container logs (%s)", containerName))
	app.TextView.Clear()
	app.TextView.ScrollToBeginning()
	app.TextView.SetText(logs)
	app.Pages.SwitchToPage("text")
}

func (app *App) showContainerInfo(containerID string, containerName string) {
	inspect, err := app.DockerClient.InspectContainer(containerID)
	if err != nil {
		app.showError("Error inspecting container: " + err.Error())
		return
	}
	app.TextView.Clear()
	app.TextView.ScrollToBeginning()
	app.TextView.SetText(inspect)
	app.TextView.SetTitle(fmt.Sprintf("Container information (%s)", containerName))
	app.Pages.SwitchToPage("text")
}

func (app *App) handleContainerAction(action string) {
	if app.Table.GetTitle() != "Containers" {
		return
	}

	row, _ := app.Table.GetSelection()
	if row == 0 {
		return // Header row selected
	}

	containerID := app.Table.GetCell(row, 0).Text
	var err error

	switch action {
	case "start":
		err = app.DockerClient.StartContainer(containerID)
	case "stop":
		err = app.DockerClient.StopContainer(containerID)
	case "pause":
		err = app.DockerClient.PauseContainer(containerID)
	case "unpause":
		err = app.DockerClient.UnpauseContainer(containerID)
	case "kill":
		err = app.DockerClient.KillContainer(containerID)
	case "delete":
		app.showConfirmationDialog("Delete Container", "Are you sure you want to delete this container?", func() {
			err = app.DockerClient.DeleteContainer(containerID)
			if err != nil {
				app.showError(fmt.Sprintf("Error deleting container: %v", err))
			} else {
				app.updateContainers()
			}
		})
	default:
		app.showError("Unknown action: " + action)
		return
	}

	if err != nil {
		app.showError(fmt.Sprintf("Error performing %s action: %v", action, err))
	} else {
		app.updateContainers() // Refresh the container list
	}
}
