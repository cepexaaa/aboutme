package ru.itmo.wp.service;

import org.springframework.stereotype.Service;
import ru.itmo.wp.domain.TicTacToe;

import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

@Service
public class TicTacToeService {
    private final Map<String, TicTacToe> games = new HashMap<>();

    public TicTacToe startNewGame() {
        TicTacToe game = new TicTacToe();
        game.setId(UUID.randomUUID().toString());
        game.setCurrentPlayer("X"); // X всегда ходит первым
        game.setBoard(new String[3][3]); // Пустое поле
        games.put(game.getId(), game);
        return game;
    }

    public TicTacToe makeMove(String gameId, int row, int col) {
        TicTacToe game = games.get(gameId);
        if (game == null) {
            throw new RuntimeException("Game not found");
        }

        if (game.getBoard()[row][col] != null) {
            throw new RuntimeException("Cell is already occupied");
        }

        game.getBoard()[row][col] = game.getCurrentPlayer();

        if (checkWin(game.getBoard(), game.getCurrentPlayer())) {
            game.setWinner(game.getCurrentPlayer());
        } else if (checkDraw(game.getBoard())) {
            game.setDraw(true);
        } else {
            game.setCurrentPlayer(game.getCurrentPlayer().equals("X") ? "O" : "X");
        }

        return game;
    }

    public TicTacToe getGameStatus(String gameId) {
        return games.get(gameId);
    }

    private boolean checkWin(String[][] board, String player) {
        // Проверка строк, столбцов и диагоналей
        for (int i = 0; i < 3; i++) {
            if (board[i][0] != null && board[i][0].equals(player) && board[i][1].equals(player) && board[i][2].equals(player)) {
                return true; // Проверка строк
            }
            if (board[0][i] != null && board[0][i].equals(player) && board[1][i].equals(player) && board[2][i].equals(player)) {
                return true; // Проверка столбцов
            }
        }
        if (board[0][0] != null && board[0][0].equals(player) && board[1][1].equals(player) && board[2][2].equals(player)) {
            return true; // Главная диагональ
        }
        if (board[0][2] != null && board[0][2].equals(player) && board[1][1].equals(player) && board[2][0].equals(player)) {
            return true; // Побочная диагональ
        }
        return false;
    }

    private boolean checkDraw(String[][] board) {
        for (int i = 0; i < 3; i++) {
            for (int j = 0; j < 3; j++) {
                if (board[i][j] == null) {
                    return false; // Есть пустые клетки
                }
            }
        }
        return true; // Все клетки заполнены
    }
}