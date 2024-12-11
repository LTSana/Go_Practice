/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Shows a list of all your tasks.",
	Long:  `Shows a list of all your tasks that are stored in the CSV file.`,
	Run: func(cmd *cobra.Command, args []string) {

		showAll, _ := cmd.Flags().GetBool("all")

		if showAll {
			listTasks(true)
		} else {
			listTasks(false)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().Bool("all", false, "Shows the done status of the tasks.")
	listCmd.Flags().BoolP("all", "a", false, "Shows the done status of the tasks.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
