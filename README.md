# 4U

Фан сайт с базой заданий и ответов к ним
Функционал - (LaTeX-формулы, фильтрацию по типу задания и управление через встроенную админ-панель)

 Бэк - Go 1.22, Gin 
 База данных - MongoDB 7 
 Фронт - Vanilla JS, HTML, CSS 
 Формулы - KaTeX (CDN) 
 Контейнеры - Docker, Docker Compose 

## Структура проекта

```
4U/
├── .env                        # переменные для docker-compose (не коммитить)
├── .gitignore
├── docker-compose.yml          # MongoDB + бэкенд 
├── backup.sh                   # создать бэкап БД   
├── restore.sh                  # восстановить БД из бэкапа
└── backend/
    ├── .env                    # переменные для локального запуска
    ├── Dockerfile
    ├── main.go                 # точка входа
    ├── config/                 # загрузка конфигурации из env
    ├── database/               # подключение к MongoDB
    ├── models/                 # Task, Admin
    ├── dto/                    # структуры запросов/ответов
    ├── repository/             # работа с коллекциями MongoDB
    ├── service/                # бизнес-логика
    ├── handlers/               # HTTP-хендлеры
    ├── middleware/             # JWT auth, роли, CORS
    ├── jwt/                    # генерация и валидация токенов
    ├── response/               # единый формат ответа API
    ├── router/                 # сборка роутов
    └── web/                    # фронтенд
        ├── css/style.css
        ├── html/               # main.html, login.html, admin.html
        └── js/                 # shared-data.js, Js.js, admin.js
```

## Локальный запуск (без Docker)

**Требования:** Go 1.22+, MongoDB запущена на `localhost:27017`

```bash
cd 4U/backend
go run main.go
```

Сайт доступен на `http://localhost:8080`

Переменные окружения берутся из `backend/.env`:

```env
PORT=8080
MONGO_URI=mongodb://localhost:27017
MONGO_DB=exam_tasks_db
JWT_SECRET=ваш-секрет
JWT_EXPIRY_HOURS=72
ADMIN_LOGIN=admin
ADMIN_PASSWORD=ваш-пароль
CORS_ORIGINS=*
WEB_ROOT=web
```

## Запуск через Docker

```bash
cd 4U

# Заполнить переменные (1 раз)
nano .env

# Запустить
docker compose up --build -d

# Логи
docker compose logs -f backend

# Остановить
docker compose down
```

Переменные в `4U/.env`:

```env
MONGO_ROOT_PASSWORD=пароль-для-монго-(либо задать самому, либо при )
JWT_SECRET=минимум-32-символа
ADMIN_LOGIN=admin
ADMIN_PASSWORD=любой-пароль-для-входа-в-лк-(по умолчанию admin123)
```

Сгенерировать JWT_SECRET:
```bash
openssl rand -hex 32
```

## API

Все ответы оборачиваются в `{ "status": "OK", "message": "...", "data": ... }`

| Метод | Путь | Auth | Описание |
|---|---|---|---|
| GET | `/api/v1/tasks` | нет | Список заданий (фильтр: `?subject=&type=`) |
| POST | `/api/v1/auth/login` | нет | Логин администратора |
| GET | `/api/v1/admin/tasks` | JWT | Список заданий в админке (поиск, пагинация) |
| POST | `/api/v1/admin/tasks` | JWT | Создать задание |
| PUT | `/api/v1/admin/tasks/:id` | JWT | Обновить задание |
| DELETE | `/api/v1/admin/tasks/:id` | JWT | Удалить задание |

## Страницы

| URL | Описание |
|---|---|
| `/` | Главная — выбор предмета, типа, список заданий |
| `/login` | Вход в админ-панель |
| `/admin` | Админ-панель — таблица заданий, создание, редактирование, удаление |

## Предметы и типы заданий

**Математика (профиль)**
- Теория вер./чисел
- Преобразования
- Планиметрия
- Стереометрия
- Уравнения
- Текстовое задание
- Графики функций
- Неравенства
- Производные/первообразные

**Физика**
- Магнетизм
- Молекулярная физика и термодинамика
- Перевод в СИ
- Квантовая/Ядерная физика
- Механика
- Механические колебания и волны
- Электродинамика
- Оптика

## LaTeX-формулы

В полях условия и решения поддерживается синтаксис LaTeX через KaTeX.

 **Для ознакомления**
https://ru.wikibooks.org/wiki/LaTeX/%D0%9C%D0%B0%D1%82%D0%B5%D0%BC%D0%B0%D1%82%D0%B8%D1%87%D0%B5%D1%81%D0%BA%D0%B8%D0%B5_%D1%84%D0%BE%D1%80%D0%BC%D1%83%D0%BB%D1%8B#
https://katex.org/docs/supported#operators


## Бэкап БД

**Создать бэкап:**
```bash
cd 4U
chmod +x backup.sh
./backup.sh
```
Файл сохраняется в `4U/backups/backup_YYYY-MM-DD_HH-MM.gz`. Бэкапы старше 30 дней удаляются автоматически.

**Восстановить из бэкапа:**
```bash
./restore.sh ./backups/backup_2025-01-15_03-00.gz
```

**Автоматический бэкап (cron, каждый день в 03:00):**
```bash
crontab -e
# Добавить строку:
0 3 * * * /opt/exam-tasks/4U/backup.sh >> /var/log/exam-backup.log 2>&1
```

## Деплой на сервер (если на Ubuntu 22.04)

```bash
# 1. Установить Docker на сервере
curl -fsSL https://get.docker.com | sh

# 2. Скопировать проект
scp -r ./4U root@IP_СЕРВЕРА:/opt/exam-tasks/

# 3. Создать .env на сервере
nano /opt/exam-tasks/4U/.env

# 4. Запустить
cd /opt/exam-tasks/4U
docker compose up --build -d
```

**Когда домен будет готов** обновить `CORS_ORIGINS` в `4U/.env`:
```env
CORS_ORIGINS=https://ваш-домен.ru
```
И перезапустить: `docker compose up -d`

# После изменения  ADMIN_PASSWORD в .env
docker compose up --build -d
```

Если задания уже есть — изменить через mongosh:
```bash
docker exec -it exam-mongo mongosh \
  -u root -p $MONGO_ROOT_PASSWORD --authenticationDatabase admin

use exam_tasks_db
db.admins.findOne()   # посмотреть текущего админа
```
