#include <iostream>
#include <vector>
#include <cmath>

using namespace std;

struct Node {
    int max_val;
    int index;
};

class SegmentTree {
private:
    vector<Node> tree;
    vector<int> arr;
    int n;

    Node combine(Node a, Node b) {
        if (arr[a.index] > arr[b.index]) {
            return a;
        } else {
            return b;
        }
    }

    void build(int v, int tl, int tr) {
        if (tl == tr) {
            tree[v] = {arr[tl], tl};
        } else {
            int tm = (tl + tr) / 2;
            build(v*2, tl, tm);
            build(v*2+1, tm+1, tr);
            tree[v] = combine(tree[v*2], tree[v*2+1]);
        }
    }

    Node query(int v, int tl, int tr, int l, int r) {
        if (l > r) {
            return {INT_MIN, -1};
        }
        if (l == tl && r == tr) {
            return tree[v];
        }
        int tm = (tl + tr) / 2;
        return combine(query(v*2, tl, tm, l, min(r, tm)), query(v*2+1, tm+1, tr, max(l, tm+1), r));
    }

public:
    SegmentTree(vector<int>& a) {
        arr = a;
        n = a.size();
        tree.resize(4*n);
        build(1, 0, n-1);
    }

    pair<int, int> query(int l, int r) {
        Node result = query(1, 0, n-1, l-1, r-1);
        return {result.max_val, result.index+1};
    }
};

int main() {
    int N;
    cin >> N;

    vector<int> arr(N);
    for (int i = 0; i < N; ++i) {
        cin >> arr[i];
    }

    SegmentTree segTree(arr);

    int K;
    cin >> K;

    for (int i = 0; i < K; ++i) {
        int l, r;
        cin >> l >> r;
        pair<int, int> result = segTree.query(l, r);
        cout << result.first << " " << result.second << endl;
    }

    return 0;
}
