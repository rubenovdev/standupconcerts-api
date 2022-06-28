package repo

import (
	"comedians/src/core/usersConcerts/model"
	"comedians/src/db"
	"errors"
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
	usersDB = db.DBS.Preload("UsersLikes").Preload("UsersDislikes").Preload("Concerts").Preload("Subscriptions").Preload("FavoriteComedians").Preload("Roles").Preload("Roles.Permissions").Preload("FavoriteConcerts").Table(usersTable)
	rolesDB = db.DBS.Table(rolesTable)
}

func GetUser(id uint64) (model.User, error) {
	lazyInit()

	var user model.User

	err := usersDB.First(&user, "id = ?", id).Error

	return user, err
}

func GetUsers() ([]*model.User, error) {
	lazyInit()

	var users []*model.User

	err := usersDB.Find(&users).Error

	return users, err
}

func GetUserByEmail(email string) (model.User, error) {
	lazyInit()

	var user model.User
	var count int64

	if usersDB.Where("email = ?", email).Find(&user).Count(&count); count == 0 {
		return user, errors.New("not found")
	}

	return user, nil
}

func CreateUser(user model.User) (model.User, error) {
	lazyInit()
	log.Print("repo user", user)
	err := usersDB.Create(&user).Error

	return user, err
}

func UpdateUser(user model.User) error {
	lazyInit()

	return usersDB.Where("id = ?", user.Id).Updates(&user).Error
}

func UpdateFavoriteComedians(user model.User, favoriteComedians []*model.User) error {
	lazyInit()

	db.DBS.Model(&user).Association("FavoriteComedians").Clear()
	db.DBS.Save(&user)
	user.FavoriteComedians = favoriteComedians
	return db.DBS.Save(&user).Error
}

func UpdateFavoriteConcerts(user model.User, favoriteConcerts []*model.Concert) error {
	lazyInit()

	db.DBS.Model(&user).Association("FavoriteConcerts").Clear()
	db.DBS.Save(&user)
	user.FavoriteConcerts = favoriteConcerts
	return db.DBS.Save(&user).Error
}

func UpdateSubscripions(user model.User, subscripions []*model.Concert) error {
	lazyInit()

	db.DBS.Model(&user).Association("Subscriptions").Clear()
	db.DBS.Save(&user)
	user.Subscriptions = subscripions
	return db.DBS.Save(&user).Error
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

func UpdateUsersLikes(user model.User, usersLikes []*model.User) (model.User, error) {
	lazyInit()

	log.Print("usersLikes", usersLikes)

	err := db.DBS.Model(&user).Association("UsersLikes").Clear()

	if err != nil {
		return model.User{}, err
	}

	err = usersDB.Save(&user).Error

	if err != nil {
		return model.User{}, err
	}

	user.UsersLikes = usersLikes

	err = usersDB.Save(&user).Error
	return user, err
}

func UpdateUsersDislikes(user model.User, usersDislikes []*model.User) (model.User, error) {
	lazyInit()

	log.Print("usersDislikes", usersDislikes)

	err := db.DBS.Model(&user).Association("UsersDislikes").Clear()

	if err != nil {
		return model.User{}, err
	}

	err = usersDB.Save(&user).Error

	if err != nil {
		return model.User{}, err
	}

	user.UsersDislikes = usersDislikes

	err = usersDB.Save(&user).Error
	return user, err
}
