#!/usr/bin/bash

first_script="$(pwd)/1_datatime.bash"

(crontab -l; echo "5 * * * 3 $first_script") | crontab -

echo "Crontab successfully."

exit 0
