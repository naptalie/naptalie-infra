package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// TUI command
var tuiCmd = &cobra.Command{
	Use:   "purr-tui",
	Short: "Start the interactive MeowTUI",
	Long:  "Start MeowTUI - An interactive pink cat-themed terminal UI for managing your K3s cluster",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("? Meow! Starting interactive terminal UI...")
		fmt.Println("? [Controls]")
		fmt.Println("- Arrow keys: Navigate resources")
		fmt.Println("- Tab: Cycle through panels")
		fmt.Println("- q or Ctrl+C: Quit")
		fmt.Println("- h or ?: Show help")
		fmt.Println("- :: Enter command")
		
		tui := NewMeowTUI()
		if err := tui.Run(); err != nil {
			fmt.Println("? Hiss! Error running TUI:", err)
		}
	},
}
