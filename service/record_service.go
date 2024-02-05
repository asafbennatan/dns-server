package service

import (
	"dnsServer/daos"
	"dnsServer/data"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RecordService struct {
	db *gorm.DB
}

func NewRecordService(db *gorm.DB) *RecordService {
	return &RecordService{db: db}
}

func (zs *RecordService) CreateRecord(zoneId string, create daos.DNSRecordCreate) daos.DNSRecord {

	record := data.Record{
		Base: data.Base{
			ID: uuid.NewString(),
		},
		Name:   create.Name,
		Type:   create.Type,
		Value:  create.Value,
		TTL:    create.TTL,
		ZoneID: zoneId,
	}
	zs.db.Create(&record).Commit()
	return record.ToDNSRecord()
}

func (zs *RecordService) UpdateRecord(update daos.DNSRecordUpdate) daos.DNSRecord {

	record := data.Record{
		Base: data.Base{
			ID: update.ID,
		},
		Name:  update.Name,
		Type:  update.Type,
		Value: update.Value,
		TTL:   update.TTL,
	}
	zs.db.Updates(&record).Commit()
	return record.ToDNSRecord()

}

func (zs *RecordService) DeleteRecord(recordId string) {
	record := data.Record{
		Base: data.Base{
			ID: recordId,
		},
	}
	record.ID = recordId
	zs.db.Delete(&record).Commit()
}

func (zs *RecordService) GetRecord(recordId string) (*daos.DNSRecord, error) {
	var record data.Record
	res := zs.db.Where("id = ?", recordId).First(&record) // Corrected line
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			// Optionally handle not found error specifically
			return nil, res.Error
		}
		// Handle other possible errors
		return nil, res.Error
	}
	dnsRecord := record.ToDNSRecord() // Presumably converts data.Record to daos.DNSRecord
	return &dnsRecord, nil
}

func (zs *RecordService) GetRecords(zoneId string, nameLike string) []daos.DNSRecord {
	var records []data.Record

	zs.db.Where("name ILIKE ? and zone_id = ?", "%"+nameLike+"%", zoneId).Find(&records)
	var toRet []daos.DNSRecord
	for _, record := range records {
		toRet = append(toRet, record.ToDNSRecord())
	}
	return toRet
}
