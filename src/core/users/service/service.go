package service

import (
	concertsRepo "comedians/src/core/concerts/repo"
	"comedians/src/core/users/repo"
	"comedians/src/core/usersConcerts/model"
	"comedians/src/utils"
	"log"
)

func UpdateUser(id uint64, user model.User) error {
	user.Password = utils.HashPassword(user.Password)
	user.Id = id

	return repo.UpdateUser(user)
}

func DeleteUser(id uint64) error {
	return repo.DeleteUser(id)
}

func GetUser(id uint64) (model.User, error) {
	return repo.GetUser(id)
}

func GetUsers() ([]model.User, error) {
	return repo.GetUsers()
}

func AppendFavoriteConcert(userId uint64, concertId uint64) error {
	user, err := GetUser(userId)

	if err != nil {
		return err
	}

	concert, err := concertsRepo.GetConcert(concertId)

	if err != nil {
		return err
	}

	user.FavoriteConcerts = append(user.FavoriteConcerts, &concert)

	return repo.UpdateUser(user)
}

func AppendFavoriteComedian(userId uint64, comedianId uint64) error {
	user, err := GetUser(userId)

	if err != nil {
		return err
	}

	comedian, err := GetUser(comedianId)

	if err != nil {
		return err
	}
	
	user.FavoriteComedians = append(user.FavoriteComedians, &comedian)

	return repo.UpdateUser(user)
}

func AppendSubscription(userId uint64, concertId uint64) error {
	user, err := GetUser(userId)

	if err != nil {
		return err
	}

	concert, err := concertsRepo.GetConcert(concertId)

	if err != nil {
		return err
	}
	
	user.Subscriptions = append(user.Subscriptions, &concert)

	return repo.UpdateUser(user)
}

func DeleteFavoriteConcert(userId uint64, concertId uint64) error {
	user, err := GetUser(userId)

	if err != nil {
		return err
	}

	concert, err := concertsRepo.GetConcert(concertId)

	if err != nil {
		return err
	}

	user.FavoriteConcerts = utils.Filter(user.FavoriteConcerts, func(elem *model.Concert) bool {
		log.Print(elem.Id, concert.Id, elem.Id != concert.Id)

		return elem.Id != concert.Id
	})

	return repo.UpdateUser(user)
}

func DeleteFavoriteComedian(userId uint64, comedianId uint64) error {
	user, err := GetUser(userId)

	if err != nil {
		return err
	}

	comedian, err := GetUser(comedianId)

	if err != nil {
		return err
	}
	
	user.FavoriteComedians = utils.Filter(user.FavoriteComedians, func(elem *model.User) bool {
		return elem.Id != comedian.Id
	})

	return repo.UpdateUser(user)
}

func DeleteSubscription(userId uint64, concertId uint64) error {
	user, err := GetUser(userId)

	if err != nil {
		return err
	}

	concert, err := concertsRepo.GetConcert(concertId)

	if err != nil {
		return err
	}
	
	user.Subscriptions = utils.Filter(user.Subscriptions, func(elem *model.Concert) bool {
		return elem.Id != concert.Id
	})

	return repo.UpdateUser(user)
}