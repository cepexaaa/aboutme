#include "mainwindow.h"
#include "cellwidget.h"
#include "gameboard.h"
#include "gamesettingsdialog.h"
#include <QApplication>
#include <QDebug>
#include <QFileDialog>
#include <QInputDialog>
#include <QMessageBox>
#include <QTime>
#include <QTimer>
#include <QTranslator>

MainWindow::MainWindow(bool debugMode, QWidget *parent)
    : QMainWindow(parent), openedCellsCount(0), debugMode(debugMode),
      peekMode(false), playing(true) {

  QWidget *centralWidget = new QWidget(this);
  centralWidget->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
  QVBoxLayout *mainLayout = new QVBoxLayout(centralWidget);

  QHBoxLayout *buttonLayout = new QHBoxLayout();
  settingsButton = new QPushButton(this);
  settingsButton->setIcon(QIcon(":/resourses_picture/New_Game.png"));
  settingsButton->setIconSize(QSize(50, 50));
  settingsButton->setFixedSize(QSize(50, 50));
  connect(settingsButton, &QPushButton::clicked, this,
          &MainWindow::onSettingsButtonClicked);
  buttonLayout->addWidget(settingsButton);

  if (debugMode) {
    peekButton = new QPushButton(this);
    peekButton->setIcon(QIcon(":/resourses_picture/Peek.png"));
    peekButton->setIconSize(QSize(50, 50));
    peekButton->setFixedSize(QSize(50, 50));
    connect(peekButton, &QPushButton::clicked, this,
            &MainWindow::onPeekButtonClicked);
    buttonLayout->addWidget(peekButton);
  }

  mainLayout->addLayout(buttonLayout);

  gridLayout = new QGridLayout();
  gridLayout->setSpacing(0);
  gridLayout->setContentsMargins(0, 0, 0, 0);
  mainLayout->addLayout(gridLayout);

  setCentralWidget(centralWidget);
  loadGameState();

  setMinimumSize(QSize(300, 200));
  setMaximumSize(QSize(1950, 1040));
}

MainWindow::~MainWindow() {
  for (auto &row : cells) {
    for (auto &cell : row) {
      delete cell;
    }
  }
}

void MainWindow::onSettingsButtonClicked() { showGameSettingsDialog(); }

void MainWindow::onPeekButtonClicked() {
  if (playing) {
    togglePeekMode();
  }
}

void MainWindow::togglePeekMode() {
  peekMode = !peekMode;
  cellStates.resize(rows);
  for (int i = 0; i < rows; ++i) {
    cellStates[i].resize(columns);
  }
  for (int i = 0; i < cells.size(); ++i) {
    for (int j = 0; j < cells[i].size(); ++j) {
      if (peekMode) {
        cellStates[i][j] = cells[i][j]->state();
        cells[i][j]->setState(CellWidget::Revealed);
      } else {
        cells[i][j]->setState(cellStates[i][j]);
      }
    }
  }
}

void MainWindow::showGameSettingsDialog() {
  firstReveal = true;
  peekMode = false;
  if (playing) {
    delete gameBoard;
    cleanupUi();
    playing = false;
  }
  openedCellsCount = 0;

  GameSettingsDialog dialog(this);
  connect(&dialog, &GameSettingsDialog::languageChanged, this,
          &MainWindow::changeLanguage);

  if (dialog.exec() == QDialog::Accepted) {
    rows = dialog.getRows();
    columns = dialog.getColumns();
    mines = dialog.getMines();
    gameBoard = new GameBoard(this);
    gameBoard->initializeBoard(rows, columns);

    setupUi();
    playing = true;
  }
}

void MainWindow::changeLanguage(const QString &language) {
  if (translator.load(":/translator/minesweeper_" + language + ".qm")) {
    qApp->installTranslator(&translator);
  }
}

void MainWindow::setupUi() {
  cells.resize(rows);
  for (int i = 0; i < rows; ++i) {
    cells[i].resize(columns);
    for (int j = 0; j < columns; ++j) {
      cells[i][j] = new CellWidget(this);
      gridLayout->addWidget(cells[i][j], i, j);
      connect(cells[i][j], &CellWidget::cellRevealed, this,
              &MainWindow::onCellRevealed);
      connect(cells[i][j], &CellWidget::cellFlagged, this,
              &MainWindow::onCellFlagged);
      connect(cells[i][j], &CellWidget::cellInspected, this,
              &MainWindow::onCellInspected);
    }
  }
  resizeCells();
}

