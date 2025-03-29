package main

import (
	"fmt"
//	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// MeowTUI is the terminal UI component for Meowkube
type MeowTUI struct {
	app           *tview.Application
	pages         *tview.Pages
	resourceList  *tview.List
	detailView    *tview.TextView
	statusBar     *tview.TextView
	commandInput  *tview.InputField
	currentNS     string
	currentView   string
	resources     []string
	selectedIndex int
}

// Colors for the pink theme
var (
	colorBackground = tcell.GetColor("#2e0a1a") // Dark pink background
	colorPrimary    = tcell.GetColor("#ff6ec7") // Bright pink
	colorSecondary  = tcell.GetColor("#ff9fe5") // Light pink
	colorHighlight  = tcell.GetColor("#ffcbf2") // Very light pink
	colorText       = tcell.GetColor("#ffffff") // White text
	colorWarning    = tcell.GetColor("#ffca28") // Warning (amber)
	colorError      = tcell.GetColor("#ff5252") // Error (red)
	colorSuccess    = tcell.GetColor("#69f0ae") // Success (green)
)

// NewMeowTUI creates a new terminal UI
func NewMeowTUI() *MeowTUI {
	tui := &MeowTUI{
		app:         tview.NewApplication(),
		currentNS:   "default",
		currentView: "pods",
		resources:   []string{"pods", "deployments", "services", "nodes", "namespaces"},
	}

	tui.initUI()
	return tui
}

// initUI sets up the UI components
func (t *MeowTUI) initUI() {
	// Main layout
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Title bar with cat ears
	titleBar := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("/?????\\ MeowTUI - Purr-fect Kubernetes Dashboard /?????\\").
		SetTextColor(colorPrimary)
	titleBar.SetBackgroundColor(colorBackground)

	// Resource list (left panel)
	t.resourceList = tview.NewList().
		SetHighlightFullLine(true).
		SetMainTextColor(colorText).
		SetSelectedTextColor(colorHighlight).
		SetSelectedBackgroundColor(colorPrimary)
	
	t.resourceList.SetBorder(true).
		SetTitle(" ? Resources ").
		SetTitleColor(colorPrimary).
		SetBorderColor(colorPrimary)
	t.resourceList.SetBackgroundColor(colorBackground)

	// Detail view (right panel)
	t.detailView = tview.NewTextView().
		SetDynamicColors(true).
		SetTextColor(colorText)
	t.detailView.SetBorder(true).
		SetTitle(" ? Details ").
		SetTitleColor(colorPrimary).
		SetBorderColor(colorPrimary)
	t.detailView.SetBackgroundColor(colorBackground)
	t.detailView.SetScrollable(true)

	// Middle layout (list and details)
	middleFlex := tview.NewFlex().
		AddItem(t.resourceList, 0, 1, true).
		AddItem(t.detailView, 0, 3, false)

	// Status bar
	t.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetText("[" + colorToHex(colorPrimary) + "]Namespace:[-] " + t.currentNS + 
			" | [" + colorToHex(colorPrimary) + "]View:[-] " + t.currentView + 
			" | [" + colorToHex(colorPrimary) + "]Meow for help[-]")
	t.statusBar.SetBackgroundColor(colorBackground)

	// Command input
	t.commandInput = tview.NewInputField().
		SetLabel("Command: ").
		SetLabelColor(colorPrimary).
		SetFieldTextColor(colorText)
	t.commandInput.SetBackgroundColor(colorBackground)
	t.commandInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			t.handleCommand(t.commandInput.GetText())
			t.commandInput.SetText("")
		}
	})

	// Help page
	helpText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText(helpContent())
	helpText.SetBorder(true).
		SetTitle(" ? MeowTUI Help ").
		SetTitleColor(colorPrimary).
		SetBorderColor(colorPrimary)
	helpText.SetBackgroundColor(colorBackground)
	helpText.SetTextColor(colorText)

	// Add everything to the pages
	t.pages = tview.NewPages().
		AddPage("main", flex.
			AddItem(titleBar, 1, 0, false).
			AddItem(middleFlex, 0, 1, true).
			AddItem(t.statusBar, 1, 0, false).
			AddItem(t.commandInput, 1, 0, false), true, true).
		AddPage("help", helpText, true, false)

	// Initial data load
	t.populateResourceList()
}

