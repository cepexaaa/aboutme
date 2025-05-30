# 2 Semestr Algo and Sistem of Data
### 
#### Задание:

Решение:
```

```
### 10А - Пещеры и туннели
#### Задание:
После посадки на Марс учёные нашли странную систему пещер, соединённых туннелями. И учёные начали исследовать эту систему, используя управляемых роботов. Было обнаружено, что существует ровно один путь между каждой парой пещер. Но потом учёные обнаружили специфическую проблему. Иногда в пещерах происходят небольшие взрывы. Они вызывают выброс радиоактивных изотопов и увеличивают уровень радиации в пещере. К сожалению, роботы плохо выдерживают радиацию. Но для исследования они должны переместиться из одной пещеры в другую. Учёные поместили в каждую пещеру сенсор для мониторинга уровня радиации. Теперь они каждый раз при движении робота хотят знать максимальный уровень радиации, с которым придётся столкнуться роботу во время его перемещения. Как вы уже догадались, программу, которая это делает, будете писать вы.
Входные данные

Первая строка входного файла содержит одно целое число n
(1≤n≤100000

) — количество пещер.

Следующие n−1
строк описывают туннели. Каждая из этих строк содержит два целых числа — ai и bi (1≤ai,bi≤N), описывающие туннель из пещеры с номером ai в пещеру с номером bi

.

Следующая строка содержит целое число q
(1≤q≤100000

), означающее количество запросов.

Далее идут q
запросов, по одному на строку. Каждый запрос имеет вид «c u v», где c

— символ «I» либо «G», означающие тип запроса (кавычки только для ясности).

В случае запроса «I» уровень радиации в u
-й пещере (1≤u≤n) увеличивается на v (0≤v≤10000). В случае запроса «G» ваша программа должна вывести максимальный уровень радиации на пути между пещерами с номерами u и v (1≤u,v≤N

) после всех увеличений радиации (запросов «I»), указанных ранее.

Предполагается, что изначальный уровень радиации равен 0

во всех пещерах, и он никогда не уменьшается со временем (потому что период полураспада изотопов много больше времени наблюдения).
Выходные данные

Для каждого запроса «G» выведите одну строку, содержащую максимальный уровень радиации.

Решение:
```
    #include <iostream>
    #include <vector>
    #include <algorithm>
    using namespace std;
     
    const int MAXN = 100005;
    vector<int> lstS[MAXN];
    int parsent[MAXN], deep[MAXN], hvv[MAXN], head[MAXN], pos[MAXN], cur_pos = 0;
    int rid[MAXN];
     
    struct SmTres {
        int tres[4 * MAXN];
     
        void build(int v, int lft, int rght) {
            if (lft == rght) {
                tres[v] = rid[lft];
            } else {
                int tm = (lft + rght) / 2;
                build(v * 2, lft, tm);
                build(v * 2 + 1, tm + 1, rght);
                tres[v] = max(tres[v * 2], tres[v * 2 + 1]);
            }
        }
     
        void update(int v, int lft, int rght, int pos, int new_val) {
            if (lft == rght) {
                tres[v] = new_val;
            } else {
                int tm = (lft + rght) / 2;
                if (pos <= tm)
                    update(v * 2, lft, tm, pos, new_val);
                else
                    update(v * 2 + 1, tm + 1, rght, pos, new_val);
                tres[v] = max(tres[v * 2], tres[v * 2 + 1]);
            }
        }
     
        int query(int v, int lft, int rght, int l, int r) {
            if (l > r)
                return 0;
            if (l == lft && r == rght)
                return tres[v];
            int tm = (lft + rght) / 2;
            return max(query(v * 2, lft, tm, l, min(r, tm)),
                       query(v * 2 + 1, tm + 1, rght, max(l, tm + 1), r));
        }
    } segtree;
     
    int dfs(int v, int p) {
        parsent[v] = p;
        deep[v] = deep[p] + 1;
        int size = 1, Maxl = 0;
        for (int l : lstS[v]) {
            if (l != p) {
                int l_size = dfs(l, v);
                size += l_size;
                if (l_size > Maxl) {
                    Maxl = l_size;
                    hvv[v] = l;
                }
            }
        }
        return size;
    }
     
    void decompose(int v, int h) {
        head[v] = h;
        pos[v] = cur_pos++;
        if (hvv[v] != -1)
            decompose(hvv[v], h);
        for (int l : lstS[v]) {
            if (l != parsent[v] && l != hvv[v])
                decompose(l, l);
        }
    }
     
    int query(int a, int b) {
        int res = 0;
        for (; head[a] != head[b]; b = parsent[head[b]]) {
            if (deep[head[a]] > deep[head[b]])
                swap(a, b);
            int cur_heavy_path_max = segtree.query(1, 0, cur_pos - 1, pos[head[b]], pos[b]);
            res = max(res, cur_heavy_path_max);
        }
        if (deep[a] > deep[b])
            swap(a, b);
        int lhpvm = segtree.query(1, 0, cur_pos - 1, pos[a], pos[b]);
        res = max(res, lhpvm);
        return res;
    }
     
    int main() {
        int n;
        cin >> n;
        for (int i = 0; i < n - 1; ++i) {
            int a, b;
            cin >> a >> b;
            lstS[a].push_back(b);
            lstS[b].push_back(a);
        }
     
        fill(hvv, hvv + MAXN, -1);
        dfs(1, 0);
        decompose(1, 1);
     
        int h;
        cin >> h;
        while (h--) {
            char type;
            int u, v;
            cin >> type >> u >> v;
            if (type == 'I') {
                rid[pos[u]] += v;
                segtree.update(1, 0, cur_pos - 1, pos[u], rid[pos[u]]);
            } else if (type == 'G') {
                cout << query(u, v) << "\n";
            }
        }
     
        return 0;
    }
```
### 11D/1N - Различные числа
#### Задание:
Сколько различных чисел на отрезке массива?
Входные данные

На первой строке длина массива n
(1≤n≤300000). На второй строке n целых чисел от 0 до 109−1 . На третьей строке количество запросов q (1≤q≤300000). Следующие q строк содержат описание запросов, по одному на строке. Каждый запрос задаётся парой целых чисел l,r (1≤l≤r≤n

).
Выходные данные

Выведите ответы на запросы по одному в строке.

