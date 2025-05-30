package ru.itmo.wp.domain;

public class TicTacToe {
    private String id; // Уникальный идентификатор игры
    private String[][] board = new String[3][3]; // Игровое поле 3x3
    private String currentPlayer; // Текущий игрок (X или O)
    private String winner; // Победитель (X, O или null, если игра продолжается)
    private boolean draw; // Ничья

    // Геттеры и сеттеры
    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String[][] getBoard() {
        return board;
    }

    public void setBoard(String[][] board) {
        this.board = board;
    }

    public String getCurrentPlayer() {
        return currentPlayer;
    }

    public void setCurrentPlayer(String currentPlayer) {
        this.currentPlayer = currentPlayer;
    }

    public String getWinner() {
        return winner;
    }

    public void setWinner(String winner) {
        this.winner = winner;
    }

    public boolean isDraw() {
        return draw;
    }

    public void setDraw(boolean draw) {
        this.draw = draw;
    }
}