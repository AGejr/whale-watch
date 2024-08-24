package ui

import "fmt"

func (app *App) showImages() {
	app.updateImages()
	app.SetFocus(app.Table)
}

func (app *App) updateImages() {
	images, err := app.DockerClient.ListImages()
	if err != nil {
		app.showError("Error listing images: " + err.Error())
		return
	}
	app.updateTable("Images", []string{"REPOSITORY", "TAG", "IMAGE ID", "CREATED", "SIZE"}, images)
}

func (app *App) showImageInfo(imageID string, imageName string, imageTag string) {

	imageLongName := fmt.Sprintf("%s:%s", imageName, imageTag)
	inspect, err := app.DockerClient.InspectImage(imageID)
	if err != nil {
		app.showError("Error inspecting image: " + err.Error())
		return
	}
	app.TextView.Clear()
	app.TextView.ScrollToBeginning()
	app.TextView.SetText(inspect)
	app.TextView.SetTitle(fmt.Sprintf("Image information (%s)", imageLongName))
	app.Pages.SwitchToPage("text")
}

func (app *App) handleImageAction(action string) {
	if app.Table.GetTitle() != "Images" {
		return
	}

	row, _ := app.Table.GetSelection()
	if row == 0 {
		return // Header row selected
	}

	imageID := app.Table.GetCell(row, 2).Text // Assuming the image ID is in the third column
	var err error

	switch action {
	case "delete":
		app.showConfirmationDialog("Delete Image", "Are you sure you want to delete this image?", func() {
			err = app.DockerClient.DeleteImage(imageID)
			if err != nil {
				app.showError(fmt.Sprintf("Error deleting image: %v", err))
			} else {
				app.updateImages()
			}
		})
	default:
		app.showError("Unknown action: " + action)
		return
	}
}
