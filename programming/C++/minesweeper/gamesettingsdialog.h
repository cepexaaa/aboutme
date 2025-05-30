#ifndef GAMESETTINGSDIALOG_H
#define GAMESETTINGSDIALOG_H

#include <QComboBox>
#include <QDialog>
#include <QDialogButtonBox>
#include <QFormLayout>
#include <QMessageBox>
#include <QSpinBox>
#include <QTranslator>

class GameSettingsDialog : public QDialog {
  Q_OBJECT

public:
  explicit GameSettingsDialog(QWidget *parent = nullptr);

  int getRows() const;
  int getColumns() const;
  int getMines() const;

signals:
  void languageChanged(const QString &language);

private slots:
  void updateMines(int value);
  void updateRows(int value);
  void updateColumns(int value);
  void updateMinesRange();
  void changeLanguage(int index);

private:
  QComboBox *languageComboBox;
  QTranslator translator;
  QSpinBox *minesSpinBox;
  QSpinBox *rowsSpinBox;
  QSpinBox *columnsSpinBox;
  void initializeUi();
  int rows;
  int columns;
  int mines;
};

#endif // GAMESETTINGSDIALOG_H
