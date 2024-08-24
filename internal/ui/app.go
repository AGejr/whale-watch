package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
	"time"
)

type App struct {
	*tview.Application
	stopUpdateChan chan struct{}
	Pages          *tview.Pages
	Table          *tview.Table
	InputField     *tview.InputField
	TextView       *tview.TextView
	HelpModal      *tview.Modal
	ListLayout     *tview.Flex
	TextLayout     *tview.Flex
	searchInput    *tview.InputField
	searchMode     bool
	searchTerm     string
	searchResults  []string
	currentResult  int
	DockerClient   DockerClient
}

type DockerClient interface {
	ListContainers() ([][]string, error)
	ListImages() ([][]string, error)
	ListNetworks() ([][]string, error)
	DeleteNetwork(string) error
	ListVolumes() ([][]string, error)
	DeleteVolume(string) error
	InspectContainer(string) (string, error)
	InspectImage(string) (string, error)
	DeleteImage(string) error
	GetContainerLogs(string) (string, error)
	StartContainer(string) error
	StopContainer(string) error
	PauseContainer(string) error
	UnpauseContainer(string) error
	KillContainer(string) error
	DeleteContainer(string) error
}

func NewApp(dockerClient DockerClient) *App {
	tview.Borders.HorizontalFocus = tcell.RuneHLine
	tview.Borders.VerticalFocus = tcell.RuneVLine
	tview.Borders.TopLeftFocus = tcell.RuneULCorner
	tview.Borders.TopRightFocus = tcell.RuneURCorner
	tview.Borders.BottomLeftFocus = tcell.RuneLLCorner
	tview.Borders.BottomRightFocus = tcell.RuneLRCorner
	tview.Styles.BorderColor = tcell.ColorGrey

	app := &App{
		Application:  tview.NewApplication(),
		Pages:        tview.NewPages(),
		Table:        tview.NewTable(),
		InputField:   tview.NewInputField(),
		TextView:     tview.NewTextView(),
		DockerClient: dockerClient,
	}

	app.Table.SetBorder(true)

	app.HelpModal = app.createHelpModal()

	app.Table.SetSelectedFunc(app.selectTableItem)

	app.Table.SetBackgroundColor(tcell.NewRGBColor(10, 22, 34))

	app.InputField.SetAutocompleteStyles(tcell.NewRGBColor(10, 22, 34), tcell.Style{}.Foreground(tcell.ColorWhite), tcell.Style{}.Background(tcell.ColorLightBlue).Foreground(tcell.NewRGBColor(10, 22, 34)))
	app.InputField.SetFieldBackgroundColor(tcell.NewRGBColor(10, 22, 34))
	app.InputField.SetLabelColor(tcell.ColorWhite)

	app.InputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			command := strings.TrimSpace(strings.ToLower(app.InputField.GetText()))
			app.handleCommand(command)
			app.InputField.SetText("")
		}
	})

	app.InputField.SetAutocompleteFunc(func(currentText string) (entries []string) {
		keywords := []string{"container", "image", "volume", "network", "help", "q"}
		for _, keyword := range keywords {
			if strings.Contains(keyword, strings.ToLower(currentText)) {
				entries = append(entries, keyword)
			}
		}
		return entries
	})

	app.ListLayout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(app.Table, 0, 1, true).
		AddItem(app.InputField, 1, 0, false)

	app.TextView = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Application.Draw()
		})
	app.TextView.SetBorder(true)

	app.searchInput = tview.NewInputField().
		SetLabel("/").
		SetFieldWidth(0).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				app.performSearch()
			}
		})
	app.searchInput.SetFieldBackgroundColor(tcell.NewRGBColor(10, 22, 34))
	app.searchInput.SetLabelColor(tcell.ColorWhite)

	app.TextLayout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(app.TextView, 0, 1, true).
		AddItem(app.searchInput, 1, 0, false)

	app.Pages.AddPage("text", app.TextLayout, true, true)
	app.Pages.AddPage("help", app.HelpModal, true, true)
	app.Pages.AddPage("main", app.ListLayout, true, true)
	app.SetRoot(app.Pages, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if app.GetFocus() == app.TextView {
			switch event.Rune() {
			case '/':
				app.searchInput.SetText("")
				app.toggleSearchMode()
				return nil
			case 'n':
				app.navigateSearchResult(1)
				return nil
			case 'N':
				app.navigateSearchResult(-1)
				return nil
			}
		}
		if app.searchMode && event.Key() == tcell.KeyEscape {
			app.toggleSearchMode()
			return nil
		}
		if event.Rune() == ':' && app.GetFocus() != app.InputField && app.GetFocus() == app.Table {
			app.SetFocus(app.InputField)
			return nil
		}
		if event.Key() == tcell.KeyEscape && app.GetFocus() == app.TextView || app.GetFocus() == app.HelpModal {
			app.Pages.SwitchToPage("main")
			app.SetFocus(app.Table)
			return nil
		}
		if event.Rune() == '?' && app.GetFocus() == app.Table {
			app.showHelpModal()
			return nil
		}
		if app.GetFocus() == app.Table {
			switch app.Table.GetTitle() {
			case "Containers":
				switch event.Rune() {
				case 'l':
					app.showContainerLogs()
					return nil
				case 's':
					app.handleContainerAction("start")
					return nil
				case 'S':
					app.handleContainerAction("stop")
					return nil
				case 'p':
					app.handleContainerAction("pause")
					return nil
				case 'P':
					app.handleContainerAction("unpause")
					return nil
				case 'k':
					app.handleContainerAction("kill")
					return nil
				case 'd':
					app.handleContainerAction("delete")
					return nil
				}
			case "Images":
				switch event.Rune() {
				case 'd':
					app.handleImageAction("delete")
					return nil
				}
			case "Volumes":
				switch event.Rune() {
				case 'd':
					app.handleVolumeAction("delete")
					return nil
				}
			case "Networks":
				switch event.Rune() {
				case 'd':
					app.handleNetworkAction("delete")
					return nil
				}
			}

		}
		return event
	})

	app.showContainers()

	app.startAsyncUpdates()

	return app
}

