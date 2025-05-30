#include <iostream>
#include <vector>

using namespace std;

void print(const vector<int>& p) {
    for (const int i : p) {
        cout << i << " ";
    }
}

bool condition(const int k, const int s1, const int s2) {
    return k > 0 && s1 != s2;
}

vector<int> N() {
    string s;
    cin >> s;
    vector z_function_fo_calculating_very_importaint_data(s.length(), 0);

   int left = 0, right = 0;
    for (int i = 0; i < s.length(); ++i) {
        z_function_fo_calculating_very_importaint_data[i] = max(0, min(right - i, z_function_fo_calculating_very_importaint_data[i - left]));
        while (i + z_function_fo_calculating_very_importaint_data[i] < s.length() && s[z_function_fo_calculating_very_importaint_data[i]] == s[i + z_function_fo_calculating_very_importaint_data[i]]) {
            z_function_fo_calculating_very_importaint_data[i]++;
        } if (i + z_function_fo_calculating_very_importaint_data[i] > right) {
            left = i;
            right = i + z_function_fo_calculating_very_importaint_data[i];
        }
    }

    return z_function_fo_calculating_very_importaint_data;;
    // print(z_function_fo_calculating_very_importaint_data);
}

vector<int> P() {
    string s;
    cin >> s;
    vector p(s.length(), 0);

    for (int i = 1; i < s.length(); i++) {
        int k = p[i - 1];
        while (condition(k, s[k], s[i])) {
            k = p[k - 1];
        }
        if (s[i] == s[k]) {
            k++;
        }
        p[i] = k;
    }

    return p;
    // print(p);
}

void substringsInStr() {
    vector<int> p(N());
    print(p);
    // string text;
    string s;
    // cin >> text;
    cin >> s;
    int len = s.length();
    for (int i = 0; i < p.size(); i++) {
        if (p[i] == len) {
            cout << i << " ";
        }
    }
}

int main() {
    ios::sync_with_stdio(0);
    cin.tie(0);
    cout.tie(0);

    substringsInStr();
    return 0;
}



