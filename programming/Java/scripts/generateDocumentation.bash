#!/bin/bash

IMPLEMENTOR_PATH="../java-solutions/info/kgeorgiy/ja/kubesh/implementor/Implementor.java"
TEST_PATH="../../java-advanced-2025/modules/info.kgeorgiy.java.advanced.implementor"
IMPLER_PATH="$TEST_PATH/info/kgeorgiy/java/advanced/implementor/Impler.java"
IMPLERJAR_PATH="$TEST_PATH/info/kgeorgiy/java/advanced/implementor/ImplerException.java"
IMPLERECX_PATH="$TEST_PATH.tools/info/kgeorgiy/java/advanced/implementor/tools/JarImpler.java"

OUTPUT_DIR="../javadoc"

mkdir -p $OUTPUT_DIR

javadoc -d "$OUTPUT_DIR" "$IMPLEMENTOR_PATH" "$IMPLER_PATH" "$IMPLERJAR_PATH" "$IMPLERECX_PATH"

xdg-open "$OUTPUT_DIR/index.html"







