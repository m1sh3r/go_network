# Создать нового
curl -u admin:secret -X POST -H "Content-Type: application/json" \
  -d '{"name":"Charlie"}' http://localhost:8080/users
