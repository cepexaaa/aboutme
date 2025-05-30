#include "cellwidget.h"
#include <QMouseEvent>
#include <QPainter>
#include <QPushButton>

CellWidget::CellWidget(QWidget *parent)
    : QWidget(parent), m_isMine(false), m_adjacentMines(0), m_state(Hidden) {
  for (int i = 0; i < 9; ++i) {
    numberPixmaps[i] =
        QPixmap(QString(":/resourses_picture/number_%1.png").arg(i));
  }
  numberPixmaps[10] = QPixmap(":/resourses_picture/unopened_square.png");
  numberPixmaps[11] = QPixmap(":/resourses_picture/flag.png");
  numberPixmaps[9] = QPixmap(":/resourses_picture/mine.png");
  numberPixmaps[12] = QPixmap(":/resourses_picture/mine_lastPressed.png");
  setMinimumSize(40, 40);
}

void CellWidget::setMine(bool mine) { m_isMine = mine; }

bool CellWidget::isMine() const { return m_isMine; }

void CellWidget::setAdjacentMines(int count) { m_adjacentMines = count; }

int CellWidget::adjacentMines() const { return m_adjacentMines; }

void CellWidget::setState(CellState state) {
  m_state = state;
  update();
}

CellWidget::CellState CellWidget::state() const { return m_state; }

void CellWidget::setLastPressed(bool lastPressed) {
  m_lastPressed = lastPressed;
  update();
}

bool CellWidget::isLastPressed() const { return m_lastPressed; }
void CellWidget::setRow(int row) { this->row = row; }

void CellWidget::setColumn(int column) { this->column = column; }

int CellWidget::getRow() const { return row; }

int CellWidget::getColumn() const { return column; }

void CellWidget::paintEvent(QPaintEvent *event) {
  QPainter painter(this);
  if (m_state == Revealed) {
    if (m_isMine) {
      if (m_lastPressed) {
        painter.drawPixmap(rect(), numberPixmaps[12]);
      } else {
        painter.drawPixmap(rect(), numberPixmaps[9]);
      }
    } else {
      painter.drawPixmap(rect(), numberPixmaps[m_adjacentMines]);
    }
  } else if (m_state == Flagged) {
    painter.drawPixmap(rect(), numberPixmaps[11]);
  } else {
    painter.drawPixmap(rect(), numberPixmaps[10]);
  }
}

void CellWidget::mousePressEvent(QMouseEvent *event) {
  if (event->button() == Qt::LeftButton) {
    emit cellRevealed(this);
  } else if (event->button() == Qt::RightButton) {
    emit cellFlagged(this);
  } else if (event->button() == Qt::MiddleButton) {
    emit cellInspected(this);
  }
  QWidget::mousePressEvent(event);
}
