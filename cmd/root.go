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
	UpdatedAt string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mypass",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
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

func init() {
	initializeDatabase()
	checkAndSetFirstPasswordKey()
	fmt.Println(getPasswordKey())

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mypass.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func checkFirstPasswordKey() (bool, string) {
	//fmt.Print("Enter Password Key: ")
	//passwordBytes, err := term.ReadPassword(0) // 0 通常代表标准输入
	//if err != nil {
	//	log.Fatal(err)
	//}
	//password_key := string(passwordBytes)

	fmt.Printf("Enter Password Key: ")
	passwordBytes, err := gopass.GetPasswdMasked() // 显示星号
	if err != nil {
		fmt.Println("Error:", err)
		return false, ""
	}
	password_key := string(passwordBytes)
	//fmt.Println("\nPassword you entered is: ", password_key) // 注意：出于安全考虑，实际应用中不应打印密码
	password_key_db, _ := getPasswordKey()
	if tools.DefaultDecryptString(password_key_db) != strings.TrimSpace(password_key) {
		log.Fatal("Password key is not correct!")
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
	} else {
		fmt.Printf("Table %s already exists\n", tableName)
	}
}

// listPasswords retrieves and displays all passwords
func listPasswords() {
	rows, err := app.Sqlite.DB().Query("SELECT id, account, password, url, email, note, updated_at FROM passwords")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	items := make([]PasswordItem, 0)
	for rows.Next() {
		var item PasswordItem
		if err := rows.Scan(&item.ID, &item.Account, &item.Password, &item.URL, &item.Email, &item.Note, &item.UpdatedAt); err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
		//fmt.Printf("%d: %s - %s -%s - %s - %s\n", item.ID, item.Account, item.Password, item.URL, item.Email, item.Note)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	tools.Output(items)
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
		passwordKey = tools.DefaultEncryptString(passwordKey)
		// 获取当前时间作为 created_at 和 updated_at
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		// 将新的 password key 插入到数据库中
		_, err = app.Sqlite.DB().Exec("INSERT INTO password_keys (key, created_at, updated_at) VALUES (?, ?, ?)", passwordKey, currentTime, currentTime)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Password key set successfully.")
	} else {
		fmt.Println("Password key already exists.")
	}
}

func getPasswordKey() (string, error) {
	var passwordKey string
	err := app.Sqlite.DB().QueryRow("SELECT key FROM password_keys ORDER BY id ASC LIMIT 1").Scan(&passwordKey)
	if err != nil {
		if err == sql.ErrNoRows {
			// 表中没有记录
			return "", fmt.Errorf("no password key found")
		}
		// 处理其他错误
		return "", err
	}
	return passwordKey, nil
}