Решение:
```
        #include <iostream>
        #include <vector>
        #include <algorithm>
        #include <unordered_map>
         
        using namespace std;
         
        vector<int> a;
        vector<int> tree;
        vector<int> last_occurrence;
         
        void updating(int node, int start, int end, int idx, int val) {
            if (start == end) {
                tree[node] = val;
            } else {
                int mid = (start + end) / 2;
                if (start <= idx && idx <= mid) {
                    updating(2 * node, start, mid, idx, val);
                } else {
                    updating(2 * node + 1, mid + 1, end, idx, val);
                }
                tree[node] = tree[2 * node] + tree[2 * node + 1];
            }
        }
         
        int query(int node, int start, int end, int l, int r) {
            if (r < start || end < l) {
                return 0;
            }
            if (l <= start && end <= r) {
                return tree[node];
            }
            int mid = (start + end) / 2;
            int p1 = query(2 * node, start, mid, l, r);
            int p2 = query(2 * node + 1, mid + 1, end, l, r);
            return p1 + p2;
        }
         
        int main() {
            ios::sync_with_stdio(false);
            cin.tie(nullptr);
         
            int n;
            cin >> n;
            a.resize(n);
            vector<int> sorted_a(n);
            for (int i = 0; i < n; ++i) {
                cin >> a[i];
                sorted_a[i] = a[i];
            }
         
            sort(sorted_a.begin(), sorted_a.end());
            sorted_a.erase(unique(sorted_a.begin(), sorted_a.end()), sorted_a.end());
            unordered_map<int, int> compressed;
            for (int i = 0; i < sorted_a.size(); ++i) {
                compressed[sorted_a[i]] = i;
            }
         
            vector<int> compressed_a(n);
            for (int i = 0; i < n; ++i) {
                compressed_a[i] = compressed[a[i]];
            }
         
            tree.resize(4 * n, 0);
            last_occurrence.resize(sorted_a.size(), -1);
         
            int q;
            cin >> q;
            vector<vector<int>> queries(q, vector<int>(2));
            vector<int> answers(q);
            for (int i = 0; i < q; ++i) {
                cin >> queries[i][0] >> queries[i][1];
                queries[i][0]--;
                queries[i][1]--;
            }
         
            vector<vector<int>> buckets(n);
            for (int i = 0; i < q; ++i) {
                buckets[queries[i][1]].push_back(i);
            }
         
            for (int i = 0; i < n; ++i) {
                if (last_occurrence[compressed_a[i]] != -1) {
                    updating(1, 0, n - 1, last_occurrence[compressed_a[i]], 0);
                }
                last_occurrence[compressed_a[i]] = i;
                updating(1, 0, n - 1, i, 1);
         
                for (int j : buckets[i]) {
                    answers[j] = query(1, 0, n - 1, queries[j][0], queries[j][1]);
                }
            }
         
            for (int i = 0; i < q; ++i) {
                cout << answers[i] << '\n';
            }
         
            return 0;
        }
```

### 8В - Дуумвират
#### Задание:
Вам дано дерево. В вершинах записаны числа. Нужно научиться находить сумму чисел на пути из v в u

.
Входные данные

В первой строке записано число n
— количество вершин дерева (1≤n≤105). Во сторой сроке записаны через пробел n чисел vi (|vi|<109), задающие значения в вершинах. В следующих n−1 строках описаны ребра дерева. В (i+2)-й строке записаны номера вершин ai, bi (1≤ai,bi≤n), означающие, что в дереве есть ребро из вершины ai в вершину bi

.

Далее на отдельной строке записано число m
— количество запросов (1≤m≤105). После этого идут m строк с описанием запросов, в (n+2+i)-й строке записаны через пробел числа xi и yi (1≤xi,yi≤n

).
Выходные данные

Для каждого запроса на отдельной строке требуется вывести сумму всех значений vi
по всем вершинам на пути из xi в yi.

Решение:
```
    #include <iostream>
    #include <vector>
    #include <cmath>
    using namespace std;
     
    struct TresStr {
        int value;
        vector<int> children;
    };
     
    struct Tres {
        int n;
        vector<TresStr> vars;
        vector<int> deep;
        vector<vector<int>> parent;
        vector<long long> sum;
        int logn;
     
        Tres(int n, const vector<int>& vales) : n(n), vars(n), deep(n, 0), sum(n, 0) {
            for (int i = 0; i < n; ++i) {
                vars[i].value = vales[i];
            }
            logn = log2(n) + 1;
            parent.assign(n, vector<int>(logn, -1));
        }
     
        void addEdge(int u, int v) {
            vars[u].children.push_back(v);
            vars[v].children.push_back(u);
        }
     
        void dfs(int node, int par) {
            parent[node][0] = par;
            if (par != -1) {
                deep[node] = deep[par] + 1;
                sum[node] = sum[par] + vars[node].value;
            } else {
                sum[node] = vars[node].value;
            }
     
            for (int child : vars[node].children) {
                if (child != par) {
                    dfs(child, node);
                }
            }
        }
     
        void preprocess() {
            dfs(0, -1);
            for (int j = 1; j < logn; ++j) {
                for (int i = 0; i < n; ++i) {
                    if (parent[i][j - 1] != -1) {
                        parent[i][j] = parent[parent[i][j - 1]][j - 1];
                    }
                }
            }
        }
     
        int getLCA(int u, int v) {
            if (deep[u] < deep[v]) swap(u, v);
     
            for (int i = logn - 1; i >= 0; --i) {
                if (deep[u] - (1 << i) >= deep[v]) {
                    u = parent[u][i];
                }
            }
     
            if (u == v) return u;
     
            for (int i = logn - 1; i >= 0; --i) {
                if (parent[u][i] != parent[v][i]) {
                    u = parent[u][i];
                    v = parent[v][i];
                }
            }
     
            return parent[u][0];
        }
     
        long long query(int u, int v) {
            int lca = getLCA(u, v);
            return sum[u] + sum[v] - 2 * sum[lca] + vars[lca].value;
        }
    };
     
    int main() {
        ios::sync_with_stdio(false);
        cin.tie(0);
     
        int n;
        cin >> n;
        vector<int> vales(n);
        for (int i = 0; i < n; ++i) {
            cin >> vales[i];
        }
     
        Tres Tres(n, vales);
     
        for (int i = 0; i < n - 1; ++i) {
            int u, v;
            cin >> u >> v;
            --u; --v;
            Tres.addEdge(u, v);
        }
     
        Tres.preprocess();
     
        int m;
        cin >> m;
        for (int i = 0; i < m; ++i) {
            int u, v;
            cin >> u >> v;
            --u; --v;
            cout << Tres.query(u, v) << '\n';
        }
     
        return 0;
    }
```

### 8С - Самое дешёвое ребро
#### Задание:
Дано подвешенное дерево с корнем в первой вершине. Все ребра имеют веса (стоимости). Вам нужно ответить на M

запросов вида "найти у двух вершин минимум среди стоимостей ребер пути между ними".
Входные данные

В первой строке файла записано одно числ — n

(количество вершин).

В следующих n−1
строках записаны два числа — x и y. Число x на строке i означает, что x — предок вершины i, y

