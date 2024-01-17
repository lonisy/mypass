/*
Copyright © 2024 lonisy@163.com

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"mypass/app"
	"mypass/tools"
	"time"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		addPassword()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func addPassword() {
	//var password_key, account, password, url, email, note string
	var account, password, url, email, note string
	//fmt.Print("Enter Password Key: ")
	//fmt.Scanln(&password_key)
	// 获取用户输入
	fmt.Print("Enter account: ")
	fmt.Scanln(&account)
	fmt.Print("Enter password: ")
	fmt.Scanln(&password)
	fmt.Print("Enter URL: ")
	fmt.Scanln(&url)
	fmt.Print("Enter email: ")
	fmt.Scanln(&email)
	fmt.Print("Enter note: ")
	fmt.Scanln(&note)
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	// 加密密码
	password_key, _ = getPasswordKey()
	password, _ = tools.EncryptString(password, tools.AdjustTo16Characters(password_key))
	// 插入到数据库
	statement, err := app.Sqlite.DB().Prepare("INSERT INTO passwords (account, password, url, email, note,updated_at,created_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	_, err = statement.Exec(account, password, url, email, note, currentTime, currentTime)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Password added successfully!")
}
