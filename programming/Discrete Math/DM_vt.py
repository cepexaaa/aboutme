# f = [int(i) for i in open('task1-22336033773525194629.txt').read().split()]
# n = f.pop(0)
# a = []
# for i in range(n):
#     a.append([f[i*2], f[i*2+1]])
# sm = sum([a[i][1] for i in range(n)])
# rez = 0
# for i in range(n):
#     rez += a[i][0] * (a[i][1]/sm)
# print(rez)

# f = [int(i) for i in open('task2-32097988411719732113.txt').read().split()]
# n = f.pop(0)
# a = []
# for i in range(n):
#     a.append([f[i*2], f[i*2+1]])
# sm = sum([a[i][1] for i in range(n)])
# rezD = 0
# for i in range(n):
#     rez += a[i][0] * (a[i][1]/sm)
# for i in range(n):
#     rezD += a[i][0]**2 * (a[i][1] / sm)
# print(rezD-rez**2)
# a = [27.23425811063423, 42.194341889365774, 19.754216221268457, 49.674383778731546, 12.274174331902685, 57.15442566809732, 4.794132442536913, 64.6344675574631]
# rez = 0
# sm = 8
# n = 8
# rezD = 0
# # v = []
# # for i in range(n):
# #     v.append(a[i]/8)
# for i in range(n):
#     rez += a[i] * (1/sm)
# for i in range(n):
#     rezD += a[i]**2 * (1/ sm)
# print(rezD-rez**2)


#
# import  math
# # f = [int(i) for i in open('task3-57475313946029606269.txt').read().split()]
# # n = f.pop(0)
# f = [1 ,28 ,104, 210]
# n = 4
# sm = sum([f[i] for i in range(n)])
# a = []
# rez = 0
# for i in range(n):
#     a.append(f[i]/sm)
# for i in range(n):
#     rez += a[i] * math.log2(a[i])
# print(-rez)

# f = open('task4-44679138530509391636.txt', "r")
# n, k = map(int, f.readline().split())
# a = []
# for i in range(n):
#     a.append(list(map(int, f.readline().split())))
# cnt = 0
# for i in range(k):
#     for j in range(i+1, k):
#         E1 = 0
#         E2 = 0
#         E12 = 0
#         for l in range(n):
#             E1 += (a[l][i])
#             E2 += (a[l][j])
#             E12 += (a[l][j] * a[l][i])
#         if (E1 * E2 == E12*n):
#             cnt += 1
#             if i == 0:
#                 print(j + 1)
# print()
# print(cnt)
# f.close()

# n = 8
# E = 34.7143
# D = 419.6327
# step = (D * n / 60)**0.5
# for i in range(1, 1 + n//2):
#     print(E - step * i, 1)
#     print(E + step * i, 1)






# 1 28 104 210
print("a"*1 + "b"*28 + "c"*104 + "d"*210)
print(len("abbbbbbbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")== 1 +28+ 104 +210)












#№1, 2
# f = [int(x) for x in open('task1-22336033773525194629.txt').read().split()]
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


# №4
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


#number 5
# n = 9
# E = 11.8827
# D = 56.8397
# print(sum([i**2 for i in range(n//2+1)])*2)
# x = (D * n / (sum([i**2 for i in range(n//2+1)])*2))**0.5
# # x = (D * n / 60)**0.5
# for i in range(1, n//2 + 1):
#     print(E + i * x)
#     print(E - i * x)


#number 6
# import math
# H = 1.5890
# eps = 0.0001
# n = 333
# for a in range(1, n):
#     for b in range(a + 1, n):
#         for c in range(b + 1, n):
#             for d in range(c + 1, n):
#                 sm = a + b + c + d
#                 fa = a/sm
#                 fb = a/sm
#                 fc = a/sm
#                 fd = a/sm
#                 h = -fa*math.log2(fa) - fb * math.log2(fb) - fb * math.log2(fc) - fb * math.log2(fd)
#                 if abs(H-h) < eps:
#                     print(a, b, c, d)
#                     exit(0)
