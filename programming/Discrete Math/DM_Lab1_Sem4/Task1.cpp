#include <iostream>
#include <vector>
#include <algorithm>

using namespace std;

constexpr int MOD = 998244353;

vector<int> add_poly(const vector<int>& p, const vector<int>& q) {
    const int n = max(p.size(), q.size());
    vector<int> res(n);
    for (int i = 0; i < n; ++i) {
        const int a = (i < p.size()) ? p[i] : 0;
        const int b = (i < q.size()) ? q[i] : 0;
        res[i] = (a + b) % MOD;
    }
    while (res.size() > 1 && res.back() == 0) res.pop_back();
    return res;
}

vector<int> multiply_poly(const vector<int>& p, const vector<int>& q) {
    const int n = p.size(), m = q.size();
    vector res(n + m - 1, 0);
    for (int i = 0; i < n; ++i) {
        for (int j = 0; j < m; ++j) {
            res[i + j] = (res[i + j] + 1LL * p[i] * q[j]) % MOD;
        }
    }
    while (res.size() > 1 && res.back() == 0) res.pop_back();
    return res;
}

vector<int> inverse_series(const vector<int>& q, int k) {
    vector inv(k, 0);
    inv[0] = 1; // since q[0] = 1
    for (int i = 1; i < k; ++i) {
        for (int j = 1; j <= min(i, (int)q.size() - 1); ++j) {
            inv[i] = (inv[i] - 1LL * q[j] * inv[i - j] % MOD + MOD) % MOD;
        }
    }
    return inv;
}

vector<int> divide_series(const vector<int>& p, const vector<int>& q, int k) {
    vector res(k, 0);
    vector<int> inv_q = inverse_series(q, k);
    for (int i = 0; i < k; ++i) {
        for (int j = 0; j <= i; ++j) {
            int p_val = (j < p.size()) ? p[j] : 0;
            res[i] = (res[i] + 1LL * p_val * inv_q[i - j]) % MOD;
        }
    }
    return res;
}

void print_poly(const vector<int>& p) {
    int degree = p.size() - 1;
    while (degree > 0 && p[degree] == 0) degree--;
    cout << degree << endl;
    for (int i = 0; i <= degree; ++i) {
        cout << p[i] << (i < degree ? " " : "\n");
    }
}

int main() {
    ios_base::sync_with_stdio(false);
    cin.tie(nullptr);

    int n, m;
    cin >> n >> m;

    vector<int> p(n + 1), q(m + 1);
    for (int i = 0; i <= n; ++i) cin >> p[i];
    for (int i = 0; i <= m; ++i) cin >> q[i];

    // P + Q
    vector<int> sum = add_poly(p, q);
    print_poly(sum);

    // P * Q
    vector<int> product = multiply_poly(p, q);
    print_poly(product);

    // P / Q
    vector<int> division = divide_series(p, q, 1000);
    for (int i = 0; i < 1000; ++i) {
        cout << division[i] << (i < 999 ? " " : "\n");
    }

    return 0;
}



/// it is task 2

