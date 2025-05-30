#include "return_codes.h"
#include <inttypes.h>
#include <stdio.h>

#define MAX_SINGLE_EXP 255
#define MAX_HALF_EXP 31
#define WHOLE_SINGLE_MANTISSA 8388607
#define WHOLE_HALF_MANTISSA 1023
#define MAX(a, b) ((a) > (b) ? (a) : (b))
#define MIN(a, b) ((a) < (b) ? (a) : (b))

typedef struct {
    uint8_t signBit;
    uint8_t exp;
    uint8_t len_exp;
    uint8_t len_mantissa;
    uint32_t mantissa;
    uint64_t value;
    char type;
} Number;

char check_exception(Number *number) {
    if (number->exp == MAX_SINGLE_EXP && number->len_exp == 8 || number->exp == MAX_HALF_EXP && number->len_exp == 5) {
        if (number->mantissa == 0) {
            return 'i';
        } else {
            return 'n';
        }
    }
    if (number->exp == 0) {
        if (number->mantissa == 0) {
            return '0';
        } else {
            return 'd';
        }
    }
    return 'f';
}

void isBigger(Number *number1, Number *number2) {
    if (number2->exp > number1->exp || number1->exp == number2->exp && number1->mantissa < number2->mantissa) {
        Number mid = *number1;
        *number1 = *number2;
        *number2 = mid;
    }
}

void rounding(Number *result, const int8_t *round, const uint8_t *num_remnant) {
    if (*round == 1) {
        result->mantissa += (((*num_remnant & 1) || (result->mantissa & 1)) * (*num_remnant >> 1));
    } else if (*round == 2) {
        result->mantissa += ((!(result->signBit & 1)) && *num_remnant != 0);
    } else if (*round == 3) {
        result->mantissa += ((result->signBit & 1) && *num_remnant != 0);
    }
    if (result->mantissa & (1 << result->len_mantissa)) {
        result->exp += 1;
        result->mantissa = 0;
    }
    result->type = check_exception(result);
}

Number multiply(Number *ns1, Number *ns2, int8_t *round) {
    Number rez;
    ns1->type = check_exception(ns1);
    ns2->type = check_exception(ns2);
    if (ns1->type == '0' && ns2->type == 'i' || ns1->type == 'i' && ns2->type == '0') {
        rez.type = 'n';
    } else if (ns1->type == 'i' || ns2->type == 'i') {
        rez.type = 'i';
    } else if (ns1->type == '0' || ns2->type == '0') {
        rez.type = '0';
    } else {
        uint8_t remnant2;
        uint64_t norm_val_rez = ns1->value * ns2->value;
        uint64_t dop_bit = ((norm_val_rez & ((uint64_t) 1 << (ns1->len_mantissa * 2 + 1))) != 0);
        rez.exp = ns1->exp + ns2->exp - (1 << (ns1->len_exp - 1)) + 1 + dop_bit;
        if (!(ns1->type == 'f' && ns2->type == 'f')) {
            while ((norm_val_rez & ((uint64_t) 1 << (ns1->len_mantissa * 2))) == 0) {
                norm_val_rez = norm_val_rez << 1;
            }
        }
        uint16_t maskRemnant = 1 << (ns1->len_mantissa - 1 + dop_bit);
        remnant2 = (((norm_val_rez & maskRemnant) != 0) << 1) +
                   ((norm_val_rez & (maskRemnant - 1)) != 0);
        rez.mantissa = (norm_val_rez >> (dop_bit + ns1->len_mantissa)) & ((1 << ns1->len_mantissa) - 1);
        rez.len_mantissa = ns1->len_mantissa;
        rez.signBit = ns1->signBit ^ ns2->signBit;
        rounding(&rez, round, &remnant2);
        rez.type = check_exception(&rez);
    }
    rez.len_exp = ns1->len_exp;
    return rez;
}

Number add(Number *ns1, Number *ns2, int8_t *round) {
    int16_t shiftLeft = MIN(ns1->exp - ns2->exp, 63 - ns1->len_mantissa);
    uint64_t norm_val_1 = (uint64_t)(ns1->value) << shiftLeft;
    uint64_t norm_val_2 = (uint64_t)(ns2->value) >> MAX(ns1->exp - ns2->exp - 63 + ns1->len_mantissa, 0);
    uint64_t norm_val_res = norm_val_1 + norm_val_2;
    uint8_t dop_bit = norm_val_res >> (ns1->len_mantissa + shiftLeft + 1);
    uint16_t maskRemnant = 1 << (shiftLeft + dop_bit - 1);
    uint8_t remnant = (((norm_val_res & maskRemnant) != 0) << 1) +
                      ((norm_val_res & (maskRemnant - 1)) != 0);
    ns1->mantissa = (norm_val_res >> (shiftLeft + dop_bit)) & ((1 << ns1->len_mantissa) - 1);
    ns1->exp = MAX(ns1->exp, ns2->exp) + dop_bit;
    rounding(ns1, round, &remnant);
    return *ns1;
}

