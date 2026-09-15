# internal/entity

Сущности базы данных.

## Назначение

Эта директория содержит:
- Структуру Item, соответствующую таблице в PostgreSQL
- Методы для сериализации/десериализации
- Связи между сущностями

## Структура

- `item.go` - основная сущность Item

## Поля сущности

- ID, Type, Title, URL, Score
- Text (для комментариев)
- ByUser, ParentID, PollID
- CreatedAt, UpdatedAt
- Deleted, Dead, Descendants