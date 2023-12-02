package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"test/models"

	"github.com/google/uuid"
)

func CreateUser(in models.User) {
	url := "http://user-service:3000/user/"
	jsonData, err := json.Marshal(in)
	if err != nil {
		fmt.Println("marshal error: ", err)
		return
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("request error: ", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("response error: ", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("1. Пользователь успешно создан")
	}
}

func Login(in models.UserAuth) string {
	url := "http://user-service:3000/login/"
	jsonData, err := json.Marshal(in)
	if err != nil {
		return ""
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return ""
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return ""
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("2. Авторизация прошла успешно")
		return response.Header.Get("Authorization")
	}
	return ""
}

func DepositMoney(authHeader string) {
	url := "http://billing-service:8080/balance/"
	data := models.UpdateBalance{
		Amount: 1000,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", authHeader)
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("3. Деньги успешно зачислены (1000 рублей)")
	}
}

func CreateSuccessfulOrder(authHeader string) {
	url := "http://order-service:9000/createOrder/"
	data := models.Order{
		ProductID: 200,
		Quantity:  2,
		Price:     470,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", authHeader)
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("4. Создали заказ, на который хватило денег (Стоимость 940 рублей)")
		return
	}
}

func CheckBalance(authHeader string) {
	url := "http://billing-service:8080/balance/"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	request.Header.Set("Authorization", authHeader)
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var balance int64
		json.NewDecoder(response.Body).Decode(&balance)
		fmt.Println("5. Баланс успешно получен (ожидается 60): ", balance)
	}
}

func CreateFailedOrder(authHeader string) uuid.UUID {
	url := "http://order-service:9000/createOrder/"
	data := models.Order{
		ProductID: 200,
		Quantity:  1,
		Price:     150,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return uuid.Nil
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return uuid.Nil
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", authHeader)
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return uuid.Nil
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var data models.CreateOrderResponse
		decoder := json.NewDecoder(response.Body)
		decoder.Decode(&data)
		if data.Success == "order creation failed: not enough funds to write off" {
			fmt.Println("6. Создали заказ, на который не хватило денег (Стоимость 150 рублей)")
		}
		return data.OrderID
	}
	return uuid.Nil
}

func GetOrderByID(authHeader string, orderID uuid.UUID) (models.OrderInfo, error) {
	url := fmt.Sprintf("http://order-service:9000/getOrderByID/%s", orderID.String())
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.OrderInfo{}, err
	}
	request.Header.Set("Authorization", authHeader)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return models.OrderInfo{}, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var orderInfo models.OrderInfo
		json.NewDecoder(response.Body).Decode(&orderInfo)
		fmt.Printf("7. Проверить, что заказ сохранился как неуспешный (status = failed):\n \tUsername: %s\n \tOrderID: %s\n \tProductID: %d\n \tQuantity: %d\n \tPrice: %d\n \tTotalCost: %d\n \tStatus: %s\n",
			orderInfo.Username, orderInfo.OrderID.String(), orderInfo.ProductID, orderInfo.Quantity, orderInfo.Price, orderInfo.TotalCost, orderInfo.Status)
		return orderInfo, nil
	}
	return models.OrderInfo{}, nil
}

func CheckBalanceNoChanged(authHeader string) {
	url := "http://billing-service:8080/balance/"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	request.Header.Set("Authorization", authHeader)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var balance int64
		json.NewDecoder(response.Body).Decode(&balance)
		fmt.Println("8. Ожидается, что баланс не изменился (остался равным 60). Баланс: ", balance)
	}
}