означает стоимость ребра.

x<i, |y|≤106

.

Далее m
запросов вида (x,y) — найти минимум на пути из x в y (x≠y

).

Ограничения: 2≤n≤5⋅104
, 0≤m≤5⋅104

.
Выходные данные

Выведите m
ответов на запросы.

Решение:
```
    #include <iostream>
    #include <vector>
    #include <cmath>
    #include <climits>
    using namespace std;
     
    struct Edd {
        int now;
        int rub;
        Edd* next;
    };
     
    void addEdd(Edd** head, int now, int rub) {
        Edd* newEdge = new Edd;
        newEdge->now = now;
        newEdge->rub = rub;
        newEdge->next = *head;
        *head = newEdge;
    }
     
    void dfs(int var, int par, Edd** lstS, vector<vector<int>>& verh, vector<vector<int>>& minRub, vector<int>& deep) {
        Edd* current = lstS[var];
        while (current != nullptr) {
            int child = current->now;
            int c = current->rub;
            if (child != par) {
                verh[child][0] = var;
                minRub[child][0] = c;
                deep[child] = deep[var] + 1;
                for (int j = 1; j < verh[0].size(); ++j) {
                    verh[child][j] = verh[verh[child][j - 1]][j - 1];
                    minRub[child][j] = min(minRub[child][j - 1], minRub[verh[child][j - 1]][j - 1]);
                }
                dfs(child, var, lstS, verh, minRub, deep);
            }
            current = current->next;
        }
    }
     
    int lca(int u, int v, vector<vector<int>>& verh, vector<vector<int>>& minRub, vector<int>& deep) {
        if (deep[u] < deep[v]) swap(u, v);
        int diff = deep[u] - deep[v];
        int minC = INT_MAX;
        for (int j = verh[0].size() - 1; j >= 0; --j) {
            if (diff & (1 << j)) {
                minC = min(minC, minRub[u][j]);
                u = verh[u][j];
            }
        }
        if (u == v) return minC;
        for (int j = verh[0].size() - 1; j >= 0; --j) {
            if (verh[u][j] != verh[v][j]) {
                minC = min(minC, min(minRub[u][j], minRub[v][j]));
                u = verh[u][j];
                v = verh[v][j];
            }
        }
        return min(minC, min(minRub[u][0], minRub[v][0]));
    }
     
    int main() {
        int n;
        cin >> n;
        Edd** lstS = new Edd*[n + 1];
        for (int i = 1; i <= n; ++i) {
            lstS[i] = nullptr;
        }
     
        for (int i = 2; i <= n; ++i) {
            int x, y;
            cin >> x >> y;
            addEdd(&lstS[x], i, y);
        }
     
        int logN = log2(n) + 1;
        vector<vector<int>> verh(n + 1, vector<int>(logN, 0));
        vector<vector<int>> minRub(n + 1, vector<int>(logN, INT_MAX));
        vector<int> deep(n + 1, 0);
     
        dfs(1, 0, lstS, verh, minRub, deep);
     
        int m;
        cin >> m;
        for (int i = 0; i < m; ++i) {
            int x, y;
            cin >> x >> y;
            cout << lca(x, y, verh, minRub, deep) << "\n";
        }
     
        for (int i = 1; i <= n; ++i) {
            Edd* current = lstS[i];
            while (current != nullptr) {
                Edd* next = current->next;
                delete current;
                current = next;
            }
        }
        delete[] lstS;
        return 0;
    }
```

### 
#### Задание: 2А - Максимум на подотрезках с добавлением на отрезках
Реализуйте эффективную структуру данных для хранения массива и выполнения следующих операций: увеличение всех элементов данного интервала на одно и то же число; поиск максимума на интервале.
Входные данные

В первой строке вводится одно натуральное число N(1≤N≤100000)

— количество чисел в массиве.

Во второй строке вводятся N
чисел от 0 до 100000

— элементы массива.

В третьей строке вводится одно натуральное число M(1≤M≤30000)

— количество запросов.

Каждая из следующих M
строк представляет собой описание запроса. Сначала вводится одна буква, кодирующая вид запроса (m — найти максимум, a

— увеличить все элементы на отрезке).

Следом за m

вводятся два числа — левая и правая граница отрезка.

Следом за a
вводятся три числа — левый и правый концы отрезка и число add, на которое нужно увеличить все элементы данного отрезка массива (0≤add≤100000)

.
Выходные данные

Выведите в одну строку через пробел ответы на каждый запрос m
.

Решение:
```
    #include <iostream>
    #include <vector>
    #include <cmath>
    #include <algorithm>
     
    using namespace std;
     
    class SegmentTree {
    public:
        SegmentTree(const vector<int>& array) {
            n = array.size();
            int height = (int)ceil(log2(n));
            size = 2 * (1 << height) - 1;
            tree.resize(size, 0);
            lazy.resize(size, 0);
            build(array, 0, n - 1, 0);
        }
     
        void updateRange(int update_left, int update_right, int add_value) {
            updateRangeUtil(0, 0, n - 1, update_left, update_right, add_value);
        }
     
        int queryMax(int query_left, int query_right) {
            return queryMaxUtil(0, 0, n - 1, query_left, query_right);
        }
     
        int n;
        int size;
        vector<int> tree;
        vector<int> lazy;
     
        void build(const vector<int>& array, int left, int right, int node) {
            if (left == right) {
                tree[node] = array[left];
            } else {
                int mid = (left + right) / 2;
                build(array, left, mid, 2 * node + 1);
                build(array, mid + 1, right, 2 * node + 2);
                tree[node] = max(tree[2 * node + 1], tree[2 * node + 2]);
            }
        }
     
        void propagate(int node, int left, int right) {
            if (lazy[node] != 0) {
                tree[node] += lazy[node];
                if (left != right) {
                    lazy[2 * node + 1] += lazy[node];
                    lazy[2 * node + 2] += lazy[node];
                }
                lazy[node] = 0;
            }
        }
     
        void updateRangeUtil(int node, int left, int right, int update_left, int update_right, int add_value) {
            propagate(node, left, right);
            if (update_left > right || update_right < left) {
                return;
            }
            if (update_left <= left && update_right >= right) {
                tree[node] += add_value;
                if (left != right) {
                    lazy[2 * node + 1] += add_value;
                    lazy[2 * node + 2] += add_value;
                }
                return;
            }
            int mid = (left + right) / 2;
            updateRangeUtil(2 * node + 1, left, mid, update_left, update_right, add_value);
            updateRangeUtil(2 * node + 2, mid + 1, right, update_left, update_right, add_value);
            tree[node] = max(tree[2 * node + 1], tree[2 * node + 2]);
        }
     
        int queryMaxUtil(int node, int left, int right, int query_left, int query_right) {
            propagate(node, left, right);
            if (query_left > right || query_right < left) {
                return -1e9;
            }
            if (query_left <= left && query_right >= right) {
                return tree[node];
            }
            int mid = (left + right) / 2;
            int left_max = queryMaxUtil(2 * node + 1, left, mid, query_left, query_right);
            int right_max = queryMaxUtil(2 * node + 2, mid + 1, right, query_left, query_right);
            return max(left_max, right_max);
        }
    };
     
    int main() {
        ios_base::sync_with_stdio(false);
        cin.tie(nullptr);
        int n, m;
        cin >> n;
        vector<int> array(n);
        for (int i = 0; i < n; ++i) {
            cin >> array[i];
        }
        SegmentTree segment_tree(array);
        cin >> m;
        vector<int> results;
        for (int i = 0; i < m; ++i) {
            char query_type;
            cin >> query_type;
            if (query_type == 'm') {
                int left, right;
                cin >> left >> right;
                results.push_back(segment_tree.queryMax(left - 1, right - 1));
            } else if (query_type == 'a') {
                int left, right, add_value;
                cin >> left >> right >> add_value;
                segment_tree.updateRange(left - 1, right - 1, add_value);
            }
        }
        for (int result : results) {
            cout << result << " ";
        }
        cout << endl;
     
        return 0;
    }
```

