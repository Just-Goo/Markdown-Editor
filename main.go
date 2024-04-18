package main

import (
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type config struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI
	SaveMenuItem  *fyne.MenuItem
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
	cfg.createMenuItems(window)

	// set the content of the window
	window.SetContent(container.NewHSplit(edit, preview))

	// resize the window
	window.Resize(fyne.Size{Width: 800, Height: 500})

	// center the window on screen
	window.CenterOnScreen()

	// show window and run app
	window.ShowAndRun()
}

func (c *config) makeUI() (*widget.Entry, *widget.RichText) {
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("")

	c.EditWidget = edit
	c.PreviewWidget = preview

	edit.OnChanged = preview.ParseMarkdown

	return edit, preview
}

func (c *config) createMenuItems(win fyne.Window) {
	openMenuItem := fyne.NewMenuItem("Open...", c.openFunc(win))
	saveMenuItem := fyne.NewMenuItem("Save", c.saveFunc(win))
	saveAsMenuItem := fyne.NewMenuItem("Save as...", c.saveAsFunc(win))
	c.SaveMenuItem = saveMenuItem
	c.SaveMenuItem.Disabled = true

	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)
	menu := fyne.NewMainMenu(fileMenu)

	win.SetMainMenu(menu)
}

var myFilter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

func (c *config) saveFunc(win fyne.Window) func()  {
	return func() {
		if c.CurrentFile != nil { // this function should only run if there's a file currently open
			write, err := storage.Writer(c.CurrentFile)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			write.Write([]byte(c.EditWidget.Text))
			defer write.Close()
		}
	}
}

func (c *config) openFunc(win fyne.Window) func()  {
	return func() {
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			// user cancelled
			if read == nil {
				return
			}

			defer read.Close() // to avoid resource leak

			data, err := io.ReadAll(read)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			c.EditWidget.SetText(string(data))
			c.CurrentFile = read.URI() // keep track of what file is currently open

			win.SetTitle(read.URI().Name()) // change the window's title to the opened file name
			c.SaveMenuItem.Disabled = false

		}, win)
		openDialog.SetFilter(myFilter)
		openDialog.Show()
	}
}

func (c *config) saveAsFunc(win fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			// user cancelled
			if write == nil {
				return
			}

			// check if file name is saved with '.md'
			if !strings.HasSuffix(strings.ToLower(write.URI().String()), ".md") {
				dialog.ShowInformation("Error", "file must end with .md extension", win)
				return
			}

			write.Write([]byte(c.EditWidget.Text))// save file
			c.CurrentFile = write.URI() // keep track of what file is currently open

			defer write.Close() // clean up resources to avoid memory leak

			win.SetTitle(write.URI().Name()) // change the window's title to the saved file name

			c.SaveMenuItem.Disabled = false // enable 'save menu item' after the file has been saved the first time
		}, win)

		saveDialog.SetFileName("untitled.md") // default name
		saveDialog.SetFilter(myFilter)
		saveDialog.Show()
	}
}
