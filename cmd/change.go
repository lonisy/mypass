/*
Copyright © 2024 lonisy@163.com

*/
package cmd

import (
	"fmt"
	"log"
	"mypass/app"
	"time"

	"github.com/spf13/cobra"
)

// changeCmd represents the change command
var changeCmd = &cobra.Command{
	Use:   "change",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		changePassword()
	},
}

func init() {
	rootCmd.AddCommand(changeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// changeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// changeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func changePassword() {
	if ok, _ := checkFirstPasswordKey(); !ok {
		return
	}

	var id int
	fmt.Print("Enter the ID of the password to change: ")
	fmt.Scanln(&id)
	// 从数据库中搜索账号或URL
	rows, err := app.Sqlite.DB().Query("SELECT id, account, password, url, email, note FROM passwords WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		found = true
		var item PasswordItem
		if err := rows.Scan(&item.ID, &item.Account, &item.Password, &item.URL, &item.Email, &item.Note); err != nil {
			log.Fatal(err)
		}
	}

	if !found {
		fmt.Println("No matching records found.")
		return
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// 获取新的密码信息
	var newAccount, newPassword, newURL, newEmail, newNote string
	fmt.Print("Enter new account (leave blank to keep current): ")
	fmt.Scanln(&newAccount)
	fmt.Print("Enter new password (leave blank to keep current): ")
	fmt.Scanln(&newPassword)
	fmt.Print("Enter new URL (leave blank to keep current): ")
	fmt.Scanln(&newURL)
	fmt.Print("Enter new email (leave blank to keep current): ")
	fmt.Scanln(&newEmail)
	fmt.Print("Enter new note (leave blank to keep current): ")
	fmt.Scanln(&newNote)
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// 在数据库中更新记录
	statement, err := app.Sqlite.DB().Prepare("UPDATE passwords SET account = COALESCE(NULLIF(?, ''), account), password = COALESCE(NULLIF(?, ''), password), url = COALESCE(NULLIF(?, ''), url), email = COALESCE(NULLIF(?, ''), email), note = COALESCE(NULLIF(?, ''), note), updated_at = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	_, err = statement.Exec(newAccount, newPassword, newURL, newEmail, newNote, currentTime, id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Password record updated successfully!")
}
