package tasks

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	path2 "path"
	"path/filepath"
	"strings"
)

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func ToMap(object interface{}) (*map[string]interface{}, error) {
	var dict map[string]interface{}
	serializsdObject, _ := json.Marshal(object)
	err := json.Unmarshal(serializsdObject, &dict)
	if err != nil {
		return nil, err
	}
	return &dict, nil
}

func CreateTempDirectory(prefix string) string {
	dir, err := ioutil.TempDir("/tmp", prefix+"-*")
	if err != nil {
		log.Fatal("Cannot create temp directory: ", err)
		panic(err)
	}
	return dir
}

func ExpandPath(path string) string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	if path == "~" {
		return homeDir
	}
	if strings.HasPrefix(path, "~") {
		path = path2.Join(homeDir, path[2:])
	}
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Cannot expand path of: %v. Original error: %v\n", path, err.Error())
	}
	return absolutePath
}
