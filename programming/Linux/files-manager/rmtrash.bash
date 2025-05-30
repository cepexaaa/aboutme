#!/usr/bin/bash

if [ "$#" -ne 1 ]; then
    echo "Использование: $0 <имя_файла>"
    exit 1
fi

FILE="$1"
TRASH_DIR="$HOME/.trash"
TRASH_LOG="$HOME/.trash.log"

if [ ! -f "$FILE" ]; then
    echo "Файл '$FILE' не найден."
    exit 1
fi

if [ ! -d "$TRASH_DIR" ]; then
    mkdir "$TRASH_DIR"
fi

LINK_NUM=1
while [ -e "$TRASH_DIR/$LINK_NUM" ]; do
    LINK_NUM=$((LINK_NUM + 1))
done

ln "$FILE" "$TRASH_DIR/$LINK_NUM"

rm "$FILE"

echo "$(pwd)/$FILE -> $LINK_NUM" >> "$TRASH_LOG"

echo "Файл '$FILE' перемещен в корзину как '$LINK_NUM'."

exit 0
