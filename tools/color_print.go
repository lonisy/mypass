package tools

import "github.com/fatih/color"

func ColorPrintDemo() {
	red := color.New(color.FgRed).PrintfFunc()
	green := color.New(color.FgGreen).PrintfFunc()
	blue := color.New(color.FgBlue).PrintfFunc()
	// 使用颜色打印文本
	red("这是红色文本\n")
	green("这是绿色文本\n")
	blue("这是蓝色文本\n")

	color.Cyan("Prints text in cyan.")
	color.Blue("Prints %s in blue.", "text")
	color.Red("We have red")
	color.Magenta("And many others ..")
	// Create a new color object
	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("Prints cyan text with an underline.")
	// Or just add them to New()
	d := color.New(color.FgCyan, color.Bold)
	d.Printf("This prints bold cyan %s\n", "too!.")
	// Mix up foreground and background colors, create new mixes!
	red2 := color.New(color.FgRed)
	boldRed := red2.Add(color.Bold)
	boldRed.Println("This will print text in bold red.")
	whiteBackground := red2.Add(color.BgWhite)
	whiteBackground.Println("Red text with white background.")
}