### 
#### Задание: 1N - Различные числа
Сколько различных чисел на отрезке массива?
Входные данные

На первой строке длина массива n (1 ≤ n ≤ 300 000). На второй строке n целых чисел от 0 до 109. На третьей строке количество запросов q (1 ≤ q ≤ 300 000). Следующие q строк содержат описание запросов, по одному на строке. Каждый запрос задаётся парой целых чисел l, r (1 ≤ l ≤ r ≤ n).
Выходные данные

Выведите ответы на запросы по одному в строке.

Решение:
```
    #include <iostream>
    #include <vector>
    #include <algorithm>
    #include <unordered_map>
     
    using namespace std;
     
    vector<int> a;
    vector<int> tree;
    vector<int> last_occurrence;
     
    void updating(int node, int start, int end, int idx, int val) {
        if (start == end) {
            tree[node] = val;
        } else {
            int mid = (start + end) / 2;
            if (start <= idx && idx <= mid) {
                updating(2 * node, start, mid, idx, val);
            } else {
                updating(2 * node + 1, mid + 1, end, idx, val);
            }
            tree[node] = tree[2 * node] + tree[2 * node + 1];
        }
    }
     
    int query(int node, int start, int end, int l, int r) {
        if (r < start || end < l) {
            return 0;
        }
        if (l <= start && end <= r) {
            return tree[node];
        }
        int mid = (start + end) / 2;
        int p1 = query(2 * node, start, mid, l, r);
        int p2 = query(2 * node + 1, mid + 1, end, l, r);
        return p1 + p2;
    }
     
    int main() {
        ios::sync_with_stdio(false);
        cin.tie(nullptr);
     
        int n;
        cin >> n;
        a.resize(n);
        vector<int> sorted_a(n);
        for (int i = 0; i < n; ++i) {
            cin >> a[i];
            sorted_a[i] = a[i];
        }
     
        sort(sorted_a.begin(), sorted_a.end());
        sorted_a.erase(unique(sorted_a.begin(), sorted_a.end()), sorted_a.end());
        unordered_map<int, int> compressed;
        for (int i = 0; i < sorted_a.size(); ++i) {
            compressed[sorted_a[i]] = i;
        }
     
        vector<int> compressed_a(n);
        for (int i = 0; i < n; ++i) {
            compressed_a[i] = compressed[a[i]];
        }
     
        tree.resize(4 * n, 0);
        last_occurrence.resize(sorted_a.size(), -1);
     
        int q;
        cin >> q;
        vector<vector<int>> queries(q, vector<int>(2));
        vector<int> answers(q);
        for (int i = 0; i < q; ++i) {
            cin >> queries[i][0] >> queries[i][1];
            queries[i][0]--;
            queries[i][1]--;
        }
     
        vector<vector<int>> buckets(n);
        for (int i = 0; i < q; ++i) {
            buckets[queries[i][1]].push_back(i);
        }
     
        for (int i = 0; i < n; ++i) {
            if (last_occurrence[compressed_a[i]] != -1) {
                updating(1, 0, n - 1, last_occurrence[compressed_a[i]], 0);
            }
            last_occurrence[compressed_a[i]] = i;
            updating(1, 0, n - 1, i, 1);
     
            for (int j : buckets[i]) {
                answers[j] = query(1, 0, n - 1, queries[j][0], queries[j][1]);
            }
        }
     
        for (int i = 0; i < q; ++i) {
            cout << answers[i] << '\n';
        }
     
        return 0;
    }
```

### 1М - Звёзды
#### Задание:
Астрономы часто изучают звездные карты, на которых звезды обозначены точками на плоскости, и каждая звезда задана своими декартовыми координатами. Назовем уровнем звезды количество звезд, которые расположены не выше и не правее данной звезды. Астрономы хотят знать распределение уровней звезд.

В качестве примера рассмотрим карту, показанную на рисунке выше. Уровень звезды 5
равен трем (уровень сформирован звездами с номерами 1,2 и 4). Уровни звезд с номерами 2 и 4

равны одному. На данной карте есть единственная звезда с уровнем, равным нулю, две звезды с уровнем, равным одному, одна звезда с уровнем два и одна звезда с уровнем три.

От вас требуется написать программу, которая посчитает количество звезд каждого уровня на заданной карте.
Входные данные

Первая строка содержит число N
— количество звезд на карте (1≤N≤150000)

.

Следующие N
строк описывают координаты звезд. В каждой строке находятся два целых числа, разделенные пробелом — Xi и Yi (0≤Xi,Yi≤200000). В каждой точке плоскости находится не более одной звезды. Звезды перечислены в порядке увеличения координаты Y. Звезды с равными Y-координатами перечислены в порядке увеличения координаты X

.
Выходные данные

Ваша программа должна вывести N
строк. В i-й строке должно быть записано одно целое число — количество звезд уровня i−1.

