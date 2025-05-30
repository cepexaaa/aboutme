#!/usr/bin/bash

echo $$ > .pid
echo $((RANDOM % 20 + 1)) > .secret

attempts=10

show() {
    echo "Осталось попыток: $attempts"
}

less() {
    let attempts--
    if [ $attempts -le 0 ]; then
        echo "Вы проиграли. Загаданное число было $(cat .secret)"
        rm .pid .secret
        exit 0
    else
        echo "Загаданное число больше"
    fi
}

more() {
    let attempts--
    if [ $attempts -le 0 ]; then
        echo "Вы проиграли. Загаданное число было $(cat .secret)"
        rm .pid .secret
        exit 0
    else
        echo "Загаданное число меньше"
    fi
}

win() {
    echo "Поздравляем! Вы угадали число $(cat .secret)"
    rm .pid .secret
    exit 0
}

trap 'show' SIGHUP
trap 'less' SIGUSR2
trap 'more' SIGUSR1
trap 'win' SIGCONT
trap 'echo "Обработчик остановлен"; rm .pid .secret; exit 0' SIGTERM

while true; do
    sleep 1
done
exit 0

