package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

type config struct {
	EditWidget *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile fyne.URI
	SaveMenuItem *fyne.MenuItem
}

// var cfg config

func main() {
	
	// create a fyne app
	myApp := app.New()

	// create a new window with a title
	window := myApp.NewWindow("Markdown Editor")

	cfg := config{}
	// get the user interface
	edit, preview := cfg.makeUI()

	// set the content of the window
	window.SetContent(container.NewHSplit(edit, preview))	

	// resize the window
	window.Resize(fyne.Size{Width: 800, Height: 500})
	
	// center the window on screen
	window.CenterOnScreen()

	// show window and run app
	window.ShowAndRun()
}

func (c *config) makeUI() (*widget.Entry, *widget.RichText)  {
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("")

	c.EditWidget = edit
	c.PreviewWidget = preview
	
	edit.OnChanged = preview.ParseMarkdown

	return edit, preview
}