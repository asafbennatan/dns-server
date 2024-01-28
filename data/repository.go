package data

import (
	"dnsServer/api"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

type Zone struct {
	gorm.Model
	Name      string   `gorm:"unique"`
	Records   []Record // One-to-many relationship
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Record struct {
	gorm.Model
	Name      string
	Type      string
	Value     string
	TTL       int
	ZoneID    uint // Foreign key for Zone
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (zs *Zone) ToDNSZone() api.DNSZone {
	return api.DNSZone{
		ID:   zs.ID,
		Name: zs.Name,
	}
}

func (zs *Record) ToDNSRecord() api.DNSRecord {
	return api.DNSRecord{
		ID:    zs.ID,
		Name:  zs.Name,
		Type:  zs.Type,
		Value: zs.Value,
		TTL:   zs.TTL,
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
