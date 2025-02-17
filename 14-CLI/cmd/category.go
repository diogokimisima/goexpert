/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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
		db := GetDb()
		category := GetCategoryDB(db)

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		category.Create(name, description)

	},
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("category called with name ", category)
	// 	name, _ := cmd.Flags().GetString("name")
	// 	fmt.Println("category called", name)
	// 	exists, _ := cmd.Flags().GetBool("exists")
	// 	fmt.Println("category exists", exists)
	// },
	// PreRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("Chamado antes do Run")
	// },
	// PostRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("Chamado depois do Run")
	// },
}

func init() {
	rootCmd.AddCommand(categoryCmd)
	createCmd.Flags().StringP("name", "n", "", "Category name")
	createCmd.Flags().StringP("description", "n", "", "Description of the category")
	createCmd.MarkFlagsRequiredTogether("name", "description")
	// categoryCmd.PersistentFlags().StringVarP(&category, "name", "n", "", "Category name")

	// categoryCmd.PersistentFlags().BoolP("exists", "e", false, "Check if category exists")
	// categoryCmd.PersistentFlags().Int16P("id", "i", 0, "Category ID")
}
