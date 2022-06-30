package service

import (
	commonService "comedians/src/common/service"
	authModel "comedians/src/core/auth/model"
	"comedians/src/core/users/repo"
	"comedians/src/core/users/service"
	usersModel "comedians/src/core/usersConcerts/model"
	"comedians/src/utils"
	"errors"
	"fmt"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"log"
	"os"
	"time"
)

const (
	tokenTTL = 7 * 24 * time.Hour
)

func CreateUser(user usersModel.User) (usersModel.User, error) {
	_, err := GetUserByEmail(user.Email)

	if err == nil {
		return usersModel.User{}, errors.New("email already taken")
	}

	user.Password = utils.HashPassword(user.Password)
	roleUser, err := repo.GetRoleByName("user")

	if err != nil {
		return usersModel.User{}, err
	}

	user.Roles = append(user.Roles, roleUser)
	log.Print("user", user, user.Email)

	user, err = repo.CreateUser(user)
	log.Print("user id !!", user.Id)
	return user, err
}

func AuthGoogle(dto authModel.AuthGoogleDto) (string, error) {
	if user, err := repo.GetUserByGoogleId(dto.Id); err == nil {
		signedToken, err := GenerateTokenJWT(user.Id, user.Roles)

		if err != nil {
			return "", err
		}

		return signedToken, nil
	}

	pwd, err := utils.GeneratePassword()

	if err != nil {
		return "", err
	}

	dir := os.Getenv("USERS_IMAGES_DIR")
	commonService.MakeDirIfNotExists(dir)

	filepath := dir + "/" + fmt.Sprint(time.Now().Unix()) + ".jpg"

	err = commonService.DownloadFile(dto.ImgUrl, filepath)

	if err != nil {
		return "", err
	}

	user, err := CreateUser(usersModel.User{Email: dto.Email, Password: pwd, Name: dto.Name, ImgUrl: &filepath, GoogleId: dto.Id})

	if err != nil {
		commonService.DeleteFile(filepath)
		return "", err
	}

	signedToken, err := GenerateTokenJWT(user.Id, user.Roles)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func AuthVk(dto authModel.AuthVkDto) (string, error) {
	if user, err := repo.GetUserByVkId(dto.Id); err == nil {
		signedToken, err := GenerateTokenJWT(user.Id, user.Roles)

		if err != nil {
			return "", err
		}

		return signedToken, nil
	}

	pwd, err := utils.GeneratePassword()

	if err != nil {
		return "", err
	}

	dir := os.Getenv("USERS_IMAGES_DIR")
	commonService.MakeDirIfNotExists(dir)

	filepath := dir + "/" + fmt.Sprint(time.Now().Unix()) + ".jpg"

	err = commonService.DownloadFile(dto.ImgUrl, filepath)

	if err != nil {
		return "", err
	}

	user, err := CreateUser(usersModel.User{Email: dto.Email, Password: pwd, Name: dto.Name, ImgUrl: &filepath, VkId: dto.Id})

	if err != nil {
		commonService.DeleteFile(filepath)
		return "", err
	}

	signedToken, err := GenerateTokenJWT(user.Id, user.Roles)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func AuthYandex(dto authModel.AuthYandexDto) (string, error) {
	if user, err := repo.GetUserByYandexId(dto.Id); err == nil {
		signedToken, err := GenerateTokenJWT(user.Id, user.Roles)

		if err != nil {
			return "", err
		}

		return signedToken, nil
	}

	pwd, err := utils.GeneratePassword()

	if err != nil {
		return "", err
	}

	dir := os.Getenv("USERS_IMAGES_DIR")
	commonService.MakeDirIfNotExists(dir)

	filepath := dir + "/" + fmt.Sprint(time.Now().Unix()) + ".jpg"

	err = commonService.DownloadFile(dto.ImgUrl, filepath)

	if err != nil {
		return "", err
	}

	user, err := CreateUser(usersModel.User{Email: dto.Email, Password: pwd, Name: dto.Name, ImgUrl: &filepath, YandexId: dto.Id})

	if err != nil {
		commonService.DeleteFile(filepath)
		return "", err
	}

	signedToken, err := GenerateTokenJWT(user.Id, user.Roles)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GetUserByEmail(email string) (usersModel.User, error) {
	user, err := repo.GetUserByEmail(email)

	if err != nil {
		return user, err
	}

	return user, nil
}

func RecoveryUserPassword(email string) (string, error) {
	newPassword, _ := utils.GeneratePassword()
	user, err := GetUserByEmail(email)

	if err != nil {
		return "", err
	}

	err = service.UpdateUserPassword(user.Id, newPassword)

	if err != nil {
		return "", err
	}

	var messageData commonService.MessageData

	messageData.To = "e-shvedov@list.ru"
	messageData.Subject = "Password recovery"
	messageData.Body = fmt.Sprintf("Ваш новый пароль: %s", newPassword)

	err = commonService.SendMail(messageData)

	if err != nil {
		log.Panic(err)
		return "", err
	}
	return newPassword, nil
}

func GenerateTokenJWT(userId uint64, roles []usersModel.Role) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &authModel.Token{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Id:    userId,
		Roles: roles,
	})

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))

	return signedToken, err
}

func ParseTokenJWT(accessToken string) (*authModel.Token, error) {
	token, err := jwt.ParseWithClaims(accessToken, &authModel.Token{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		return &authModel.Token{}, err
	}

	claims, ok := token.Claims.(*authModel.Token)

	if !ok {
		return &authModel.Token{}, err
	}

	return claims, nil
}
