# Merch Store

## Описание

Онлайн-магазин мерча с возможностью покупки товаров, перевода средств между пользователями и просмотра информации о своих покупках и переводах.

## Установка и запуск

Для запуска проекта вам потребуется [Docker](https://www.docker.com/get-started) и [Docker Compose](https://docs.docker.com/compose/install/).

### 1. Клонирование репозитория

```bash
git clone https://github.com/yourusername/avito-shop.git
cd avito-shop
```

### 2. Запуск проекта

```bash
docker-compose up --build
```

### 3. Запуск тестов

Запуск всех тестов в контейнере:

```bash
docker exec -it avito-shop-service  bash
go test -v ./...
```
### 4. Вопросы и решения

почему используется метод GET для /api/buy/{item}?
Решение: метод был изменен на метод POST, так как он используется в HTTP для отправки данных на сервер с целью создания или обновления ресурса (создаём или меняем данные о покупах в таблице)?