void MainWindow::resizeCells() {
  int cellSize = qMin(width() / columns, height() / rows);
  for (int i = 0; i < rows; ++i) {
    for (int j = 0; j < columns; ++j) {
      cells[i][j]->setFixedSize(cellSize, cellSize);
    }
  }
}

void MainWindow::resizeEvent(QResizeEvent *event) {
  QMainWindow::resizeEvent(event);
  resizeCells();
}

void MainWindow::cleanupUi() {
  for (int i = 0; i < rows; ++i) {
    for (int j = 0; j < columns; ++j) {
      delete cells[i][j];
      cells[i][j] = nullptr;
    }
  }
  cells.clear();
}

void MainWindow::onCellRevealed(CellWidget *cell) {
  if (cell->state() == CellWidget::Revealed) {
    return;
  }
  cell->setState(CellWidget::Revealed);
  if (cell->isMine()) {
    cell->setLastPressed(true);
    QMessageBox::information(this, tr("Game Over"), tr("You hit a mine!"));
    revealAllMines();
  } else {
    cell->update();
    int row = gridLayout->indexOf(cell) / columns;
    int column = gridLayout->indexOf(cell) % columns;
    if (firstReveal) {
      firstReveal = false;
      onFirstCellRevealed(row, column);
    }
    if (cell->adjacentMines() == 0) {
      cell->setState(CellWidget::Hidden);
      revealAdjacentCells(row, column);
      openedCellsCount--;
    }
    openedCellsCount++;
    checkWinCondition(cell);
  }
}

void MainWindow::revealAdjacentCells(int row, int column) {
  if (row < 0 || row >= rows || column < 0 || column >= columns) {
    return;
  }
  CellWidget *cell = cells[row][column];
  if (cell->state() == CellWidget::Revealed || cell->isMine()) {
    return;
  }
  openedCellsCount++;
  cell->setState(CellWidget::Revealed);
  if (cell->adjacentMines() == 0) {
    for (int di = -1; di <= 1; ++di) {
      for (int dj = -1; dj <= 1; ++dj) {
        revealAdjacentCells(row + di, column + dj);
      }
    }
  }
}

void MainWindow::onFirstCellRevealed(int row, int column) {
  generateMines(row, column);
  calculateAdjacentMines();
}

void MainWindow::generateMines(int rowFirst, int columnFirst) {
  QTime time = QTime::currentTime();
  qsrand((uint)time.msec());
  int placedMines = 0;
  while (placedMines < mines) {
    int i = qrand() % rows;
    int j = qrand() % columns;
    if ((!cells[i][j]->isMine()) && (!(i == rowFirst && j == columnFirst))) {
      cells[i][j]->setMine(true);
      ++placedMines;
    }
  }
}

void MainWindow::calculateAdjacentMines() {
  for (int i = 0; i < rows; ++i) {
    for (int j = 0; j < columns; ++j) {
      int count = 0;
      for (int di = -1; di <= 1; ++di) {
        for (int dj = -1; dj <= 1; ++dj) {
          int ni = i + di;
          int nj = j + dj;
          if (ni >= 0 && ni < rows && nj >= 0 && nj < columns &&
              cells[ni][nj]->isMine()) {
            ++count;
          }
        }
      }
      cells[i][j]->setAdjacentMines(count);
    }
  }
}

void MainWindow::checkWinCondition(CellWidget *cell) {
  int cellsToOpen = rows * columns - mines;
  if (openedCellsCount >= cellsToOpen) {
    revealAllMines();
    cell->setLastPressed(true);
    QMessageBox::information(
        this, tr("Victory"),
        tr("You have opened all the cells without mines!"));
  }
}

void MainWindow::revealAllMines() {
  openedCellsCount = rows * columns;
  for (int i = 0; i < rows; ++i) {
    for (int j = 0; j < columns; ++j) {
      cells[i][j]->setState(CellWidget::Revealed);
      cells[i][j]->update();
    }
  }
}

void MainWindow::onCellFlagged(CellWidget *cell) {
  if (cell->state() == CellWidget::Hidden) {
    cell->setState(CellWidget::Flagged);
  } else if (cell->state() == CellWidget::Flagged) {
    cell->setState(CellWidget::Hidden);
  }
  cell->update();
}