Решение:
```
    #include <iostream>
    #include <vector>
    #include <algorithm>
    #include <map>
     
    using namespace std;
     
    class FenwickTree {
    public:
        FenwickTree(int size) : tree(size + 1, 0) {}
     
        void update(int index, int value) {
            while (index < tree.size()) {
                tree[index] += value;
                index += index & -index;
            }
        }
     
        int query(int index) {
            int sum = 0;
            while (index > 0) {
                sum += tree[index];
                index -= index & -index;
            }
            return sum;
        }
     
        vector<int> tree;
    };
     
    int main() {
        int N;
        cin >> N;
     
        vector<pair<int, int>> stars(N);
        vector<int> x_coords;
     
        for (int i = 0; i < N; ++i) {
            cin >> stars[i].first >> stars[i].second;
            x_coords.push_back(stars[i].first);
        }
     
        sort(x_coords.begin(), x_coords.end());
        x_coords.erase(unique(x_coords.begin(), x_coords.end()), x_coords.end());
     
        map<int, int> compressed_x;
        for (int i = 0; i < x_coords.size(); ++i) {
            compressed_x[x_coords[i]] = i + 1;
        }
     
        FenwickTree fenwick(x_coords.size());
        vector<int> levels(N, 0);
     
        for (int i = 0; i < N; ++i) {
            int compressed_coord = compressed_x[stars[i].first];
            int level = fenwick.query(compressed_coord);
            levels[level]++;
            fenwick.update(compressed_coord, 1);
        }
     
        for (int i = 0; i < N; ++i) {
            cout << levels[i] << endl;
        }
     
        return 0;
    }
```

### 3А - AVL-дерево
#### Задание:
Данную задачу возможно решать только на языке C++.

В данной задаче от вас требуется реализовать AVL-дерево. Вам дан заголовочный файл avl.h, в котором определена структура вершины дерева, а также описаны сигнатуры функций, которые вы должны реализовать.

В качестве решения вы должны отправить .cpp файл, который будет содержать реализации описанных функций. Вам дан шаблон файла с решением avl.cpp.

Ваше решение не должно содержать функции main. Также ваше решение не должно использовать ввод-вывод. Функции будут вызываться из вспомогательной программы автоматически. При необходимости вы можете реализовать произвольное количество вспомогательных функций, однако решение должно состоять из одного файла.

В любой момент времени вспомогательная программа может проверить текущее дерево на соблюдение инварианта двоичного дерева поиска, а также на соблюдение инварианта AVL-дерева.

Решение:
```
    #include "avl.h"
    #include <algorithm>
     
    size_t height(node* n) {
        return n ? n->h : 0;
    }
     
    void update_height(node* n) {
        if (n) {
            n->h = 1 + std::max(height(n->l), height(n->r));
        }
    }
     
    node* rotate_left(node* n) {
        node* new_root = n->r;
        n->r = new_root->l;
        new_root->l = n;
        update_height(n);
        update_height(new_root);
        return new_root;
    }
     
    node* rotate_right(node* n) {
        node* new_root = n->l;
        n->l = new_root->r;
        new_root->r = n;
        update_height(n);
        update_height(new_root);
        return new_root;
    }
     
    node* balance_tree(node* n) {
        if (!n) return nullptr;
        update_height(n);
        if (height(n->l) - height(n->r) == 2) {
            if (height(n->l->r) > height(n->l->l)) {
                n->l = rotate_left(n->l);
            }
            return rotate_right(n);
        } else if (height(n->r) - height(n->l) == 2) {
            if (height(n->r->l) > height(n->r->r)) {
                n->r = rotate_right(n->r);
            }
            return rotate_left(n);
        }
        return n;
    }
     
    node* insert(node *root, int key) {
        if (!root) {
            node* new_node = new node{nullptr, nullptr, key, 1};
            return new_node;
        }
        if (key < root->key) {
            root->l = insert(root->l, key);
        } else if (key > root->key) {
            root->r = insert(root->r, key);
        }
        return balance_tree(root);
    }
     
    node* min_value_node(node* n) {
        node* current = n;
        while (current->l) {
            current = current->l;
        }
        return current;
    }
     
    node* remove(node *root, int key) {
        if (!root) return nullptr;
        if (key < root->key) {
            root->l = remove(root->l, key);
        } else if (key > root->key) {
            root->r = remove(root->r, key);
        } else {
            if (!root->l || !root->r) {
                node* temp = root->l ? root->l : root->r;
                if (!temp) {
                    temp = root;
                    root = nullptr;
                } else {
                    *root = *temp;
                }
                delete temp;
            } else {
                node* temp = min_value_node(root->r);
                root->key = temp->key;
                root->r = remove(root->r, temp->key);
            }
        }
        return balance_tree(root);
    }
     
    bool exists(node *root, int key) {
        while (root) {
            if (key == root->key) {
                return true;
            } else if (key < root->key) {
                root = root->l;
            } else {
                root = root->r;
            }
        }
        return false;
    }
```

### 1Н - Сережа и скобочки
#### Задание:
У Сережи есть строка s длины n

, состоящая из символов «(» и «)».

Сереже нужно ответить на m
запросов, каждый из которых характеризуется двумя целыми числами li,ri. Ответом на i-ый запрос является длина наибольшей правильной скобочной подпоследовательности последовательности sli,sli+1,…,sri

. Помогите Сереже ответить на все запросы.
Входные данные

Первая строка содержит последовательность символов без пробелов s1,s2,…,sn(1≤n≤106)
. Каждый символ это либо «(», либо «)». Вторая строка содержит целое число m(1≤m≤105) количество запросов. Каждая из следующих m строк содержит пару целых чисел. В i-ой строке записаны числа li,ri, (1≤li≤ri≤n) — описание i

-го запроса.
Выходные данные

Выведите ответ на каждый запрос в отдельной строке. Ответы выводите в порядке следования запросов во входных данных.

Решение:
```
    #include <iostream>
    #include <vector>
    #include <string>
     
    using namespace std;
     
    struct Node {
        int open, close, max_len;
        Node() : open(0), close(0), max_len(0) {}
    };
     
    class SegmentTree {
    public:
        SegmentTree(const string& s) {
            int n = s.size();
            tree.resize(4 * n);
            building(s, 0, 0, n - 1);
        }
     
        int query(int l, int r) {
            Node result = query(0, 0, tree.size() / 4 - 1, l, r);
            return result.max_len;
        }
        
        vector<Node> tree;
     
        Node query(int node, int start, int end, int l, int r) {
            if (l > end || r < start) {
                return Node();
            }
            if (l <= start && r >= end) {
                return tree[node];
            }
            int mid = (start + end) / 2;
            Node left_result = query(2 * node + 1, start, mid, l, r);
            Node right_result = query(2 * node + 2, mid + 1, end, l, r);
            Node result;
            int matched = min(left_result.open, right_result.close);
            result.max_len = left_result.max_len + right_result.max_len + 2 * matched;
            result.open = left_result.open + right_result.open - matched;
            result.close = left_result.close + right_result.close - matched;
            return result;
        }
     
        void building(const string& s, int node, int start, int end) {
            if (start == end) {
                if (s[start] == '(') {
                    tree[node].open = 1;
                } else {
                    tree[node].close = 1;
                }
            } else {
                int mid = (start + end) / 2;
                building(s, 2 * node + 1, start, mid);
                building(s, 2 * node + 2, mid + 1, end);
                merge(node);
            }
        }
     
        void merge(int node) {
            int left = 2 * node + 1;
            int matched = min(tree[left].open, tree[left + 1].close);
            tree[node].max_len = tree[left].max_len + tree[left + 1].max_len + 2 * matched;
            tree[node].open = tree[left].open + tree[left + 1].open - matched;
            tree[node].close = tree[left].close + tree[left + 1].close - matched;
        }
    };
     
    int main() {
        ios::sync_with_stdio(false);
        cin.tie(nullptr);
     
        string s;
        cin >> s;
        int m;
        cin >> m;
        SegmentTree st(s);
        int l, r;
        for (int i = 0; i < m; ++i) {
            cin >> l >> r;
            cout << st.query(l - 1, r - 1) << '\n';
        }
        return 0;
    }
```

