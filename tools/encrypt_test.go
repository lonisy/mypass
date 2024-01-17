package tools

import (
	"fmt"
	"testing"
)

func TestPasswd(t *testing.T) {
	var myPasswd string
	myPasswd = "7Pto1K5dA5nvEK9fi asd..123 adf123lmamfm @#)()(#)$*)!@$"
	saf, _ := AAesEncrypt([]byte(myPasswd))
	ron, _ := AAesDecrypt(saf)
	fmt.Println(string(ron))

	m, _ := EncryptString(myPasswd, AdjustTo16Characters("asd..123"))
	fmt.Println(m)
	ss, _ := DecryptString(m, AdjustTo16Characters("asd..123"))
	fmt.Println(ss)
}
