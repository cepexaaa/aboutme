package search;

public class BinarySearch {
    // Pred: args.length >= 1 && forall i=1..(args.length - 1): args[i] >= args[i + 1]
    // Post: min(ind : 0 <= ind < args.length && args[0] >= args[ind + 1])
    public static void main(String[] args) {
        int x = Integer.parseInt(args[0]);
        int[] arr = new int[args.length - 1];
        for (int i = 1; i < args.length; i++) {
            arr[i - 1] = Integer.parseInt(args[i]);
        }
        int sum = 0;
        for (int i = 0; i < arr.length; i++) {
            sum += arr[i] % 2;
        }
        sum %= 2;
        if (sum == 1) {
            // forall i=0..(args.length - 2): arr[i] >= arr[i + 1]
            System.out.println(iterativeBinarySearch(arr, x));
        } else {
            // forall i=0..(args.length - 2): arr[i] >= arr[i + 1]
            System.out.println(recursiveBinarySearch(arr, x, 0, arr.length - 1));
        }
    }

    // Pred: forall i=0..(arr.length - 1): arr[i] >= arr[i + 1] && arr[-1] == +inf && arr[arr.length] == -inf
    // Post: R = r && 0 <= r < arr.length && arr[r-1] > x >= arr[r]
    public static int iterativeBinarySearch(int[] a, int x) {
        int left = 0;
        int right = a.length - 1;
        // Inv: -1 <= l <= r <= arr.length && arr[l] > x >= arr[r]
        while (left <= right) {
            // l < r
            int mid = (right + left) / 2;
            // l <= mid < r
            if (a[mid] <= x) {
                // arr[l] > x >= arr[mid]
                if (mid == 0 || a[mid - 1] > x) {
                    //a[mid - 1] > x > a[mid + 1] || x >= a[0] = max in a[]
                    return mid;
                } else {
                    right = mid - 1;
                    // r' == mid - 1  && l' == l
                    // l' < r' < r && arr[l'] > x >= arr[r']
                }
            } else {
                // arr[mid] > x >= arr[r]
                left = mid + 1;
                // l' == mid + 1 && r' == r
                // l < l' < r' && arr[l'] > x >= arr[r']
            }
            //r' - l' = 1/2(r - l)
            // 0 <= l' < r' <= arr.length - 1 && arr[l'] > x >= arr[r']
        }
        // l == r - 1 && arr[l] > x >= arr[r]
        // l == r - 1 => R = r && arr[r-1] > x >= arr[r]
        return right + 1;
    }

    // Pred: forall i=0..(arr.length - 2): arr[i] >= arr[i + 1] && arr[-1] == +inf && arr[arr.length] == -inf
    // Post: R = r && 0 <= r < arr.length && arr[r-1] > x >= arr[r]
    // Inv: 0 <= l < r <= arr.length - 1 && arr[l] > x >= arr[r]
    public static int recursiveBinarySearch(int[] a, int x, int left, int right) {
        if (left > right) {
            // l == r - 1 => R = r && arr[r-1] > x >= arr[r]
            return right + 1;
        }
        // l < r - 1
        int mid = (right + left) / 2;
        // l <= mid < r
        if (a[mid] <= x && (mid == 0 || a[mid - 1] > x)) {
            //a[mid - 1] > x > a[mid + 1] || x >= a[0] = max in a[]
            return mid;
        } else if (a[mid] > x) {
            // arr[mid] > x >= arr[r]
            // l' == mid && r' == r
            // l < l' < r' && arr[l'] > x >= arr[r']
            // -1 <= l' < r' <= arr.length && arr[l'] > x >= arr[r']
            return recursiveBinarySearch(a, x, mid + 1, right);
        } else {
            // arr[l] > x >= arr[mid]
            // r' == mid && l' == l
            // l' < r' < r && arr[l'] > x >= arr[r']
            // -1 <= l' < r' <= arr.length && arr[l'] > x >= arr[r']
            return recursiveBinarySearch(a, x, left, mid - 1);
        }
    }
}

