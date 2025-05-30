#!/usr/bin/bash

at now + 2 minutes <<EOF
./1_datetime.bash
EOF

tail -n 0 -f "$HOME/report"
exit 0