Number subtract(Number *ns1, Number *ns2, int8_t *round) {
    int16_t shiftLeft = MIN(ns1->exp - ns2->exp, 63 - ns1->len_mantissa);
    uint64_t norm_val_1 = (uint64_t)(ns1->value) << shiftLeft;
    uint64_t norm_val_2 = (uint64_t)(ns2->value) >> MAX(ns1->exp - ns2->exp - 63 + ns1->len_mantissa, 0);
    uint64_t norm_val_res = norm_val_1 - norm_val_2;
    uint64_t differ = 1 << (ns1->len_mantissa + shiftLeft);
    uint8_t count = 0;
    while (norm_val_res < differ) {
        differ >>= 1;
        count++;
    }
    uint8_t remnant;
    ns1->exp = count > ns1->exp ? 0 : ns1->exp - count;
    uint16_t lenRemnant = shiftLeft - count - 1;
    if (shiftLeft > count) {
        remnant = (((norm_val_res & (1 << lenRemnant)) != 0) << 1) +
                  ((norm_val_res & ((1 << lenRemnant) - 1)) != 0);
        ns1->mantissa = (norm_val_res >> (lenRemnant + 1)) & ((1 << ns1->len_mantissa) - 1);
    } else {
        ns1->mantissa = norm_val_res << (count - shiftLeft);
        remnant = 0;
    }
    rounding(ns1, round, &remnant);
    return *ns1;
}

Number divide(Number *ns1, Number *ns2, int8_t *round) {
    Number rez;
    ns1->type = check_exception(ns1);
    ns2->type = check_exception(ns2);
    if ((ns1->type == ns2->type) && (ns1->type == '0' || ns1->type == 'i')) {
        rez.type = 'n';
    } else if (ns1->type == '0' || ns2->type == 'i') {
        rez.type = '0';
    } else if (ns1->type == 'i' || ns2->type == '0') {
        rez.type = 'i';
    } else {
        uint64_t norm_val_rez = (ns1->value << (63 - ns1->len_mantissa)) / ns2->value;
        rez.exp = ns1->exp - ns2->exp + (1 << (ns1->len_exp - 1));
        uint8_t index = 63;
        uint64_t valIndex = (uint64_t) 1 << 63;
        while (!(norm_val_rez & valIndex)) {
            valIndex >>= 1;
            index--;
        }
        uint16_t lenRemnant = 63 - ns1->len_mantissa - index;
        uint8_t remnant2 = (((norm_val_rez & (1 << lenRemnant)) != 0) << 1) +
                           ((norm_val_rez & ((1 << lenRemnant) - 1)) != 0);
        rez.mantissa = (norm_val_rez >> lenRemnant) & ((1 << ns1->len_mantissa) - 1);
        rez.len_mantissa = ns1->len_mantissa;
        rez.signBit = ns1->signBit ^ ns2->signBit;
        rounding(&rez, round, &remnant2);
        rez.type = check_exception(&rez);
    }
    rez.len_exp = ns1->len_exp;
    return rez;
}

void determine(const char *fh, const uint32_t *numb, Number *number) {
    if (*fh == 'f') {
        number->signBit = (*numb >> 31);
        number->exp = (*numb >> 23) & MAX_SINGLE_EXP;
        number->len_exp = 8;
        number->len_mantissa = 23;
        number->mantissa = (*numb & WHOLE_SINGLE_MANTISSA);
        number->type = check_exception(number);
        number->value = number->type == 'f' ? number->mantissa + (1 << number->len_mantissa) : number->mantissa;
    } else if (*fh == 'h') {
        number->signBit = *numb >> 15;
        number->exp = (*numb >> 10) & MAX_HALF_EXP;
        number->len_exp = 5;
        number->len_mantissa = 10;
        number->mantissa = (*numb & WHOLE_HALF_MANTISSA);
        number->type = check_exception(number);
        number->value = number->type == 'f' ? number->mantissa + (1 << number->len_mantissa) : number->mantissa;
    }
}

