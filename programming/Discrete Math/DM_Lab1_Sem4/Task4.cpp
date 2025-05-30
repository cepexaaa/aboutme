#include <iostream>
#include <vector>
#include <algorithm>

using namespace std;

long long gcd(long long a, long long b) {
    while (b != 0) {
        long long temp = b;
        b = a % b;
        a = temp;
    }
    return a;
}

void simplify_fraction(long long &num, long long &den) {
    long long common_divisor = gcd(abs(num), abs(den));
    num /= common_divisor;
    den /= common_divisor;
    if (den < 0) {
        num *= -1;
        den *= -1;
    }
}

int main() {
    int r, k;
    cin >> r >> k;

    vector<long long> P(k + 1);
    for (int i = 0; i <= k; ++i) {
        cin >> P[i];
    }

    // Инициализация коэффициентов f_k(n) = c_0 + c_1*n + ... + c_k*n^k
    vector<long long> numerator(k + 1, 0);
    vector<long long> denominator(k + 1, 1);

    // Для каждого коэффициента P_i в P(t)
    for (int i = 0; i <= k; ++i) {
        if (P[i] == 0) continue;

        // Коэффициент при t^i в P(t) даёт вклад в c_j для j от 0 до k
        for (int j = 0; j <= k; ++j) {
            // Вычисляем коэффициент при n^j в разложении (n choose i) * r^{-i}
            // (n choose i) = n*(n-1)*...*(n-i+1)/i! = полином от n степени i
            if (j > i) continue;

            // Вычисляем коэффициент полинома (n choose i) при n^j
            long long coeff_num = 1;
            long long coeff_den = 1;

            // Вычисляем [n^j] (n choose i)
            // Это коэффициент при n^j в n*(n-1)*...*(n-i+1)/i!
            // Можно вычислить через Стирлинга или явно
            // Здесь упрощённый вариант для небольших k
            vector<long long> poly(i + 1, 0);
            poly[0] = 1;
            for (int m = 0; m < i; ++m) {
                // Умножаем на (n - m)
                for (int l = i; l >= 1; --l) {
                    poly[l] = poly[l - 1] - m * poly[l];
                }
                poly[0] = -m * poly[0];
            }
            // Делим на i!
            long long fact = 1;
            for (int m = 1; m <= i; ++m) {
                fact *= m;
            }
            coeff_num = poly[j];
            coeff_den = fact;

            // Умножаем на P_i и r^{-i}
            coeff_num *= P[i];
            coeff_den *= 1; // r^i уже учтено в общем множителе r^n

            // Приводим к общему знаменателю и складываем
            long long new_num = coeff_num * denominator[j] + numerator[j] * coeff_den;
            long long new_den = coeff_den * denominator[j];
            simplify_fraction(new_num, new_den);
            numerator[j] = new_num;
            denominator[j] = new_den;
        }
    }

    // Упрощаем дроби и выводим
    for (int j = 0; j <= k; ++j) {
        simplify_fraction(numerator[j], denominator[j]);
        cout << numerator[j] << "/" << denominator[j];
        if (j < k) cout << " ";
    }
    cout << endl;

    return 0;
}




