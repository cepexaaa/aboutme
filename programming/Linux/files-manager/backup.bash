#!/usr/bin/bash

SOURCE_DIR="$HOME/source"
BACKUP_REPORT="$HOME/backup-report"
CURRENT_DATE=$(date +%Y-%m-%d)
BACKUP_DIR="$HOME/Backup-$CURRENT_DATE"

if [ ! -d "$SOURCE_DIR" ]; then
    echo "The $SOURCE_DIR directory does not exist. Create it and add files for backup."
    exit 1
fi

find_active_backup() {
    find "$HOME" -maxdepth 1 -type d -name "Backup-*" | sort -r | head -n 1
}

create_backup_dir() {
    mkdir "$BACKUP_DIR"
    echo "A new backup directory has been created: $BACKUP_DIR ($CURRENT_DATE)" >> "$BACKUP_REPORT"
    echo "Files copied from $SOURCE_DIR:" >> "$BACKUP_REPORT"
    cp -v "$SOURCE_DIR"/* "$BACKUP_DIR/" >> "$BACKUP_REPORT"
}

update_backup_dir() {
    local active_backup="$1"
    echo "Making changes to the current backup directory: $active_backup ($CURRENT_DATE)" >> "$BACKUP_REPORT"

    for file in "$SOURCE_DIR"/*; do
        local filename=$(basename "$file")
        local backup_file="$active_backup/$filename"

        if [ -f "$backup_file" ]; then
            local source_size=$(stat -c %s "$file")
            local backup_size=$(stat -c %s "$backup_file")

            if [ "$source_size" -ne "$backup_size" ]; then
                local new_name="${filename}.${CURRENT_DATE}"
                mv "$backup_file" "$active_backup/$new_name"
                echo "New version of file: $filename -> $new_name" >> "$BACKUP_REPORT"
                cp -v "$file" "$active_backup/" >> "$BACKUP_REPORT"
            fi
        else
            cp -v "$file" "$active_backup/" >> "$BACKUP_REPORT"
        fi
    done
}

ACTIVE_BACKUP=$(find_active_backup)

if [ -z "$ACTIVE_BACKUP" ]; then
    create_backup_dir
else
    ACTIVE_BACKUP_DATE=$(basename "$ACTIVE_BACKUP" | cut -d'-' -f2)
    ACTIVE_BACKUP_DATE_SECONDS=$(date -d "$ACTIVE_BACKUP_DATE" +%s)
    CURRENT_DATE_SECONDS=$(date -d "$CURRENT_DATE" +%s)
    SEVEN_DAYS_SECONDS=$((7 * 24 * 60 * 60))

    if [ $((CURRENT_DATE_SECONDS - ACTIVE_BACKUP_DATE_SECONDS)) -lt $SEVEN_DAYS_SECONDS ]; then
        update_backup_dir "$ACTIVE_BACKUP"
    else
        create_backup_dir
    fi
fi

echo "Backup is complete."
