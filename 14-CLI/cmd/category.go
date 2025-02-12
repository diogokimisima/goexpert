/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// categoryCmd represents the category command
var categoryCmd = &cobra.Command{
	Use:   "category",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		fmt.Println("category called", name)
		exists, _ := cmd.Flags().GetBool("exists")
		fmt.Println("category exists", exists)
	},
}

func init() {
	rootCmd.AddCommand(categoryCmd)
	categoryCmd.PersistentFlags().StringP("name", "n", "Y", "Category name")
	categoryCmd.PersistentFlags().BoolP("exists", "e", false, "Check if category exists")
	categoryCmd.PersistentFlags().Int16P("id", "i", 0, "Category ID")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// categoryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// categoryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
