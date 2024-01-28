package service

import (
	"dnsServer/api"
	"dnsServer/data"
	"errors"
	"gorm.io/gorm"
	"time"
)

type ZoneService struct {
	db *gorm.DB
}

func NewZoneService(db *gorm.DB) *ZoneService {
	return &ZoneService{db: db}
}

func (zs *ZoneService) CreateZone(create api.DNSZoneCreate) api.DNSZone {
	zone := data.Zone{
		Model:     gorm.Model{},
		Name:      create.Name,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	zs.db.Create(&zone).Commit()
	return zone.ToDNSZone()
}

func (zs *ZoneService) UpdateZone(update api.DNSZoneUpdate) api.DNSZone {
	zone := data.Zone{
		Model:     gorm.Model{},
		Name:      update.Name,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	zs.db.Updates(&zone).Commit()
	return zone.ToDNSZone()
}

func (zs *ZoneService) DeleteZone(zoneId uint) {
	zone := data.Zone{}
	zone.ID = zoneId
	zs.db.Delete(&zone).Commit()
}

func (zs *ZoneService) GetZone(zoneId uint) api.DNSZone {
	zone := data.Zone{}
	res := zs.db.First(zone, zoneId)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			// Record not found
		} else {
			// Other error
		}
	}
	return zone.ToDNSZone()
}

func (zs *ZoneService) GetZones(nameLike string) []api.DNSZone {
	var records []data.Zone

	zs.db.Where("name ILIKE ?", "%"+nameLike+"%").Find(&records)
	var toRet []api.DNSZone
	for _, record := range records {
		toRet = append(toRet, record.ToDNSZone())
	}
	return toRet
}