func (app *App) selectTableItem(row, column int) {
	switch app.Table.GetTitle() {
	case "Containers":
		containerID := app.Table.GetCell(row, 0).Text
		containerName := app.Table.GetCell(row, 1).Text
		app.showContainerInfo(containerID, containerName)
	case "Images":
		imageID := app.Table.GetCell(row, 2).Text
		imageName := app.Table.GetCell(row, 0).Text
		imageTag := app.Table.GetCell(row, 1).Text
		app.showImageInfo(imageID, imageName, imageTag)
	default:
		return
	}
}

func (app *App) startAsyncUpdates() {
	app.stopUpdateChan = make(chan struct{})
	go func() {
		ticker := time.NewTicker(1 * time.Second) // Update every 5 seconds
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				app.Application.QueueUpdateDraw(func() {
					if app.GetFocus() == app.Table {
						switch app.Table.GetTitle() {
						case "Containers":
							app.updateContainers()
						case "Images":
							app.updateImages()
						case "Volumes":
							app.updateVolumes()
						case "Networks":
							app.updateNetworks()
						}
					}
				})
			case <-app.stopUpdateChan:
				return
			}
		}
	}()
}

func (app *App) stopAsyncUpdates() {
	if app.stopUpdateChan != nil {
		close(app.stopUpdateChan)
	}
}

func (app *App) Run() error {
	defer app.stopAsyncUpdates()
	return app.Application.Run()
}

func (app *App) handleCommand(command string) {
	switch command {
	case "container", "containers":
		app.showContainers()
	case "image", "images":
		app.showImages()
	case "network", "networks":
		app.showNetworks()
	case "volume", "volumes":
		app.showVolumes()
	case "help", "h":
		app.showHelpModal()
	case "q":
		app.Application.Stop()
	default:
		app.InputField.SetFieldBackgroundColor(tcell.ColorRed)
		app.SetFocus(app.Table)
	}
}

func (app *App) updateTable(title string, headers []string, data [][]string) {
	selectedRow, selectedColumn := app.Table.GetSelection()
	app.Table.Clear()
	app.Table.SetTitle(title)

	// Set header
	for col, header := range headers {
		cell := tview.NewTableCell(header).
			SetSelectable(false).
			SetAttributes(tcell.AttrBold).
			SetTextColor(tcell.ColorWhite)
		app.Table.SetCell(0, col, cell)
	}

	// Populate rows
	lightBlue := tcell.ColorLightBlue
	for row, rowData := range data {
		for col, cellData := range rowData {
			app.Table.SetCell(row+1, col, tview.NewTableCell(cellData).SetTextColor(lightBlue))
		}
	}

	app.Table.Select(selectedRow, selectedColumn).SetFixed(1, 0).SetSelectable(true, false)
}

func (app *App) showError(message string) {
	app.Table.Clear()
	app.Table.SetCell(0, 0, tview.NewTableCell(message))
}
