module ternary_min(a, b, out);
  input [1:0] a;
  input [1:0] b;
  output [1:0] out;
  
  assign out[1] = (a[1] & b[1]);
  assign out[0] = (a[0] & b[0]) | (a[1] & b[0] & ~a[0]) | (a[0] & b[1] & ~b[0]);
endmodule

module ternary_max(a, b, out);
  input [1:0] a;
  input [1:0] b;
  output [1:0] out;
  
  assign out[1] = (a[1] | b[1]);
  assign out[0] = ~(a[1] | b[1]) & (a[0] | b[0]);
endmodule

module ternary_any(a, b, out);
  input [1:0] a;
  input [1:0] b;
  output [1:0] out;
  
  assign out[1] = (a[1] & b[1]) | (((~a[1] & a[0]) | (~b[1] & b[0])) & (a[1] | b[1]));
  assign out[0] = ((a[0] & b[0]) | (a[1] & ~b[1] & ~b[0]) | (b[1] & ~a[1] & ~a[0]));//((~(a[0] | b[0])) & (a[1] | b[1]));
  
endmodule

module ternary_consensus(a, b, out);
  input [1:0] a;
  input [1:0] b;
  output [1:0] out;
  
  assign out[1] = (a[1] & b[1]);
  assign out[0] = (a[0] | b[0]) | (a[1] & ~b[1] & ~ b[0]) | (b[1] & ~a[1] & ~a[0]);
  
endmodule

module nand_gate(a, b, out);
  input wire a, b;
  output out;

  supply1 pwr;
  supply0 gnd;

  wire nmos1_out;

  pmos pmos1(out, pwr, a);
  pmos pmos2(out, pwr, b);

  nmos nmos1(nmos1_out, gnd, b);
  nmos nmos2(out, nmos1_out, a);
endmodule

module and_gate(a, b, out);
  input wire a, b;
  output wire out;

  wire nand_out;

  nand_gate nand_gate1(a, b, nand_out);
  not_gate not_gate1(nand_out, out);
endmodule

module not_gate(a, out);
  input wire a;
  output out;

  supply1 pwr;
  supply0 gnd;

  pmos pmos1(out, pwr, a);
  nmos nmos1(out, gnd, a);
endmodule

module or_gate(a, b, out);
  input wire a, b;
  output wire out;

  wire not_a, not_b, nand_out;

  not_gate not_gate1(a, not_a);
  not_gate not_gate2(b, not_b);

  nand_gate nand_gate1(not_a, not_b, nand_out);

  not_gate not_gate3(nand_out, out);
endmodule