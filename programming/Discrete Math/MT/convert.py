import json

def convert_turing_table(input_json, output_file):
    # Загружаем JSON
    data = json.loads(input_json)

    # Открываем файл для записи
    with open(output_file, 'w') as f:
        # Записываем заголовки
        f.write("start: s\n")
        f.write("accept: ac\n")
        f.write("reject: rj\n")
        f.write("blank: _\n")

        # Обрабатываем каждое состояние
        for state_name, state_data in data['states'].items():
            for symbol, action in state_data.items():
                # Пропускаем комментарии
                if symbol == "comment":
                    continue

                # Обрабатываем пустой символ (λ)
                if symbol == "λ":
                    symbol = "_"

                # Разбираем действие
                if isinstance(action, str):
                    parts = action.split()
                    if len(parts) == 1:
                        # Просто движение (R/L) - оставляем символ и состояние прежними
                        move = parts[0]
                        if move == "R":
                            new_state = state_name
                            new_symbol = symbol
                            direction = "->"
                        elif move == "L":
                            new_state = state_name
                            new_symbol = symbol
                            direction = "<-"
                        elif move == "N":
                            new_state = state_name
                            new_symbol = symbol
                            direction = "^"
                    elif len(parts) == 2:
                        # Движение и новое состояние (например "R q1")
                        move = parts[0]
                        new_state = parts[1]
                        new_symbol = symbol
                        if move == "R":
                            direction = "->"
                        elif move == "L":
                            direction = "<-"
                        elif move == "N":
                            direction = "^"
                    elif len(parts) == 3:
                        # Новый символ, движение и состояние (например "1 R q1")
                        new_symbol = parts[0]
                        move = parts[1]
                        new_state = parts[2]
                        if move == "R":
                            direction = "->"
                        elif move == "L":
                            direction = "<-"
                        elif move == "N":
                            direction = "^"
                    else:
                        continue

                    # Заменяем специальные символы для конечных состояний
                    if new_state == "!":
                        new_state = "ac"  # допускающее состояние

                    # Записываем строку
                    f.write(f"{state_name} {symbol} -> {new_state} {new_symbol} {direction}\n")


inoutFile = input()
outputFile = input()
convert_turing_table(inoutFile, outputFile)