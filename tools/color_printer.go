package tools

import "fmt"

type ColorPrinterStruct struct {
}

var ColorPrinter ColorPrinterStruct

func (pr *ColorPrinterStruct) Primary(a ...interface{}) {
	fmt.Print(a...)
}

func (pr *ColorPrinterStruct) Success(a ...interface{}) {
	fmt.Print(a...)
}

func (pr *ColorPrinterStruct) Info(a ...interface{}) {
	fmt.Print(a...)
}

func (pr *ColorPrinterStruct) Warning(a ...interface{}) {
	fmt.Print(a...)
}

func (pr *ColorPrinterStruct) Danger(a ...interface{}) {
	fmt.Print(a...)
}
