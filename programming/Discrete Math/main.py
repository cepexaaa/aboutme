#№1, 2
# f = [int(x) for x in open('task2-72488783264421512023.txt').read().split()]
# n = f.pop(0)
# a = []
# for i in range(n):
#     a.append([f[i*2], f[i*2+1]])
# s = sum([a[i][1] for i in range(n)])
# rez = 0
# for i in range(n):
#     rez += a[i][0] * (a[i][1] / s)
# print("Мат ожидание:", rez)
# rezD = 0
# for i in range(n):
#     rezD += a[i][0]**2 * (a[i][1] / s)
# print("dispersiya:", rezD - rez**2)

#№3
# import math
# f = [int(x) for x in open('task3-63153217817442523443.txt').read().split()]
# n = f.pop(0)
# sum_value = 0
# sm = sum([f[i] for i in range(n)])
# a = []
# for i in range(n):
#     a.append(f[i] / sm)
# for p in a:
#     sum_value += p * math.log(p) / math.log(2)
#
# print(sum_value)

#must: u[i, count/n] :
# a = []
# n = 9
# D = 56.8397
# e = 0
# rez = 0
# for i in range(n):
#     rez += a[i] * (1/n)
# rezD = 0
# for i in range(n):
#     rezD += a[i]**2 * (1/n)
# rez == e and rezD - rez ** 2 == d
# E = 11.8827
# v = 8.96277619962438
# print(8.96277619962438 + 14.802623800375619)
# print(6.042852399248761 + 17.722547600751238)
# print(11.8827 * 2)
# print(2 * E - v)
# print( - E**2)
# print((v**2 + (2 * E - v)**2) - E**2)
# print(D/n)
#n//2 пар таких, что v1 + v2 = 2*E and
# E == среднее арифметическое
# print(2.9 * 8)
# print("step =", 17.722547600751238 - 14.802623800375619)
# step = D / (E + n)
# print(step, "= step my")
# step = 17.722547600751238 - 14.802623800375619
# print('d/step =', D / step)
# print(step, "= step right")
# print("d * e =", D / (E+n-1/n))
# print("d/n =", D / n)
# step = D /(E + n)
# E =
# D =
# N =
# for i in range(n//2 + 1):
#     print(E-step*i)
#     print(E+step*i)

# №4
# print("number 4")
#
# file = open("task4-34340075695174195750(1).txt")
# n, k = map(int,file.readline().split())
# arr=[]
# for i in range(n):
#     arr.append(list(map(int,file.readline().split())))
# # arr = [int(x) for x in open('task3-63153217817442523443.txt').read().split()]
# # n = arr.pop(0)
# # k = arr.pop(0)
# cnt = 0
# for i in range(k):
#     for j in range(i+1,k):
#         E1 = 0
#         E2 = 0
#         E12 = 0
#         for l in range(n):
#             E1+=(arr[l][i])
#             E2+=(arr[l][j])
#             E12+= (arr[l][i] * arr[l][j])
#         if(E1*E2 == E12 * n):
#             cnt +=1
#             if i == 0:
#                 print(j+1)
# print(cnt)
# file.close()
# print("number 5")

#number 5
# import math
#
# n = 9
# E = 11.8827
# D = 56.8397
#
# E_2 = D + E*E
# arr = [E for i in range(n)]
# print(E_2)
# a = math.sqrt((E_2 - E*E)*9/8)
# print(a)
# for i in range(4):
#     arr[i]=E + a
# arr[4]=E
# for i in range(5,n):
#     arr[i]=E - a
# print(arr)
# E_2_new = 0
# E_new = 0
# for i in arr:
#     E_2_new += (1/n)*i*i
#     E_new += (1/n)*i
# print(E_2_new- E_new*E_new)
# print(E_new)
#
# #number 6
# import math
# print("number 6")
# f = 1
# for a in range(1,500):
#     for b in range(1,500-a):
#         for c in range(1,500-a-b):
#             for c in range(1, 500 - a - b):
#                 s = a+b+c
#                 p1 = a/s
#                 p2 = b/s
#                 p3 = c/s
#                 H = -p1*math.log2(p1) -p2*math.log2(p2) -p3*math.log2(p3)
#                 d = H - 0.1361
#                 if(abs(d)<0.00005):
#                     print(H, a, b, c)
#                     f= 0
#                     break


# for k in range(1, 4):
#     for t in range(1, 4):
#         for r in range(1, 4):
#             for m in range(1, 4):
#                 if (r+m==6):
#                     print(-2*t+2*m-3*r+2*k, end="; ")
#                 else:
#                     print(-2*t+2*m-3*r+2*k, end=", ")
# x = 1
# while x != 0:
#     x = int(input())
#     b = int(input())
#     print(2*x+5*b)