// #include <iostream>
// #include <vector>
//
// using namespace std;
//
// const int MOD = 998244353;
//
// // Быстрое возведение в степень по модулю
// int mod_pow(int a, int b) {
//     int res = 1;
//     while (b > 0) {
//         if (b % 2 == 1) {
//             res = (1LL * res * a) % MOD;
//         }
//         a = (1LL * a * a) % MOD;
//         b /= 2;
//     }
//     return res;
// }
//
// // Обратный элемент по модулю
// int mod_inv(int a) {
//     return mod_pow(a, MOD - 2);
// }
//
// // Умножение многочленов по модулю
// vector<int> multiply_poly(const vector<int>& p, const vector<int>& q, int m) {
//     vector<int> res(m, 0);
//     for (int i = 0; i < m; ++i) {
//         for (int j = 0; j <= i; ++j) {
//             int p_val = (j < p.size()) ? p[j] : 0;
//             int q_val = (i - j < q.size()) ? q[i - j] : 0;
//             res[i] = (res[i] + 1LL * p_val * q_val) % MOD;
//         }
//     }
//     return res;
// }
//
// // Вычисление производной многочлена
// vector<int> derivative(const vector<int>& p) {
//     vector<int> res(p.size() - 1);
//     for (int i = 1; i < p.size(); ++i) {
//         res[i - 1] = (1LL * p[i] * i) % MOD;
//     }
//     return res;
// }
//
// // Интегрирование многочлена (без константы)
// vector<int> integrate(const vector<int>& p) {
//     vector<int> res(p.size() + 1, 0);
//     for (int i = 0; i < p.size(); ++i) {
//         res[i + 1] = (1LL * p[i] * mod_inv(i + 1)) % MOD;
//     }
//     return res;
// }
//
// // Вычисление sqrt(1 + P(t)) до m коэффициентов
// vector<int> sqrt_series(const vector<int>& p, int m) {
//     vector<int> res(m, 0);
//     res[0] = 1; // sqrt(1 + 0) = 1
//     vector<int> current_pow = {1}; // (1 + P(t))^{1/2}
//     for (int k = 1; k < m; ++k) {
//         vector<int> term = multiply_poly(current_pow, p, m);
//         int coeff = (1LL * (MOD - 1) * mod_inv(2 * k)) % MOD;
//         coeff = (1LL * coeff * (2 * k - 3)) % MOD;
//         coeff = (1LL * coeff * mod_inv(k)) % MOD;
//         for (int i = 0; i < m; ++i) {
//             if (i < term.size()) {
//                 res[i] = (res[i] + 1LL * coeff * term[i]) % MOD;
//             }
//         }
//         current_pow = multiply_poly(current_pow, p, m);
//     }
//     return res;
// }
//
// // Вычисление exp(P(t)) до m коэффициентов
// vector<int> exp_series(const vector<int>& p, int m) {
//     vector<int> res(m, 0);
//     res[0] = 1; // exp(0) = 1
//     vector<int> current_pow = {1}; // P(t)^k / k!
//     for (int k = 1; k < m; ++k) {
//         current_pow = multiply_poly(current_pow, p, m);
//         int inv_k = mod_inv(k);
//         for (int i = 0; i < m; ++i) {
//             if (i < current_pow.size()) {
//                 res[i] = (res[i] + 1LL * current_pow[i] * inv_k) % MOD;
//             }
//         }
//     }
//     return res;
// }
//
// // Вычисление ln(1 + P(t)) до m коэффициентов
// vector<int> ln_series(const vector<int>& p, int m) {
//     vector<int> p_deriv = derivative(p);
//     vector<int> p_plus_one = p;
//     p_plus_one.insert(p_plus_one.begin(), 1);
//     vector<int> inv_p_plus_one(m, 0);
//     inv_p_plus_one[0] = 1; // 1 / (1 + P(t)) = 1 - P(t) + P(t)^2 - ...
//     vector<int> current_pow = p;
//     for (int k = 1; k < m; ++k) {
//         for (int i = 0; i < m; ++i) {
//             if (i < current_pow.size()) {
//                 inv_p_plus_one[i] = (inv_p_plus_one[i] + ((k % 2 == 1) ? MOD - current_pow[i] : current_pow[i])) % MOD;
//             }
//         }
//         current_pow = multiply_poly(current_pow, p, m);
//     }
//     vector<int> integrand = multiply_poly(p_deriv, inv_p_plus_one, m);
//     vector<int> res = integrate(integrand);
//     res.resize(m);
//     return res;
// }
//
// void dump(const vector<int>& p, int m) {
//     for (int i = 0; i < m; ++i) {
//         cout << p[i] << (i < m - 1 ? " " : "\n");
//     }
// }
//
// int main() {
//     ios_base::sync_with_stdio(false);
//     cin.tie(nullptr);
//
//     int n, m;
//     cin >> n >> m;
//
//     vector<int> p(n + 1);
//     for (int i = 0; i <= n; ++i) {
//         cin >> p[i];
//     }
//
//     // Убедимся, что p0 = 0
// //    p[0] = 0;
//
//     // Вычисляем sqrt(1 + P(t))
//     vector<int> sqrt_res = sqrt_series(p, m);
//     dump(sqrt_res, m);
// //    for (int i = 0; i < m; ++i) {
// //        cout << sqrt_res[i] << (i < m - 1 ? " " : "\n");
// //    }
//
//     // Вычисляем exp(P(t))
//     vector<int> exp_res = exp_series(p, m);
//     dump(exp_res, m);
// //    for (int i = 0; i < m; ++i) {
// //        cout << exp_res[i] << (i < m - 1 ? " " : "\n");
// //    }
//
//     // Вычисляем ln(1 + P(t))
//     vector<int> ln_res = ln_series(p, m);
//     dump(ln_res, m);
// //    for (int i = 0; i < m; ++i) {
// //        cout << ln_res[i] << (i < m - 1 ? " " : "\n");
// //    }
//
//     return 0;
// }












