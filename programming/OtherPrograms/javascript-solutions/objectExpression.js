"use struct"

const pseudoClass = (constructor, evaluate, toString, prefix, postfix) => {
    constructor.prototype.evaluate = evaluate;
    constructor.prototype.toString = toString;
    constructor.prototype.prefix = prefix;
    constructor.prototype.postfix = postfix;
    return constructor;
}

const Const = pseudoClass(
    function (value) {
        this.value = value
    },
    function () {
        return this.value
    },
    function () {
        return this.value.toString()
    },
    function () {
        return this.value.toString()
    },
    function () {
        return this.value.toString()
    }
);

const mapXYZ = {"x": 0, "y": 1, "z": 2};
const Variable = pseudoClass(
    function (value) {
        this.value = value
    },
    function (...args) {
        return args[mapXYZ[this.value]]
    },
    function () {
        return this.value
    },
    function () {
        return this.value
    },
    function () {
        return this.value
    }
);

const Operation = pseudoClass(
    function (...args) {
        this.args = args
    },
    function (...vars) {
        return this.certainOperation(...this.args.map(val => val.evaluate(...vars)))
    },
    function () {
        return this.args.join(" ") + " " + this.certainOperationSign
    },
    function () {
        return "(" + this.certainOperationSign + " " + this.args.map(arg => arg.prefix()).join(" ") + ")"
    },
    function () {
        return "(" + this.args.map(arg => arg.postfix()).join(" ") + " " + this.certainOperationSign + ")"
    }
);

const map2 = {}
const map3 = {}

function calculus(certainOperation, certainOperationSign) {
    const constructor = function (...args) {
        Operation.call(this, ...args)
    }
    constructor.prototype = Object.create(Operation.prototype)
    constructor.prototype.certainOperation = certainOperation
    constructor.prototype.certainOperationSign = certainOperationSign
    map2[certainOperationSign] = (...args) => new constructor(...args.reverse());
    map3[certainOperationSign] = (...args) => new constructor(...args);
    return constructor;
}

const Add = calculus((a, b) => a + b, "+")
const Subtract = calculus((a, b) => a - b, "-")
const Multiply = calculus((a, b) => a * b, "*")
const Divide = calculus((a, b) => a / b, "/")
const Negate = calculus(a => -a, "negate")
const Sinh = calculus(a => Math.sinh(a), "sinh")
const Cosh = calculus(a => Math.cosh(a), "cosh")
const Product = calculus((...args) => args.reduce((acc, x) => acc * x, 1), "product");
const Geom = calculus((...args) => Math.pow(Math.abs(Product.prototype.certainOperation(...args)), 1 / args.length), "geom");
const Sum = calculus((...args) => args.reduce((acc, x) => acc + x, 0), "")
const Mean = calculus((...args) => Sum.prototype.certainOperation(...args) / args.length, "mean")
const Sumsq = calculus((...args) => args.reduce((acc, x) => acc + x * x, 0) / args.length, "")
const Var = calculus((...args) => Sumsq.prototype.certainOperation(...args) - Math.pow(Mean.prototype.certainOperation(...args), 2), "var")
const ArithMean = calculus((...args) => Sum.prototype.certainOperation(...args) / args.length, "arithMean");
const GeomMean = calculus((...args) => Math.pow(Math.abs(Product.prototype.certainOperation(...args)), 1 / args.length), "geomMean");
const HarmMean = calculus((...args) => (args.length / args.reduce((acc, x) => acc + 1 / x, 0)), "harmMean")

function isOb(a, val) {
    if (!isNaN(a)) {
        val.push(new Const(parseInt(a)));
        return true;
    } else if (a in mapXYZ) {
        val.push(new Variable(a));
        return true;
    }
    return false;
}

function parse(expr) {
    let expressions = expr.split(/\s+/).filter(t => t !== "");
    let val = [];
    for (let index = 0; index < expressions.length; index++) {
        if (isOb(expressions[index], val)) {
        } else {
            const operation2 = val.pop();
            if (expressions[index].length === 1) {
                const operation1 = val.pop();
                val.push(map2[expressions[index]](operation2, operation1));
            } else {
                val.push(map2[expressions[index]](operation2));
            }
        }
    }
    return val.pop();
}

function SomeError(message) {
    this.message = message
}

SomeError.prototype = Object.create(Error.prototype);

function BracketError(message) {
    SomeError.call(this, message)
}

BracketError.prototype = Object.create(SomeError.prototype)

function EmptyError(message) {
    SomeError.call(this, message)
}

EmptyError.prototype = Object.create(SomeError.prototype);

function ArgumentsError(message) {
    SomeError.call(this, message)
}

ArgumentsError.prototype = Object.create(SomeError.prototype);

function parsePostfix(expr) {
    return parser(expr, 0);
}

function parsePrefix(expr) {
    return parser(expr, 1);
}