/*#include <iostream>
#include <vector>
#include <cmath>
#include <algorithm>
/// DO NOT WORK. NEED TO FIX
using namespace std;

long long gcd(long long a, long long b) {
    return b == 0 ? a : gcd(b, a % b);
}

struct Fraction {
    long long num, den;

    Fraction(long long n = 0, long long d = 1) {
        long long g = gcd(abs(n), abs(d));
        num = n / g;
        den = d / g;
        if (den < 0) {
            num *= -1;
            den *= -1;
        }
    }

    Fraction operator+(const Fraction& other) const {
        return Fraction(num * other.den + other.num * den, den * other.den);
    }

    Fraction operator*(const Fraction& other) const {
        return Fraction(num * other.num, den * other.den);
    }
    Fraction operator^(const Fraction& other) const {
        return Fraction(pow(num, other.num), pow(den, other.den));
    }
};

vector<Fraction> solve(int r, int k, const vector<int>& P) {
    // Вычисляем α коэффициенты разложения P(t) по степеням (1-rt)
    vector<Fraction> alpha(k+1, Fraction(0,1));
    for (int i = 0; i <= k; ++i) {
        Fraction sum(0,1);
        for (int j = i; j <= k; ++j) {
            long long binom = 1;
            for (int t = 1; t <= i; ++t) {
                binom = binom * (j - t + 1) / t;
            }
            sum = sum + Fraction(P[j],1) * Fraction(binom * (i%2 ? -1 : 1), 1) * Fraction(1,1)^(j-i);
        }
        alpha[i] = sum * Fraction(1,1)^(i);
    }

    // Вычисляем коэффициенты C_j
    vector<Fraction> C(k+1);
    vector<long long> fact(k+1);
    fact[0] = 1;
    for (int i = 1; i <= k; ++i) {
        fact[i] = fact[i-1] * i;
    }

    for (int j = 0; j <= k; ++j) {
        long long comb = fact[k] / (fact[j] * fact[k-j]);
        C[j] = alpha[j] * Fraction(comb, 1) * Fraction(1,1)^(j);
    }

    return C;
}

int main() {
    int r, k;
    cin >> r >> k;

    vector<int> P(k+1);
    for (int i = 0; i <= k; ++i) {
        cin >> P[i];
    }

    vector<Fraction> res = solve(r, k, P);

    for (int i = 0; i <= k; ++i) {
        cout << res[i].num << "/" << res[i].den;
        if (i < k) cout << " ";
    }
    cout << endl;

    return 0;
}*/


/*
 //first wersion
 #include <iostream>
#include <vector>
#include <algorithm>

using namespace std;

long long gcd(long long a, long long b) {
    return b == 0 ? a : gcd(b, a % b);
}

struct Fraction {
    long long num, den;

    Fraction(long long n = 0, long long d = 1) {
        long long g = gcd(n, d);
        num = n / g;
        den = d / g;
        if (den < 0) {
            num *= -1;
            den *= -1;
        }
    }

    Fraction operator+(const Fraction& other) const {
        return Fraction(num * other.den + other.num * den, den * other.den);
    }

    Fraction operator-(const Fraction& other) const {
        return Fraction(num * other.den - other.num * den, den * other.den);
    }

    Fraction operator*(const Fraction& other) const {
        return Fraction(num * other.num, den * other.den);
    }
};

vector<Fraction> solve(int r, int k, const vector<int>& P) {
    vector<vector<Fraction>> dp(k + 1, vector<Fraction>(k + 1));
    dp[0][0] = Fraction(1, 1);
    for (int i = 1; i <= k; ++i) {
        for (int j = 0; j <= k; ++j) {
            if (j > 0) {
                dp[i][j] = dp[i - 1][j - 1] * Fraction(1, r);
            }
            dp[i][j] = dp[i][j] + dp[i - 1][j] * Fraction(i, 1);
        }
    }

    vector<Fraction> alpha(k + 1);
    for (int i = 0; i <= k; ++i) {
        alpha[i] = Fraction(0, 1);
        for (int j = 0; j <= k; ++j) {
            alpha[i] = alpha[i] + Fraction(P[j], 1) * dp[j][i];
        }
    }

    vector<Fraction> fact(k + 1);
    fact[0] = Fraction(1, 1);
    for (int i = 1; i <= k; ++i) {
        fact[i] = fact[i - 1] * Fraction(i, 1);
    }

    vector<Fraction> res(k + 1);
    for (int i = 0; i <= k; ++i) {
        res[i] = alpha[i] * Fraction(1, 1);
        for (int j = 0; j < i; ++j) {
            res[i] = res[i] - res[j] * fact[i] * Fraction(1, fact[j].num * fact[i - j].num);
        }
        res[i] = res[i] * Fraction(1, fact[i].num);
    }

    return res;
}

int main() {
    int r, k;
    cin >> r >> k;

    vector<int> P(k + 1);
    for (int i = 0; i <= k; ++i) {
        cin >> P[i];
    }

    vector<Fraction> res = solve(r, k, P);

    for (int i = 0; i <= k; ++i) {
        cout << res[i].num << "/" << res[i].den;
        if (i < k) {
            cout << " ";
        }
    }
    cout << endl;

    return 0;
}
 */