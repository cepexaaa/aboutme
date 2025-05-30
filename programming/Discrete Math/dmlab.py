
def dfs(u, v, transitions1, transitions2, visited, associations):
    visited[u] = True

    if u in transitions1 and v in transitions2:
        if transitions1[u] != transitions2[v]:
            return False
    else:
        return False

    associations[u] = v
    result = True
    for c, q in transitions1[u].items():
        t1 = q
        t2 = transitions2[v][c]
        if t1 is None or t2 is None:
            return False
        if visited[t1]:
            result = result and t2 == associations[t1]
        else:
            result = result and dfs(t1, t2, transitions1, transitions2, visited, associations)

    return result

with open('problem5.in', 'r') as file:
    n1, m1, k1 = map(int, file.readline().split())
    accepting_states1 = set(map(int, file.readline().split()))
    transitions1 = {}
    for _ in range(m1):
        a, b, c = file.readline().split()
        a, b = int(a), int(b)
        if a not in transitions1:
            transitions1[a] = {}
        transitions1[a][c] = b
    n2, m2, k2 = map(int, file.readline().split())
    accepting_states2 = set(map(int, file.readline().split()))
    transitions2 = {}
    for _ in range(m2):
        a, b, c = file.readline().split()
        a, b = int(a), int(b)
        if a not in transitions2:
            transitions2[a] = {}
        transitions2[a][c] = b

    if n1 != n2 or len(accepting_states1) != len(accepting_states2) or m1 != m2 or k1 != k2:
        print("first")
        result = "NO"

    visited = [False] * (n1 + 1)
    associations = [None] * (n1 + 1)

    if not dfs(1, 1, transitions1, transitions2, visited, associations):
        print("srcond")
        result = "NO"

    for state in range(1, n1 + 1):
        if not visited[state]:
            print("third")
            result = "NO"

    with open('problem5.out', 'w') as fileout:
         print(result)
         fileout.write(str(result))

