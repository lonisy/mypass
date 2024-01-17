/*
Copyright © 2024 lonisy@163.com

*/
package cmd

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"log"
	"mypass/app"
	"mypass/tools"
	"os"
	"strings"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		deletePassword()
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func deletePassword() {
	if ok, _ := checkFirstPasswordKeyByScanInput(); !ok {
		return
	}
	var id int
	fmt.Print("Enter the ID of the password to delete: ")
	fmt.Scanln(&id)
	// 从数据库中搜索账号或URL
	rows, err := app.Sqlite.DB().Query("SELECT id, account, password, url, email, note FROM passwords WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	items := make([]PasswordItem, 0)
	found := false
	for rows.Next() {
		found = true
		var item PasswordItem
		if err := rows.Scan(&item.ID, &item.Account, &item.Password, &item.URL, &item.Email, &item.Note); err != nil {
			log.Fatal(err)
		}
		item.Password = "******"
		items = append(items, item)
		fmt.Printf("ID: %d, Account: %s, Password: %s, URL: %s, Email: %s, Note: %s\n", item.ID, item.Account, item.Password, item.URL, item.Email, item.Note)
	}
	if !found {
		tools.ColorPrinter.Warning("No matching records found.")
		app.Log.Error("No matching records found.")
		return
	}
	tools.Output(items)

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(os.Stdin)
	redPrint := color.New(color.FgRed).PrintfFunc()
	redPrint("Are you sure you want to delete this password? [y/n]: \n")
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	response = strings.TrimSpace(response)
	if strings.ToLower(response) == "y" {
		//执行删除操作
		statement, err := app.Sqlite.DB().Prepare("DELETE FROM passwords WHERE id = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer statement.Close()

		_, err = statement.Exec(id)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Password deleted successfully!")
	} else {
		fmt.Println("Delete operation cancelled.")
	}
}