Number operations(Number number1, const char sign, Number number2, int8_t round) {
    number1.type = check_exception(&number1);
    number2.type = check_exception(&number2);
    if (number1.type == 'n' || number2.type == 'n') {
        isBigger(&number1, &number2);
        number1.type = 'n';
        return number1;
    }
    if (sign == '+' || sign == '-') {
        uint8_t flag_grow_module = (number1.signBit ^ number2.signBit) && sign == '-' ||
                                   (!number1.signBit ^ number2.signBit) && sign == '+';
        if (number1.type == 'i' && number2.type == 'i' && !flag_grow_module) {
            number1.type = 'n';
            return number1;
        } else if (number1.type == 'i' || number2.type == 'i') {
            if (number1.type == 'i' && number2.type == 'i') {
                return number1;
            } else {
                if (number1.type == 'i') {
                    return number1;
                } else {
                    number2.signBit = sign == '+' != 0 == number2.signBit;
                    return number2;
                }
            }
        } else if (number2.type == '0') {
            return number1;
        } else if (number1.type == '0') {
            return number2;
        } else {
            Number forCheck = number1;
            isBigger(&number1, &number2);
            if (flag_grow_module) {
                number1.signBit = forCheck.signBit;
                return add(&number1, &number2, &round);
            } else {
                if (number1.signBit == number2.signBit && sign == '-' &&
                    !(number1.signBit == forCheck.signBit && number1.mantissa == forCheck.mantissa &&
                      number1.exp == forCheck.exp)) {
                    number1.signBit = ~number1.signBit;
                }
                return subtract(&number1, &number2, &round);
            }
        }
    } else if (sign == '*') {
        return multiply(&number1, &number2, &round);
    } else if (sign == '/') {
        return divide(&number1, &number2, &round);
    }
    return number1;
}

int main(int argc, char *argv[]) {
    char fh;
    int8_t round;
    uint32_t number1;
    uint32_t number2;
    char sign;
    Number numb_struct_1;
    Number numb_struct_2;
    Number result;
    if (argc != 4 && argc != 6) {
        fprintf(stderr, "Incorrect count of arguments %i\n", argc);
        return ERROR_ARGUMENTS_INVALID;
    }
    sscanf(argv[1], "%c", &fh);
    sscanf(argv[2], "%"SCNx8, &round);
    sscanf(argv[3], "%"SCNx32, &number1);
    if (!(fh == 'f' || fh == 'h')) {
        fprintf(stderr, "Incorrect format\n");
        return ERROR_ARGUMENTS_INVALID;
    }
    if (round > 3 || round < 0) {
        fprintf(stderr, "Incorrect count of arguments\n");
        return ERROR_ARGUMENTS_INVALID;
    }
    determine(&fh, &number1, &numb_struct_1);
    if (argc == 6) {
        sscanf(argv[4], "%c", &sign);
        sscanf(argv[5], "%"SCNx32, &number2);
        determine(&fh, &number2, &numb_struct_2);
    }
    if (argc == 6) {
        if (sign != '*' && sign != '-' && sign != '+' && sign != '/') {
            for (int i = 0; i < argc; i++) {
                printf("%s \n", argv[i]);
            }
            fprintf(stderr, "Incorrect operation '%c'\n", sign);
            return ERROR_ARGUMENTS_INVALID;
        }
        result = operations(numb_struct_1, sign, numb_struct_2, round);
    } else {
        result = numb_struct_1;
    }
    if (result.type == 'f') {
        if (fh == 'f') {
            if (result.signBit == 0) {
                printf("0x1.%06xp%+i\n", result.mantissa << ((result.len_exp & 1) + 1),
                       ((int16_t) result.exp) - ((1 << (result.len_exp - 1)) - 1));
            } else {
                printf("-0x1.%06xp%+i\n", result.mantissa << ((result.len_exp & 1) + 1),
                       ((int16_t) result.exp) - ((1 << (result.len_exp - 1)) - 1));
            }
        } else {
            if (!(result.signBit & 1)) {
                printf("0x1.%03xp%+i\n", result.mantissa << ((result.len_exp & 1) + 1),
                       ((int8_t) result.exp) - ((1 << (result.len_exp - 1)) - 1));
            } else {
                printf("-0x1.%03xp%+i\n", result.mantissa << ((result.len_exp & 1) + 1),
                       ((int8_t) result.exp) - ((1 << (result.len_exp - 1)) - 1));
            }
        }
        return SUCCESS;
    }
    switch (result.type) {
        case 'n':
            printf("nan");
            break;
        case 'i':
            if (result.signBit) {
                printf("-inf");
            } else {
                printf("inf");
            }
            break;
        case '0':
            if (result.len_exp == 5) {
                printf("0x0.000p+0");
            } else {
                printf("0x0.000000p+0");
            }
    }
    return SUCCESS;
}