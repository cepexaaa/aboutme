#include "mainwindow.h"

#include <QApplication>
#include <QCommandLineParser>
#include <QTranslator>

int main(int argc, char *argv[]) {
  QApplication a(argc, argv);
  QCommandLineParser parser;
  parser.addOption({"dbg", "Enable debug mode with peek button"});
  parser.process(a);
  MainWindow w(parser.isSet("dbg"));
  w.show();
  return a.exec();
}
