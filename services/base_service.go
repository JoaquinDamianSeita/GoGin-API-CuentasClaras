package services

import (
	dao "GoGin-API-CuentasClaras/dao"
	"GoGin-API-CuentasClaras/repository"
	"strconv"
	"time"
)

var utcLocation, _ = time.LoadLocation("UTC")

func RetrieveCurrentUser(userRepository repository.UserRepository, userID string) (user dao.User, recordError error) {
	intUserID, _ := strconv.Atoi(userID)
	user, recordError = userRepository.FindUserById(intUserID)
	return
}
