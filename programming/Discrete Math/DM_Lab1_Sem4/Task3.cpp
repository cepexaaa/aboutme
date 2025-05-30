#include <iostream>
#include <vector>

using namespace std;

vector<int> compute_P(const vector<int>& a, const vector<int>& c, int k) {
    vector P(k, 0);
    for (int n = 0; n < k; ++n) {
        P[n] = a[n];
        for (int i = 1; i <= min(n, k); ++i) {
            if (n - i >= 0) {
                P[n] -= c[i - 1] * a[n - i];
            }
        }
    }
    while (P.size() > 1 && P.back() == 0) {
        P.pop_back();
    }
    return P;
}

vector<int> compute_Q(const vector<int>& c, int k) {
    vector Q(c.size() + 1, 0);
    Q[0] = 1;
    for (int i = 1; i <= k; ++i) {
        Q[i] = -c[i - 1];
    }
    while (Q.size() > 1 && Q.back() == 0) {
        Q.pop_back();
    }
    return Q;
}

void print_poly(const vector<int>& poly) {
    int degree = poly.size() - 1;
    cout << degree << endl;
    for (int i = 0; i <= degree; ++i) {
        cout << poly[i];
        if (i < degree) {
            cout << " ";
        }
    }
    cout << endl;
}

int main() {
    int k;
    cin >> k;

    vector<int> a(k);
    for (int i = 0; i < k; ++i) {
        cin >> a[i];
    }

    vector<int> c(k);
    for (int i = 0; i < k; ++i) {
        cin >> c[i];
    }

    vector<int> P = compute_P(a, c, k);
    vector<int> Q = compute_Q(c, k);

    print_poly(P);
    print_poly(Q);

    return 0;
}