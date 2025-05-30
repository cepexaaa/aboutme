#include <iostream>
#include <vector>

using namespace std;

constexpr long long MOD = 104857601;

vector<long long> multiply_poly(const vector<long long>& a, const vector<long long>& b) {
    vector<long long> res(a.size() + b.size() - 1, 0);
    for (size_t i = 0; i < a.size(); ++i) {
        for (size_t j = 0; j < b.size(); ++j) {
            res[i + j] = (res[i + j] + a[i] * b[j]) % MOD;
        }
    }
    return res;
}

vector<long long> Q_from_C(const vector<long long>& C) {
    vector<long long> Q(C.size() + 1);
    Q[0] = 1;
    for (size_t i = 0; i < C.size(); ++i) {
        Q[i + 1] = (-C[i] + MOD) % MOD;
    }
    return Q;
}

vector<long long> Q_minus_t(const vector<long long>& Q) {
    vector<long long> res = Q;
    for (size_t i = 1; i < res.size(); i += 2) {
        res[i] = (-res[i] + MOD) % MOD;
    }
    return res;
}

vector<long long> sqrt_poly(const vector<long long>& p) {
    vector<long long> res((p.size() + 1) / 2);
    for (size_t i = 0; i < res.size(); ++i) {
        res[i] = p[2 * i];
    }
    return res;
}

long long get_nth(long long n, vector<long long> a, vector<long long> Q) {
    long long k = Q.size() - 1;
    while (n >= k) {
        for (long long i = k; i < 2 * k; ++i) {
            a.push_back(0);
            for (long long j = 1; j <= k; ++j) {
                a[i] = (a[i] - Q[j] * a[i - j] % MOD + MOD) % MOD;
            }
        }

        vector<long long> Q_neg = Q_minus_t(Q);
        vector<long long> R = multiply_poly(Q, Q_neg);

        vector<long long> new_a;
        for (long long i = n % 2; i < 2 * k; i += 2) {
            new_a.push_back(a[i]);
        }
        a = new_a;
        Q = sqrt_poly(R);
        n /= 2;
    }
    return a[n];
}

int main() {
    ios_base::sync_with_stdio(false);
    cin.tie(nullptr);

    long long k, n;
    cin >> k >> n;
    n--;

    vector<long long> A(k);
    for (long long i = 0; i < k; ++i) {
        cin >> A[i];
    }

    vector<long long> C(k);
    for (long long i = 0; i < k; ++i) {
        cin >> C[i];
    }

    vector<long long> Q = Q_from_C(C);
    long long result = get_nth(n, A, Q);

    cout << result << endl;

    return 0;
}