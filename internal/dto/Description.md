# internal/dto

DTO (Data Transfer Objects) для API.

## Назначение

Эта директория содержит:
- Структуры запросов для API (ItemListRequest, CommentTreeRequest)
- Структуры ответов для API (ItemListResponse, CommentTreeResponse)
- Валидация входных параметров
- Преобразование между DTO и внутренними сущностями

## Структура

- `requests.go` - структуры запросов
- `responses.go` - структуры ответов

## Особенности

- Валидация параметров через теги `binding`
- Чистое разделение между входными и выходными данными