# YozaTune Support Bot

Telegram-бот поддержки: пересылает вопросы пользователей в группу операторов и отправляет ответы обратно.

## Переменные окружения

| Переменная | Обязательная | Описание |
|---|---|---|
| `BOT_TOKEN` | ✅ | Токен от @BotFather |
| `SUPPORT_CHAT_ID` | ✅ | ID группы поддержки (отрицательное число) |
| `DATABASE_URL` | ✅ | PostgreSQL DSN |
| `WEBHOOK_URL` | — | Публичный URL бота. Если не задан — бот работает в режиме polling |
| `WEBHOOK_PATH` | — | Путь вебхука (по умолчанию `/webhook`) |
| `WEBHOOK_LISTEN_ADDR` | — | Адрес прослушивания (по умолчанию `:8080`) |
| `WEBHOOK_SECRET_TOKEN` | — | Секрет для верификации запросов от Telegram |

Получить `SUPPORT_CHAT_ID`: добавить бота в группу, написать любое сообщение, открыть `https://api.telegram.org/bot<TOKEN>/getUpdates` и найти `chat.id`.

---

## Локальная разработка

### Зависимости
- Go 1.24+
- Docker + Docker Compose
- [mage](https://magefile.org): `go install github.com/magefile/mage@latest`
- [ngrok](https://ngrok.com) для вебхука (или убрать `WEBHOOK_URL` из `.env` для polling)

### Запуск

```bash
# 1. Скопировать и заполнить конфиг
cp .env.example .env

# 2. Поднять БД
mage db

# 3. Запустить бота (с polling — без ngrok)
go run ./cmd/bot

# 4. Или с вебхуком через ngrok:
ngrok http 8080
# Прописать полученный URL в .env → WEBHOOK_URL=https://xxxx.ngrok-free.app
go run ./cmd/bot
```

### Команды mage

| Команда | Описание |
|---|---|
| `mage start` | Собрать образы и запустить всё в фоне |
| `mage stop` | Остановить и удалить контейнеры |
| `mage restart` | Пересобрать и перезапустить бота |
| `mage logs` | Стримить логи бота |
| `mage status` | Статус контейнеров |
| `mage db` | Запустить только локальную БД |

---

## Деплой на сервер (рядом с основным ботом)

Предполагается, что основной бот уже запущен через Docker Compose с Caddy и PostgreSQL.

### 1. Создать общую Docker-сеть (один раз)

```bash
docker network create shared
```

### 2. Подключить основной бот к сети

В `docker-compose.yml` основного бота добавить `shared` сеть к сервисам `postgres` и `caddy`:

```yaml
services:
  postgres:
    # ... существующий конфиг ...
    networks:
      - default
      - shared

  caddy:
    # ... существующий конфиг ...
    networks:
      - default
      - shared

networks:
  shared:
    external: true
    name: shared
```

Применить:

```bash
docker compose up -d
```

### 3. Добавить маршрут в Caddyfile

```caddy
supportbot.yourdomain.com {
    reverse_proxy support-bot:8080
}
```

Перезагрузить Caddy:

```bash
docker exec caddy caddy reload --config /etc/caddy/Caddyfile
```

### 4. Подготовить `.env` на сервере

```env
BOT_TOKEN=
SUPPORT_CHAT_ID=

WEBHOOK_URL=https://supportbot.yourdomain.com
WEBHOOK_PATH=/webhook
WEBHOOK_LISTEN_ADDR=:8080
WEBHOOK_SECRET_TOKEN=

# DATABASE_URL не нужен — задаётся в docker-compose.yml через postgres основного бота
```

### 5. Запустить

```bash
git clone <repo> /opt/support-bot
cd /opt/support-bot
cp .env.example .env && nano .env  # заполнить токен и SUPPORT_CHAT_ID

docker compose up -d --build
docker compose logs -f bot          # убедиться что запустился
```

---

## Как работает

```
Пользователь → пишет боту
    Бот → пересылает в группу поддержки (текст — одним сообщением, медиа — форвардом)
    Бот → сохраняет связку "сообщение в группе ↔ пользователь" в БД

Оператор → делает Reply на сообщение в группе
    Бот → находит пользователя по ID сообщения в группе
    Бот → пересылает ответ пользователю (без плашки "переслано")
```

### Закрытие тикетов

Оператор делает **Reply** на сообщение в группе и пишет `/done` — тикет удаляется из БД, последующие реплаи на это сообщение уже не дойдут пользователю.

Тикеты старше 30 дней удаляются автоматически раз в 24 часа.
