#ifndef CELLWIDGET_H
#define CELLWIDGET_H

#include <QIcon>
#include <QMouseEvent>
#include <QPixmap>
#include <QWidget>

class CellWidget : public QWidget {
  Q_OBJECT

public:
  explicit CellWidget(QWidget *parent = nullptr);

  enum CellState { Hidden, Revealed, Flagged };

  void setRow(int row);
  void setColumn(int column);
  int getRow() const;
  int getColumn() const;
  void setMine(bool mine);
  void setLastPressed(bool lastPressed);
  bool isLastPressed() const;
  bool isMine() const;
  void setAdjacentMines(int count);
  int adjacentMines() const;
  void setState(CellState state);
  CellState state() const;

protected:
  void paintEvent(QPaintEvent *event) override;
  void mousePressEvent(QMouseEvent *event) override;

signals:
  void cellRevealed(CellWidget *cell);
  void cellFlagged(CellWidget *cell);
  void cellInspected(CellWidget *cell);

private:
  int row;
  int column;
  bool m_isMine;
  int m_adjacentMines;
  bool m_lastPressed = false;
  CellState m_state;
  QPixmap numberPixmaps[13];
};

#endif // CELLWIDGET_H
