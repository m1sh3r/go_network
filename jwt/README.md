# JWT пример (Задание 7)

Генерация JWT и сервер с /login и /protected.

Запуск:
- Подготовка: `./prepare.sh`
- Генерация токена: `./run_token.sh`
- Сервер: `./run_server.sh`

Проверка:
- Валидный токен: `./curl_success.sh`
- Невалидный токен: `./curl_invalid.sh`
- Без токена: `./curl_no_token.sh`
