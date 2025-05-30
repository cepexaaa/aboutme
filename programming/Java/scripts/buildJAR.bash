#!/bin/bash

JAVA_SOLUTIONS="../java-solutions"
IMPLENEMTOR_PATH="info/kgeorgiy/java/advanced/implementor/"
IMPLER_TEST_PATH="../../java-advanced-2025/modules/info.kgeorgiy.java.advanced.implementor"
SOURCE_IMPLEMENTOR="info/kgeorgiy/ja/kubesh/implementor/"

javac -encoding UTF-8 -d "$JAVA_SOLUTIONS" \
"$JAVA_SOLUTIONS"/"$SOURCE_IMPLEMENTOR"Implementor.java \
-cp "$IMPLER_TEST_PATH":\
"$IMPLER_TEST_PATH".tools

jar cfm "$JAVA_SOLUTIONS"/"$SOURCE_IMPLEMENTOR"Implementor.jar \
./MANIFEST.MF \
-C "$JAVA_SOLUTIONS" "$SOURCE_IMPLEMENTOR"Implementor.class \
-C "$JAVA_SOLUTIONS" "$IMPLENEMTOR_PATH"Impler.class \
-C "$JAVA_SOLUTIONS" "$IMPLENEMTOR_PATH"ImplerException.class \
-C "$JAVA_SOLUTIONS" "$IMPLENEMTOR_PATH"tools/JarImpler.class

java -jar "$JAVA_SOLUTIONS"/"$SOURCE_IMPLEMENTOR"Implementor.jar
