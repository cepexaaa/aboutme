"use struct"
//let expressions = expr.split(" ").filter(t => t !== " " && t !== "");
const cnst = (value) => () => value;
let lst = ["x", "y", "z"];
const variable = name => (...args) => {
    // if (!lst.includes(name)) {
    //     lst.push(name);
    // }
    return args[lst.indexOf(name)];
};
//const variable = (name) => (...args) => args.indexOf(charAt());
//
const add = (...args) => (...vars) => args.reduce((acc, func) => acc + func(...vars), 0);
// const multiply = (...args) => (...vars) => args.reduce((acc, func) => acc * func(...vars), 1);
//const subtract = (...args) => (...vars) => args.reduce((acc, func) => func(...vars) - acc, 0);
//const multiply = (...args) => (...vars) => args.reduce((acc, func) => acc * func(...vars), 0);
const calculus = operate => (...args) => (...vars) => operate(...args.map(arg => arg(...vars)))
// const divide = (...args) => (...vars) => args.reduce((acc, func) => acc / func(...vars), 1);
// const operate1 = (...args) => (...vars) => args.reduce((acc, func) => acc / func(...vars), 1);
// const divide = (...args) => (...vars) => args/vars;//.reduce((acc, func) => acc / func(...vars), 1);
//const divide  = operate1((a, b) => a / b)
//const negate = (...args) => (...vars) => args.reduce((acc, func) => func(acc, ...vars), 0);


// const subtract = (...args) => (...vars) => - args.reduce((acc, func) => acc - func(...vars), 0);//

const negate = calculus(a => -a)
const cube = calculus(a => a * a * a);
const cbrt = calculus(a => Math.cbrt(a));

// const add = operation => (...args) => (...vars) => args.reduce((acc, func) => acc + func(...vars), 0);
// const subtract = operation => (...args) => (...vars) => args.reduce((acc, func) => acc - func(...vars), 0);
// const multiply = operation => (...args) => (...vars) => args.reduce((acc, func) => acc * func(...vars), 0);
// const divide = operation => (...args) => (...vars) => args.reduce((acc, func) => acc / func(...vars), 0);
// const negate = operation => (...args) => args.reduce(acc => -acc, 0);


//const add = calculus((a, b) => a + b)
const subtract = calculus((a, b) => a - b)
const multiply = calculus((a, b) => a * b)
const divide = calculus((a, b) => a / b)

const pi = cnst(Math.PI);
const e = cnst(Math.E);


//const operationSign = {"+": 1, "-": 1, "/": 2, "*": 2}

// function doOperation(index) {
//     const operation2 = val.pop();
//     const operation1 = val.pop();
//     const sign = expressions[index];
//     switch (sign) {
//         case "+" :
//             val.push(new Add(operation1, operation2));
//             break;
//         case "-" :
//             val.push(new Subtract(operation1, operation2));
//             break;
//         case "/" :
//             val.push(new Divide(operation1, operation2));
//             break;
//         case "*" :
//             val.push(new Multiply(operation1, operation2));
//     }
// }


// } else if (expressions[index] === 'negate') {
//     let a = val.pop();
//     val.push(new Negate(a));
// } else if (expressions[index] in operationSign) {
//     doOperation(index);

// const someOperation = (constructor,  operation , evaluate,  toString) => {
//     constructor.prototype.operation = operation;
//     constructor.prototype.evaluate = evaluate;
//     constructor.prototype.toString = toString;
//     return constructor;
// }

//.match(/(\d+|\+|\-|\*|\/|\(|\)|\x|\y|\z|\negate)/g);//split(" ").filter(t => t.length > 0);

//console.log(new Add(new Const(2), new Const(3)));
// for (let i = 0; i < expressions.length; i++) {
//     console.log(expressions[i]);
// }

// if (expressions[index] === 'negate') {
//     console.log(expr)
//     a = val.pop();
//     console.log(a)
//     val.push(new Negate(a));
// } else

//console.log(parse("x"))

// } else if (index > 0 && expressions[index - 1] in operationSign && expressions[index] === '-') {
//     val.push(new Negate(negateOperation()));

// function negateOperation() {
//     if (expressions[index++] === "(") {
//         return processing();
//     } else if (isOb()) {
//         return val.pop();
//     } else if (expressions[index] === '-') {
//         return new Negate(negateOperation());
//     }
// }


