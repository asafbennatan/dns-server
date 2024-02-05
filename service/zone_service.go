package service

import (
	"dnsServer/daos"
	"dnsServer/data"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ZoneService struct {
	db *gorm.DB
}

func NewZoneService(db *gorm.DB) *ZoneService {
	return &ZoneService{db: db}
}

func (zs *ZoneService) CreateZone(create daos.DNSZoneCreate) daos.DNSZone {
	zone := data.Zone{
		Base: data.Base{
			ID: uuid.NewString(),
		},
		Name: create.Name,
	}
	zs.db.Create(&zone).Commit()
	return zone.ToDNSZone()
}

func (zs *ZoneService) UpdateZone(update daos.DNSZoneUpdate) daos.DNSZone {
	zone := data.Zone{
		Base: data.Base{
			ID: update.ID,
		},
		Name: update.Name,
	}
	zs.db.Updates(&zone).Commit()
	return zone.ToDNSZone()
}

func (zs *ZoneService) DeleteZone(zoneId string) {
	zone := data.Zone{
		Base: data.Base{
			ID: zoneId,
		},
	}
	zs.db.Delete(&zone).Commit()
}

func (zs *ZoneService) GetZone(zoneId string) (*daos.DNSZone, error) {
	zone := data.Zone{
		Base: data.Base{
			ID: zoneId,
		},
	}
	res := zs.db.First(&zone)
	if res.Error != nil {
		return nil, res.Error // Return the error to the caller
	}
	dnsZone := zone.ToDNSZone() // Presumably converts data.Zone to api.DNSZone
	return &dnsZone, nil
}

func (zs *ZoneService) GetZones(nameLike string) []daos.DNSZone {
	var records []data.Zone

	zs.db.Where("name ILIKE ?", "%"+nameLike+"%").Find(&records)
	var toRet []daos.DNSZone
	for _, record := range records {
		toRet = append(toRet, record.ToDNSZone())
	}
	return toRet
}