void MainWindow::onCellInspected(CellWidget *cell) {
  int row = gridLayout->indexOf(cell) / columns;
  int column = gridLayout->indexOf(cell) % columns;
  int flaggedNeighbors = 0;
  for (int di = -1; di <= 1; ++di) {
    for (int dj = -1; dj <= 1; ++dj) {
      int ni = row + di;
      int nj = column + dj;
      if (ni >= 0 && ni < rows && nj >= 0 && nj < columns) {
        if (cells[ni][nj]->state() == CellWidget::Flagged) {
          ++flaggedNeighbors;
        }
      }
    }
  }

  if (flaggedNeighbors == cell->adjacentMines()) {
    for (int di = -1; di <= 1; ++di) {
      for (int dj = -1; dj <= 1; ++dj) {
        int ni = row + di;
        int nj = column + dj;
        if (ni >= 0 && ni < rows && nj >= 0 && nj < columns) {
          if (cells[ni][nj]->state() == CellWidget::Hidden) {
            onCellRevealed(cells[ni][nj]);
          }
        }
      }
    }
  } else {
    for (int di = -1; di <= 1; ++di) {
      for (int dj = -1; dj <= 1; ++dj) {
        int ni = row + di;
        int nj = column + dj;
        if (ni >= 0 && ni < rows && nj >= 0 && nj < columns) {
          cells[ni][nj]->setLastPressed(true);
          cells[ni][nj]->update();
        }
      }
    }
    QTimer::singleShot(500, this, [this, row, column]() {
      for (int di = -1; di <= 1; ++di) {
        for (int dj = -1; dj <= 1; ++dj) {
          int ni = row + di;
          int nj = column + dj;
          if (ni >= 0 && ni < rows && nj >= 0 && nj < columns) {
            cells[ni][nj]->setLastPressed(false);
            cells[ni][nj]->update();
          }
        }
      }
    });
  }
}

void MainWindow::closeEvent(QCloseEvent *event) {
  if ((openedCellsCount < rows * columns - mines) && (openedCellsCount != 0)) {
    saveGameState();
    QMainWindow::closeEvent(event);
  }
}

void MainWindow::saveGameState() {
  QSettings settings("MyProject", "MyApp");
  settings.beginGroup("GameState");
  settings.setValue("rows", rows);
  settings.setValue("columns", columns);
  settings.setValue("mines", mines);
  settings.setValue("openedCellsCount", openedCellsCount);
  for (int i = 0; i < rows; i++) {
    for (int j = 0; j < columns; j++) {
      settings.setValue(QString("isMine i=%0, j=%1").arg(i).arg(j),
                        cells[i][j]->isMine());
      settings.setValue(QString("state i=%0, j=%1").arg(i).arg(j),
                        static_cast<int>(cells[i][j]->state()));
    }
  }
  settings.endGroup();
}
void MainWindow::loadGameState() {
  QSettings settings("MyProject", "MyApp");
  if (settings.childGroups().contains("GameState")) {
    settings.beginGroup("GameState");
    rows = settings.value("rows", 10).toInt();
    columns = settings.value("columns", 10).toInt();
    mines = settings.value("mines", 10).toInt();
    openedCellsCount = settings.value("openedCellsCount", 0).toInt();
    cells.resize(rows);
    firstReveal = false;
    for (int i = 0; i < rows; i++) {
      cells[i].resize(columns);
      for (int j = 0; j < columns; j++) {
        cells[i][j] = new CellWidget(this);
        gridLayout->addWidget(cells[i][j], i, j);
        cells[i][j]->setMine(
            settings.value(QString("isMine i=%0, j=%1").arg(i).arg(j), false)
                .toBool());
        cells[i][j]->setState(static_cast<CellWidget::CellState>(
            settings
                .value(QString("state i=%0, j=%1").arg(i).arg(j),
                       CellWidget::Hidden)
                .toInt()));
        connect(cells[i][j], &CellWidget::cellRevealed, this,
                &MainWindow::onCellRevealed);
        connect(cells[i][j], &CellWidget::cellFlagged, this,
                &MainWindow::onCellFlagged);
        connect(cells[i][j], &CellWidget::cellInspected, this,
                &MainWindow::onCellInspected);
      }
    }
    settings.remove("");
    settings.endGroup();
    calculateAdjacentMines();
  } else {
    rows = 10;
    columns = 10;
    mines = 10;
    gameBoard = new GameBoard(this);
    gameBoard->initializeBoard(rows, columns);
    firstReveal = true;
    setupUi();
  }
}
