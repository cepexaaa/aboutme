/**
 * В теле класса решения разрешено использовать только переменные делегированные в класс RegularInt.
 * Нельзя volatile, нельзя другие типы, нельзя блокировки, нельзя лазить в глобальные переменные.
 *
 * @author Кубеш Сергей
 */
class Solution : MonotonicClock {
    private var c1d1 by RegularInt(0)
    private var c1d2 by RegularInt(0)
    private var c1d3 by RegularInt(0)
    private var c2d1 by RegularInt(0)
    private var c2d2 by RegularInt(0)
    private var c2d3 by RegularInt(0)

    override fun write(time: Time) {
        // write right-to-left
        c2d1 = time.d1
        c2d2 = time.d2
        c2d3 = time.d3
        c1d3 = time.d3
        c1d2 = time.d2
        c1d1 = time.d1
    }

    override fun read(): Time {
        // read left-to-right
        val r1d1 = c1d1
        val r1d2 = c1d2
        val r1d3 = c1d3
        val r2d3 = c2d3
        val r2d2 = c2d2
        val r2d1 = c2d1
        val r1 = Time(r1d1, r1d2, r1d3)
        val r2 = Time(r2d1, r2d2, r2d3)

        if (r1 == r2) {
            return r1
        }
        if (r1.d1 != r2.d1) {
            return Time(r1.d1, 9999, 9999)
        }
        if (r1.d2 != r2.d2) {
            return Time(r1.d1, r1.d2, 9999)
        }
        return Time(r1.d1, r1.d2, maxOf(r1.d3, r2.d3))
    }
}