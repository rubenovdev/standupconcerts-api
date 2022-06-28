package repo

import (
	concertsModel "comedians/src/core/concerts/model"
	usersConcertsModel "comedians/src/core/usersConcerts/model"
	"comedians/src/db"
	"fmt"
	"log"

	// "log"

	"gorm.io/gorm"
)

const (
	concertsTable = "concerts"
)

var concertsDB *gorm.DB

func lazyInit() {
	concertsDB = db.DBS.Preload("User").Preload("UsersLikes").Table(concertsTable)
}

func CreateConcert(concert usersConcertsModel.Concert) error {
	lazyInit()

	return concertsDB.Create(&concert).Error
}

func GetConcerts(filters concertsModel.Filters) ([]usersConcertsModel.Concert, error) {
	lazyInit()

	var concerts []usersConcertsModel.Concert

	query := concertsDB.Order(filters.SortBy + " DESC")

	if (filters.ComedianId != 0) {
		query = query.Where("user_id  = ?", filters.ComedianId)
	}

	if filters.Year != 0 {
		yearAt := fmt.Sprint(filters.Year) + "-01-01"
		yearTo := fmt.Sprint(filters.Year) + "-12-31"

		query = query.Where("created_at >= ? AND created_at <= ?", yearAt, yearTo)
	}

	if filters.Limit != 0 {
		query.Limit(int(filters.Limit))
	}

	if filters.ExcludedId != 0 {
		query = query.Where("id != ?", filters.ExcludedId)
	}


	err := query.Find(&concerts).Error
	return concerts, err
}

func GetConcert(id uint64) (usersConcertsModel.Concert, error) {
	lazyInit()

	var concert usersConcertsModel.Concert

	err := concertsDB.First(&concert, "id = ?", id).Error
	return concert, err
}

func DeleteConcert(id uint64) error {
	lazyInit()

	return concertsDB.Delete(&usersConcertsModel.Concert{}, id).Error
}

func UpdateConcert(concert usersConcertsModel.Concert) error {
	lazyInit()

	return concertsDB.Where("id = ?", concert.Id).Updates(&concert).Error
}

func UpdateUsersLikes(concert usersConcertsModel.Concert, usersLikes []*usersConcertsModel.User) (usersConcertsModel.Concert, error) {
	lazyInit()

	log.Print("usersLikes", usersLikes)

	err := db.DBS.Model(&concert).Association("UsersLikes").Clear()

	if err != nil {
		return usersConcertsModel.Concert{}, err
	}

	err = concertsDB.Save(&concert).Error

	if err != nil {
		return usersConcertsModel.Concert{}, err
	}

	concert.UsersLikes = usersLikes

	err = concertsDB.Save(&concert).Error
	return concert, err
}
