n = 8 # количество лампочек
A = 15 # заданная высота первой лампочки
tolerance = 0.0001 # точность
B = 10 # начальное значение B

def calculate_B(B):
    h = [0] * (n + 1)
    h[1] = A
    h[n] = B

    # итеративно вычисляем высоты лампочек
    for i in range(2, n):
        h[i] = (h[i-1] + h[i+1]) / 2 - 1
    
    # проверяем что все высоты положительные
    for i in range(1, n+1):
        if h[i] <= 0:
            return False

    return True

while True:
    # вычисляем новое значение B
    new_B = B - 0.1
    if calculate_B(new_B):
        B = new_B
    else:
        new_B = B + 0.1
        if calculate_B(new_B):
            B = new_B
        
    # проверяем достижение точности
    if abs(new_B - B) < tolerance:
        break

print("Минимальная высота второго конца гирлянды:", round(B, 2))







"""def binSearch(a, key):  
    l = -1
    r = len(a)    
    while l < r - 1:    
        m = (l + r) // 2
        if a[m] < key:
            l = m
        else: 
            r = m      
    return r
n, m = map(int, input().split())
list_n = list(map(int, input().split()))
list_m = list(map(int, input().split()))
for i in range(len(list_m)):
    element = list_m[i]
    index = binSearch(list_n, element)
    if (n > index) and (element == list_n[index]):
        index1 = index
        while True:
            if (index1 > 0) and list_n[index1-1] == element:
                index1 -= 1
            else:
                break
        print(index1 + 1, end = " ")
        index1 = index
        while True:
            if (index1+1 < n) and list_n[index1+1] == element:
                index1 += 1
            else:
                break
        print(index1 + 1)
    else:
        print(0)
        continue"""


"""n, k = map(int, input().split())
lst1 = sorted(map(int, input().split()))
lst2 = sorted(map(int, input().split()))
i, j = 0, 0
count = 0
while True:
    if lst2[i] == lst1[j]:
        print("YES")
        count += 1
        i += 1
        if i < n:
            continue
        else:
            break
    elif lst2[i] > lst1[j]:
        j+=1
        if j < n:
            continue
        else:
            break
    elif lst2[i] < lst1[j]:
        print("NO")
        count += 1
        i+=1
        continue
for i in range(k-count):
    print("NO")"""
