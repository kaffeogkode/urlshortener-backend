package main

type link struct {
	ID        int64       `gorm:"primaryKey" json:"id"`
	URL       string      `json:"url"`
	Hash      string      `gorm:"unique" json:"hash"`
	CreatedAt int64       `json:"created_at"`
	Visits    []linkVisit `gorm:"foreignKey:LinkID" json:"-"`
}

type linkPost struct {
	URL string `json:"url"`
}

type linkVisit struct {
	ID        int64  `gorm:"primaryKey" json:"id"`
	LinkID    int64  `json:"-"`
	VisitTime int64  `json:"visit_time"`
	VisitorIP string `json:"visitor_ip"`
}
