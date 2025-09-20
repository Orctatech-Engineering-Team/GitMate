/*
Copyright © 2025 Bernard Katamanso <bernard@orctatech.com>
*/
package cmd

import (
	"github.com/Orctatech-Engineering-Team/GitMate/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"log"

	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up commits interactively (squash/fixup)",
	Long:  `This command will clean up commits interactively (squash/fixup).`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(tui.NewModel(tui.CmdClean))
		if err, _ := p.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
