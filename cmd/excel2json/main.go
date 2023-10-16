package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/FerdinaKusumah/excel2json"
)

func main() {
	var (
		result []*map[string]interface{}
		err    error
		path   = "./combined_score.xlsx"
		// select sheet name
		sheetName = "combined"
		// select only selected field
		// if you want to show all headers just passing nil or empty list
		headers = []string{"Profit", "Shipping Cost", "Unit Price"}
	)
	if result, err = excel2json.GetExcelFilePath(path, sheetName, headers); err != nil {
		log.Fatalf(`unable to parse file, error: %s`, err)
	}
	for _, val := range result {
		result, _ := json.Marshal(val)
		fmt.Println(string(result))
	}
}