// colorToHex converts a tcell.Color to a hex string for tview
func colorToHex(color tcell.Color) string {
	r, g, b := color.RGB()
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// populateResourceList fills the resource list with available resources
func (t *MeowTUI) populateResourceList() {
	t.resourceList.Clear()

	for i, resource := range t.resources {
		t.resourceList.AddItem(fmt.Sprintf("? %s", resource), "", rune('a'+i), func() {
			idx := t.resourceList.GetCurrentItem()
			selectedResource := t.resources[idx]
			t.currentView = selectedResource
			t.loadResourceDetails(selectedResource)
			t.updateStatusBar()
		})
	}
}

// loadResourceDetails loads details for the selected resource
func (t *MeowTUI) loadResourceDetails(resource string) {
	// Show loading message
	t.detailView.SetText("[" + colorToHex(colorPrimary) + "]? Fetching " + resource + "...[" + colorToHex(colorText) + "]\n")
	t.app.Draw()

	// Prepare command
	var cmd *exec.Cmd
	if resource == "namespaces" {
		cmd = exec.Command("kubectl", "get", resource, "-o", "wide")
	} else if resource == "nodes" {
		cmd = exec.Command("kubectl", "get", resource, "-o", "wide")
	} else {
		cmd = exec.Command("kubectl", "get", resource, "-n", t.currentNS, "-o", "wide")
	}

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.detailView.SetText("[" + colorToHex(colorError) + "]? Error fetching " + resource + ": " + err.Error() + "[" + colorToHex(colorText) + "]\n\n" + string(output))
		return
	}

	// Format output with colors
	formattedOutput := formatKubectlOutput(string(output))
	t.detailView.SetText(formattedOutput)
}

// formatKubectlOutput adds colors to kubectl output
func formatKubectlOutput(output string) string {
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return output
	}

	// Color the header row
	headerColor := colorToHex(colorPrimary)
	lines[0] = "[" + headerColor + "]" + lines[0] + "[-]"

	// Color alternating rows
	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			continue
		}

		// Check for status in the line to color appropriately
		lineColor := colorToHex(colorText)
		if strings.Contains(strings.ToLower(lines[i]), "running") || 
		   strings.Contains(strings.ToLower(lines[i]), "active") || 
		   strings.Contains(strings.ToLower(lines[i]), "ready") {
			lineColor = colorToHex(colorSuccess)
		} else if strings.Contains(strings.ToLower(lines[i]), "pending") || 
				  strings.Contains(strings.ToLower(lines[i]), "unknown") {
			lineColor = colorToHex(colorWarning)
		} else if strings.Contains(strings.ToLower(lines[i]), "failed") || 
				  strings.Contains(strings.ToLower(lines[i]), "error") || 
				  strings.Contains(strings.ToLower(lines[i]), "crashloopbackoff") {
			lineColor = colorToHex(colorError)
		}

		lines[i] = "[" + lineColor + "]" + lines[i] + "[-]"
	}

	return strings.Join(lines, "\n")
}

// handleCommand processes user commands
func (t *MeowTUI) handleCommand(cmd string) {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return
	}

	parts := strings.Fields(cmd)
	switch parts[0] {
	case "ns", "namespace":
		if len(parts) > 1 {
			t.switchNamespace(parts[1])
		}
	case "view":
		if len(parts) > 1 && contains(t.resources, parts[1]) {
			t.currentView = parts[1]
			t.loadResourceDetails(parts[1])
		}
	case "describe":
		if len(parts) > 2 {
			t.describeResource(parts[1], parts[2])
		}
	case "logs":
		if len(parts) > 1 {
			t.showPodLogs(parts[1])
		}
	case "refresh", "r":
		t.loadResourceDetails(t.currentView)
	case "help", "h", "?":
		t.pages.SwitchToPage("help")
	case "quit", "q", "exit":
		t.app.Stop()
	default:
		t.detailView.SetText("[" + colorToHex(colorError) + "]? Unknown command: " + cmd + "[" + colorToHex(colorText) + "]\n" +
			"Try 'help' for a list of commands")
	}

	t.updateStatusBar()
}

// switchNamespace changes the current namespace
func (t *MeowTUI) switchNamespace(namespace string) {
	// Verify namespace exists
	cmd := exec.Command("kubectl", "get", "namespace", namespace)
	if err := cmd.Run(); err != nil {
		t.detailView.SetText("[" + colorToHex(colorError) + "]? Namespace not found: " + namespace + "[" + colorToHex(colorText) + "]")
		return
	}

	t.currentNS = namespace
	t.loadResourceDetails(t.currentView)
	t.detailView.SetText("[" + colorToHex(colorSuccess) + "]? Switched to namespace: " + namespace + "[" + colorToHex(colorText) + "]\n\n" + t.detailView.GetText(false))
}

// describeResource shows detailed info about a resource
func (t *MeowTUI) describeResource(resourceType, name string) {
	cmd := exec.Command("kubectl", "describe", resourceType, name, "-n", t.currentNS)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.detailView.SetText("[" + colorToHex(colorError) + "]? Error describing " + resourceType + "/" + name + ": " + err.Error() + "[" + colorToHex(colorText) + "]\n\n" + string(output))
		return
	}

	t.detailView.SetText("[" + colorToHex(colorPrimary) + "]? Description of " + resourceType + "/" + name + ":[" + colorToHex(colorText) + "]\n\n" + string(output))
}

