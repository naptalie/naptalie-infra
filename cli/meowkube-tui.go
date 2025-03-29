package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Global variables for application state
type AppState struct {
	currentNamespace string
	namespaces       []string
}

// SimpleMeowTUI creates a basic but functional MeowTUI
func SimpleMeowTUI() error {
	// Initialize app state
	state := AppState{
		currentNamespace: "default",
	}

	// Load available namespaces
	loadNamespaces(&state)

	// Create the application
	app := tview.NewApplication()

	// Create a flex layout for the main screen
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Add a title
	title := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("? MeowTUI - Purr-fect K3s Dashboard ?")
	title.SetTextColor(tcell.ColorPink)

	// Create three panels - left for resources, right for content, top for namespace selection
	topFlex := tview.NewFlex()
	midFlex := tview.NewFlex()

	// Namespace selector (top panel)
	namespaceSelector := tview.NewDropDown().
		SetLabel("Namespace: ").
		SetFieldWidth(20)
	
	for _, ns := range state.namespaces {
		namespaceSelector.AddOption(ns, nil)
	}
	
	// Set the default namespace
	for i, ns := range state.namespaces {
		if ns == state.currentNamespace {
			namespaceSelector.SetCurrentOption(i)
			break
		}
	}

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
		SetText(fmt.Sprintf("Press q to quit | ? for help | Namespace: %s", state.currentNamespace))
	statusBar.SetTextColor(tcell.ColorPink)

	// Resources to display
	resources := []string{"pods", "deployments", "services", "nodes", "namespaces"}
	
	// Update namespaceSelector's selected function to use statusBar
	namespaceSelector.SetSelectedFunc(func(text string, index int) {
		state.currentNamespace = text
		statusBar.SetText(fmt.Sprintf("Press q to quit | ? for help | Namespace: %s", state.currentNamespace))
		// Reload current resource with new namespace
		if resourceList.GetItemCount() > 0 {
			idx := resourceList.GetCurrentItem()
			if idx >= 0 && idx < len(resources) {
				loadResource(resources[idx], details, state.currentNamespace)
			}
		}
	})
	
	// Add items to the resource list
	for _, res := range resources {
		r := res // Local copy to use in closure
		resourceList.AddItem("? "+r, "", 0, func() {
			loadResource(r, details, state.currentNamespace)
		})
	}

	// Add refresh button
	refreshButton := tview.NewButton("? Refresh").
		SetSelectedFunc(func() {
			if resourceList.GetItemCount() > 0 {
				idx := resourceList.GetCurrentItem()
				if idx >= 0 && idx < len(resources) {
					loadResource(resources[idx], details, state.currentNamespace)
				}
			}
		})

	// Build the layout
	topFlex.AddItem(namespaceSelector, 0, 1, false).
		AddItem(refreshButton, 10, 0, false)

	midFlex.AddItem(resourceList, 0, 1, true).
		AddItem(details, 0, 3, false)
	
	flex.AddItem(title, 1, 0, false).
		AddItem(topFlex, 1, 0, false).
		AddItem(midFlex, 0, 1, true).
		AddItem(statusBar, 1, 0, false)

	// Show welcome message
	details.SetText("? Welcome to MeowTUI!\n\n" +
		"Select a resource from the list to view it.\n\n" +
		"[pink]Keyboard shortcuts:[-]\n" +
		"- q: Quit\n" +
		"- n: Change namespace\n" +
		"- r: Refresh view\n" +
		"- Tab: Switch panels\n" +
		"- Arrows: Navigate\n" +
		"- Enter: Select")

	// Set key handler
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Handle global keys
		switch event.Key() {
		case tcell.KeyTab:
			// Cycle focus between panels: resource list -> details -> namespace selector -> back to resource list
			if app.GetFocus() == resourceList {
				app.SetFocus(details)
			} else if app.GetFocus() == details {
				app.SetFocus(namespaceSelector)
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
			case 'n':
				// Focus on namespace selector
				app.SetFocus(namespaceSelector)
				return nil
			case 'r':
				// Refresh current view
				if resourceList.GetItemCount() > 0 {
					idx := resourceList.GetCurrentItem()
					if idx >= 0 && idx < len(resources) {
						loadResource(resources[idx], details, state.currentNamespace)
					}
				}
				return nil
			}
		}

		return event
	})

	// Load the first resource on startup if there are any
	if resourceList.GetItemCount() > 0 {
		resourceList.SetCurrentItem(0)
		loadResource(resources[0], details, state.currentNamespace)
	}

	// Start the application
	return app.SetRoot(flex, true).EnableMouse(true).Run()
}

// loadNamespaces gets all available namespaces and updates the app state
func loadNamespaces(state *AppState) {
	cmd := exec.Command("kubectl", "get", "namespaces", "-o", "jsonpath={.items[*].metadata.name}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If can't get namespaces, default to just "default"
		state.namespaces = []string{"default"}
		return
	}
	
	// Split the output by spaces
	namespaces := strings.Fields(string(output))
	if len(namespaces) == 0 {
		state.namespaces = []string{"default"}
	} else {
		state.namespaces = namespaces
	}
}

// loadResource loads a kubernetes resource and displays it
func loadResource(resource string, view *tview.TextView, namespace string) {
	view.SetText(fmt.Sprintf("? Loading %s...", resource))
	
	// Run kubectl to get the resource
	var cmd *exec.Cmd
	if resource == "namespaces" || resource == "nodes" {
		// These are cluster-scoped resources, not namespace-scoped
		cmd = exec.Command("kubectl", "get", resource)
	} else {
		cmd = exec.Command("kubectl", "get", resource, "-n", namespace)
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
- [pink]Tab[-]: Cycle between panels
- [pink]Arrow keys[-]: Navigate the list
- [pink]Enter[-]: Select an item
- [pink]n[-]: Focus on namespace selector
- [pink]r[-]: Refresh current view
- [pink]q[-]: Quit the application
- [pink]?[-]: Show this help

[pink]Namespace Selection:[-]
- Use the dropdown at the top to select a namespace
- Resources will be filtered by the selected namespace
- Some resources (nodes, namespaces) are cluster-wide

[pink]Tips:[-]
- The left panel shows Kubernetes resources
- The right panel shows details
- Select a resource to view its instances
- Make sure your K3s cluster is running

[pink]Meow and enjoy! ?[-]`

	view.SetText(help)
}