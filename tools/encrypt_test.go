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
	m, _ := EncryptString(myPasswd, "PfHRR48%sFhw8*K1")
	fmt.Println(m)
	ss, _ := DecryptString(m, "PfHRR48%sFhw8*K1")
	fmt.Println(ss)
}
