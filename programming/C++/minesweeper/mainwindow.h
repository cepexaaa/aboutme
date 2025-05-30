#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include "cellwidget.h"
#include "gameboard.h"
#include <QAction>
#include <QCloseEvent>
#include <QDataStream>
#include <QGridLayout>
#include <QMainWindow>
#include <QMessageBox>
#include <QPushButton>
#include <QResizeEvent>
#include <QSettings>
#include <QTime>
#include <QToolBar>
#include <QTranslator>
#include <QVBoxLayout>
#include <QVector>

class MainWindow : public QMainWindow {
  Q_OBJECT

public:
  explicit MainWindow(bool debugMode, QWidget *parent = nullptr);
  ~MainWindow();
  void loadGame(QDataStream &stream);
  void saveGameState();
  void loadGameState();

private slots:
  void onCellRevealed(CellWidget *cell);
  void revealAdjacentCells(int row, int column);
  void onFirstCellRevealed(int row, int column);
  void onSettingsButtonClicked();
  void onPeekButtonClicked();
  void changeLanguage(const QString &language);

private:
  void setupUi();
  void cleanupUi();
  void generateMines(int rowFirst, int columnFirst);
  void calculateAdjacentMines();
  void showGameSettingsDialog();
  void checkWinCondition(CellWidget *cell);
  void revealAllMines();
  void onCellFlagged(CellWidget *cell);
  void onCellInspected(CellWidget *cell);
  void togglePeekMode();
  void resizeCells();

  QPushButton *saveButton;
  QPushButton *settingsButton;
  QPushButton *peekButton;
  QVBoxLayout *mainLayout;
  QGridLayout *gridLayout;
  QVector<QVector<CellWidget *>> cells;
  GameBoard *gameBoard = nullptr;
  QTranslator translator;
  int rows;
  int columns;
  int mines;
  bool firstReveal = true;
  int openedCellsCount = 0;
  QVector<QVector<CellWidget::CellState>> cellStates;
  bool debugMode;
  bool peekMode;
  bool playing;

protected:
  void closeEvent(QCloseEvent *event) override;
  void resizeEvent(QResizeEvent *event) override;
};

#endif // MAINWINDOW_H
