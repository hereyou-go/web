package web

type I18N interface {
	Lang(local, name string, args ...interface{}) (lang string)
}
