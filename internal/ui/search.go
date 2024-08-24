package ui

import (
	"fmt"
	"strconv"
	"strings"
)

func addSearchResultRegion(src string, replace string, id string) string {
	regionClose := "[\"\"]"
	startIndex := strings.Index(strings.ToLower(src), strings.ToLower(replace))
	stopIndex := startIndex + len(replace)
	regionedSrc := strings.Join([]string{src[:startIndex], id, src[startIndex:stopIndex], regionClose, src[stopIndex:]}, "")
	return regionedSrc
}

func (app *App) performSearch() {
	app.searchTerm = app.searchInput.GetText()
	app.searchResults = []string{}
	app.currentResult = -1

	if app.searchTerm == "" {
		return
	}

	text := app.TextView.GetText(false)
	lines := strings.Split(text, "\n")

	regionID := "[\"0\"]"
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), strings.ToLower(app.searchTerm)) {
			app.searchResults = append(app.searchResults, fmt.Sprintf("%d:%s", i, line))
			newLine := addSearchResultRegion(line, app.searchTerm, regionID)
			text = strings.Replace(text, line, newLine, 1)
		}
	}

	app.TextView.SetText(text)

	if len(app.searchResults) > 0 {
		app.TextView.Highlight("0")
		app.navigateSearchResult(1)
	}
	app.SetFocus(app.TextView)
}

func (app *App) toggleSearchMode() {
	app.searchMode = !app.searchMode
	if app.searchMode {
		app.SetFocus(app.searchInput)
	} else {
		app.TextView.SetText(app.TextView.GetText(true))
		app.searchInput.SetText("")
		app.SetFocus(app.TextView)
	}
}

func (app *App) navigateSearchResult(direction int) {
	if len(app.searchResults) == 0 {
		return
	}

	app.currentResult += direction
	if app.currentResult >= len(app.searchResults) {
		app.currentResult = 0
	} else if app.currentResult < 0 {
		app.currentResult = len(app.searchResults) - 1
	}

	parts := strings.SplitN(app.searchResults[app.currentResult], ":", 2)
	if len(parts) == 2 {
		lineNumber, _ := strconv.Atoi(parts[0])
		app.TextView.ScrollTo(lineNumber, len(parts[1]))
	}
}
