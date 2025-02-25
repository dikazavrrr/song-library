# Song Library

## Описание
Этот проект представляет REST API для управления песнями, а также получения информации о песнях из внешнего API.

## Установка и запуск

### 1. Клонирование репозитория
```sh
git clone https://github.com/dikazavrrr/song-library.git
```

### . Конфигурация
Замените `EXTERNAL_API` на реальный адрес API в `config/config.yaml`:
```yaml
api:
  external_url: "http://EXTERNAL_API"
```

### 4. Запуск сервера
```sh
make run
```

Сервер будет работать по адресу `http://localhost:8080`.

## API Эндпоинты

### Получить все песни
```http
GET /songs/
```
**Ответ:**
```json
[
  { "id": 1, "song": "Song 1", "group": "Artist 1", "releaseDate": "2025-01-01", "lyrics": "", "link": ""},
  { "id": 2, "song": "Song 2", "group": "Artist 2", "releaseDate": "2025-01-01", "lyrics": "", "link": ""},

]
```

### Получить песню по ID
```http
GET /songs/:id
```
**Ответ:**
```json
{ "id": 1, "song": "Song 1", "group": "Artist 1", "releaseDate": "2025-01-01", "lyrics": "", "link": ""}
```

### Добавить песню
```http
POST /songs/
```
**Тело запроса:**
```json
{
    "group": "Metallica",
    "song": "Nothing Else Matters"
}
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
{
    "group": "Metallicaa",
    "song": "Nothing Else Matters"
}
```
**Ответ:**
```json
{
    "group": "Metallicaa",
    "song": "Nothing Else Matters",
    "releaseDate": "2025-01-01", 
    "lyrics": "", 
    "link": ""
}
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
GET /songs/info?group={group}&song={song}
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

## Swagger
API документация доступна по адресу:
```
http://localhost:8080/swagger/index.html
```

