package main

import (
	"log"
	"strconv"
)

type Hero struct {
	Name   string
	Race   bool
	Gender bool
	Level  int
}
type mountsInfo struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Database string `json:"database"`
	Data     []struct {
		ID      string `json:"entry"`
		Name    string `json:"name"`
		SpellID string `json:"spellid_2"`
	}
}
type CompanionsInfo struct {
	SqlExpr string `json:"sql_expr"`
	Data    []struct {
		ID      string `json:"entry"`
		SpellID string `json:"spellid_2"`
	}
}

func main() {
	//// Open our jsonFile
	//jsonFile, err := os.Open("WoWUtils/companions_info.json")
	//// if we os.Open returns an error then handle it
	//if err != nil {
	//	fmt.Println(err)
	//}
	//// defer the closing of our jsonFile so that we can parse it later on
	//defer jsonFile.Close()
	//byteValue, _ := ioutil.ReadAll(jsonFile)
	//
	//var mounts CompanionsInfo
	//_ = json.Unmarshal(byteValue, &mounts)
	//expr := "("
	//for _, dt := range mounts.Data {
	//	if dt.SpellID != "0" {
	//		expr += dt.SpellID + ", "
	//	}
	//}
	//expr = expr[:len(expr)-2]
	//expr += ")"
	//log.Print(expr)
	log.Print(strconv.FormatFloat(3.12345678, 'f', 3, 64))

}
