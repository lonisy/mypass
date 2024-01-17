package translations

import (
	"os"
	"strings"
)

var LocalLanguage string

var Translations = map[string]map[string]string{
	"en_US": {
		"hello":              "Hello",
		"no_passwords_found": "No passwords found.",
	},
	"fr": {
		"hello":   "Bonjour",
		"goodbye": "Au revoir",
	},
	"zh": {
		"hello":              "你好",
		"goodbye":            "再见",
		"no_passwords_found": "没有找到任何密码哦",
	},
}

func init() {
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = os.Getenv("LC_ALL")
	}
	if lang == "" {
		lang = "en_US.UTF-8" // 或你想要的默认语言
	}
	LocalLanguage = strings.Split(lang, ".")[0]
}

func GetLocalTranslation(key string) string {
	return GetTranslation(LocalLanguage, key)
}

func GetTranslation(currentLanguage, key string) string {
	translationsForCurrentLanguage, ok := Translations[currentLanguage]
	if !ok {
		return ""
	}
	translation, ok := translationsForCurrentLanguage[key]
	if !ok {
		return ""
	}
	return translation
}