### 1G - Противник слаб
#### Задание:
Римляне снова наступают. На этот раз их гораздо больше чем персов, но Шапур готов победить их. Он говорит: «Лев никогда не испугается сотни овец».

Не смотря на это, Шапур должен найти слабость римской армии чтобы победить ее. Как вы помните, Шапур — математик, поэтому он определяет насколько слаба армии как число — степень слабости.

Шапур считает, что степень слабости армии равна количеству таких троек i,j,k
, что i<j<k и ai>aj>ak, где ax  — сила человека, стоящего в строю на месте с номером x

.

Помогите Шапуру узнать, насколько слаба армия римлян.
Входные данные

В первой строке записано одно целое число n
(3≤n≤106) — количество солдат в римской армии. Следующая строка содержит n целых чисел ai(1≤i≤n,1≤ai≤109)

— силы людей в римской армии.
Выходные данные

Выведите одно число — степень слабости римской армии.

Решение:
```
    #include <iostream>
    #include <vector>
    #include <algorithm>
     
    using namespace std;
     
    class FenwickTree {
    public:
        FenwickTree(int size) : tree(size + 1, 0) {}
     
        void update(int index, int value) {
            while (index < tree.size()) {
                tree[index] += value;
                index += index & -index;
            }
        }
     
        int query(int index) {
            int sum = 0;
            while (index > 0) {
                sum += tree[index];
                index -= index & -index;
            }
            return sum;
        }
        vector<int> tree;
    };
     
    int main() {
        ios::sync_with_stdio(false);
        cin.tie(nullptr);
     
        int n, left, right, index;
        long long result = 0;
        cin >> n;
        vector<int> a(n);
        vector<pair<int, int>> elements(n);
     
        for (int i = 0; i < n; ++i) {
            cin >> a[i];
            elements[i] = {a[i], i};
        }
     
        sort(elements.begin(), elements.end(), greater<pair<int, int>>());
     
        FenwickTree leftCount(n);
        FenwickTree rightCount(n);
     
        for (int i = 0; i < n; ++i) {
            rightCount.update(i + 1, 1);
        }
     
        for (const auto& elem : elements) {
            index = elem.second;
            rightCount.update(index + 1, -1);
            left = leftCount.query(index);
            right = rightCount.query(n) - rightCount.query(index + 1);
            result += static_cast<long long>(left) * right;
            leftCount.update(index + 1, 1);
        }
        cout << result << endl;
        return 0;
    }
```

### 1Е - Ближайшее большее число справа
#### Задание:
Дан массив a из n

чисел. Нужно обрабатывать запросы:

    set(i, x) – присвоить новое значение элементу массива a[i] = x;
    get(i, x) – найти mink:k≥i

и ak≥x

    . 

Входные данные

Первая строка входных данных содержит два числа: длину массива n
и количество запросов m (1≤n,m≤200000

).

Во второй строке записаны n
целых чисел – элементы массива a (0≤ai≤200000

).

Следующие m
строк содержат запросы, каждый запрос содержит три числа t, i, x. Первое число t равно 0 или 1 – тип запроса. t=0 означает запрос типа set, t=1 соответствует запросу типа get, 1≤i≤n, 0≤x≤200000

. Элементы массива нумеруются с единицы.
Выходные данные

На каждой запрос типа get на отдельной строке выведите соответствующее значение k
. Если такого k не существует, выведите −1.

Решение:
```
    #include <iostream>
    #include <vector>
    #include <algorithm>
     
    using namespace std;
     
    class SegmentTree {
    private:
        int n;
        vector<int> tree;
     
        void build(const vector<int>& arr, int node, int start, int end) {
            if (start == end) {
                tree[node] = arr[start];
            } else {
                int mid = (start + end) / 2;
                build(arr, 2 * node + 1, start, mid);
                build(arr, 2 * node + 2, mid + 1, end);
                tree[node] = max(tree[2 * node + 1], tree[2 * node + 2]);
            }
        }
     
        void update(int node, int start, int end, int idx, int val) {
            if (start == end) {
                tree[node] = val;
            } else {
                int mid = (start + end) / 2;
                if (start <= idx && idx <= mid) {
                    update(2 * node + 1, start, mid, idx, val);
                } else {
                    update(2 * node + 2, mid + 1, end, idx, val);
                }
                tree[node] = max(tree[2 * node + 1], tree[2 * node + 2]);
            }
        }
     
        int query(int node, int start, int end, int l, int r, int x) {
            if (end < l || start > r) {
                return -1;
            }
            if (start == end) {
                return (tree[node] >= x) ? start : -1;
            }
            if (l <= start && end <= r && tree[node] < x) {
                return -1;
            }
            int mid = (start + end) / 2;
            int left_result = query(2 * node + 1, start, mid, l, r, x);
            if (left_result != -1) {
                return left_result;
            }
            return query(2 * node + 2, mid + 1, end, l, r, x);
        }
     
    public:
        SegmentTree(const vector<int>& arr) {
            n = arr.size();
            tree.resize(4 * n);
            build(arr, 0, 0, n - 1);
        }
     
        void set(int idx, int val) {
            update(0, 0, n - 1, idx, val);
        }
     
        int get(int l, int x) {
            return query(0, 0, n - 1, l, n - 1, x);
        }
    };
     
    int main() {
        ios_base::sync_with_stdio(false);
        cin.tie(nullptr);
     
        int n, m;
        cin >> n >> m;
        vector<int> arr(n);
        for (int i = 0; i < n; ++i) {
            cin >> arr[i];
        }
     
        SegmentTree st(arr);
     
        for (int i = 0; i < m; ++i) {
            int t, idx, x;
            cin >> t >> idx >> x;
            if (t == 0) {
                st.set(idx - 1, x);
            } else {
                int result = st.get(idx - 1, x);
                cout << (result != -1 ? result + 1 : -1) << '\n';
            }
        }
     
        return 0;
    }
```

