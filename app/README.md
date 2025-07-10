# K8s Echo App

Простое веб-приложение на Go для Kubernetes, которое возвращает отправленные данные и системную информацию.

## Функциональность

Приложение имеет один эндпоинт `/` который при любом запросе возвращает:

- Полученные данные (echo)
- Информацию о запросе (метод, заголовки, URI)
- Системную информацию (директория, hostname, пользователь)
- Сетевую информацию (DNS серверы, сетевые адаптеры)
- Переменные окружения
- Время работы сервера

## Использование

```bash
# GET запрос
curl http://localhost:8080/

# POST запрос с JSON данными
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, K8s!"}'

# POST запрос с текстовыми данными
curl -X POST http://localhost:8080/ \
  -d "Hello from K8s"
```

## Локальный запуск

```bash
# Через Docker
docker build -t k8s-app .
docker run -p 8080:8080 k8s-app

# Если Go установлен
go run main.go
```

## Сборка и публикация образа

```bash
# Вход в Docker Hub
docker login

# Сборка и публикация
./build-and-push.sh
```

## Развертывание в Kubernetes

```bash
# Применить манифесты
kubectl apply -f k8s-deployment.yaml

# Проверить статус
kubectl get pods -l app=k8s-echo-app
kubectl get svc k8s-echo-app-service

# Проверить логи
kubectl logs -l app=k8s-echo-app
```

## Docker образ

Образ доступен на Docker Hub: `matvey5686/k8s-app:latest`
