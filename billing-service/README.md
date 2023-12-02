![License](https://img.shields.io/badge/license-MIT-green)

# Billing service

<img align="right" width="45%" src="../images/billing.jpg">
Сервис, хранящий информацию об остатках на счетах пользователей.  

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

После запуска сервиса доку можно посмотреть по адресу `http://localhost:8080/swagger/index.html`.

## Parametrs formats
  * Authorization header: `Bearer <access_token>;<refresh_token>` - получаем после выполнения команды /login.
  * Amount - значение может быть положительным для пополнения баланса или отрицательным для снятия средств.
## Examples

Некоторые примеры запросов
- [Проверка доступности сервиса](#health)
- [Получение баланса](#get)
- [Пополнение баланса](#update_add)
- [Снятие средств](#update_del)

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


### Получение баланса <a name="get"></a>

```curl
curl -X 'GET' \
  'http://localhost:8080/balance' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer <access_token>;<refresh_token>'
```
Пример ответа:
```json
{
  1000
}
```

### Пополнение баланса <a name="update_add"></a>

```curl
curl -X 'POST' \
  'http://localhost:8080/balance' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer <access_token>;<refresh_token>' \
  -H 'Content-Type: application/json' \
  -d '{
  "amount": 1000
}'
```
Пример ответа:
```json
{
  "success": "information for user with username 'IvanIvanov2000' updated"
}
```

### Снятие средств <a name="update_del"></a>

```curl
curl -X 'POST' \
  'http://localhost:8080/balance' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer <access_token>;<refresh_token>' \
  -H 'Content-Type: application/json' \
  -d '{
  "amount": -600
}'
```
Пример ответа:
```json
{
  "success": "information for user with username 'IvanIvanov2000' updated"
}
```