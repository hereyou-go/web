package i18n

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/go-ini/ini"
)

// type Localization interface {
// 	Local() string
// 	Lang(local string, name string, args ...interface{}) string
// }

// type Bundles map[string]*Bundle

// func (b Bundles) Lang(local string, name string, args ...interface{}) string {
// 	local = strings.ToUpper(local)

// 	bundle, _ := b[local]
// 	if bundle != nil {
// 		index := strings.Index(local, "-")
// 		if index < 0 {
// 			index = strings.Index(local, "_")
// 		}
// 		if index > 0 {
// 			bundle, _ = b[local[:index]]
// 			if bundle == nil {
// 				bundle, _ = b[local[index+1:]]
// 			}
// 		}
// 	}

// 	if bundle == nil {
// 		return ""
// 	}

// 	return bundle.Lang(local, name, args...)
// }

var reg, _ = regexp.Compile(`\w+`)
var localReg, _ = regexp.Compile(`\w{1,4}(\-\w{1,4})?`)

func readFile(file os.FileInfo, path string, section string, entry BundleEntry) error {
	if file.IsDir() {
		files, err := ioutil.ReadDir(path + "/" + file.Name())
		if err != nil {
			return err
		}
		for _, child := range files {
			sec := file.Name()
			if section != "" {
				sec = section + "." + sec
			}
			readFile(child, path, sec, entry)
		}
		return err
	}

	name := strings.ToUpper(file.Name())
	if strings.HasSuffix(name, ".PROPERTIES") {
		name = name[:len(name)-11]
	}

	if !reg.MatchString(name) {
		return nil
	}

	if section != "" {
		name = section + "." + name
	}

	cfg, err := ini.Load(path + "/" + file.Name())
	if err != nil {
		return err
	}
	for _, key := range cfg.Section("").Keys() {
		entry[strings.ToUpper(name+"."+key.Name())] = key.Value()
	}

	return nil

}

func readLang(path string, local string) (entry BundleEntry, err error) {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	entry = make(BundleEntry)

	for _, file := range files {
		readFile(file, path, "", entry)
	}
	return
}

func Load(dir string) (bundle Bundle, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	bundle = make(Bundle)
	for _, file := range files {
		if file.IsDir() {
			local := strings.ToUpper(file.Name())
			if !localReg.MatchString(local) {
				continue
			}

			entry, err := readLang(dir+"/"+file.Name(), local)
			if err != nil {
				return nil, err
			}
			bundle[local] = entry
		}
	}
	return
}
