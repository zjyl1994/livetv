package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/zjyl1994/livetv/services/parser"
)

func main() {
	fmt.Println("Under Development")
	p := parser.NewYoutubeLiveParser(os.Getenv("LIVETV_DATADIR"))
	fmt.Println("Can Update?", p.ParserUpgradeable())
	needUpdate, err := p.ParserNeedUpdate()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Need Update?", needUpdate)
	if needUpdate {
		fmt.Println("Update!", p.ParserUpdate())
	}
	info, err := p.GetLiveStream("https://www.youtube.com/watch?v=63RmMXCd_bQ")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Title", info.Title)
	for _, v := range info.Formats {
		fmt.Println(v.ID, "=>", v.Format, "=>", v.URL)
	}
}
