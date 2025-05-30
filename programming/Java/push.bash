#!/bin/bash

cd Java/java-advanced

# Получаем текущую дату и время
current_date=$(date +"%Y-%m-%d %H:%M:%S")

# Добавляем все изменения
git add .
if [ $? -ne 0 ]; then
    echo "Ошибка при выполнении 'git add .'"
    exit 1
fi

# Коммитим с текущей датой и временем
git commit -m "Update: $current_date"
if [ $? -ne 0 ]; then
    echo "Ошибка при выполнении 'git commit'."
    exit 1
fi 

# Пушим изменения
expect << EOF
spawn git push
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

if [ $? -ne 0 ]; then
    echo "Ошибка при выполнении 'git push'."
    exit 1
fi

read -p "Нажмите Enter, чтобы закрыть окно..."
