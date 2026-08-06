#!/bin/bash

set -euo pipefail

BACKUP_FILE="${1:-}"
if [ -z "$BACKUP_FILE" ] || [ ! -f "$BACKUP_FILE" ]; then
    echo "Использование: $0 <путь-к-файлу-бэкапа.gz>" >&2
    echo "Пример:        $0 ./backups/backup_2025-01-15_03-00.gz" >&2
    exit 1
fi

ENV_FILE="$(dirname "$0")/.env"
if [ -f "$ENV_FILE" ]; then

    source "$ENV_FILE"
fi

MONGO_PASSWORD="${MONGO_ROOT_PASSWORD:-}"
if [ -z "$MONGO_PASSWORD" ]; then
    echo "ERROR: MONGO_ROOT_PASSWORD не задан в .env" >&2
    exit 1
fi

echo "[$(date)] Восстанавливаю из: $BACKUP_FILE"

cat "$BACKUP_FILE" | docker exec -i exam-mongo mongorestore \
    --uri="mongodb://root:${MONGO_PASSWORD}@localhost:27017/?authSource=admin" \
    --nsInclude="exam_tasks_db.*" \
    --archive \
    --gzip \
    --drop

echo "[$(date)] Восстановление завершено"
