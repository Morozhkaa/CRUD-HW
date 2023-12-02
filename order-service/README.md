![License](https://img.shields.io/badge/license-MIT-green)

# Order service
<img align="right" width="45%" src="../images/orders.jpg">

Сервис, хранящий информацию о заказах пользователей.  

__Используемые технологии__:
- PostgreSQL (в качестве хранилища данных)
- Docker (для запуска сервиса)
- Swagger (для документации API)
- gin-gonic/gin (веб фреймворк)
- pgx (драйвер для работы с PostgreSQL)
- slog (для логирования)

Сервис был написан с Clean Architecture, что позволяет легко расширять функционал сервиса и тестировать его. Также был реализован Graceful Shutdown для корректного завершения работы сервиса.


# Usage
Для использования требуется запустить docker-compose.yml файл, расположенный на директорию выше.

Документация была сгенерирована с помощью команды: ```swag init --dir ./internal/adapters/http --generalInfo swagger.go --output ./api/swagger/public --parseDepth 1 --parseDependency```.

После запуска сервиса доку можно посмотреть по адресу `http://localhost:9000/swagger/index.html`.


## Parametrs formats
  * Authorization header: `Bearer <access_token>;<refresh_token>` - получаем после выполнения команды /login.
  * Amount - значение может быть положительным для пополнения баланса или отрицательным для снятия средств.
## Examples

Некоторые примеры запросов
- [Проверка доступности сервиса](#health)
- [Создать заказ](#createOrder)
- [Посмотреть все заказы](#getAllOrders)
- [Посмотреть все заказы пользователя](#getUserOrders)
- [Посмотреть конкретный заказ](#getOrderByID)
- [Удалить заказ](#deleteOrder)

### Проверка доступности сервиса <a name="health"></a>

```curl
curl -X 'GET' \
  'http://localhost:3000/health' \
  -H 'accept: application/json'
```
Ответ:
```json
{
  "success": "service available"
}
```


### Создать заказ <a name="createOrder"></a>

```curl
curl -X 'POST' \
  'http://localhost:9000/createOrder' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer <access_token>;<refresh_token>' \
  -H 'Content-Type: application/json' \
  -d '{
  "price": 560,
  "product_id": 312,
  "quantity": 2
}'
```
Пример ответа:
```json
{
  "order_id": "3f8f0d05-0c59-4e7b-a7b6-48e0d5c11f71",
  "success": "order was successfully created"
}
```

### Посмотреть все заказы <a name="getAllOrders"></a>

```curl
curl -X 'GET' \
  'http://localhost:9000/getAllOrders' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer <access_token>;<refresh_token>'
```
Пример ответа:
```json
[
  {
    "order_id": "3f8f0d05-0c59-4e7b-a7b6-48e0d5c11f71",
    "price": 560,
    "product_id": 312,
    "quantity": 2,
    "status": "success",
    "total_cost": 1120,
    "username": "Maria"
  },
  {
    "order_id": "458f0d05-0c59-4e7b-a7b6-48e0d5c11f67",
    "price": 450,
    "product_id": 123,
    "quantity": 1,
    "status": "success",
    "total_cost": 450,
    "username": "Olga"
  }
]
```

### Посмотреть все заказы пользователя <a name="getUserOrders"></a>

```curl
curl -X 'GET' \
  'http://localhost:9000/getUserOrders' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer <access_token>;<refresh_token>'
```
Пример ответа:
```json
[
  {
    "order_id": "3f8f0d05-0c59-4e7b-a7b6-48e0d5c11f71",
    "price": 560,
    "product_id": 312,
    "quantity": 2,
    "status": "success",
    "total_cost": 1120,
    "username": "Maria"
  }
]
```

### Посмотреть конкретный заказ <a name="getOrderByID"></a>

```curl
curl -X 'GET' \
  'http://localhost:9000/getOrderByID/3f8f0d05-0c59-4e7b-a7b6-48e0d5c11f71' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer <access_token>;<refresh_token>'
```
Пример ответа:
```json
{
  "order_id": "3f8f0d05-0c59-4e7b-a7b6-48e0d5c11f71",
  "price": 560,
  "product_id": 312,
  "quantity": 2,
  "status": "success",
  "total_cost": 1120,
  "username": "Maria"
}
```


### Удалить заказ <a name="deleteOrder"></a>

```curl
curl -X 'DELETE' \
  'http://localhost:9000/deleteOrder/3f8f0d05-0c59-4e7b-a7b6-48e0d5c11f71' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer <access_token>;<refresh_token>'
```

Пример ответа:
```json
{
  "success": "order successfully deleted"
}
```