# class SegmentTree:
#     def __init__(self, arr):
#         self.N = len(arr)
#         self.tree = [None] * (4 * self.N)
#         self.build(arr, 0, 0, self.N - 1)
#
#     def build(self, arr, v, l, r):
#         if l == r:
#             self.tree[v] = (arr[l], l)
#         else:
#             m = (l + r) // 2
#             self.build(arr, 2 * v + 1, l, m)
#             self.build(arr, 2 * v + 2, m + 1, r)
#             self.tree[v] = max(self.tree[2 * v + 1], self.tree[2 * v + 2])
#
#     def query(self, v, l, r, ql, qr):
#         if ql > qr:
#             return (-float('inf'), -1)
#         if l == ql and r == qr:
#             return self.tree[v]
#         m = (l + r) // 2
#         left = self.query(2 * v + 1, l, m, ql, min(qr, m))
#         right = self.query(2 * v + 2, m + 1, r, max(ql, m + 1), qr)
#         return max(left, right)
#
#
# N = int(input())
# arr = list(map(int, input().split()))
# tree = SegmentTree(arr)
#
# K = int(input())
# for _ in range(K):
#     l, r = map(int, input().split())
#     res = tree.query(0, 0, N - 1, l - 1, r - 1)
#     print(res[0], res[1] + 1)


# import numpy as np
#
# matrix = [
#     [0.00, 0.11, 0.00, 0.03, 0.03, 0.01, 0.02, 0.66, 0.03, 0.11],
#     [0.10, 0.08, 0.11, 0.09, 0.10, 0.12, 0.10, 0.11, 0.08, 0.11],
#     [0.05, 0.13, 0.05, 0.11, 0.09, 0.11, 0.14, 0.14, 0.10, 0.08],
#     [0.10, 0.13, 0.06, 0.09, 0.15, 0.08, 0.10, 0.10, 0.08, 0.11],
#     [0.08, 0.04, 0.11, 0.10, 0.11, 0.12, 0.16, 0.10, 0.09, 0.09],
#     [0.01, 0.00, 0.10, 0.00, 0.66, 0.00, 0.12, 0.02, 0.03, 0.06],
#     [0.14, 0.09, 0.14, 0.05, 0.07, 0.12, 0.08, 0.06, 0.12, 0.13],
#     [0.08, 0.15, 0.15, 0.09, 0.06, 0.10, 0.10, 0.07, 0.07, 0.13],
#     [0.13, 0.14, 0.10, 0.06, 0.13, 0.09, 0.09, 0.07, 0.14, 0.05],
#     [0.00, 0.04, 0.73, 0.05, 0.01, 0.02, 0.10, 0.02, 0.03, 0]
# ]
#
# k = int(input("Введите степень: "))
#
# # Преобразование матрицы в numpy массив
# matrix_np = np.array(matrix)
#
# # Возведение матрицы в степень k
# result_matrix = np.linalg.matrix_power(matrix_np, k)
#
# print("Результат:")
# print(result_matrix)

# import numpy as np
#
# # Ввод матрицы из консоли
# n = int(input("Введите размерность матрицы: "))
# print("Введите элементы матрицы построчно:")
# matrix = []
# for _ in range(n):
#     a = input().replace(",", ".")
#     row = list(map(float, a.split()))
#     matrix.append(row)
#
# k = int(input("Введите степень: "))
#
# # Преобразование матрицы в numpy массив
# matrix_np = np.array(matrix)
#
# # Возведение матрицы в степень k
# result_matrix = np.linalg.matrix_power(matrix_np, k)
#
# print("Результат:")
# print(result_matrix)


#Number 5
# import numpy as np
#
# # Ввод матрицы
# n = int(input("Введите размерность матрицы: "))
# matrix = np.zeros((n, n))
# print("Введите элементы матрицы построчно:")
# for i in range(n):
#     matrix[i] = list(map(float, input().replace(",", ".").split()))
#
# # Функция для умножения матриц
# def multiply_matrices(matrix1, matrix2):
#     return np.dot(matrix1, matrix2)
#
# # Возводим матрицу в степень k
# k = 321633955
# result_matrix = np.linalg.matrix_power(matrix, k)
#
# print(result_matrix)



# st = input()
# n = int(input())
# for i in range(n):
#     l, r = map(int, input().split())
#     ps = 0
#     pl = 0
#     res = 0
#     sts = st[l-1:r]
#     while "()" in sts:
#         sts = sts.replace("()", "", 1)
#         res+=2
#     # for j in range(l-1, r):
#     #     if st[j] == "(":
#     #         ps+=1
#     #     else:
#     #         ps-=1
#     #     if ps < 0:
#     #         pl, ps = 0, 0
#     #         continue
#     #     pl += 1
#     #     if ps == 0 and (j < (r - 1) and st[j+1] != "(" or j == r - 1):
#     #         res += pl #max(res, pl)
#
#     print(res)




import math

# Функция для построения дерева отрезков
def build_segment_tree(arr, n):
    height = math.ceil(math.log2(n))
    size = 2 * (2 ** height) - 1
    segment_tree = [float('inf')] * size
    return segment_tree

# Функция для обновления дерева отрезков
def update_segment_tree(segment_tree, index, value, node, start, end):
    if start == end:
        segment_tree[node] = value
    else:
        mid = (start + end) // 2
        if start <= index <= mid:
            update_segment_tree(segment_tree, index, value, 2 * node + 1, start, mid)
        else:
            update_segment_tree(segment_tree, index, value, 2 * node + 2, mid + 1, end)
        segment_tree[node] = min(segment_tree[2 * node + 1], segment_tree[2 * node + 2])

def query_segment_tree(segment_tree, query_start, query_end, node, start, end):
    if query_start <= start and query_end >= end:
        return segment_tree[node]
    if end < query_start or start > query_end:
        return float('inf')
    mid = (start + end) // 2
    return min(query_segment_tree(segment_tree, query_start, query_end, 2 * node + 1, start, mid),
               query_segment_tree(segment_tree, query_start, query_end, 2 * node + 2, mid + 1, end))


n = int(input())
arr = []
segment_tree = build_segment_tree(arr, n)
for _ in range(n):
    operation = input().split()
    if operation[0] == '+':
        i, x = int(operation[1]), int(operation[2])
        arr.insert(i, x)
        update_segment_tree(segment_tree, i, x, 0, 0, len(arr) - 1)
    elif operation[0] == '?':
        i, j = int(operation[1]), int(operation[2])
        result = query_segment_tree(segment_tree, i, j, 0, 0, len(arr) - 1)
        print(result)

