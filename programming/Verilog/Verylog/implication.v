// module implication_mos(a, b, out);
//   input a, b;
//   output out;
 
// wire notout;

// not_gate not_gate1(a, notout);
// or_gate or_gate1(notout, b, out);

// endmodule

// module not_gate(a, out);
//     input a;
//     output out;

//     supply1 pwr;
//     supply0 gnd;

//     pmos pmos1(out, pwr, a);
//     nmos nmos1(out, gnd, a);
// endmodule

// module nor_gate(a, b, out);
//     input a, b;
//     output out;

//     supply1 pwr;
//     supply0 gnd;

//     wire pmos1_out;

//     pmos pmos1(pmos1_out, pwr, a);
//     pmos pmos2(out, pmos1_out, b);

//     nmos nmos1(out, gnd, a);
//     nmos nmos2(out, gnd, b);
// endmodule

// module or_gate(a,b,out);
//     input a, b;
//     output out;

//     wire nor_out;
    
//     nor_gate nor_gate1(a, b, nor_out);
//     not_gate not_gate1(nor_out, out);
// endmodule

// module implication_mos(a, b, out);
//   input a, b;
//   output out;

//   wire notout;

//   not_gate not_gate1(a, notout);
//   or_gate or_gate1(notout, b, out);

// endmodule

// module not_gate(a, out);
//     input a;
//     output out;

//     supply1 pwr;
//     supply0 gnd;

//     pmos pmos1(out, pwr, a);
//     nmos nmos1(out, gnd, a);
// endmodule

// module nor_gate(a, b, out);
//     input a, b;
//     output out;

//     supply1 pwr;
//     supply0 gnd;

//     wire pmos1_out;

//     pmos pmos1(pmos1_out, pwr, a);
//     pmos pmos2(out, pwr, pmos1_out);

//     nmos nmos1(out, gnd, a);
//     nmos nmos2(out, gnd, b);
// endmodule

// module or_gate(a,b,out);
//     input a, b;
//     output out;

//     wire nor_out;

//     nor_gate nor_gate1(a, b, nor_out);
//     not_gate not_gate1(nor_out, out);
// endmodule

module implication_mos(a, b, out);
  input a, b;
  output out;

  wire notout;

  not_gate not_gate1(a, notout);
  or_gate or_gate1(notout, b, out);

endmodule

module not_gate(a, out);
    input a;
    output out;

    supply1 pwr;
    supply0 gnd;

    pmos pmos1(out, pwr, a);
    nmos nmos1(out, gnd, a);
endmodule

module nor_gate(a, b, out);
    input a, b;
    output out;

    supply1 pwr;
    supply0 gnd;

    wire pmos1_out;

    pmos pmos1(pmos1_out, pwr, a);
    pmos pmos2(out, pwr, b);

    nmos nmos1(out, gnd, a);
    nmos nmos2(out, gnd, b);
endmodule

module or_gate(a,b,out);
    input a, b;
    output out;

    wire nor_out;

    nor_gate nor_gate1(a, b, nor_out);
    not_gate not_gate1(nor_out, out);
endmodule