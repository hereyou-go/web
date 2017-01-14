package i18n

// Locale 指定为每个区域，根据 RFC 4646 的唯一名称。
// 该名称是与语言相关的 ISO 639 两个字母小写区域性代码和与国家或地区相关 ISO 3166 双字母大写子区域性代码的组合。
type Locale struct {
}

// TODO:
var (
	English = &Locale{}
	Chinese = &Locale{}
)
