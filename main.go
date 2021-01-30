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
	info, err := p.GetLiveStream("https://www.youtube.com/watch?v=63RmMXCd_bQ")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Title", info.Title)
	for _, v := range info.Formats {
		fmt.Println(v.ID, "=>", v.Format, "=>", v.URL)
	}
}
