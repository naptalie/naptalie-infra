package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// TUI command
var tuiCmd = &cobra.Command{
	Use:   "meow",
	Short: "Start the interactive MeowTUI",
	Long:  "Start MeowTUI - An interactive pink cat-themed terminal UI for managing your K3s cluster",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("? Meow! Starting interactive terminal UI...")
		tui := NewMeowTUI()
		if err := tui.Run(); err != nil {
			fmt.Println("? Hiss! Error running TUI:", err)
		}
	},
}
