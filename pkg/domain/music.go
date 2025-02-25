package domain

type Songs struct {
	ID          uint64 `gorm:"primaryKey" json:"id"`
	SongName    string `gorm:"not null" json:"song"`
	GroupName   string `gorm:"not null" json:"group"`
	ReleaseDate string `json:"releaseDate"`
	Lyrics      string `json:"lyrics"`
	Link        string `json:"link"`
}
