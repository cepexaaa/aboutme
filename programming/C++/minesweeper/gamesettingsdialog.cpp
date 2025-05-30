#include "gamesettingsdialog.h"
#include <QApplication>
#include <QDialogButtonBox>
#include <QFormLayout>
#include <QLocale>
#include <QSpinBox>
#include <QTranslator>

GameSettingsDialog::GameSettingsDialog(QWidget *parent)
    : QDialog(parent), rows(10), columns(10), mines(10) {
  initializeUi();
}

void GameSettingsDialog::initializeUi() {
  QFormLayout *formLayout = new QFormLayout(this);

  rowsSpinBox = new QSpinBox(this);
  rowsSpinBox->setRange(1, 50);
  rowsSpinBox->setValue(rows);
  connect(rowsSpinBox, QOverload<int>::of(&QSpinBox::valueChanged), this,
          &GameSettingsDialog::updateRows);
  connect(rowsSpinBox, QOverload<int>::of(&QSpinBox::valueChanged), this,
          &GameSettingsDialog::updateMinesRange);

  columnsSpinBox = new QSpinBox(this);
  columnsSpinBox->setRange(1, 50);
  columnsSpinBox->setValue(columns);
  connect(columnsSpinBox, QOverload<int>::of(&QSpinBox::valueChanged), this,
          &GameSettingsDialog::updateColumns);
  connect(columnsSpinBox, QOverload<int>::of(&QSpinBox::valueChanged), this,
          &GameSettingsDialog::updateMinesRange);

  minesSpinBox = new QSpinBox(this);
  minesSpinBox->setRange(1, rows * columns - 1);
  minesSpinBox->setValue(mines);
  connect(minesSpinBox, QOverload<int>::of(&QSpinBox::valueChanged), this,
          &GameSettingsDialog::updateMines);
  connect(minesSpinBox, QOverload<int>::of(&QSpinBox::valueChanged), this,
          &GameSettingsDialog::updateMinesRange);

  languageComboBox = new QComboBox(this);
  languageComboBox->addItem("English");
  languageComboBox->addItem("Русский");
  languageComboBox->addItem("Français");
  connect(languageComboBox, QOverload<int>::of(&QComboBox::currentIndexChanged),
          this, &GameSettingsDialog::changeLanguage);

  formLayout->addRow(tr("Rows:"), rowsSpinBox);
  formLayout->addRow(tr("Columns:"), columnsSpinBox);
  formLayout->addRow(tr("Mines:"), minesSpinBox);
  formLayout->addRow(tr("Language:"), languageComboBox);

  QDialogButtonBox *buttonBox = new QDialogButtonBox(
      QDialogButtonBox::Ok | QDialogButtonBox::Cancel, this);
  connect(buttonBox, &QDialogButtonBox::accepted, this, &QDialog::accept);
  connect(buttonBox, &QDialogButtonBox::rejected, this, &QDialog::reject);

  formLayout->addWidget(buttonBox);

  setLayout(formLayout);
}

void GameSettingsDialog::updateMinesRange() {
  int maxMines = rows * columns - 1;
  minesSpinBox->setRange(1, maxMines);
  if (mines > maxMines) {
    mines = maxMines;
    minesSpinBox->setValue(mines);
  }
}

void GameSettingsDialog::updateRows(int value) { rows = value; }

void GameSettingsDialog::updateColumns(int value) { columns = value; }

void GameSettingsDialog::updateMines(int value) { mines = value; }

int GameSettingsDialog::getRows() const { return rows; }

int GameSettingsDialog::getColumns() const { return columns; }

int GameSettingsDialog::getMines() const { return mines; }

void GameSettingsDialog::changeLanguage(int index) {
  QString language = languageComboBox->itemText(index);
  if (language == "Русский") {
    emit languageChanged("ru_RU");
  } else if (language == "English") {
    emit languageChanged("en_EN");
  } else {
    emit languageChanged("fr_FR");
  }
}
