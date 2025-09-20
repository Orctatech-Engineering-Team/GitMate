/*
Copyright © 2025 Bernard Katamanso <bernard@orctatech.com>
*/
package cmd

import (
	"GitMate/internal/tui"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync your branch with main",
	Long:  `This command will sync your branch with main.`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	p := tea.NewProgram(tui.NewModel(tui.CmdClean))
	//	if err, _ := p.Run(); err != nil {
	//		log.Fatal(err)
	//	}
	//},
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.RunSyncTUI()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
