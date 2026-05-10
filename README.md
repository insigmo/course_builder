# course-builder

Go-порт оригинального Python-скрипта. Конвертирует папку с `.docx` и видеофайлами в единый интерактивный HTML-курс. Нулевые внешние зависимости.

## Использование

```bash
course-builder ./МойКурс
```

Рядом с бинарником создастся `МойКурс.html`.

## Сборка

```bash
go build -o course-builder ./cmd/course-builder
```

## Релиз через GitHub Actions

```bash
git tag v1.0.0
git push --tags
```

Actions автоматически собирает бинарники для всех платформ и публикует их в Releases.

## Поддерживаемые платформы

| OS      | Arch  |
|---------|-------|
| Linux   | amd64 |
| Linux   | arm64 |
| macOS   | amd64 |
| macOS   | arm64 (Apple Silicon) |
| Windows | amd64 |
| Windows | arm64 |

## Структура

```
course-builder/
├── cmd/course-builder/main.go
├── internal/
│   ├── config/config.go          ← константы
│   ├── prefix/prefix.go          ← детект/удаление [X] префиксов
│   ├── docx/parser.go            ← парсинг .docx → HTML (stdlib zip+xml)
│   ├── video/converter.go        ← конвертация через ffmpeg
│   ├── builder/builder.go        ← рекурсивный сбор уроков/шагов
│   ├── builder/sort.go           ← числово-алфа сортировка
│   └── htmlrender/render.go      ← финальный HTML (//go:embed)
├── .github/workflows/release.yml
└── go.mod
```
