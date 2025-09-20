/*
Copyright © 2025 Bernard Katamanso <bernard@orctatech.com>
*/
package cmd

import (
	"GitMate/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"log"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new feature branch workflow",
	Long:  ` This command will start a new workflow.`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(tui.NewModel(tui.CmdStart)) // pass initial command
		if err, _ := p.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
