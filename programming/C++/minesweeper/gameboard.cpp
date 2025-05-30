#include "gameboard.h"
#include "cellwidget.h"

GameBoard::GameBoard(QObject *parent) : QObject(parent), rows(0), columns(0) {}

GameBoard::~GameBoard() {
  for (auto &row : board) {
    for (auto &cell : row) {
      delete cell;
    }
  }
  board.clear();
}

void GameBoard::initializeBoard(int rows, int columns) {
  this->rows = rows;
  this->columns = columns;
  board.resize(rows);
  for (auto &row : board) {
    row.resize(columns);
  }
  for (int i = 0; i < rows; ++i) {
    for (int j = 0; j < columns; ++j) {
      if (board[i][j] == nullptr) {
        board[i][j] = new CellWidget();
      }
    }
  }
}

CellWidget *GameBoard::getCell(int row, int column) {
  if (row >= 0 && row < rows && column >= 0 && column < columns) {
    return board[row][column];
  }
  return nullptr;
}

int GameBoard::getRows() const { return rows; }

int GameBoard::getColumns() const { return columns; }
