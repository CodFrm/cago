package i18n

var DefaultLang = "zh-cn"

var langs = map[string]map[int]string{}

func Register(lang string, code map[int]string) {
	if _, ok := langs[lang]; ok {
		// append
		for k, v := range code {
			langs[lang][k] = v
		}
	} else {
		langs[lang] = code
	}
}
