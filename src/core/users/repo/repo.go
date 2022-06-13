package repo

import (
	"comedians/src/core/usersConcerts/model"
	"comedians/src/db"
	"log"

	"gorm.io/gorm"
)

const (
	usersTable       = "users"
	rolesTable       = "roles"
	permissionsTable = "permissions"
)

var usersDB *gorm.DB
var rolesDB *gorm.DB

func lazyInit() {
	if usersDB == nil {
		usersDB = db.DBS.Preload("Subscriptions").Preload("FavoriteComedians").Preload("Roles").Preload("Roles.Permissions").Preload("FavoriteConcerts").Table(usersTable)
	}

	if rolesDB == nil {
		rolesDB = db.DBS.Table(rolesTable)
	}
}

func GetUser(id uint64) (model.User, error) {
	lazyInit()

	var user model.User

	err := usersDB.First(&user, "id = ?", id).Error

	return user, err
}

func GetUsers() ([]model.User, error) {
	lazyInit()

	var users []model.User

	err := usersDB.Find(&users).Error

	return users, err
}

func GetUserByEmail(email string) (model.User, error) {
	lazyInit()

	var user model.User

	if err := usersDB.Where("email = ?", email).Find(&user).Error; err != nil {
		log.Panic(err)
		return user, err
	}

	return user, nil
}

func CreateUser(user model.User) error {
	lazyInit()

	err := usersDB.Create(&user).Error

	return err
}

func UpdateUser(user model.User) error {
	lazyInit()

	log.Print(user)
	favoriteComedians := user.FavoriteComedians
	favoriteConcerts := user.FavoriteConcerts
	subscriptions := user.Subscriptions

	db.DBS.Model(&user).Association("FavoriteComedians").Clear()
	db.DBS.Model(&user).Association("FavoriteConcerts").Clear()
	db.DBS.Model(&user).Association("Subscriptions").Clear()

	user.Subscriptions = subscriptions
	user.FavoriteComedians = favoriteComedians
	user.FavoriteConcerts = favoriteConcerts

	return usersDB.Where("id = ?", user.Id).Updates(&user).Error
}

func DeleteUser(id uint64) error {
	lazyInit()

	return usersDB.Delete(&model.User{}, id).Error
}

func GetRoleByName(roleName string) (model.Role, error) {
	lazyInit()

	var role model.Role

	err := rolesDB.First(&role, "title = ?", roleName).Error
	return role, err
}
