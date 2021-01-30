package main

import (
	"fmt"

	_ "github.com/joho/godotenv/autoload"
	"github.com/zjyl1994/livetv/services/parser"
)

func main() {
	fmt.Println("Under Development")
	parser.Init()
	p := parser.ParserRegistry["RTHK-31/32"]
	fmt.Println("Can Update?", p.ParserUpgradeable())
	needUpdate, err := p.ParserNeedUpdate()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Need Update?", needUpdate)
	if needUpdate {
		fmt.Println("Update!", p.ParserUpdate())
	}
	info, err := p.GetLiveStream("https://www.rthk.hk/feeds/dtt/rthktv31_https.m3u8")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Title", info.Title)
	for _, v := range info.Formats {
		fmt.Println(v.ID, "=>", v.Format, "=>", v.URL)
	}
}
