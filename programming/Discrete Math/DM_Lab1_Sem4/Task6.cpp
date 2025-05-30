#include <iostream>
#include <vector>

using namespace std;

const long long MOD = 1e9 + 7;

int main() {
    int k, m;
    cin >> k >> m;

    vector<long long> dp(m + 1, 0);
    for (int i = 0; i < k; ++i) {
        int weight;
        cin >> weight;
        dp[weight]++;
    }

    vector<long long> prefixSum(m + 1, 0);
    vector<long long> trees(m + 1, 0);
    trees[0] = 1;
    prefixSum[0] = 1;

    for (int i = 1; i <= m; ++i) {
        for (int j = 1; j <= i; ++j) {
            trees[i] = ((trees[i] + dp[j] * prefixSum[i - j])% MOD + MOD) % MOD;
        }

        for (int j = 0; j <= i; ++j) {
            prefixSum[i] = ((prefixSum[i] + trees[j] * trees[i - j])% MOD + MOD) % MOD;
        }
    }

    for (int i = 1; i <= m; ++i) {
        cout << trees[i] << " ";
    }
    cout << endl;

    return 0;
}