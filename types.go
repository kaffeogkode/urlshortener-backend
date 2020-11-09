package main

type link struct {
	ID        int64 `gorm:"primaryKey"`
	URL       string
	Hash      string `gorm:"unique"`
	CreatedAt int64
	Visits    []linkVisit `gorm:"foreignKey:LinkID"`
}

type linkVisit struct {
	ID        int64 `gorm:"primaryKey"`
	LinkID    int64
	VisitTime int64
	VisitorIP string
}
