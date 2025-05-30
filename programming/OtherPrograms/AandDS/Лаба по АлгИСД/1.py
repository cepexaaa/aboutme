##def quicksort(array):
##    if len(array) <= 1:
##        return array
##    p = array[len(array) // 2]
##    l, e, g = [], [], []
##    for x in array:
##        if x < p: l.append(x)
##        elif x > p: g.append(x)
##        else: e.append(x)
##    return quicksort(l) + e + quicksort(g)
##
##n = int(input())
##List = list(map(int, input().split()))
##sorted_List = quicksort(List)
##print(*sorted_List)
##def quicksort(a, l, r):
##    if l >= r:
##        return
##    v = a[r]
##    i = l
##    j = r - 1
##    p = l - 1
##    q = r
##    
##    while i <= j:
##        while a[i] < v:
##            i += 1
##        while a[j] > v and j >= l:
##            j -= 1
##        if i >= j:
##            break
##        a[i], a[j] = a[j], a[i]
##        if a[i] == v:
##            p += 1
##            a[p], a[i] = a[i], a[p]
##        i += 1
##        if j >= l and a[j] == v:
##            q -= 1
##            a[q], a[j] = a[j], a[q]
##        j -= 1
##    
##    a[i], a[r] = a[r], a[i]
##    j = i - 1
##    i += 1
##    
##    k = l
##    while k <= p:
##        a[k], a[j] = a[j], a[k]
##        k += 1
##        j -= 1
##    
##    k = r - 1
##    while k >= q:
##        a[k], a[i] = a[i], a[k]
##        k -= 1
##        i += 1
##
##    quicksort(a, l, j)
##    quicksort(a, i, r)
##
##n = int(input())
##array = list(map(int, input().split()))
##quicksort(array, 0, n - 1)
##print(*array)
def quicksort(a, l, r):
    if l >= r:
        return
    v = a[r]
    i = l
    j = r - 1
    p = l - 1
    q = r
    
    while i <= j:
        while i <= j and a[i] < v:
            i += 1
        while j >= l and a[j] > v:
            j -= 1
        if i >= j:
            break
        a[i], a[j] = a[j], a[i]
        if a[i] == v:
            p += 1
            a[p], a[i] = a[i], a[p]
        i += 1
        if j >= l and a[j] == v:
            q -= 1
            a[q], a[j] = a[j], a[q]
        j -= 1
    
    a[i], a[r] = a[r], a[i]
    j = i - 1
    i += 1
    
    k = l
    while k <= p:
        a[k], a[j] = a[j], a[k]
        k += 1
        j -= 1
    
    k = r - 1
    while k >= q:
        a[k], a[i] = a[i], a[k]
        k -= 1
        i += 1

    quicksort(a, l, j)
    quicksort(a, i, r)

n = int(input())
array = list(map(int, input().split()))
quicksort(array, 0, n - 1)
print(*array)