function parser(expr, prefOrPost) {
    let parentCounter = 0;
    if (expr === "") {
        throw new EmptyError("Empty expression");
    }
    let expressions = expr.match(/[^()\s]+|\(|\)|\s/g).filter(t => t !== " ");

    let index = 0;
    var result;
    if (prefOrPost) {
        const mapMuch = ["harmMean", "geomMean", "arithMean", "var", "mean", "product"];
        if (expressions[0] === "(" && expressions[expressions.length - 1] === ")") {
            expressions.shift();
            expressions.pop();
            if ((expressions.length <= 1 || ((!(expressions[0] in map3)) && expressions.length !== 1)) && (!(mapMuch.includes(expressions[0])))) {
                if (expressions[0] in map3) {
                    throw new ArgumentsError("Operator without some arguments")
                }
                if (expressions.length === 0) {
                    throw new EmptyError("Empty Expression")
                }
                if (isOb(expressions[index], [])) {
                    throw new ArgumentsError("Arguments without some operations");
                }
                throw new ArgumentsError("Unknown symbol -->'" + expressions[0] + "'");
            }
        }
        if (expressions.length === 0 || ((!(expressions[0] in map3)) && expressions.length !== 1)) {
            throw new ArgumentsError("Either there is a missing operator or an extra argument");
        }
        result = parsePrefRec(prefOrPost);
    } else {
        if (expressions[0] === "(" && expressions[expressions.length - 1] === ")") {
            expressions.shift();
            parentCounter++;
        }
        if (expressions.length === 0 || (!isOb(expressions[0], [])) && (expressions[0] !== "(")) {
            if (expressions[0] in map3) {
                throw new ArgumentsError("Operator without some arguments")
            }
            if (expressions[0] === ")") {
                throw new EmptyError("Empty expression.");
            }
            throw new ArgumentsError("Unknown symbol -->'" + expressions[0] + "'");
        }
        result = parsePostRec(prefOrPost);
        if (index + 1 < expressions.length) {
            throw new BracketError("Brackets is more then must");
        }
    }
    if (parentCounter) {
        throw new BracketError("There isn't open bracket.");
    } else {
        return result;
    }

    function certainOp(x, flagPrefOrPost) {
        let val = [];
        if (isOb(x, val)) {
            return val.pop()
        } else if (x === "(") {
            parentCounter += 1;
            index++;
            if (flagPrefOrPost) {
                return parsePrefRec(flagPrefOrPost);
            } else {
                return parsePostRec(flagPrefOrPost);
            }
        } else {
            throw new ArgumentsError("There isn't operation");
        }
    }

    function parsePostRec(flagPrefOrPost) {
        let value = [];
        index--;
        while ((index + 1 < expressions.length) && !(expressions[index + 1] in map3)) {
            value.push(certainOp(expressions[++index], flagPrefOrPost));
        }
        checkError(0);
        const sign = expressions[++index];
        if (!(sign in map3) && value.length !== 1) {
            throw new ArgumentsError("There isn't sign of operation on " + index + " position");
        }
        if ((sign !== undefined) && value.length !== 2 && sign.length === 1 || sign === "negate" && value.length !== 1) {
            throw new ArgumentsError("Uncorrect count of arguments. There is/are " + value.length + " arguments. Should to be " + ((sign.length === 1) + 1))
        }
        if ((sign === undefined) && value.length === 1) {
            return value.pop();
        }
        if (expressions[++index] === ")") {
            parentCounter--;
            return map3[sign](...value);
        } else {
            throw new BracketError("There isn't close bracket in " + index + " position.");
        }

    }

    function parsePrefRec(flagPrefOrPost) {
        if (expressions[index] in map3) {
            const sign = expressions[index];
            let value = [];
            while ((index + 1 < expressions.length) && expressions[index + 1] !== ")") {
                value.push(certainOp(expressions[++index], flagPrefOrPost));
            }
            checkError(1);
            if (value.length !== 2 && sign.length === 1 || sign === "negate" && value.length !== 1) {
                throw new SomeError("Uncorrect count of arguments. There is/are " + value.length + " arguments. Should to be " + ((sign.length === 1) + 1))
            }
            return map3[sign](...value);
        } else {
            let val = [];
            const symbol = expressions[index];
            if (isOb(symbol, val)) {
                checkError(1);
                return val.pop();
            } else if (symbol === ")" || (symbol === "(")) {
                throw new EmptyError("null expression");
            }
            if (!isOb(symbol) && !(symbol in map3) && !(symbol === ")" || symbol === "(")) {
                throw new ArgumentsError("Unknown symbol -->'" + expressions[index] + "'");
            }
        }

    }

    function checkError(flag) {
        if (!isOb(expressions[index], []) && !(expressions[index] in map3) && !(expressions[index] === ")" || expressions[index] === "(")) {
            throw new ArgumentsError("Unknown symbol");
        }
        if (flag === 1) {
            let fl = true;
            if (index + 1 < expressions.length && expressions[index + 1] === ")") {
                fl = false;
                parentCounter -= 1;
                index += 1;
            } else if (((index + 1) < expressions.length)) {
                throw new BracketError("There isn't close bracket");
            }
            if (parentCounter && fl) {
                throw new BracketError("no Bracket");
            }
        }
    }
}
