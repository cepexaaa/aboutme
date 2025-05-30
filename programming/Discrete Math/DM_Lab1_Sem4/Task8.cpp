// #include <iostream>
// #include <vector>
//
// using namespace std;
//
// const int MOD = 998244353;
//
// int main() {
//     int k, n;
//     cin >> k >> n;
//
//     vector<long long> dp(n + 1, 0);
//     vector<long long> left(n + 1, 0);
//
//     dp[1] = 1; // Дерево из одного листа
//
//     for (int i = 2; i <= n; ++i) {
//         for (int j = 1; j < i; ++j) {
//             // Учитываем только те left[j], которые не создают расчёску порядка k
//             if (j < k) {
//                 dp[i] = (dp[i] + dp[j] * dp[i - j]) % MOD;
//                 left[i] = (left[i] + left[j] * dp[i - j]) % MOD;
//             }
//         }
//         // Добавляем случай, когда корень имеет левого ребёнка
//         if (i < k) {
//             left[i] = (left[i] + dp[i]) % MOD;
//         }
//     }
//
//     for (int i = 1; i <= n; ++i) {
//         cout << dp[i] << " ";
//     }
//     cout << endl;
//
//     return 0;
// }
#include <iostream>
#include <vector>

using namespace std;

const int MOD = 998244353;

int main() {
    int k, n;
    cin >> k >> n;

    vector<long long> dp(n + 1, 0);
    vector<long long> sum(n + 1, 0);
    dp[1] = 1;
    sum[1] = 1;

    for (int i = 2; i <= n; ++i) {
        // Максимальная глубина левой расчески, которую мы можем допустить
        int max_depth = min(k - 2, i - 1);

        // Сумма dp[j] * dp[i-j] для j от 1 до i-1
        // Но с ограничением на глубину левого поддерева
        for (int j = 1; j <= i - 1; ++j) {
            if (j <= max_depth || (i - j) <= max_depth) {
                dp[i] = (dp[i] + dp[j] * dp[i - j]) % MOD;
            }
        }

        sum[i] = (sum[i - 1] + dp[i]) % MOD;
    }

    for (int i = 1; i <= n; ++i) {
        cout << dp[i] << " ";
    }
    cout << endl;

    return 0;
}