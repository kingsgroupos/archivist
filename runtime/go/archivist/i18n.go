package archivist

import "github.com/kingsgroupos/misc/wlang"

var (
	ResolveI18n func(lang, key string) string
)

type I18n string

func (s I18n) I18n(lang string) string {
	return ResolveI18n(lang, string(s))
}

func (s I18n) Sprintf(lang string, a ...interface{}) (string, error) {
	return wlang.Sprintf(lang, ResolveI18n(lang, string(s)), a...)
}
