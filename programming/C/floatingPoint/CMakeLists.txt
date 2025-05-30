cmake_minimum_required(VERSION 3.27)
project(Cpp)

set(CMAKE_CXX_STANDARD 14)

include_directories(ct-c24-lw-floating-point-cepexaaa)

add_executable(Cpp
        ct-c24-lw-floating-point-cepexaaa/.clang-format
        ct-c24-lw-floating-point-cepexaaa/main.c
        ct-c24-lw-floating-point-cepexaaa/return_codes.h
        ct-c24-lw-floating-point-cepexaaa/tests.py)
