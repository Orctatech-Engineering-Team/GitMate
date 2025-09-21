/*
Copyright © 2025 Bernard Katamanso <bernard@orctatech.com>
*/
package cmd

import (
	"github.com/Orctatech-Engineering-Team/GitMate/internal/tui"
	"github.com/spf13/cobra"
)

// tutorCmd represents the tutor command
var tutorCmd = &cobra.Command{
	Use:   "tutor",
	Short: "Interactive tutorial for Git workflows",
	Long:  `This command will guide you through the basics of Git workflows.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.RunCleanTUI()
	},
}

func init() {
	rootCmd.AddCommand(tutorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tutorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tutorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
