package persist

import (
	"encoding/json"
	"os"
)

func Save(data interface{}) {
	file, _ := os.Create("data.tmp")
	defer file.Close()
	jsonData, _ := json.MarshalIndent(data, "", " ")
	file.Write(jsonData)
}

func Load(data interface{}) {
	file, _ := os.ReadFile("data.tmp")
	json.Unmarshal(file, data)
}
