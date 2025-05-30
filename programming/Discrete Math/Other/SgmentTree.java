package Other;

import java.util.Scanner;

public class SgmentTree {
    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);

        int N = scanner.nextInt();
        int[] arr = new int[N];
        for (int i = 0; i < N; i++) {
            arr[i] = scanner.nextInt();
        }
        SegmentTree tree = new SegmentTree(arr);

        int K = scanner.nextInt();
        for (int i = 0; i < K; i++) {
            int l = scanner.nextInt();
            int r = scanner.nextInt();
            int[] res = tree.query(0, 0, N-1, l-1, r-1);
            System.out.println(res[0] + " " + (res[1]+1));
        }
    }
    public static class SegmentTree {
        int N;
        int[] tree;

        public SegmentTree(int[] arr) {
            N = arr.length;
            tree = new int[4*N];
            build(arr, 0, 0, N-1);
        }

        private void build(int[] arr, int v, int l, int r) {
            if (l == r) {
                tree[v] = arr[l];
            } else {
                int m = (l + r) / 2;
                build(arr, 2*v+1, l, m);
                build(arr, 2*v+2, m+1, r);
                tree[v] = Math.max(tree[2*v+1], tree[2*v+2]);
            }
        }

        public int[] query(int v, int l, int r, int ql, int qr) {
            if (ql > qr) {
                return new int[]{-Integer.MAX_VALUE, -1};
            }
            if (l == ql && r == qr) {
                return new int[]{tree[v], l};
            }
            int m = (l + r) / 2;
            int[] left = query(2*v+1, l, m, ql, Math.min(qr, m));
            int[] right = query(2*v+2, m+1, r, Math.max(ql, m+1), qr);
            return new int[]{Math.max(left[0], right[0]), left[0] >= right[0] ? left[1] : right[1]};
        }
    }
}
