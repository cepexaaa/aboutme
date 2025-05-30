#!/usr/bin/bash

RESTORE_DIR="$HOME/restore"

find_latest_backup() {
    find "$HOME" -maxdepth 1 -type d -name "Backup-*" | sort -r | head -n 1
}

if [ ! -d "$RESTORE_DIR" ]; then
    mkdir "$RESTORE_DIR"
fi

LATEST_BACKUP=$(find_latest_backup)

if [ -z "$LATEST_BACKUP" ]; then
    echo "The current backup directory was not found."
    exit 1
fi

for file in "$LATEST_BACKUP"/*; do
    if [[ "$(basename "$file")" != *.*([0-9])-*[0-9] ]]; then
        cp "$file" "$RESTORE_DIR/"
    fi
done

echo "Files have been successfully copied to the directory $RESTORE_DIR."
