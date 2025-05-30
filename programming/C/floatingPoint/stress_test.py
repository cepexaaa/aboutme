from os import system
from sys import argv
from re import sub

N = 20000
OPERATIONS = ['+', '-', '*', '/']

if (len(argv) > 1):
    OPERATIONS = []
    for opeartion in argv[1:]:
        OPERATIONS.append(opeartion)
        if (opeartion != '*' and opeartion != '/' and opeartion != '+' and opeartion != '-'):
            print(f"Invalid operation: {opeartion}")
            exit(1)

system("clang -O2 main.c -o main.exe")
system("clang -O2 tester.c -o tester.exe -lm")
open("out.txt", "w")

for operation in OPERATIONS:
    correct = 0
    for testNumber in range(N):
        system(f"tester.exe '{operation}' > out.txt")

        with open("out.txt", "r") as f:
            test = f.readline().strip()
            expected = f.readline().strip()
            decimal = f.readline().strip()

        system(f"main.exe {test} > out.txt")
        test = test.replace("'", "", 2)
        with open("out.txt", "r") as f:
            actual = f.readline().strip()

        if sub("0+p", "p", actual) != sub("0p", "0.p", expected):
            print()
            print(f"==== test failed => | {testNumber}")
            print(f" ==== expected ==> | {expected}")
            print(f"  === actual ====> | {actual}")
            print(f"   == test case => | {test}")
            if input(f"    = decimal: {decimal}\n") == "q":
                break
        else:
            correct += 1
            if (testNumber % 100 == 0):
                print(f"\n======== ONE HUNDRED MORE TESTS ========\n")
            elif (testNumber % 10 == 0):
                i = testNumber % 3
                print("=" * i + f"===> test passed - {testNumber}")

    print("\n========== ALL TESTS COMPLETED ==========")
    print(f"             operation was {operation}")
    print(f"                 {correct / N * 100}%")
    input(f"           press any button...\n")
