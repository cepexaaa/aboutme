n, m, p = map(int, input().split())
lst = list(map(int, input().split()))
count = 0
for i in range(n):
    if lst[i] > p:
        a = lst[i]
        while a > p:
            count+=1
            a-=p
        lst[i] = a
if count < m:
    print(p * m)
else:
    print(count * p + sum(sorted(lst, reverse=True)[:m - count]))









##n, T = map(int, input().split())
##a = list(map(int, input().split()))
##f = [0] * (200000 + 13)
##
##def upd(x):
##    i = x
##    while i < len(f):
##        f[i] += 1
##        i |= i + 1
##
##def get(x):
##    res = 0
##    i = x
##    while i >= 0:
##        res += f[i]
##        i = (i & (i + 1)) - 1
##    return res
##
##sums = [0]
##pr = 0
##for i in range(n):
##    pr += a[i]
##    sums.append(pr)
##
##sums.sort()
##i = 1
##while i < len(sums):
##    if sums[i] == sums[i-1]:
##        sums.pop(i)
##    else:
##        i += 1
##
##ans = 0
##pr = 0
##idx = sums.index(0)
##upd(idx)
##for i in range(n):
##    pr += a[i]
##    npos = bisect_right(sums, pr - T) - 1
##    ans += (i + 1 - get(npos))
##    k = bisect_left(sums, pr)
##    upd(k)
##
##print(ans)




















####def count_pairs(A):
####    count = 0
####    stack = []
####    
####    for num in A:
####        while stack and num >= stack[-1]:
####            stack.pop()
####            count += len(stack)
####        stack.append(num)
####    
####    return count
####
##### Чтение входных данных
####n = int(input())
####A = list(map(int, input().split()))
####
##### Вызов функции и вывод результатов
####result = count_pairs(A)
####print(result)
##def count_inversions(arr):
##    n = len(arr)
##
##    # Создаем пустые словари для отслеживания количества инверсий
##    smaller = {}
##    greater = {}
##
##    count = 0
##
##    # Проходим по массиву справа налево
##    for i in range(n-1, -1, -1):
##        num = arr[i]
##
##        # Увеличиваем счетчик инверсий, считая количество элементов в greater,
##        # соответствующих условию ai>aj, где i<j
##        count += len(smaller.get(num, []))
##
##        # Обновляем словари smaller и greater, добавляя текущий элемент num
##        if num in smaller:
##            smaller[num].append(i)
##        else:
##            smaller[num] = [i]
##
##        if num in greater:
##            for j in greater[num]:
##                count += len(smaller.get(j, []))
##            greater[num].append(i)
##        else:
##            greater[num] = [i]
##
##    return count
##
### Считываем входные данные
##n = int(input())
##arr = list(map(int, input().split()))
##
### Вызываем функцию и выводим результат
##result = count_inversions(arr)
##print(result)
##def merge(left, right):
##    return sorted(left + right)
##
##
##def recursion(arr):
##    if len(arr) <= 1:
##        return 0, arr
##    half = len(arr) // 2
##    left_ans, left = recursion(arr[:half])
##    right_ans, right = recursion(arr[half:])
##    cross_ans = 0
##    j = 0
##    for i in range(len(left)):
##        while j < len(right) and left[i] >= 2 * right[j]:
##            j += 1
##        cross_ans += j
##    return left_ans + cross_ans + right_ans, merge(left, right)
##
##n = int(input())
##arr = list(map(int, input().split()))
##print(recursion(arr)[0])