// vector<int> sqrt_series(const vector<int>& p, int m) {
//     vector<int> res(m, 0);
//     res[0] = 1; // sqrt(1 + 0) = 1
//     vector<int> current_pow = {1}; // (1 + P(t))^{1/2}
//     for (int k = 1; k < m; ++k) {
//         vector<int> term = multiply_poly(current_pow, p, m);
//         int coeff = (1LL * (MOD - 1) * mod_inv(2 * k)) % MOD;
//         coeff = (1LL * coeff * (2 * k - 3)) % MOD;
//         coeff = (1LL * coeff * mod_inv(k)) % MOD;
//         for (int i = 0; i < m; ++i) {
//             if (i < term.size()) {
//                 res[i] = (res[i] + 1LL * coeff * term[i]) % MOD;
//             }
//         }
//         current_pow = multiply_poly(current_pow, p, m);
//     }
//     return res;
// }

// vector<int> sqrt_series(const vector<int>& p, int m) {
//     vector<int> res(m, 0);
//     res[0] = 1; // sqrt(1 + 0) = 1
//
//     vector<int> binom_coeff(m, 0);
//     binom_coeff[0] = 1;
//     for (int k = 1; k < m; ++k) {
//         int prev = binom_coeff[k-1];
//         int numerator = (3 - 2 * k + MOD) % MOD; // (3 - 2k) mod MOD
//         int denominator = (2LL * k) % MOD;
//         binom_coeff[k] = (1LL * prev * numerator) % MOD;
//         binom_coeff[k] = (1LL * binom_coeff[k] * mod_inv(denominator)) % MOD;
//     }
//
//     vector<int> current_pow(m, 0);
//     current_pow[0] = 1; // P(t)^0 = 1
//
//     for (int k = 0; k < m; ++k) {
//         int coeff = binom_coeff[k];
//         for (int i = 0; i < m; ++i) {
//             if (i < current_pow.size()) {
//                 res[i] = (res[i] + 1LL * coeff * current_pow[i]) % MOD;
//             }
//         }
//
//         if (k < m - 1) {
//             vector<int> new_pow(m, 0);
//             for (int i = 0; i < m; ++i) {
//                 for (int j = 0; j <= i; ++j) {
//                     int p_val = (j < p.size()) ? p[j] : 0;
//                     int pow_val = (i - j < current_pow.size()) ? current_pow[i - j] : 0;
//                     new_pow[i] = (new_pow[i] + 1LL * p_val * pow_val) % MOD;
//                 }
//             }
//             current_pow = new_pow;
//         }
//     }
//
//     return res;
// }