### 1С - нолики
#### Задание:
Дедус любит давть своим ученикам сложные задачки. На этот раз он придумал такую задачу:

Рейтинг всех его учеников записан в массив A

. Запросы Дедуса таковы:

    Изменить рейтинг i

-го ученика на число x
Найти максимальную последовательность подряд идущих ноликов в массиве A
на отрезке [l,r]

    . 

Помогите бедным фиксикам избежать зверского наказания за нерешение задачи на этот раз.
Входные данные

В первой строке входного файла записано число N
(1≤N≤500000) – количество учеников. Во второй строке записано N чисел – их рейтинги, числа по модулю не превосходящие 1000 (по количеству задач, которые ученик решил или не решил за время обучения). В третьей строке записано число M (1≤M≤50000) – количество запросов. Каждая из следующих M

строк содержит описания запросов:

«UPDATE i x» – обновить i
-ый элемент массива значением x (1≤i≤N, |x|≤1000

)

«QUERY l r» – найти длину максимальной последовательности из нулей на отрезке с l
по r. (1≤l≤r≤N

)
Выходные данные

В выходной файл выведите ответы на запросы «QUERY» в том же порядке, что и во входном файле

Решение:
```
    #include <iostream>
    #include <vector>
    #include <string>
    #include <algorithm>
     
    using namespace std;
     
    struct Node {
        int max_zeros;
        int left_zeros;
        int right_zeros;
    };
     
    vector<Node> tree;
     
    void build(const vector<int>& A, int v, int tl, int tr) {
        if (tl == tr) {
            tree[v] = {A[tl] == 0, A[tl] == 0, A[tl] == 0};
        } else {
            int tm = (tl + tr) / 2;
            build(A, v * 2, tl, tm);
            build(A, v * 2 + 1, tm + 1, tr);
            tree[v].left_zeros = tree[v * 2].left_zeros + (tree[v * 2].left_zeros == (tm - tl + 1) ? tree[v * 2 + 1].left_zeros : 0);
            tree[v].right_zeros = tree[v * 2 + 1].right_zeros + (tree[v * 2 + 1].right_zeros == (tr - tm) ? tree[v * 2].right_zeros : 0);
            tree[v].max_zeros = max({tree[v * 2].max_zeros, tree[v * 2 + 1].max_zeros, tree[v * 2].right_zeros + tree[v * 2 + 1].left_zeros});
        }
    }
     
    void update(int v, int tl, int tr, int pos, int new_val) {
        if (tl == tr) {
            tree[v] = {new_val == 0, new_val == 0, new_val == 0};
        } else {
            int tm = (tl + tr) / 2;
            if (pos <= tm) {
                update(v * 2, tl, tm, pos, new_val);
            } else {
                update(v * 2 + 1, tm + 1, tr, pos, new_val);
            }
            tree[v].left_zeros = tree[v * 2].left_zeros + (tree[v * 2].left_zeros == (tm - tl + 1) ? tree[v * 2 + 1].left_zeros : 0);
            tree[v].right_zeros = tree[v * 2 + 1].right_zeros + (tree[v * 2 + 1].right_zeros == (tr - tm) ? tree[v * 2].right_zeros : 0);
            tree[v].max_zeros = max({tree[v * 2].max_zeros, tree[v * 2 + 1].max_zeros, tree[v * 2].right_zeros + tree[v * 2 + 1].left_zeros});
        }
    }
     
    int query(int v, int tl, int tr, int l, int r) {
        if (l > r) return 0;
        if (l == tl && r == tr) return tree[v].max_zeros;
        int tm = (tl + tr) / 2;
        if (r <= tm) return query(v * 2, tl, tm, l, r);
        if (l > tm) return query(v * 2 + 1, tm + 1, tr, l, r);
        int left_max = query(v * 2, tl, tm, l, tm);
        int right_max = query(v * 2 + 1, tm + 1, tr, tm + 1, r);
        int cross_max = min(tree[v * 2].right_zeros, tm - l + 1) + min(tree[v * 2 + 1].left_zeros, r - tm);
        return max({left_max, right_max, cross_max});
    }
     
    int main() {
        ios::sync_with_stdio(false);
        cin.tie(nullptr);
     
        int N;
        cin >> N;
        vector<int> A(N);
        for (int i = 0; i < N; ++i) {
            cin >> A[i];
        }
     
        tree.resize(4 * N);
        build(A, 1, 0, N - 1);
     
        int M;
        cin >> M;
        for (int m = 0; m < M; ++m) {
            string query1;
            cin >> query1;
            if (query1 == "UPDATE") {
                int i, x;
                cin >> i >> x;
                update(1, 0, N - 1, i - 1, x);
            } else if (query1 == "QUERY") {
                int l, r;
                cin >> l >> r;
                cout << query(1, 0, N - 1, l - 1, r - 1) << "\n";
            }
        }
     
        return 0;
    }
```

### 1L - Инверсии
#### Задание:
Напишите программу, которая для заданного массива A=⟨a1,a2,…,an⟩ находит количество пар (i,j) таких, что i<j и ai>aj

.

Обратите внимание на то, что ответ может не влезать в int.
Входные данные

Первая строка входного файла содержит натуральное число n
(1⩽n⩽100000) — количество элементов массива. Вторая строка содержит n попарно различных элементов массива A — целых неотрицательных чисел, не превосходящих 109

.
Выходные данные

В выходной файл выведите одно число — ответ на задачу.

Решение:
```
    #include <iostream>
    #include <vector>
     
    using namespace std;
     
    long long mergeSort(vector<int>& arr, vector<int>& temp, int left, int right);
    long long merge(vector<int>& arr, vector<int>& temp, int left, int mid, int right);
     
    int main() {
        int n;
        cin >> n;
        vector<int> arr(n);
        for (int i = 0; i < n; i++) {
            cin >> arr[i];
        }
     
        vector<int> temp(n);
        long long result = mergeSort(arr, temp, 0, n - 1);
        cout << result << endl;
     
        return 0;
    }
     
    long long mergeSort(vector<int>& arr, vector<int>& temp, int left, int right) {
        long long inv_count = 0;
        if (right > left) {
            int mid = (right + left) / 2;
            inv_count += mergeSort(arr, temp, left, mid);
            inv_count += mergeSort(arr, temp, mid + 1, right);
            inv_count += merge(arr, temp, left, mid + 1, right);
        }
        return inv_count;
    }
     
    long long merge(vector<int>& arr, vector<int>& temp, int left, int mid, int right) {
        int i = left;
        int j = mid;
        int k = left;
        long long inv_count = 0;
     
        while ((i <= mid - 1) && (j <= right)) {
            if (arr[i] <= arr[j]) {
                temp[k++] = arr[i++];
            } else {
                temp[k++] = arr[j++];
                inv_count += (mid - i);
            }
        }
     
        while (i <= mid - 1) {
            temp[k++] = arr[i++];
        }
     
        while (j <= right) {
            temp[k++] = arr[j++];
        }
     
        for (i = left; i <= right; i++) {
            arr[i] = temp[i];
        }
     
        return inv_count;
    }
```

