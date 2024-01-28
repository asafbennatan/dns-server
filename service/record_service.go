package service

import (
	"dnsServer/api"
	"dnsServer/data"
	"errors"
	"gorm.io/gorm"
	"time"
)

type RecordService struct {
	db *gorm.DB
}

func NewRecordService(db *gorm.DB) *RecordService {
	return &RecordService{db: db}
}

func (zs *RecordService) CreateRecord(zoneId uint, create api.DNSRecordCreate) api.DNSRecord {

	record := data.Record{
		Model:     gorm.Model{},
		Name:      create.Name,
		Type:      create.Type,
		Value:     create.Value,
		TTL:       create.TTL,
		ZoneID:    zoneId,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	zs.db.Create(&record).Commit()
	return record.ToDNSRecord()
}

func (zs *RecordService) UpdateRecord(update api.DNSRecordUpdate) api.DNSRecord {

	record := data.Record{
		Model:     gorm.Model{},
		Name:      update.Name,
		Type:      update.Type,
		Value:     update.Value,
		TTL:       update.TTL,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	zs.db.Updates(&record).Commit()
	return record.ToDNSRecord()

}

func (zs *RecordService) DeleteRecord(recordId uint) {
	record := data.Record{}
	record.ID = recordId
	zs.db.Delete(&record).Commit()
}

func (zs *RecordService) GetRecord(recordId uint) api.DNSRecord {
	record := data.Record{}
	res := zs.db.First(record, recordId)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			// Record not found
		} else {
			// Other error
		}
	}
	return record.ToDNSRecord()
}

func (zs *RecordService) GetRecords(zoneId uint, nameLike string) []api.DNSRecord {
	var records []data.Record

	zs.db.Where("name ILIKE ? and zone_id = ", "%"+nameLike+"%", zoneId).Find(&records)
	var toRet []api.DNSRecord
	for _, record := range records {
		toRet = append(toRet, record.ToDNSRecord())
	}
	return toRet
}
