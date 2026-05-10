# Course builder
Если у вас есть курс с видео и docx, то это приложение создаст для вас единый автономный HTML-курс.  
Нулевые внешние зависимости — итогом является один `.html` файл, который работает без интернета.

---

## Установка

### macOS / Linux

```sh
curl -fsSL https://raw.githubusercontent.com/insigmo/course_builder/refs/heads/master/install.sh | sudo sh
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/insigmo/course_builder/refs/heads/master/install.ps1 | iex
```

### Сборка из исходников через go

```bash
go install github.com/insigmo/course-builder/cmd/course-builder
```

---

## Использование

```bash
course-builder ./МойКурс
```

Внутри папки создаётся файл `МойКурс.html`.

---

## Поддерживаемые платформы

| OS      | Arch                  |
|---------|-----------------------|
| Linux   | amd64                 |
| Linux   | arm64                 |
| macOS   | amd64                 |
| macOS   | arm64 (Apple Silicon) |
| Windows | amd64                 |
| Windows | arm64                 |
