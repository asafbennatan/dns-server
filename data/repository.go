package data

import (
	"dnsServer/daos"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

type Base struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Zone struct {
	Base
	Name    string `gorm:"unique"`
	Records []Record
}

type Record struct {
	Base
	Name   string
	Type   string
	Value  string
	TTL    int
	ZoneID string
}

func (zs *Zone) ToDNSZone() daos.DNSZone {
	return daos.DNSZone{
		ID:   zs.ID,
		Name: zs.Name,
	}
}

func (zs *Record) ToDNSRecord() daos.DNSRecord {
	return daos.DNSRecord{
		ID:        zs.ID,
		Name:      zs.Name,
		Type:      zs.Type,
		Value:     zs.Value,
		TTL:       zs.TTL,
		DNSZoneID: zs.ZoneID,
	}
}

func InitDB() *gorm.DB {
	dsn := "user=dns password=dns dbname=dns"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// AutoMigrate the Zone and Record structs
	err = db.AutoMigrate(&Zone{}, &Record{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}
