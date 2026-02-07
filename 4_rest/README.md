# REST пример (Задание 4)

## Запуск
```bash
cd 4_rest/server
go run main.go
```

## Проверка (curl)
```bash
# Получить всех
curl http://localhost:8080/users

# Получить по ID
curl http://localhost:8080/users/1

# Создать
curl -X POST -H "Content-Type: application/json" \
  -d '{"name":"Charlie"}' http://localhost:8080/users

# Обновить по ID
curl -X PUT -H "Content-Type: application/json" \
  -d '{"name":"Updated Name"}' http://localhost:8080/users/1

# Удалить по ID
curl -X DELETE http://localhost:8080/users/1
```
