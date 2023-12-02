package main

import (
	"fmt"
	"test/models"
	"test/requests"
	"time"
)

// Можно изменить Username и Login и прогнать скрипт еще раз
var user = models.User{
	Username:  "Maria", // change
	Email:     "mivanova@gmail.com",
	FirstName: "Maria",
	LastName:  "Ivanova",
	Password:  "qwerty",
	Phone:     "+79999999999",
}

var userAuth = models.UserAuth{
	Login:    "Maria", // change
	Password: "qwerty",
}

func main() {
	time.Sleep(5 * time.Second)
	fmt.Println("Скрипт успешен, если в выводе 8 пунктов и ожидаемые значения совпадают с реальными.")
	requests.CreateUser(user)
	authHeader := requests.Login(userAuth)
	requests.DepositMoney(authHeader)
	requests.CreateSuccessfulOrder(authHeader)
	requests.CheckBalance(authHeader)
	orderID := requests.CreateFailedOrder(authHeader)
	requests.GetOrderByID(authHeader, orderID)
	requests.CheckBalanceNoChanged(authHeader)
}
