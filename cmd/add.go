/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/jonathantrevino/resumes-cli/internal/vault"
	"github.com/spf13/cobra"
)

var (
	company string
	role    string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <file>",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		role, _ := vault.Slug(role)
		company, _ := vault.Slug(company)
		fileName := args[0]

		// add file to requested company and role location in vault folder
		resume, err := vault.Add(fileName, company, role)
		fmt.Println(err)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Printf("resume json: %v\n", resume)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	addCmd.Flags().StringVarP(&company, "company", "c", "", "Set company to manually sort resume to")
	addCmd.Flags().StringVarP(&role, "role", "r", "", "Set role to manually sort resume to")

	addCmd.MarkFlagRequired("company")
	addCmd.MarkFlagRequired("role")
}
