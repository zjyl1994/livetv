package database

type Channel struct {
	ID    uint `gorm:"primary_key"`
	Name  string
	URL   string
	Proxy bool
}

type Config struct {
	Name string `gorm:"primary_key"`
	Data string
}

// Config name："password" ，"base_url"
