![License](https://img.shields.io/badge/license-MIT-green)

# User service

<img align="right" width="45%" src="../images/users.png">
Cервис, хранящий информацию о пользователях (поддерживает базовые CRUD операции).  

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

После запуска сервиса доку можно посмотреть по адресу `http://localhost:3000/swagger/index.html`.


## Parametrs formats
  * email  `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$` - почта пользователя.
    - [a-zA-Z0-9._%+\-]+   - набор символов до @
    - @					- символ @, который разделяет имя пользователя и доменное имя.
    - [a-zA-Z0-9.\-]+		- набор символов в доменном имени
  	- \\.					- символ точки
    - [a-zA-Z]{2,}			- расширение домена (минимум 2 символа, напр: .com, .net, .org)

  * phone  `^\+\d{11}$` - номер телефона.
  * Authorization header: `Bearer <access_token>;<refresh_token>` - получаем после выполнения команды /login.

## Examples
Некоторые примеры запросов
- [Проверка доступности сервиса](#health)
- [Создание пользователя](#create)
- [Удаление пользователя](#delete)
- [Обновление информации о пользователе](#update)
- [Получение информации о пользователе](#get)
- [Авторизация пользователя](#login)
- [Проверка токенов авторизации](#verify)


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


### Создание пользователя <a name="create"></a>

```curl
curl -X 'POST' \
  'http://localhost:3000/user' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "email": "iivanov@gmail.com",
  "first_name": "Ivan",
  "last_name": "Ivanov",
  "password": "qwerty1234",
  "phone": "+79999999999",
  "username": "IvanIvanov2000"
}'
```
Пример ответа:
```json
{
  "success": "user with username 'IvanIvanov2000' created"
}
```

### Удаление пользователя <a name="delete"></a>

```curl
curl -X 'DELETE' \
  'http://localhost:3000/user/IvanIvanov2000' \
  -H 'accept: application/json'
```
Пример ответа:
```json
{
  "success": "user with username 'IvanIvanov2000' deleted"
}
```


### Обновление информации о пользователе <a name="update"></a>

```curl
curl -X 'PUT' \
  'http://localhost:3000/user/IvanIvanov2000' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "email": "iivanov@gmail.com",
  "first_name": "Ivan",
  "last_name": "Ivanov",
  "password": "qwerty1234",
  "phone": "+79999999999",
  "username": "IvanIvanov2000"
}'
```
Пример ответа:
```json
{
  "success": "information for user with username 'IvanIvanov2000' updated"
}
```


### Получение информации о пользователе <a name="get"></a>

```curl
curl -X 'GET' \
  'http://localhost:3000/user/IvanIvanov2000' \
  -H 'accept: application/json'
```
Пример ответа:
```json
{
  "email": "iivanov@gmail.com",
  "first_name": "Ivan",
  "last_name": "Ivanov",
  "phone": "+79999999999",
  "username": "IvanIvanov2000"
}
```


### Авторизация пользователя <a name="login"></a>

```curl
curl -X 'POST' \
  'http://localhost:3000/login' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "login": "IvanIvanov2000",
  "password": "qwerty1234"
}'
```
Пример ответа:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJJdmFuSXZhbm92MjAwMCIsInN1YiI6ImlpdmFub3ZAZ21haWwuY29tIiwiZXhwIjoxNzAxNTExMjM0fQ.o0sDOb2lzTMxXbsZddGpKynEcT7Zeb2DBmUgsCltzEk",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJJdmFuSXZhbm92MjAwMCIsInN1YiI6ImlpdmFub3ZAZ21haWwuY29tIiwiZXhwIjoxNzAxNTE0Nzc0fQ.xwONvJdQy1UUT4IE8orokaxAVmrW5eL07aJGnqJyaK4"
}
```


### Проверка токенов авторизации <a name="verify"></a>

```curl
curl -X 'POST' \
  'http://localhost:3000/verify' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer <access_token>;<refresh_token>' \
  -d ''
```
Пример ответа:
```json
{
  "accessToken": "<access_token>",
  "email": "iivanov@gmail.com",
  "login": "IvanIvanov2000",
  "refreshToken": "<refresh_token>"
}
```

