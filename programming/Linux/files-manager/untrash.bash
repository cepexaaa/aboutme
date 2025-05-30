#!/usr/bin/bash

TRASH_LOG="$HOME/.trash.log"
TRASH_DIR="$HOME/.trash"
RESTORE_NAME="$1"
POLICY="${2:--i}"

if [ -z "$RESTORE_NAME" ]; then
    echo "Using: $0 <name_of_file> [option]"
    echo "Options: -i (--ignore), -u (--unique), -o (--overwrite)"
    exit 1
fi

if [ ! -f "$TRASH_LOG" ]; then
    echo "File trash.log not found. Recovery is not possible."
    exit 1
fi

restore_file() {
    local full_path="$1"
    local link_name="$2"
    local target_dir
    local target_file

    target_dir="$(dirname "$full_path")"
    target_file="$(basename "$full_path")"

    if [ ! -d "$target_dir" ]; then
        echo "The '$target_dir' directory does not exist. Restoring the file to the home directory."
        target_dir="$HOME"
    fi

    local restore_path="$target_dir/$target_file"

    if [ -e "$restore_path" ]; then
        case "$POLICY" in
            -u|--unique)
                local ext="${target_file##*.}"
                local base="${target_file%.*}"
		#add to fix
		if [ "$base" = "$target_file" ]; then
                    base=""
                    ext=""
                fi
		#^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
                local unique_name="$base"
                local counter=1
                while [ -e "$target_dir/$unique_name.$ext" ]; do
		    # add this
		    if [ -z "$ext" ]; then
                        unique_name="${target_file}($counter)"
                    else
                        unique_name="$base($counter).$ext" # add in suffix here
		    fi  
		    counter=$((counter + 1))  
	        done
                restore_path="$target_dir/$unique_name" #delete add in suffix .ext
                ;;
            -o|--overwrite)
                rm -f "$restore_path"
                ;;
            -i|--ignore)
                echo "The file '$restore_path' already exists. Restoring the missing."
                return
                ;;
            *)
                echo "Unknown policy: $POLICY. The default policy is used (-i)."
                echo "The file '$restore_path' already exists. Restoring the missing."
                return
                ;;
        esac
    fi

    ln "$TRASH_DIR/$link_name" "$restore_path" && rm "$TRASH_DIR/$link_name"
    echo "Файл '$target_file' успешно восстановлен в '$restore_path'."
}

grep "/$RESTORE_NAME ->" "$TRASH_LOG" | while read -r line; do
    full_path=$(echo "$line" | awk -F' -> ' '{print $1}')
    link_name=$(echo "$line" | awk -F' -> ' '{print $2}')

    echo "File was found: $full_path"
    read -p "To restore this file? (y/n): " confirm </dev/tty

    if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
        restore_file "$full_path" "$link_name"
    else
        echo "Restoring what was missed."
    fi
done
