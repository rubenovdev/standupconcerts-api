package repo

import (
	"comedians/src/core/usersConcerts/model"
	"comedians/src/db"

	"gorm.io/gorm"
)

const (
	concertsTable = "concerts"
)

var concertsDB *gorm.DB

func lazyInit() {
	if concertsDB == nil {
		concertsDB = db.DBS.Preload("User").Table(concertsTable)
	}
}

func CreateConcert(concert model.Concert) error {
	lazyInit()

	return concertsDB.Create(&concert).Error
}

func GetConcerts() ([]model.Concert, error) {
	lazyInit()

	var concerts []model.Concert

	err := concertsDB.Find(&concerts).Error
	return concerts, err
}

func GetConcert(id uint64) (model.Concert, error) {
	lazyInit()

	var concert model.Concert

	err := concertsDB.First(&concert, "id = ?", id).Error
	return concert, err
}

func DeleteConcert(id uint64) error {
	lazyInit()

	return db.DBS.Table(concertsTable).Delete(&model.Concert{}, id).Error
}

func UpdateConcert(concert model.Concert) error {
	lazyInit()

	return db.DBS.Table(concertsTable).Where("id = ?", concert.Id).Updates(&concert).Error
}
