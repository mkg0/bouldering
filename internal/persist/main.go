package persist

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var path string

func Save(data interface{}) {
	file, _ := os.Create(path + "temp_bouldering.tmp")
	defer file.Close()
	jsonData, _ := json.MarshalIndent(data, "", " ")
	file.Write(jsonData)
}

func Load(data interface{}) {
	exec, err := os.Executable()
	execPath := filepath.Dir(exec)
	if err != nil {
		panic(err)
	}
	file, _ := os.ReadFile(execPath + "/temp_bouldering.tmp")
	if len(file) != 0 {
		path = execPath + "/"
	} else {
		fileAtCur, _ := os.ReadFile("temp_bouldering.tmp")
		if len(fileAtCur) != 0 {
			file = fileAtCur
			path = ""
		} else {
			path = execPath + "/"
		}
	}
	json.Unmarshal(file, data)
}

func GetFilePath() string {
	return path + "temp_bouldering.tmp"
}
