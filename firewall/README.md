# Firewall module

## Launch
### Debug
```bash
$ go run cmd/app/main.go
```

### Release
```bash
$ go build -o ./bin/firewall ./cmd/app/main.go
$ ./bin/firewall
```

### Usage
```bash
Usage of ./bin/firewall:
  -config string
        Path to configuration file (default "config/config.yaml")
```

### Running with PostgreSQL
Передать учетные данные от postgres можно через конфиг. Другой способ - через переменные окружения. Пример: 
```bash
$ docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres
$ export FIREWALL_dbCreds="user=postgres password=postgres dbname=postgres host=localhost port=5432 sslmode=disable" 
$ go run cmd/app/main.go
```

## APIs
* `/nexus/*` - API для общения с Nexus Repository
* `/api/*` - API для пользовательского взаимодействия с firewall
    - `GET /api/status` - получить список всей просканированных пакетов
        - Response `[{"purl":"...","state":S}]`, где 
            - `S = 1` - пакетом можно пользоваться (`healthy`)
            - `S = 2` - пакет попал в карантин (`quarantined`)
            - `S = 3` - пакет был разблокирован (`unquarantined`)
    - `POST /api/status` - получить результаты сканирования
        - Request `{"purl":"..."}`
        - Response `{<report>}`
    - `POST /api/evaluate` - просканировать пакет и получить результаты сканирования
        - Request `{"purl":"..."}`
        - Response `{<report>}`
    - `POST /api/unquarantine` - разблокировать пакет помещённый в карантин
        - Request `{"purl":"..."}`
