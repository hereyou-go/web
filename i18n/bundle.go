package i18n

import (
	"fmt"
	"strings"
)

// type MinorBundle struct {
// 	local string
// 	langs map[string]string
// }

// func (mb *MinorBundle) Local() string {
// 	return mb.local
// }

type BundleEntry map[string]string

type Bundle map[string]BundleEntry

func (b Bundle) Lang(local, name string, args ...interface{}) (lang string) {

	if !localReg.MatchString(local) {
		return
	}
	// var ok bool
	local = strings.ToUpper(local)
	name = strings.ToUpper(name)
	major := local
	entry, _ := b[local]
	if entry == nil {
		index := strings.Index(local, "-")
		if index > 0 {
			major = local[:index]
			minor := local[index+1:]

			entry, _ = b[major]
			if entry == nil {
				entry, _ = b[minor]
			}
		}
		if entry == nil {
			for loc, entryLike := range b {
				if strings.HasPrefix(loc, major) {
					entry = entryLike
					break
				}
			}
		}
	}

	var tag string
	if entry != nil {
		tag, _ = entry[name]
	}
	if tag == "" {
		for loc, entry := range b {
			if strings.HasPrefix(loc, major) {
				continue
			}
			tag, _ = entry[name]
			if tag != "" {
				break
			}
		}
	}
	// logs.Debug("tag:%s,name:%s  %s  %s    %+v", tag, name, local, major, entry)
	return fmt.Sprintf(tag, args...)
}
