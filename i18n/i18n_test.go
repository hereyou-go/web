package i18n

import "testing"

func TestLoad(t *testing.T) {
	dir := `F:\workspace\go\web\resources\lang`
	bundle, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bundle.Lang("zh-CN", "login.INVALID_USERNAME_OR_PASSWORD"))
	t.Log(bundle.Lang("zh", "login.INVALID_USERNAME_OR_PASSWORD"))
	t.Log(bundle.Lang("zh-TW", "login.INVALID_USERNAME_OR_PASSWORD"))
	t.Log(bundle.Lang("en", "login.INVALID_USERNAME_OR_PASSWORD"))
}
