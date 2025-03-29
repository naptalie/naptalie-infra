package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SimpleMeowTUI creates a basic but functional MeowTUI
func SimpleMeowTUI() error {
	// Create the application
	app := tview.NewApplication()

	// Create a flex layout for the main screen
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Add a title
	title := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("? MeowTUI - Purr-fect K3s Dashboard ?")
	title.SetTextColor(tcell.ColorPink)

	// Create two panels - left for resources, right for content
	midFlex := tview.NewFlex()

	// Resource list (left panel)
	resourceList := tview.NewList().
		SetHighlightFullLine(true).
		SetSelectedTextColor(tcell.ColorWhite).
		SetSelectedBackgroundColor(tcell.ColorPink)
	
	resourceList.SetBorder(true).
		SetTitle(" Resources ")
	
	// Details panel (right)
	details := tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	details.SetBorder(true).
		SetTitle(" Details ")
	
	// Status bar
	statusBar := tview.NewTextView().
		SetText("Press q to quit | ? for help | Namespace: default")
	statusBar.SetTextColor(tcell.ColorPink)

	// Add items to the resource list
	resources := []string{"pods", "deployments", "services", "nodes", "namespaces"}
	for _, res := range resources {
		r := res // Local copy to use in closure
		resourceList.AddItem("? "+r, "", 0, func() {
			loadResource(r, details)
		})
	}

	// Build the layout
	midFlex.AddItem(resourceList, 0, 1, true).
		AddItem(details, 0, 3, false)
	
	flex.AddItem(title, 1, 0, false).
		AddItem(midFlex, 0, 1, true).
		AddItem(statusBar, 1, 0, false)

	// Show welcome message
	details.SetText("? Welcome to MeowTUI!\n\n" +
		"Select a resource from the list to view it.\n\n" +
		"[pink]Keyboard shortcuts:[-]\n" +
		"- q: Quit\n" +
		"- Tab: Switch panels\n" +
		"- Arrows: Navigate\n" +
		"- Enter: Select")

	// Set key handler
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Handle global keys
		switch event.Key() {
		case tcell.KeyTab:
			// Switch focus between panels
			if app.GetFocus() == resourceList {
				app.SetFocus(details)
			} else {
				app.SetFocus(resourceList)
			}
			return nil
		case tcell.KeyEsc, tcell.KeyCtrlC:
			app.Stop()
			return nil
		}

		// Handle runes (individual keys)
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'q':
				app.Stop()
				return nil
			case '?':
				showHelp(details)
				return nil
			}
		}

		return event
	})

	// Load the first resource on startup if there are any
	if resourceList.GetItemCount() > 0 {
		resourceList.SetCurrentItem(0)
		loadResource(resources[0], details)
	}

	// Start the application
	return app.SetRoot(flex, true).EnableMouse(true).Run()
}

// loadResource loads a kubernetes resource and displays it
func loadResource(resource string, view *tview.TextView) {
	view.SetText(fmt.Sprintf("? Loading %s...", resource))
	
	// Run kubectl to get the resource
	var cmd *exec.Cmd
	if resource == "namespaces" || resource == "nodes" {
		cmd = exec.Command("kubectl", "get", resource)
	} else {
		cmd = exec.Command("kubectl", "get", resource, "-n", "default")
	}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		view.SetText(fmt.Sprintf("? Error: %s\n\n%s", err, string(output)))
		return
	}
	
	// Add colors to the output
	colorized := colorizeOutput(string(output))
	view.SetText(colorized)
}

// colorizeOutput adds color to kubectl output
func colorizeOutput(output string) string {
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return output
	}
	
	// Color the header
	result := "[pink]" + lines[0] + "[-]\n"
	
	// Color the rest based on status
	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			continue
		}
		
		line := lines[i]
		if strings.Contains(strings.ToLower(line), "running") || 
		   strings.Contains(strings.ToLower(line), "active") {
			result += "[green]" + line + "[-]\n"
		} else if strings.Contains(strings.ToLower(line), "pending") {
			result += "[yellow]" + line + "[-]\n"
		} else if strings.Contains(strings.ToLower(line), "error") || 
				  strings.Contains(strings.ToLower(line), "failed") {
			result += "[red]" + line + "[-]\n"
		} else {
			result += line + "\n"
		}
	}
	
	return result
}

// showHelp displays help information
func showHelp(view *tview.TextView) {
	help := `[pink]? MeowTUI Help[white]

[pink]Keyboard Shortcuts:[-]
- [pink]Tab[-]: Switch between panels
- [pink]Arrow keys[-]: Navigate the list
- [pink]Enter[-]: Select an item
- [pink]q[-]: Quit the application
- [pink]?[-]: Show this help

[pink]Tips:[-]
- The left panel shows Kubernetes resources
- The right panel shows details
- Select a resource to view its instances
- Make sure your K3s cluster is running

[pink]Meow and enjoy! ?[-]`

	view.SetText(help)
}