// showPodLogs displays logs for a pod
func (t *MeowTUI) showPodLogs(podName string) {
	cmd := exec.Command("kubectl", "logs", podName, "-n", t.currentNS)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.detailView.SetText("[" + colorToHex(colorError) + "]? Error getting logs for " + podName + ": " + err.Error() + "[" + colorToHex(colorText) + "]\n\n" + string(output))
		return
	}

	t.detailView.SetText("[" + colorToHex(colorPrimary) + "]? Logs for " + podName + ":[" + colorToHex(colorText) + "]\n\n" + string(output))
}

// updateStatusBar refreshes the status bar
func (t *MeowTUI) updateStatusBar() {
	t.statusBar.SetText("[" + colorToHex(colorPrimary) + "]Namespace:[-] " + t.currentNS + 
		" | [" + colorToHex(colorPrimary) + "]View:[-] " + t.currentView + 
		" | [" + colorToHex(colorPrimary) + "]Meow for help (Type 'help')[-]")
}

// startAutoRefresh begins the auto-refresh loop
func (t *MeowTUI) startAutoRefresh(seconds int) {
	go func() {
		for {
			time.Sleep(time.Duration(seconds) * time.Second)
			t.app.QueueUpdateDraw(func() {
				t.loadResourceDetails(t.currentView)
			})
		}
	}()
}
// Run starts the terminal UI
func (t *MeowTUI) Run() error {
	// Set up keyboard shortcuts
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Global keys
		switch event.Key() {
		case tcell.KeyEsc:
			// If in help page, go back to main
			if t.pages.HasPage("help") {
				currentPage, _ := t.pages.GetFrontPage()
				if currentPage == "help" {
					t.pages.SwitchToPage("main")
					return nil
				}
			}
		case tcell.KeyCtrlC, tcell.KeyCtrlQ:
			t.app.Stop()
			return nil
		case tcell.KeyTab:
			// Toggle focus between resource list and detail view
			if t.app.GetFocus() == t.resourceList {
				t.app.SetFocus(t.detailView)
			} else if t.app.GetFocus() == t.detailView {
				t.app.SetFocus(t.commandInput)
			} else {
				t.app.SetFocus(t.resourceList)
			}
			return nil
		case tcell.KeyCtrlR:
			t.loadResourceDetails(t.currentView)
			return nil
		}

		// Runes
		switch event.Rune() {
		case ':':
			// Focus command input if not already focused
			if t.app.GetFocus() != t.commandInput {
				t.app.SetFocus(t.commandInput)
				t.commandInput.SetText("")
				return nil
			}
		case 'q':
			// Only handle 'q' if not in an input field
			if t.app.GetFocus() != t.commandInput {
				t.app.Stop()
				return nil
			}
		case 'h', '?':
			// Only handle 'h' if not in an input field
			if t.app.GetFocus() != t.commandInput {
				t.pages.SwitchToPage("help")
				return nil
			}
		}
		return event
	})

	// Start auto-refresh (every 30 seconds)
	t.startAutoRefresh(30)

	// Start the application
	return t.app.SetRoot(t.pages, true).EnableMouse(true).Run()
}

// helpContent returns the help text
func helpContent() string {
	pink := "[" + colorToHex(colorPrimary) + "]"
	reset := "[-]"
	
	return pink + "? Welcome to MeowTUI - Your Cat-Themed Kubernetes Dashboard! ?" + reset + "\n\n" +
		pink + "Keyboard Shortcuts:" + reset + "\n" +
		"- " + pink + "Tab:" + reset + " Cycle focus between resource list, details, and command input\n" +
		"- " + pink + "Esc:" + reset + " Go back or close popup\n" +
		"- " + pink + "Ctrl+R:" + reset + " Refresh current view\n" +
		"- " + pink + "Ctrl+C/Ctrl+Q:" + reset + " Quit\n" +
		"- " + pink + "h, ?:" + reset + " Show this help\n" +
		"- " + pink + "q:" + reset + " Quit\n" +
		"- " + pink + ":"+":" + reset + " Focus command input\n\n" +
		
		pink + "Commands:" + reset + "\n" +
		"- " + pink + "ns, namespace <n>:" + reset + " Switch to namespace\n" +
		"- " + pink + "view <resource>:" + reset + " View resource type (pods, deployments, etc.)\n" +
		"- " + pink + "describe <resource> <n>:" + reset + " Describe resource\n" +
		"- " + pink + "logs <pod-name>:" + reset + " Show pod logs\n" +
		"- " + pink + "refresh, r:" + reset + " Refresh current view\n" +
		"- " + pink + "help, h, ?:" + reset + " Show this help\n" +
		"- " + pink + "quit, q, exit:" + reset + " Exit MeowTUI\n\n" +
		
		"Press " + pink + "Esc" + reset + " to return to the main view."
}

// contains checks if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
