#include <iostream>
#include <vector>

using namespace std;

constexpr int MOD = 998244353;

int mod_pow(int a, int b) {
    int res = 1;
    while (b > 0) {
        if (b % 2 == 1) {
            res = (1LL * res * a) % MOD;
        }
        a = (1LL * a * a) % MOD;
        b /= 2;
    }
    return res;
}

int mod_inv(int a) {
    return mod_pow(a, MOD - 2);
}

vector<int> multiply_poly(const vector<int>& p, const vector<int>& q, int m) {
    vector res(m, 0);
    for (int i = 0; i < m; ++i) {
        for (int j = 0; j <= i; ++j) {
            int p_val = (j < p.size()) ? p[j] : 0;
            int q_val = (i - j < q.size()) ? q[i - j] : 0;
            res[i] = (res[i] + 1LL * p_val * q_val) % MOD;
        }
    }
    return res;
}

vector<int> sqrt_series(const vector<int>& p, int m) {
    if (m == 0) return {};

    vector res(m, 0);
    res[0] = 1; // sqrt(1 + 0) = 1
    vector binom_coeff(m, 0);
    binom_coeff[0] = 1; // C(0.5, 0) = 1

    for (int k = 1; k < m; ++k) {
        int numerator = (1 - 2 * (k - 1) + MOD) % MOD; // (0.5 - (k-1)) mod MOD
        int denominator = (2LL * k) % MOD;
        binom_coeff[k] = (1LL * binom_coeff[k-1] * numerator) % MOD;
        binom_coeff[k] = (1LL * binom_coeff[k] * mod_inv(denominator)) % MOD;
    }

    vector<vector<int>> p_powers(m);
    p_powers[0] = vector(m, 0);
    p_powers[0][0] = 1; // P(t)^0 = 1

    if (m > 1) {
        p_powers[1] = p;
        while (p_powers[1].size() < m) p_powers[1].push_back(0);
    }

    for (int k = 2; k < m; ++k) {
        p_powers[k] = multiply_poly(p_powers[k-1], p, m);
    }
    for (int k = 0; k < m; ++k) {
        int coeff = binom_coeff[k];
        for (int i = 0; i < m; ++i) {
            if (i < p_powers[k].size()) {
                res[i] = (res[i] + 1LL * coeff * p_powers[k][i]) % MOD;
            }
        }
    }
    res[0] = 1;
    return res;
}

vector<int> exp_series(const vector<int>& p, int m) {
    vector res(m, 0);
    res[0] = 1;
    vector current_pow = {1};
    int fact = 1;
    for (int k = 1; k < m; ++k) {
        current_pow = multiply_poly(current_pow, p, m);
        fact = 1LL * fact * k % MOD;
        int inv_fact = mod_inv(fact);
        for (int i = 0; i < m; ++i) {
            if (i < current_pow.size()) {
                res[i] = (res[i] + 1LL * current_pow[i] * inv_fact) % MOD;
            }
        }
    }
    return res;
}

vector<int> derivative(const vector<int>& p) {
    vector<int> res(p.size() - 1);
    for (int i = 1; i < p.size(); ++i) {
        res[i - 1] = (1LL * p[i] * i) % MOD;
    }
    return res;
}

vector<int> integrate(const vector<int>& p) {
    vector res(p.size() + 1, 0);
    for (int i = 0; i < p.size(); ++i) {
        res[i + 1] = (1LL * p[i] * mod_inv(i + 1)) % MOD;
    }
    return res;
}

vector<int> ln_series(const vector<int>& p, int m) {
    vector<int> p_deriv = derivative(p);
    vector<int> p_plus_one = p;
    p_plus_one.insert(p_plus_one.begin(), 1);
    vector inv_p_plus_one(m, 0);
    inv_p_plus_one[0] = 1; // 1 / (1 + P(t)) = 1 - P(t) + P(t)^2 - ...
    vector<int> current_pow = p;
    for (int k = 1; k < m; ++k) {
        for (int i = 0; i < m; ++i) {
            if (i < current_pow.size()) {
                inv_p_plus_one[i] = (inv_p_plus_one[i] + ((k % 2 == 1) ? MOD - current_pow[i] : current_pow[i])) % MOD;
            }
        }
        current_pow = multiply_poly(current_pow, p, m);
    }
    vector<int> integrand = multiply_poly(p_deriv, inv_p_plus_one, m);
    vector<int> res = integrate(integrand);
    res.resize(m);
    return res;
}

int main() {
    ios_base::sync_with_stdio(false);
    cin.tie(nullptr);

    int n, m;
    cin >> n >> m;

    vector<int> p(n + 1);
    for (int i = 0; i <= n; ++i) {
        cin >> p[i];
    }

    p[0] = 0;

    vector<int> sqrt_res = sqrt_series(p, m);
    for (int i = 0; i < m; ++i) {
        cout << sqrt_res[i] << (i < m - 1 ? " " : "\n");
    }

    vector<int> exp_res = exp_series(p, m);
    for (int i = 0; i < m; ++i) {
        cout << exp_res[i] << (i < m - 1 ? " " : "\n");
    }

    vector<int> ln_res = ln_series(p, m);
    for (int i = 0; i < m; ++i) {
        cout << ln_res[i] << (i < m - 1 ? " " : "\n");
    }

    return 0;
}

