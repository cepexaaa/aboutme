#!/usr/bin/bash

if [ ! -d "$HOME/test" ]; then 
	mkdir "$HOME/test"
fi

current_datetime=$(date +"%Y-%m-%d_%H-%M-%S")

touch "$HOME/test/${current_datetime}"

echo "$(date +"%Y-%m-%d:%H-%M-%S") test was created successfully" >> "$HOME/report"

current_date=$(date +"%Y-%m-%d")

if [ ! -d "$HOME/test/archived" ]; then
	mkdir "$HOME/test/archived"
fi

files_to_archive=""
previous_date=""

for file in "$HOME/test/"*; do
	file_date=$(echo "$file" | grep -oP '\d{4}-\d{2}-\d{2}')

	if [[ "$file_date" < "$current_date" ]]; then
		if [ -z "$previous_date" ]; then
			previous_date="$file_date"
		fi
	
		if [ -f "$file" ] && [[ "$file_date" == "$previous_date" ]]; then
			files_to_archive="$files_to_archive $file"
		fi
	fi
done

#echo "files_to_archive: $files_to_archive"

if [ -n "$files_to_archive" ]; then
	cd "$HOME/test"
	tar -czf "archived/${previous_date}.tar.gz" $files_to_archive
	rm $(basename -a $files_to_archive)
fi

exit 0
