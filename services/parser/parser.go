package parser

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
