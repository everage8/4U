#!/bin/bash

set -euo pipefail

BACKUP_DIR="$(dirname "$0")/backups"
DATE=$(date +%Y-%m-%d_%H-%M)

ENV_FILE="$(dirname "$0")/.env"
if [ -f "$ENV_FILE" ]; then

    source "$ENV_FILE"
fi

MONGO_PASSWORD="${MONGO_ROOT_PASSWORD:-}"
if [ -z "$MONGO_PASSWORD" ]; then
    echo "ERROR: MONGO_ROOT_PASSWORD не задан в .env" >&2
    exit 1
fi

mkdir -p "$BACKUP_DIR"

echo "[$(date)] Создаю бэкап..."

docker exec exam-mongo mongodump \
    --uri="mongodb://root:${MONGO_PASSWORD}@localhost:27017/exam_tasks_db?authSource=admin" \
    --archive \
    --gzip > "$BACKUP_DIR/backup_${DATE}.gz"

echo "[$(date)] Бэкап сохранён: $BACKUP_DIR/backup_${DATE}.gz"

find "$BACKUP_DIR" -name "backup_*.gz" -mtime +30 -delete
echo "[$(date)] Старые бэкапы (>30 дней) удалены"
