# Uptime monitoring
Сервис для мониторинга сайтов и API. На данный момент поддерживает только HTTP-мониторы с проверкой статуса и содержания ответа.

### Функциональность
- Аутентификация
- Мониторинг
  - Создание и управление мониторами(http)
  - Асинхронное выполнение проверок (worker)
  - Настройка таймаутов и интервалов
  - Проверка HTTP-ответов по статусу и содержимому

## Использование
```
git clone <repo>
cd <repo>
docker compose up 
```
Если не выполняются миграции, можно повторить команду 
```
docker compose up 
```
или выполнить по шагам
```
docker compose up db -d
docker compose run migrate
docker compose up backend
```

### Конфигурация
Конфигурация хранится в configs/config.yml:
```
db:
  host: "db"
  port: "5432"
  user: "postgres"
  dbname: "postgres"
  sslmode: "disable"

log:
  level: "debug"
  format: "json" 

server:
  host: "localhost"
  port: 8080
```
### api

##### Создание HTTP-монитора
POST /monitors
```
{
  "name": "Google Homepage Monitor",
  "type": "http",
  "target": "https://www.google.com",
  "timeout": 5,
  "interval": 60,
  "is_active": true,
  "request_spec": {
    "method": "GET",
    "headers": {
      "User-Agent": "Uptime-Monitor/1.0"
    }
  },
  "expected_response": {
    "expect_status": 200,
    "expect_body_contains": ["Google"]
  }
}
```



