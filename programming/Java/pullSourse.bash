#!/bin/bash

cd Java/java-advanced

expect << EOF
spawn git pull
expect {
    "Username for 'https://www.kgeorgiy.info':" {
        send "Kubesh_Sergei\r"
        exp_continue
    }
    "Password for 'https://Kubesh_Sergei@www.kgeorgiy.info'" {
        send "ipovofuravi\r"
        exp_continue
    }
    eof {
        exit 0
    }
    timeout {
        puts "Ошибка: время ожидания истекло."
        exit 1
    }
}
EOF


read -p "Нажмите Enter, чтобы закрыть окно..."

