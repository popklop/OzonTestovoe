
# Ozon Test Project — GraphQL CRUD Service

  

##  Описание проекта


Это тестовое задание для **Ozon**.

Проект реализован на **Go** и предоставляет GraphQL API с базовыми CRUD-операциями для сущностей:

  

-  **Post** — публикации.

-  **Comment** — комментарии с поддержкой **иерархической структуры** (вложенные комментарии).

  

### Хранилище данных

  

Проект поддерживает два варианта хранения данных (указывается в `config/config.yaml`):

-  `memory` — данные хранятся в оперативной памяти.

-  `postgres` — данные сохраняются в базе PostgreSQL.

### Тестирование

Проект покрыт **unit-тестами**, которые проверяют корректность работы и взаимодействия с БД.

- Для запуска всех тестов достаточно команды: 

 `go test ./...`


##  Запуск проекта через Docker

### Требования

- Установленные [Docker](https://docs.docker.com/get-docker/) и [Docker Compose](https://docs.docker.com/compose/install/).

### Шаги запуска

1. Клонировать репозиторий:

`git clone https://github.com/popklop/OzonTestovoe`
`cd C:/.../../ozontestovoe-folder`

2. Запустить проект:

`docker compose up --build`

3. Сервисы:

GraphQL API: http://localhost:8080

PostgreSQL (для подключения снаружи): localhost:4040

База данных: ozontestovoe

Пользователь: postgres


Пароль: pass

