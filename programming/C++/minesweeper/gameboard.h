#ifndef GAMEBOARD_H
#define GAMEBOARD_H

#include "cellwidget.h"
#include <QDataStream>
#include <QObject>
#include <vector>

class GameBoard : public QObject {
  Q_OBJECT

public:
  GameBoard(QObject *parent = nullptr);
  ~GameBoard();

  void save(QDataStream &stream) const;
  void load(QDataStream &stream);
  void initializeBoard(int rows, int columns);
  CellWidget *getCell(int row, int column);
  int getRows() const;
  int getColumns() const;

private:
  QVector<QVector<CellWidget *>> board;
  int rows;
  int columns;
};

#endif // GAMEBOARD_H
