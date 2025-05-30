#include <iostream>
#include <vector>
#include <algorithm>
#include <cmath>

using namespace std;

vector<long long> multiply_poly(const vector<long long>& p, const vector<long long>& q) {
    const long long n = p.size(), m = q.size();
    vector<long long> res(n + m - 1, 0);
    for (long long i = 0; i < n; ++i) {
        for (long long j = 0; j < m; ++j) {
            res[i + j] += p[i] * q[j];
        }
    }
    return res;
}

vector<long long> power_poly(const vector<long long>& base, long long exponent) {
    vector<long long> res = {1};
    for (long long i = 0; i < exponent; ++i) {
        res = multiply_poly(res, base);
    }
    return res;
}

vector<long long> compute_P(long long r, long long d, const vector<long long>& p, vector<long long>& Q) {
    // A(t) = sum(p(n)*r^n*t^n)
    vector<long long> A(d + 1);
    for (long long n = 0; n <= d; ++n) {
        long long p_n = 0;
        for (long long i = 0; i <= d; ++i) {
            p_n += p[i] * static_cast<long long>(pow(n, i));
        }
        A[n] = p_n * static_cast<long long>(pow(r, n));
    }
    vector<long long> P = multiply_poly(A, Q);
    while (P.size() > d + 1 || P.size() > 1 && P.back() == 0) {
        P.pop_back();
    }

    return P;
}

vector<long long> compute_Q(long long r, long long d) {
    vector Q_base = {1, -r};
    vector<long long> Q = power_poly(Q_base, d + 1);
    while (Q.size() > 1 && Q.back() == 0) {
        Q.pop_back();
    }
    return Q;
}

void print_poly(const vector<long long>& poly) {
    long long degree = poly.size() - 1;
    cout << degree << endl;
    for (long long i = 0; i <= degree; ++i) {
        cout << poly[i];
        if (i < degree) {
            cout << " ";
        }
    }
    cout << endl;
}

int main() {
    long long r, d;
    cin >> r >> d;

    vector<long long> p(d + 1);
    for (long long i = 0; i <= d; ++i) {
        cin >> p[i];
    }
    vector<long long> Q = compute_Q(r, d);
    const vector<long long> P = compute_P(r, d, p, Q);

    print_poly(P);
    print_poly(Q);

    return 0;
}