// if (expressions[index] === 'negate') {
//     a = val.pop()
//     val.push(new Negate(a));
// } else {
//    doOperation();
// }
// while (stackSign.length > 0) {
//     if (operationSign[stackSign[stackSign.length - 1]] >= operationSign[expressions[index]]) {
//         doOperation();
//     } else {
//         break;
//     }
// }
//stackSign.push(expressions[index]);
// } else if (expressions[index] === ")") {
//     return val.pop();
// } else if (expressions[index] === "(") {
//     return processing();
//}
// if (expressions.length - 1 === index) {
//     return val.pop();
// }

//let a = parse("10").toString();
//console.log(parse("x negate").evaluate(2, 1, 0));

// let sign = [];
// let val = [];
// let strLen = expr.length;
// let index = 0;
//return processing(expr);





// function processing(expr) {
//     let sign = [];
//     let val = [];
//     let flagUnary = false;
//     while (index < strLen) {
//         c = take();
//         if (isNaN(c)) {
//             val.push(Const(wholeNumber()));
//             continue;
//         } else if (c in ["x", "y", "z"]) {
//             val.push((Variable(c)));
//         } else if (flagUnary && c === '-') {
//
//         } else if (c in ["+", "/", "*", "-"]) {
//
//         }
//
//
//     }
//     function take() {
//         if ((index < strLen) && (expr.charAt(index) === ' ')) {
//             return expr.charAt(index);
//         } else {
//             return expr.charAt(++index);
//         }
//     }
//     function wholeNumber() {
//         const start = index;
//         while ((index < strLen) && isNaN(expr.charAt(index))) {
//             index++;
//         }
//         if (index < strLen) {
//             index--;
//         } return expr.substring(start, index);
//     }
// }

// function skipSpace(c) {
//     if (c === ' ' && index < strLen) {
//         index++;
//     }
// }




// switch (expressions[index]) {
//     case "negate" :
//         val.push(operation1)
//         val.push(new Negate(operation2));
//         break;
//     case "+" :
//         val.push(new Add(operation1, operation2));
//         break;
//     case "-" :
//         val.push(new Subtract(operation1, operation2));
//         break;
//     case "/" :
//         val.push(new Divide(operation1, operation2));
//         break;
//     case "*" :
//         val.push(new Multiply(operation1, operation2));
//         break;
//     case "sinh" :
//         val.push(operation1)
//         val.push(new Sinh(operation2));
//         break;
//     case "cosh" :
//         val.push(operation1)
//         val.push(new Cosh(operation2));
// }















//let a = parsePostfix(""); console.log(a.postfix()); console.log(a.evaluate(1, 2, 3));
//let a = parsePrefix("(x)"); console.log(a.prefix()); console.log(a.evaluate(1, 2, 3));



// if (expressions[index] in map3) {
//     const sign = expressions[index++];
//     let value = [];
//     //const operation1 = certainOp(expressions[index]);
//     value.push(certainOp(expressions[index]));
//     if (sign === "product" || sign === "geom"){
//         let val = []
//         while (isOb(expressions[index++], val)) {}
//         //checkError(1);
//         return map3[sign](...val);
//     } else if (sign.length !== 1) {
//         checkError( 1);
//         return map3[sign](...value);
//     } else {
//         //const operation2 = certainOp(expressions[(++index) % expressions.length]);
//         value.push(certainOp(expressions[++index]));
//         checkError( 1);
//         return map3[sign](...value);
//     }
//}


// if (index + 1 <= expressions.length && expressions[index + 1] === ")") {
//     parentCounter -= 1;
// }
// index_1 = 1 + index;

// const operation1 = certanOp(expressions[index]);//maybe sweep op1 and op2
// const operation2 = certanOp(expressions[index % expressions.length]);
// return map2[sign](operation1, operation2);


// let a = parsePrefix("(+(+(* x x)(* y y))(* z z))");
// console.log(a.prefix());
// console.log(a.evaluate(1, 2, 3));



// }else if (expressions[index] === "(") {
//     parentCounter += expressions[index] === "(";
//     throw new SomeError("null expression");
// }
// if (parentCounter) {
//     throw new BracketError("no Bracket");




//const b = new Negate(new Variable("y"))
//console.log("(negate y)".match(/[^()\s]+|\(|\)|\s/g).filter(t => t !== " "));//.match(/\S+|\(|\)/g));
// const r = new Variable("x")
// console.log(r)
// console.log(b)
// console.log("(negate y)".split(/\s+|(?<=\w)(?=[()])/).filter(t => t !== ""))
// console.log("(negate y)".match(/[^()]+|\(|\)/g));
// console.log("(negate y)".match(/\S+|\(|\)/g));
// console.log("(neg)".split(/\s+|\b(?=[)(])/))
// console.log("@x".match(/\(|\w+|\)/g))
// console.log("x".match(/[a-zA-Z]+|\(|\)/g));
// const a = pars.prefix());