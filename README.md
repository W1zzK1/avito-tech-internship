# Авито. Тестовое задание для стажёра Backend (осенняя волна 2025)

## Как запустить проект

```azure
docker-compose up --build
```

из директории проекта соответственно :)

## Список ендпоинтов

| Метод |         Адрес         |
|-------|:---------------------:|
| POST  |       /team/add       |
| GET   |  /team/get/:teamName  |
| POST  |     /users/addNew     |
| GET   |  /users/getById/:id   |
| POST  |  /users/setIsActive   |
| POST  |  /pullRequest/create  |
| POST  |  /pullRequest/merge   |
| POST  | /pullRequest/reassign |
| GET   |  /stats/getAllStats   |
| GET   |       /health         |
