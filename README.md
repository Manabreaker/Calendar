# Calendar (In‑Memory Event API)

Коротко: это простой HTTP API для управления событиями (календарь) в памяти процесса без БД. Поддерживает CRUD и выборку событий за день / неделю / месяц / все.

## Возможности
- POST /create_event — создать событие
- POST /update_event — обновить (по совпадающему `id` перезаписывается весь объект)
- POST /delete_event?id=ID — удалить по ID
- GET  /events_for_day   — события на текущий день
- GET  /events_for_week  — события текущей недели (понедельник–воскресенье)
- GET  /events_for_month — события текущего месяца
- GET  /events_for_all   — все добавленные события

## Модель данных
```json
{
  "id": "string",        // обязательный, уникальный
  "title": "string",     // обязательный
  "description": "string",
  "date": "YYYY-MM-DD",  // обязательный (используется для фильтрации)
  "owner_id": "string"   // обязательный
}
```
Валидация: `id`, `title`, `date`, `owner_id` не должны быть пустыми. Формат даты — строка (парсится layout `2006-01-02`).

## Конфигурация
Файл `configs/config.yaml`:
```yaml
server:
  host: "localhost"
  port: "8080"
```
Измените при необходимости (например, порт).

## Запуск
Требования: Go (версия из `go.mod` или новее, с поддержкой модулей).  
Команды (Windows cmd):
```cmd
go mod download
go run ./cmd/calendar
```
Логи сообщат адрес: `localhost:8080`.

## Примеры запросов (curl)
Создать:
```cmd
curl -X POST http://localhost:8080/create_event ^
  -H "Content-Type: application/json" ^
  -d "{\"id\":\"1\",\"title\":\"Meet Bob\",\"description\":\"Discuss\",\"date\":\"2025-10-01\",\"owner_id\":\"user123\"}"
```
Обновить:
```cmd
curl -X POST http://localhost:8080/update_event ^
  -H "Content-Type: application/json" ^
  -d "{\"id\":\"1\",\"title\":\"New title\",\"description\":\"Updated\",\"date\":\"2025-10-01\",\"owner_id\":\"user123\"}"
```
Удалить:
```cmd
curl -X POST "http://localhost:8080/delete_event?id=1"
```
Получить события:
```cmd
curl http://localhost:8080/events_for_day
curl http://localhost:8080/events_for_week
curl http://localhost:8080/events_for_month
curl http://localhost:8080/events_for_all
```

## Ответы
Успех (создание/обновление/удаление):
```json
{"result":"success"}
```
Ошибки: JSON вида `{ "error": "описание" }` и соответствующий HTTP статус (400 / 404).

## Как работает фильтрация
- "Сегодня" — границы календарного дня локального времени.
- "Неделя" — понедельник–воскресенье текущей недели.
- "Месяц" — первый–последний день текущего месяца.

## Тесты
Запуск тестов:
```cmd
go test ./...
```
Покрывают: CRUD, ошибки валидации, некорректный JSON.

## Ограничения / Нюансы
- Данные хранятся только в памяти (пропадают при рестарте).
- Нет пагинации и аутентификации.
- Полное обновление события (partial update не реализован).
- Нет проверки пересечений событий.

## Возможные улучшения
- Хранение в БД (PostgreSQL / SQLite)
- Индексация событий по дате для ускорения выборки
- Фильтрация по `owner_id`
- Partial update (PATCH)
- Валидация формата даты через time.Parse при создании

---
Коротко: запусти `go run ./cmd/calendar`, шли JSON на эндпоинты — получай события. Удачи!

