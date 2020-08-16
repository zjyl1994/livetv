package global

import "sync"

var (
	ConfigCache sync.Map
	URLCache    sync.Map
)
