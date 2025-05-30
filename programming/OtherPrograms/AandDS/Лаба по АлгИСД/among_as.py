from math import *
t = int(input())
for i in range(t):
    imp, crew = map(int, input().split())
    s = (1 + imp) * imp // 2
    if (s >= crew):
        print("Impostors")
        print(ceil((2*imp + 1 - ((4 * imp**2 + 4 * imp + 1) - 4 * 2 * crew)**0.5) / 2))
        #print(((1 + 8 * crew)**0.5 - 1) // 2)
    else:
        print("Crewmates")
        print(imp)
