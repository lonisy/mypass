/*
Copyright © 2024 lonisy@163.com

*/
package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/howeyc/gopass"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"log"
	"mypass/app"
	"mypass/tools"
	"os"
	"strings"
	"time"
)

// PasswordItem represents a password record
type PasswordItem struct {
	ID        int
	Account   string
	Password  string
	URL       string
	Email     string
	Note      string
	UpdatedAt string `json:"updated_at"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mypass",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if use_password_key {
			var ok bool
			if ok, password_key = checkFirstPasswordKeyByScanInput(); !ok {
				return
			}
		}
		listPasswords()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var password_key string
var use_password_key bool

func init() {
	initializeDatabase()
	checkAndSetFirstPasswordKey()
	rootCmd.PersistentFlags().BoolVarP(&use_password_key, "auth", "a", false, "Description of option A")
	//rootCmd.PersistentFlags().StringVarP(&password_key, "auth", "a", "used", "password_key")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func checkFirstPasswordKeyByScanInput() (bool, string) {
	tools.ColorPrinter.Warning("Enter Password Key: ")
	passwordBytes, err := gopass.GetPasswdMasked() // 显示星号
	if err != nil {
		fmt.Println("Error:", err)
		return false, ""
	}
	password_key := string(passwordBytes)
	//检查用户是否输入了密码键
	if strings.TrimSpace(password_key) == "" {
		fmt.Println("No password key entered.")
		return false, ""
	}
	//fmt.Println("\nPassword you entered is: ", password_key) // 注意：出于安全考虑，实际应用中不应打印密码
	password_key_db, _ := getPasswordKey()
	if password_key_db != strings.TrimSpace(password_key) {
		tools.ColorPrinter.Danger("Password key is not correct!")
		return false, ""
	}
	return true, password_key
}

func checkFirstPasswordKey(password_key string) (bool, string) {
	password_key_db, err := getPasswordKey()
	if err != nil {
		app.Log.Error(err)
		return false, ""
	}
	if password_key_db != strings.TrimSpace(password_key) {
		tools.ColorPrinter.Info("Password key is not correct!")
		return false, ""
	}
	return true, password_key
}

func initializeDatabase() {
	createTableIfNotExists("passwords", `CREATE TABLE passwords (
            id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
            account TEXT,
            password TEXT,
            url TEXT,
            email TEXT,
            note TEXT,
            updated_at DATETIME,
            created_at DATETIME
        );`)

	createTableIfNotExists("password_keys", `CREATE TABLE password_keys (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    key TEXT,
    created_at DATETIME,
    updated_at DATETIME
);`)
}

func createTableIfNotExists(tableName, createTableSQL string) {
	var count int
	query := fmt.Sprintf("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='%s'", tableName)
	err := app.Sqlite.DB().QueryRow(query).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		_, err = app.Sqlite.DB().Exec(createTableSQL)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Table %s created\n", tableName)
	}
	//else {
	//	fmt.Printf("Table %s already exists\n", tableName)
	//}
}

// listPasswords retrieves and displays all passwords
func listPasswords() {
	rows, err := app.Sqlite.DB().Query("SELECT id, account, password, url, email, note, updated_at FROM passwords")
	if err != nil {
		app.Log.Fatal(err)
	}
	defer rows.Close()
	items := make([]PasswordItem, 0)
	for rows.Next() {
		var item PasswordItem
		if err := rows.Scan(&item.ID, &item.Account, &item.Password, &item.URL, &item.Email, &item.Note, &item.UpdatedAt); err != nil {
			app.Log.Fatal(err)
		}
		if password_key == "" {
			item.Password = "******"
		} else {
			fmt.Println("password_key:", password_key)
			fmt.Println("item.Password:", item.Password)
			fmt.Println(tools.AdjustTo16Characters(password_key))
			item.Password, _ = tools.DecryptString(item.Password, tools.AdjustTo16Characters(password_key))
			if err != nil {
				app.Log.Error(err)
			}
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		app.Log.Fatal(err)
	}
	//tools.ColorPrinter.Info(password_key)
	if len(items) == 0 {
		tools.ColorPrinter.Info("No passwords found.")
	} else {
		tools.GreenOutput(items)
	}
}

func checkAndSetFirstPasswordKey() {
	var count int
	err := app.Sqlite.DB().QueryRow("SELECT count(*) FROM password_keys").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	// 如果 password_keys 表中没有记录
	if count == 0 {
		fmt.Println("No password key found. Please set your first password key.")
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter password key: ")
		passwordKey, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		// 去除换行符
		passwordKey = strings.TrimSpace(passwordKey)
		// 检查用户是否输入了密码键
		if strings.TrimSpace(passwordKey) == "" {
			fmt.Println("No password key entered.")
			return
		}
		passwordKey = tools.DefaultEncryptString(passwordKey)
		// 获取当前时间作为 created_at 和 updated_at
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		// 将新的 password key 插入到数据库中
		_, err = app.Sqlite.DB().Exec("INSERT INTO password_keys (key, created_at, updated_at) VALUES (?, ?, ?)", passwordKey, currentTime, currentTime)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Password key set successfully.")
	}
	//else {
	//	fmt.Println("Password key already exists.")
	//}
}

func getPasswordKey() (string, error) {
	var passwordKey string
	err := app.Sqlite.DB().QueryRow("SELECT key FROM password_keys ORDER BY id ASC LIMIT 1").Scan(&passwordKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no password key found")
		}
		return "", err
	}
	return tools.DefaultDecryptString(passwordKey), nil
}
