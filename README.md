# WB Tech Level#0
Учебное задание для прохождения стажировки в WB Tech.  
Описание задания в файле test_task_level_0.pdf

### Использованные технологии:
- [Golang](https://go.dev)
- [Gin framework](https://github.com/gin-gonic/gin)
- [Docker](https://www.docker.com)
- [PostgreSQL](https://www.postgresql.org)
- [NATS-Streaming](https://github.com/nats-io/stan.go)
- [logrus](https://github.com/sirupsen/logrus)
- [viper](https://github.com/spf13/viper)
- [validator](https://github.com/go-playground/validator)


## Как запустить сервис:  
1. Клонировать репозиторий
```
git clone https://github.com/Headrush95/WBTech_Level0 
```  
2. Запустить PostgreSQL образ docker:
```
docker run --name=WBTech -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=WBTechDB -e POSTGRES_USER=WBTechTest -p 5436:5432 -d postgres
```
3. Создать таблицы с помощью файла:
```
.\schema\up.sql
```
4. Запустить NATS-Streaming образ docker:
```
docker run --name=WBNATS -p 4223:4223 -p 8223:8223 -d nats-streaming -p 4223 -m 8223
```
5. Запустить тестового publisher'а через команду:
```
go run .\cmd\publisher\main.go
```
6. Запустить сам сервис:
```
go run .\cmd\consumer\main.go
```
7. Открыть в браузере ```Frontend\index.html```

### <a name="up"></a>Endpoints
- [Получение информации о заказе по его UID](#GetOrderById)
- [Получение информации о всех заказах (служебный)](#GetAllOrders)
- [Сохранение информации о заказе в БД (служебный)](#CreateOrder)

## Примеры запросов
### <a name="GetOrderById">Получение информации о заказе по его UID</a> - метод GET
```localhost:8000/orders/:id"```,  
где id - номер заказа.  
Например, для заказа "b563feb7b2b84b6test" (из исходных данных) запрос будет выглядет следующим образом (при условии, что заказ уже занесли в БД):
```
localhost:8000/orders/b563feb7b2b84b6test
```
В ответе получим информацию о заказе в формате JSON:
```JSON
{
  "order_uid": "b563feb7b2b84b6test",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": "2021-11-26T09:22:07Z",
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}
```
### <a name="GetAllOrders">Получение информации о всех заказах</a> - метод GET
```localhost:8000/orders/```  
Отправив GET запрос на указанный эндпоинт получим ___список___ всех заказов в формате JSON:
```JSON
[
    {
        "order_uid": "b763feb7b2b84b6test",
        "track_number": "WBILMTESTTRACK231",
        ...
        "sm_id": 99,
        "date_created": "2021-11-26T06:22:19Z",
        "oof_shard": "1"
    },
    {
        "order_uid": "b863feb7b2b84b6test",
        "track_number": "WBILMTESTTRACK23",
        ...
        "sm_id": 99,
        "date_created": "2021-11-26T06:22:19Z",
        "oof_shard": "1"
    },
    {
        "order_uid": "b563feb7b2b84b6test",
        "track_number": "WBILMTESTTRACK",
        "entry": "WBIL",
        ...
        "sm_id": 99,
        "date_created": "2021-11-26T06:22:19Z",
        "oof_shard": "1"
    }
]
```

### <a name="CreateOrder">Сохранение информации о заказе в БД</a> - методы POST
Эндпоинт в самом сервисе:
```
localhost:8000/create
```
Эндпоинт publisher'а:
```
localhost:3000/store
```
