package ui

import "github.com/rivo/tview"

func (app *App) showHelpModal() {
	//app.Pages.SwitchToPage("help")
	//app.SetFocus(app.HelpModal)

	helpText := `General actions:
	?: Show this help message
	:: Enter command mode
	ESC: Exit text view or close modal
	Enter: Inspect selected item

Search (in TextView):
	/: Enter search mode
	n: Next search result
	N: Previous search result
	ESC: Exit search mode

Container Actions:
	s: Start container
	S: Stop container
	p: Pause container
	P: Unpause container
	k: Kill container
	d: Delete container
	l: Show logs

Image Actions:
    d: Delete image

Volume Actions:
    d: Delete volume

Network Actions:
    d: Delete network`
	app.TextView.SetText(helpText)
	app.Pages.SwitchToPage("text")
}

func (app *App) createHelpModal() *tview.Modal {
	helpText := `Keybindings:
	?: Show this help message
	:: Enter command mode
	ESC: Exit text view or close modal
	Enter: Inspect selected item

	Search (in TextView):
	/: Enter search mode
	n: Next search result
	N: Previous search result
	ESC: Exit search mode

	Container Actions:
	s: Start container
	S: Stop container
	p: Pause container
	P: Unpause container
	k: Kill container
	d: Delete container
	l: Show logs

	Image Actions:
    d: Delete image

	Volume Actions:
    d: Delete volume

    Network Actions:
    d: Delete network`

	modal := tview.NewModal().
		SetText(helpText).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.Pages.SwitchToPage("main")
			app.SetFocus(app.Table)
		})

	return modal
}

func (app *App) showConfirmationDialog(title, message string, onConfirm func()) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				onConfirm()
			}
			app.Pages.RemovePage("confirmation")
			app.SetFocus(app.Table)
		})

	app.Pages.AddPage("confirmation", modal, true, true)
	app.SetFocus(modal)
}
