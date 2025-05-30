package search;

public class BinarySearchClosestA {
    // :NOTE: string sorted
    // Pred: args.length >= 1 && forall i=1..(args.length - 1): args[i] <= args[i + 1] && args[i] can be parsed to int
    // Post: min(ind : 0 <= ind < args.length && args[0] <= args[ind + 1])
    public static void main(String[] args) {
        int x = Integer.parseInt(args[0]);
        int[] arr = new int[args.length - 1];
        for (int i = 1; i < args.length; i++) {
            arr[i - 1] = Integer.parseInt(args[i]);
        }
        //arr.length >= 1 && forall i=1..(arr.length - 1): arr[i] <= arr[i + 1]
        int sum = 0;
        for (int i = 0; i < arr.length; i++) {
            sum = (sum + arr[i]) % 2;
        }
        if (sum == 1) {
            // forall i=0..(arr.length - 1): arr[i] <= arr[i + 1]
            System.out.println(iterativeBinarySearch(arr, x));
        } else {
            // forall i=0..(arr.length - 1): arr[i] <= arr[i + 1]
            System.out.println(recursiveBinarySearch(arr, x, 0, arr.length - 1)); // do for this///////////////
        }
        // forall i=0..(args.length - 1): arr[i] <= arr[i + 1]
    }

    //Since at each iteration of the while we reduce the length of the range in question by half
    // and the length of the range is finite, the binary search will work for the finite time.
    //The length of the range is 1, where a[left] <= x <= a[right]. Therefore, there are two options to consider: a[left] or a[right]


    // Pred: forall i=0..(arr.length - 1): arr[i] <= arr[i + 1] && arr[-1] == +inf && arr[arr.length] == -inf && i <= Integer.MAX_VALUE
    // Post: R = r && 0 <= r < arr.length && arr[r-1] < x <= arr[r]
    public static int iterativeBinarySearch(int[] a, int x) {
        int left = 0;
        int right = a.length - 1;
        // Inv: -1 <= l <= r <= a.length && a[l] <= x <= a[r]
        while (right - left > 1) {
            // l < r
            int mid = (right + left) / 2;
            //l <= mid < r
            if (a[mid] >= x) {
                // a[mid] < x < a[r]
                if (a[mid] == x) {
                    // a[mid] == x
                    return a[mid];
                } else {
                    // r' == mid - 1  && l' == l
                    // l' < r' < r && a[l'] > x >= a[r']
                    right = mid;
                    //r' - l' = 1/2(r - l)
                }
                //r' - l' = 1/2(r - l) && (r' == mid && l' == l || l' == mid && r' == r)
            } else {
                // a[mid] <= x < a[r]
                left = mid;
                // l' == mid + 1 && r' == r
                // l < l' < r' && a[l'] < x <= a[r']
                //r' - l' = 1/2(r - l)
            }
            //r' - l' = 1/2(r - l) && (r' == mid && l' == l || l' == mid && r' == r)
            // 0 <= l' < r' <= a.length - 1 && a[l'] <= x < a[r']
        }
        //Since at each iteration of the while we reduce the length of the range in question by half
        // and the length of the range is finite, the binary search will work for the finite time.
        //The length of the range is 1, where a[left] <= x <= a[right]. Therefore, there are two options to consider: a[left] or a[right]
        // r == l - 1
        // l == r - 1 && a[l] < x <= a[r]
        // l == r - 1 => R = r && a[r-1] > x >= a[r]

        // x <= a[right] || x >= a[left]
        return findClosest(left, a, x);
    }


    // Pred: forall i=0..(arr.length - 2): arr[i] <= arr[i + 1] && arr[-1] == +inf && arr[arr.length] == -inf && i <= Integer.MAX_VALUE
    // Post: R = r && 0 <= r < arr.length && arr[l] <= x < arr[r]
    // Inv: 0 <= l < r <= arr.length - 1 && arr[l] <= x < arr[r]
    public static int recursiveBinarySearch(int[] a, int x, int left, int right) {
        if (right - left <= 1) {
            // l == r - 1 => R = r && arr[r-1] > x >= arr[r]
            return findClosest(left, a, x);
        }
        // l < r - 1
        int mid = (right + left) / 2;
        // l <= mid < r
        if (a[mid] == x) {
            //x == a[mid]
            return a[mid];
        } else if (a[mid] < x) {
            // arr[mid] < x <= arr[r]
            // l' == mid && r' == r
            // l < l' < r' && arr[l'] < x <= arr[r']
            // -1 <= l' < r' <= arr.length && arr[l'] < x <= arr[r']
            return recursiveBinarySearch(a, x, mid, right);
        } else {
            // arr[l] < x <= arr[mid]
            // r' == mid && l' == l
            // l' < r' < r && arr[l'] < x <= arr[r']
            // -1 <= l' < r' <= arr.length && arr[l'] < x <= arr[r']
            return recursiveBinarySearch(a, x, left, mid);
        }
        //Since at each iteration of the while we reduce the length of the range in question by half
        // and the length of the range is finite, the binary search will work for the finite time.
        //The length of the range is 1, where a[left] <= x <= a[right]. Therefore, there are two options to consider: a[left] or a[right]

    }

    // Pred: a[indexl] < x < a[indexl + 1]
    // Post: R == a[indexl] || R == a[indexl + 1]
    // :NOTE: random ans satisfies
    public static int findClosest(int indexl, int[] a, int x) {
        int valLeft;
        int valRight;
        int indexr = indexl;
        if (indexl < a.length - 1) {
            indexr++;
        }
        valRight = Math.abs(a[indexr] - x);
        valLeft = Math.abs(a[indexl] - x);
        if (x == Integer.MIN_VALUE) {
            //max value on right position && x == Integer.MIN_VALUE => x on right position
            return a[indexl];
        }
        if (x == Integer.MAX_VALUE) {
            //min value on left position && x == Integer.MIN_VALUE => x on left position
            return a[indexr];
        }
        return (valLeft <= valRight) ? a[indexl] : a[indexr];
    }
}