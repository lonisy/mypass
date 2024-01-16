/*
Copyright © 2024 lonisy@163.com

*/
package cmd

import (
	"fmt"
	"log"
	"mypass/app"
	"mypass/tools"

	"github.com/spf13/cobra"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		findPassword()
	},
}

func init() {
	rootCmd.AddCommand(findCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// findCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// findCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func findPassword() {
	var searchQuery string
	fmt.Print("Enter account or URL to search: ")
	fmt.Scanln(&searchQuery)

	// 从数据库中搜索账号或URL
	rows, err := app.Sqlite.DB().Query("SELECT id, account, password, url, email, note FROM passwords WHERE account LIKE ? OR url LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	items := []PasswordItem{}
	found := false
	for rows.Next() {
		found = true
		var item PasswordItem
		if err := rows.Scan(&item.ID, &item.Account, &item.Password, &item.URL, &item.Email, &item.Note); err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
		//fmt.Printf("ID: %d, Account: %s, Password: %s, URL: %s, Email: %s, Note: %s\n", item.ID, item.Account, item.Password, item.URL, item.Email, item.Note)
	}

	tools.Output(items)
	if !found {
		fmt.Println("No matching records found.")
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
