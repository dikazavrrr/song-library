# Music API

## Описание
Этот проект представляет REST API для управления песнями, а также получения информации о песнях из внешнего API.

## Установка и запуск

### 1. Клонирование репозитория
```sh
git clone https://github.com/dikazavrrr/song-library.git
```

### 3. Конфигурация
Создайте файл `config.yaml` в корне проекта и добавьте:
```yaml
api:
  external_url: "http://EXTERNAL_API"
```
Замените `EXTERNAL_API` на реальный адрес API.

### 4. Запуск сервера
```sh
go run main.go
```

Сервер будет работать по адресу `http://localhost:8080`.

## API Эндпоинты

### Получить все песни
```http
GET /songs
```
**Ответ:**
```json
[
  { "id": 1, "title": "Song 1", "artist": "Artist 1" },
  { "id": 2, "title": "Song 2", "artist": "Artist 2" }
]
```

### Получить песню по ID
```http
GET /songs/:id
```
**Ответ:**
```json
{ "id": 1, "title": "Song 1", "artist": "Artist 1" }
```

### Добавить песню
```http
POST /songs
```
**Тело запроса:**
```json
{ "title": "New Song", "artist": "New Artist" }
```
**Ответ:**
```json
{ "message": "Song added successfully" }
```

### Обновить песню
```http
PUT /songs/:id
```
**Тело запроса:**
```json
{ "title": "Updated Song", "artist": "Updated Artist" }
```
**Ответ:**
```json
{ "message": "Song updated successfully" }
```

### Удалить песню
```http
DELETE /songs/:id
```
**Ответ:**
```json
{ "message": "Song deleted successfully" }
```

### Получение информации о песне из внешнего API
```http
GET /info?group={group}&song={song}
```
**Пример запроса:**
```
GET /info?group=Coldplay&song=Yellow
```
**Ответ:**
```json
{
  "releaseDate": "16.07.2006",
  "text": "Ooh baby, don't you know I suffer?...",
  "link": "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
}
```

## Тестирование
Для тестирования можно использовать `curl` или Postman.
Пример теста для `CreateSong`:
```sh
curl -X POST "http://localhost:8080/songs" -H "Content-Type: application/json" -d '{"title":"New Song","artist":"New Artist"}'
```

## Swagger
API документация доступна по адресу:
```
http://localhost:8080/swagger/index.html
```

## Лицензия
MIT License

