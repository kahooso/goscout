# university

Вузовские предметы: лекции, практика, лабы, подготовка к защитам и экзаменам. Разбор материала с Claude по единой схеме — правила подачи в [`.claude/TEACHER.md`](../.claude/TEACHER.md), вход в сессию — скилл `study`.

Раздел **изолирован** от остального goscout (Go-трек, PT): ничего за пределами `university/` на него не опирается, кроме вузовских скиллов и раздела «Вуз» в `CLAUDE.md`. Отделить раздел — см. «Как отделить» внизу.

## Структура

```
university/
  README.md                этот файл: стандарт, анонимность, пути
  .gitignore               бинарники и личное — только локально
  _templates/              шаблоны трекера, банка повторения, плана лабы
  <N>-semester/
    <предмет>/
      README.md            ТРЕКЕР — читается первым: режим, где я сейчас, статусы, дедлайны
      journal.md           журнал сессий
      recall.md            банк вопросов на повторение (интервалы — TEACHER.md §5)
      private.md           ЛОКАЛЬНО: преподаватель, его стиль на защите, группа (в .gitignore)
      lectures/            конспекты теории: NN-тема.md (+ img/)
      practice/prN/        семинары, задачи, решения
      labs/labN/           plan.md · zashchita.md · otchet.md · отчёт .docx · свои файлы
```

Пустые папки не заводить: `lectures/`, `practice/`, `labs/` появляются вместе с первым материалом.

## Предметы — 3 семестр

| Папка | Предмет | Что делаем через Claude |
|---|---|---|
| [`3-semester/circuits`](3-semester/circuits/README.md) | Электроника и схемотехника | лабы (Multisim) + защиты, практика |
| [`3-semester/networks`](3-semester/networks/README.md) | Сети и системы передачи информации | не начато |
| [`3-semester/life-safety`](3-semester/life-safety/README.md) | Безопасность жизнедеятельности | лабы: теория + защита, потом измерения |
| [`3-semester/os`](3-semester/os/README.md) | Операционные системы | лабы Pintos |
| [`3-semester/pl`](3-semester/pl/README.md) | Языки программирования | лабы (оптимизация), курсовая |
| [`3-semester/probability`](3-semester/probability/README.md) | Теория вероятностей и матстатистика | лекции через Claude + задачи |

Новый предмет: папка `<N>-semester/<код>/` + `README.md` и `recall.md` из `_templates/` + строка в этой таблице.

## Материалы преподавателей — вне репозитория

```
C:\Users\ryche\Desktop\gostudy-materials\
  <N>-semester\<предмет>\    лекции, методички, требования, учебники — КАК ЕСТЬ, только чтение
D:\multisim-dump\            дистрибутив Multisim 11 (установщик, лицензии)
```

Трекер предмета (раздел «Материалы») называет конкретные файлы. Извлечение текста:
- PDF: `python` + `pypdf` (`PdfReader`) — результат писать в файл, консоль Windows рубит кириллицу. Скан или кривая кодировка — рендерить страницы `pymupdf` в PNG и читать картинкой.
- DOCX: `python` + `zipfile` → `word/document.xml`, картинки — `word/media/`.

## Анонимность (репозиторий публичный)

В git не попадает **ничего**, что указывает на человека: ФИО, преподаватели, одногруппники, номер группы, отчёты.

1. **Бинарники — только локально** (`university/.gitignore`): `.docx`, `.pdf`, `.pptx`, `.ms11`, архивы, скриншоты лаб. Отчёты лежат в папках лаб на диске, но не в git.
2. **Фильтр `anon`** (`.githooks/anon.py`, подключён в `.gitattributes`). При `git add` текстовые файлы `university/` проходят через замену: реальные имена → заглушки (`{{STUDENT}}`, `{{TEACHER_CIRCUITS}}`). На диске файлы остаются с настоящими именами, скрипты работают. Карта замен — `.anon-map` в корне репо, в `.gitignore`. Без неё `git add` по `university/` падает.
3. **Хук `pre-commit`** (`.githooks/pre-commit`) проверяет весь индекс по `.anon-map` и блокирует коммит, если имя где-то осталось (в том числе вне `university/` и в именах файлов).
4. **Детали, нужные для учёбы, но личные** (кто преподаватель, как спрашивает на защите) — в `private.md` предмета, он в `.gitignore`.
5. Появился новый человек/номер — сначала строка в `.anon-map`, потом коммит.

Настройка на новой машине (один раз):
```
git config filter.anon.clean  "python .githooks/anon.py clean"
git config filter.anon.smudge "python .githooks/anon.py smudge"
git config filter.anon.required true
git config core.hooksPath .githooks
```
и восстановить `.anon-map` из бэкапа.

## Как отделить раздел

Если вуз понадобится вынести в свой репозиторий: перенести `university/`, скиллы `study`, `lecture`, `lab`, `problems`, `defense`, `recall`, файл `.claude/TEACHER.md`, `.githooks/` и строки `filter=anon` из `.gitattributes`; удалить раздел «Вуз» из `CLAUDE.md`. Больше связей нет.
