package parser

import (
	"errors"
	"os"
)

type LiveFormat struct {
	ID     string
	Format string
	URL    string
}

type LiveInfo struct {
	Title   string
	Formats []LiveFormat
}

type Parser interface {
	GetLiveStream(url string) (info LiveInfo, err error)
	ParserUpgradeable() bool
	ParserNeedUpdate() (bool, error)
	ParserUpdate() error
}

var (
	ErrURLNotSupport = errors.New("parser: url not support")
	ParserRegistry   map[string]Parser
)

func Init() {
	ParserRegistry = make(map[string]Parser)
	ParserRegistry["youtube-live"] = NewYoutubeLiveParser(os.Getenv("LIVETV_DATADIR"))
	ParserRegistry["RTHK-31/32"] = NewRTHKLiveParser()
}