### 1В - Количество максимумов на отрезке
#### Задание:
Реализуйте структуру данных для эффективного вычисления значения максимального из нескольких подряд идущих элементов массива, а также количества элементов, равных максимальному на данном отрезке.
Входные данные

В первой строке вводится одно натуральное число N
(1≤N≤100000

) — количество чисел в массиве.

Во второй строке вводятся N
чисел от 1 до 100000

— элементы массива.

В третьей строке вводится одно натуральное число K
(1≤K≤30000

) — количество запросов на вычисление максимума.

В следующих K

строках вводится по два числа — номера левого и правого элементов отрезка массива (считается, что элементы массива нумеруются с единицы).
Выходные данные

Для каждого запроса выведите в отдельной строке через пробел значение максимального элемента на указанном отрезке массива и количество максимальных элементов на этом отрезке.

Решение:
```
    #include <iostream>
    #include <vector>
    #include <algorithm>
     
    using namespace std;
     
    class SegmentTree {
      private:
    	int N;
    	vector<pair<int, int>> tree;
     
    	void build(const vector<int>& arr, int v, int l, int r) {
    		if (l == r) {
    			tree[v] = make_pair(arr[l], 1);
    		} else {
    			int m = (l + r) / 2;
    			build(arr, 2 * v + 1, l, m);
    			build(arr, 2 * v + 2, m + 1, r);
    			if (tree[2 * v + 1].first > tree[2 * v + 2].first) {
    				tree[v] = tree[2 * v + 1];
    			} else if (tree[2 * v + 1].first < tree[2 * v + 2].first) {
    				tree[v] = tree[2 * v + 2];
    			} else {
    				tree[v] = make_pair(tree[2 * v + 1].first, tree[2 * v + 1].second + tree[2 * v + 2].second);
    			}
    		}
    	}
     
    	pair<int, int> query(int v, int l, int r, int ql, int qr) {
    		if (ql > qr) {
    			return make_pair(-1e9, 0);
    		}
    		if (l == ql && r == qr) {
    			return tree[v];
    		}
    		int m = (l + r) / 2;
    		auto left = query(2 * v + 1, l, m, ql, min(qr, m));
    		auto right = query(2 * v + 2, m + 1, r, max(ql, m + 1), qr);
    		if (left.first > right.first) {
    			return left;
    		} else if (left.first < right.first) {
    			return right;
    		} else {
    			return make_pair(left.first, left.second + right.second);
    		}
    	}
     
      public:
    	SegmentTree(const vector<int>& arr) {
    		N = arr.size();
    		tree.resize(4 * N);
    		build(arr, 0, 0, N - 1);
    	}
     
    	pair<int, int> query(int l, int r) {
    		return query(0, 0, N - 1, l - 1, r - 1);
    	}
    };
     
    int main() {
    	int N;
    	cin >> N;
    	vector<int> arr(N);
    	for (int i = 0; i < N; ++i) {
    		cin >> arr[i];
    	}
    	SegmentTree tree(arr);
     
    	int K;
    	cin >> K;
    	for (int i = 0; i < K; ++i) {
    		int l, r;
    		cin >> l >> r;
    		auto res = tree.query(l, r);
    		cout << res.first << " " << res.second << endl;
    	}
     
    	return 0;
    }
```

### 1А - RMQ
#### Задание:
Реализуйте структуру данных, которая на данном массиве из N

целых чисел позволяет узнать максимальное значение на этом массиве и индекс элемента, на котором достигается это максимальное значение.
Входные данные

В первой строке вводится натуральное число N
(1≤N≤105) – количество элементов в массиве. В следующей строке содержатся N целых чисел, не превосходящих по модулю 109 – элементы массиваб гарантируется, что в массиве нет одинаковых элементов. Далее идет число K (0≤K≤105) – количество запросов к структуре данных. Каждая из следующих K строк содержит два целых числа l и r (1≤l≤r≤N

) – левую и правую границы отрезка в массиве для данного запроса.
Выходные данные

Для каждого из запросов выведите два числа: наибольшее значение среди элементов массива на отрезке от l
до r и индекс одного из элементов массива, принадлежащий отрезку от l до r, на котором достигается этот максимум.

Решение:
```
    #include <iostream>
    #include <vector>
    #include <algorithm>
     
    using namespace std;
     
    class SegmentTree {
      private:
    	int N;
    	vector<pair<int, int>> tree;
     
    	void build(const vector<int>& arr, int v, int l, int r) {
    		if (l == r) {
    			tree[v] = make_pair(arr[l], l);
    		} else {
    			int m = (l + r) / 2;
    			build(arr, 2 * v + 1, l, m);
    			build(arr, 2 * v + 2, m + 1, r);
    			tree[v] = max(tree[2 * v + 1], tree[2 * v + 2]);
    		}
    	}
     
    	pair<int, int> query(int v, int l, int r, int ql, int qr) {
    		if (ql > qr) {
    			return make_pair(-1e9, -1);
    		}
    		if (l == ql && r == qr) {
    			return tree[v];
    		}
    		int m = (l + r) / 2;
    		auto left = query(2 * v + 1, l, m, ql, min(qr, m));
    		auto right = query(2 * v + 2, m + 1, r, max(ql, m + 1), qr);
    		return max(left, right);
    	}
     
      public:
    	SegmentTree(const vector<int>& arr) {
    		N = arr.size();
    		tree.resize(4 * N);
    		build(arr, 0, 0, N - 1);
    	}
     
    	pair<int, int> query(int l, int r) {
    		return query(0, 0, N - 1, l - 1, r - 1);
    	}
    };
     
    int main() {
    	int N;
    	cin >> N;
    	vector<int> arr(N);
    	for (int i = 0; i < N; ++i) {
    		cin >> arr[i];
    	}
    	SegmentTree tree(arr);
     
    	int K;
    	cin >> K;
    	for (int i = 0; i < K; ++i) {
    		int l, r;
    		cin >> l >> r;
    		auto res = tree.query(l, r);
    		cout << res.first << " " << res.second + 1 << endl;
    	}
     
    	return 0;
    }
```




