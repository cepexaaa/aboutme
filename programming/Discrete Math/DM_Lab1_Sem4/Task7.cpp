#include <iostream>
#include <vector>
#include <string>

using namespace std;

vector<long long> B() {
    return {0, 1, 0, 0, 0, 0, 0};
}

vector<long long> L(const vector<long long>& X) {
    vector<long long> res(7, 0);
    res[0] = 1;
    res[1] = X[1];
    res[2] = X[1]*X[1] + X[2];
    res[3] = X[1]*X[1]*X[1] + 2*X[1]*X[2] + X[3];
    res[4] = X[1]*X[1]*X[1]*X[1] + 3*X[1]*X[1]*X[2] + 2*X[1]*X[3] + X[2]*X[2] + X[4];
    res[5] = X[1]*X[1]*X[1]*X[1]*X[1] + 4*X[1]*X[1]*X[1]*X[2] + 3*X[1]*X[1]*X[3]
            + 3*X[1]*X[2]*X[2] + 2*X[1]*X[4] + 2*X[2]*X[3] + X[5];
    res[6] = X[1]*X[1]*X[1]*X[1]*X[1]*X[1] + 5*X[1]*X[1]*X[1]*X[1]*X[2]
            + 4*X[1]*X[1]*X[1]*X[3] + 6*X[1]*X[1]*X[2]*X[2] + 3*X[1]*X[1]*X[4]
            + 6*X[1]*X[2]*X[3] + 2*X[1]*X[5] + X[2]*X[2]*X[2] + 2*X[2]*X[4]
            + X[3]*X[3] + X[6];

    return res;
}

vector<long long> S(const vector<long long>& X) {
    vector<long long> res(7, 0);
    res[0] = 1;
    res[1] = X[1];
    res[2] = X[1]*(X[1]+1)/2 + X[2];
    res[3] = X[1]*(X[1]+1)*(X[1]+2)/6 + X[1]*X[2] + X[3];
    res[4] = X[1]*(X[1]+1)*(X[1]+2)*(X[1]+3)/24 + X[1]*(X[1]+1)*X[2]/2
            + X[1]*X[3] + X[2]*(X[2]+1)/2 + X[4];
    res[5] = X[1]*(X[1]+1)*(X[1]+2)*(X[1]+3)*(X[1]+4)/120 + X[1]*(X[1]+1)*(X[1]+2)*X[2]/6
            + X[1]*(X[1]+1)*X[3]/2 + X[1]*X[2]*(X[2]+1)/2 + X[1]*X[4] + X[2]*X[3] + X[5];
    res[6] = X[1]*(X[1]+1)*(X[1]+2)*(X[1]+3)*(X[1]+4)*(X[1]+5)/720 //111111
            + X[1]*(X[1]+1)*(X[1]+2)*(X[1]+3)*X[2]/24//11112
            + X[1]*(X[1]+1)*(X[1]+2)*X[3]/6 + X[1]*(X[1]+1)*X[2]*(X[2]+1)/4//1113 1122
            + X[1]*(X[1]+1)*X[4]/2 + X[1]*X[5] + X[1]*X[2]*X[3]//114 15 123
            + X[3]*(X[3]+1)/2 + X[2]*(X[2]+1)*(X[2]+2)/6//33 222
            + X[2]*X[4] + X[6];//24 6

    return res;
}

vector<long long> P(const vector<long long>& X, const vector<long long>& Y) {
    vector<long long> res(7, 0);

    res[0] = X[0] * Y[0];
    res[1] = X[0]*Y[1] + X[1]*Y[0];
    res[2] = X[0]*Y[2] + X[1]*Y[1] + X[2]*Y[0];
    res[3] = X[0]*Y[3] + X[1]*Y[2] + X[2]*Y[1] + X[3]*Y[0];
    res[4] = X[0]*Y[4] + X[1]*Y[3] + X[2]*Y[2] + X[3]*Y[1] + X[4]*Y[0];
    res[5] = X[0]*Y[5] + X[1]*Y[4] + X[2]*Y[3] + X[3]*Y[2] + X[4]*Y[1] + X[5]*Y[0];
    res[6] = X[0]*Y[6] + X[1]*Y[5] + X[2]*Y[4] + X[3]*Y[3] + X[4]*Y[2] + X[5]*Y[1] + X[6]*Y[0];

    return res;
}

int i = 0;

vector<long long> parseSets(const string& s) {
    while (i < s.size()) {
        if (s[i] == 'B') {
            i++;
            return B();
        } if (s[i] == 'L') {
            i++;
            return L(parseSets(s));
        } if (s[i] == 'S') {
            i++;
            return S(parseSets(s));
        } if (s[i] == 'P') {
            i++;
            return P(parseSets(s), parseSets(s));
        }
        i++;
    }
}
//S(S(S(S(S(S(S(S(S(S(B))))))))))
//1 1 10 55 385 2365 15367
int main() {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    string s;
    cin >> s;

    vector<long long> result = parseSets(s);

    for (int i1 = 0; i1 <= 6; ++i1) {
        cout << result[i1] << " ";
    }
    cout << endl;

    return 0;
}