#!/bin/bash

# shellcheck disable=SC2164
cd /home/sergey/Desktop/Files/CTITMO/JavaAdvanced/Java/java-advanced/java-solutions

javac -d info/kgeorgiy/ja/kubesh/bank/test/out \
  -cp "info/kgeorgiy/ja/kubesh/bank/test:info/kgeorgiy/ja/kubesh/bank/rmi:info/kgeorgiy/ja/kubesh/bank/test/lib/junit-4.13.2.jar" \
  info/kgeorgiy/ja/kubesh/bank/test/*.java \
  info/kgeorgiy/ja/kubesh/bank/rmi/*.java

jar cfm info/kgeorgiy/ja/kubesh/bank/test/bank-tests.jar info/kgeorgiy/ja/kubesh/bank/test/MANIFEST.MF -C info/kgeorgiy/ja/kubesh/bank/test/out .

java -jar info/kgeorgiy/ja/kubesh/bank/test/bank-tests.jar



# javac -cp info/kgeorgiy/ja/kubesh/bank/rmi info/kgeorgiy/ja/kubesh/bank/test/BankTest.java info/kgeorgiy/ja/kubesh/bank/test/PersonTest.java info/kgeorgiy/ja/kubesh/bank/test/LocalTest.java info/kgeorgiy/ja/kubesh/bank/test/RemoteTest.java


# javac -d info/kgeorgiy/ja/kubesh/bank/test/out -cp info/kgeorgiy/ja/kubesh/bank/test:info/kgeorgiy/ja/kubesh/bank/rmi:info/kgeorgiy/ja/kubesh/bank/test/lib/junit-4.13.2.jar info/kgeorgiy/ja/kubesh/bank/test/TestRunner.java


# javac -cp info/kgeorgiy/ja/kubesh/bank/test:info/kgeorgiy/ja/kubesh/bank/rmi:info/kgeorgiy/ja/kubesh/bank/test/lib/junit-4.13.2.jar info/kgeorgiy/ja/kubesh/bank/test/BankTest.java


# чтобы обновить тесты с pom.xml то надо выполнить:
# mvn clean compile test-compile package
# java -jar target/bank-tests-1.0-SNAPSHOT-jar-with-dependencies.jar
