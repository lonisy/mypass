/*
Copyright © 2024 lonisy@163.com

*/
package main

import (
	"mypass/app"
	"mypass/cmd"
)

func main() {
	cmd.Execute()
	app.BackupDatabase()
}
