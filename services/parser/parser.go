package parser

import (
	"errors"
	"net/url"

	"github.com/zjyl1994/livetv/utils"
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
	parserRegistry   map[string]Parser
)

func Init() {
	parserRegistry = make(map[string]Parser)
	parserRegistry["www.youtube.com"] = NewYoutubeLiveParser(utils.DataDir(""))
	parserRegistry["www.rthk.hk"] = NewRTHKLiveParser()
}

func PickParser(strUrl string) (Parser, error) {
	ui, err := url.Parse(strUrl)
	if err != nil {
		return nil, err
	}
	parser, ok := parserRegistry[ui.Host]
	if ok {
		return parser, nil
	} else {
		return nil, ErrURLNotSupport
	}
}
