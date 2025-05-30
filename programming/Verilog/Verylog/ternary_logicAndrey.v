module ternary_min(a, b, out);
    input wire [1:0] a,b;
    output wire [1:0] out;


    wire not_out_a0;
    wire not_out_a1;
    wire not_out_b0;
    wire not_out_b1;

    not_gate not_gate1(a[0], not_out_a0);
    not_gate not_gate2(b[0], not_out_b0);
    not_gate not_gate3(a[1], not_out_a1);
    not_gate not_gate4(b[1], not_out_b1);

    //out 1
    wire and_out1;
    wire and_out2;

    and_gate and_gate1(not_out_a0, a[1], and_out1);
    and_gate and_gate2(and_out1, not_out_b0, and_out2);
    and_gate and_gate3(and_out2, b[1], out[1]);

    //out 0

    wire and_out3;
    wire and_out4;
    wire and_out5;

    wire and_out6;
    wire and_out7;
    wire and_out8;

    wire and_out9;
    wire and_out10;
    wire and_out11;

    and_gate and_gate4(a[0], not_out_a1, and_out3);
    and_gate and_gate5(and_out3, b[0], and_out4);
    and_gate and_gate6(and_out4, not_out_b1, and_out5);

    and_gate and_gate7(a[0], not_out_a1, and_out6);
    and_gate and_gate8(and_out6, not_out_b0, and_out7);
    and_gate and_gate9(and_out7, b[1], and_out8);

    and_gate and_gate10(not_out_a0, a[1], and_out9);
    and_gate and_gate11(and_out9, b[0], and_out10);
    and_gate and_gate12(and_out10, not_out_b1, and_out11);

    wire or_out;
    or_gate or_gate1(and_out5, and_out8, or_out);
    or_gate or_gate2(or_out, and_out11, out[0]);
endmodule

module ternary_max(a,b,out);
    input wire [1:0] a,b;
    output wire [1:0] out;


    wire not_out_a0;
    wire not_out_a1;
    wire not_out_b0;
    wire not_out_b1;

    not_gate not_gate1(a[0], not_out_a0);
    not_gate not_gate2(b[0], not_out_b0);
    not_gate not_gate3(a[1], not_out_a1);
    not_gate not_gate4(b[1], not_out_b1);

    //out 0

    wire and_out0_1;
    wire and_out0_2;
    wire and_out0_3;

    wire and_out0_4;
    wire and_out0_5;
    wire and_out0_6;

    wire and_out0_7;
    wire and_out0_8;
    wire and_out0_9;
    
    and_gate and_gate1(not_out_a0, not_out_a1, and_out0_1);
    and_gate and_gate2(and_out0_1, b[0], and_out0_2);
    and_gate and_gate3(and_out0_2, not_out_b1, and_out0_3);

    and_gate and_gate4(a[0], not_out_a1, and_out0_4);
    and_gate and_gate5(and_out0_4, b[0], and_out0_5);
    and_gate and_gate6(and_out0_5, not_out_b1, and_out0_6);

    and_gate and_gate_add7(a[0], not_out_a1, and_out0_7);
    and_gate and_gate_add8(and_out0_7, not_out_b0, and_out0_8);
    and_gate and_gate_add9(and_out0_8, not_out_b1, and_out0_9);

    wire or_out0;
    or_gate or_gate1(and_out0_3, and_out0_6, or_out0);
    or_gate or_gate2(or_out0, and_out0_9, out[0]);

    //out 1

    wire and_out1_1;
    wire and_out1_2;
    wire and_out1_3;

    wire and_out1_4;
    wire and_out1_5;
    wire and_out1_6;

    wire and_out1_7;
    wire and_out1_8;
    wire and_out1_9;

    wire and_out1_10;
    wire and_out1_11;
    wire and_out1_12;

    wire and_out1_13;
    wire and_out1_14;
    wire and_out1_15;

    
    and_gate and_gate7(not_out_a0, not_out_a1, and_out1_1);
    and_gate and_gate8(and_out1_1, not_out_b0, and_out1_2);
    and_gate and_gate9(and_out1_2, b[1], and_out1_3);

    and_gate and_gate10(a[0], not_out_a1, and_out1_4);
    and_gate and_gate11(and_out1_4, not_out_b0, and_out1_5);
    and_gate and_gate12(and_out1_5, b[1], and_out1_6);

    and_gate and_gate13(not_out_a0, a[1], and_out1_7);
    and_gate and_gate14(and_out1_7, not_out_b0, and_out1_8);
    and_gate and_gate15(and_out1_8, b[1], and_out1_9);

    and_gate and_gate16(a[1], not_out_a0, and_out1_10);
    and_gate and_gate17(and_out1_10, not_out_b1, and_out1_11);
    and_gate and_gate18(and_out1_11, not_out_b0, and_out1_12);

    and_gate and_gate19(not_out_a0, a[1], and_out1_13);
    and_gate and_gate20(and_out1_13, not_out_b1, and_out1_14);
    and_gate and_gate21(and_out1_14, b[0], and_out1_15);

    wire or_out1_1;
    wire or_out1_2;
    wire or_out1_3;

    or_gate or_gate3(and_out1_3, and_out1_6, or_out1_1);
    or_gate or_gate4(or_out1_1, and_out1_9, or_out1_2);
    or_gate or_gate5(or_out1_2, and_out1_12, or_out1_3);
    or_gate or_gate6(or_out1_3, and_out1_15, out[1]);

endmodule

module ternary_consensus (a,b,out);
    input wire [1:0] a,b;
    output wire [1:0] out;


    wire not_out_a0;
    wire not_out_a1;
    wire not_out_b0;
    wire not_out_b1;

    not_gate not_gate1(a[0], not_out_a0);
    not_gate not_gate2(b[0], not_out_b0);
    not_gate not_gate3(a[1], not_out_a1);
    not_gate not_gate4(b[1], not_out_b1);

    //out 1

    wire and_out1_1;
    wire and_out1_2;
    
    and_gate and_gate1(not_out_a0, a[1], and_out1_1);
    and_gate and_gate2(and_out1_1, not_out_b0, and_out1_2);
    and_gate and_gate3(and_out1_2, b[1], out[1]);

    //out 0

    wire or_out1;
    wire or_out2;
    wire or_out3;

    wire or_out4;
    wire or_out5;
    wire or_out6;

    or_gate or_gate1(a[0], a[1], or_out1);
    or_gate or_gate2(or_out1, b[0], or_out2);
    or_gate or_gate3(or_out2, b[1], or_out3);

    or_gate or_gate4(a[0], not_out_a1, or_out4);
    or_gate or_gate5(or_out4, b[0], or_out5);
    or_gate or_gate6(or_out5, not_out_b1, or_out6);

    and_gate and_gate4(or_out3, or_out6, out[0]);

endmodule

module ternary_any (a,b,out);
    input wire [1:0] a,b;
    output wire [1:0] out;


    wire not_out_a0;
    wire not_out_a1;
    wire not_out_b0;
    wire not_out_b1;

    not_gate not_gate1(a[0], not_out_a0);
    not_gate not_gate2(b[0], not_out_b0);
    not_gate not_gate3(a[1], not_out_a1);
    not_gate not_gate4(b[1], not_out_b1);

    //out 0

    wire and_out0_1;
    wire and_out0_2;
    wire and_out0_3;

    wire and_out0_4;
    wire and_out0_5;
    wire and_out0_6;

    wire and_out0_7;
    wire and_out0_8;
    wire and_out0_9;

    and_gate and_gate1(not_out_a0, not_out_a1, and_out0_1);
    and_gate and_gate2(and_out0_1, not_out_b0, and_out0_2);
    and_gate and_gate3(and_out0_2, b[1], and_out0_3);

    and_gate and_gate4(a[0], not_out_a1, and_out0_4);
    and_gate and_gate5(and_out0_4, b[0], and_out0_5);
    and_gate and_gate6(and_out0_5, not_out_b1, and_out0_6);

    and_gate and_gate7(a[1], not_out_a0, and_out0_7);
    and_gate and_gate8(and_out0_7, not_out_b1, and_out0_8);
    and_gate and_gate9(and_out0_8, not_out_b0, and_out0_9);

    wire or_out0;

    or_gate or_gate1(and_out0_3, and_out0_6, or_out0);
    or_gate or_gate2(or_out0, and_out0_9, out[0]);

    //out 1

    wire and_out1_1;
    wire and_out1_2;
    wire and_out1_3;

    wire and_out1_4;
    wire and_out1_5;
    wire and_out1_6;

    wire and_out1_7;
    wire and_out1_8;
    wire and_out1_9;

    and_gate and_gate10(a[0], not_out_a1, and_out1_1);
    and_gate and_gate11(and_out1_1, not_out_b0, and_out1_2);
    and_gate and_gate12(and_out1_2, b[1], and_out1_3);

    and_gate and_gate13(not_out_a0, a[1], and_out1_4);
    and_gate and_gate14(and_out1_4, not_out_b0, and_out1_5);
    and_gate and_gate15(and_out1_5, b[1], and_out1_6);

    and_gate and_gate16(not_out_a0, a[1], and_out1_7);
    and_gate and_gate17(and_out1_7, not_out_b1, and_out1_8);
    and_gate and_gate18(and_out1_8, b[0], and_out1_9);

    wire or_out1;
    or_gate or_gate3(and_out1_3, and_out1_6, or_out1);
    or_gate or_gate4(or_out1, and_out1_9, out[1]);
endmodule


module and_gate(a, b, out);
    input wire a, b;
    output wire out;

    wire nand_out;

    nand_gate nand_gate1(a, b, nand_out);
    not_gate not_gate1(nand_out, out);
endmodule

module nand_gate(a, b, out);
    input wire a, b;
    output out;

    supply1 pwr;
    supply0 gnd;

    wire nmos1_out;

    pmos pmos1(out, pwr, a);
    pmos pmos2(out, pwr, b);

    // 1 - сток 2 - исток 3 - база
    nmos nmos1(nmos1_out, gnd, b);
    nmos nmos2(out, nmos1_out, a);
endmodule

module not_gate(a, out);
    input wire a;
    output out;

    supply1 pwr;
    supply0 gnd;

    pmos pmos1(out, pwr, a);
    nmos nmos1(out, gnd, a);
endmodule

module nor_gate(a, b, out);
    input wire a, b;
    output out;

    supply1 pwr;
    supply0 gnd;

    wire pmos1_out;

    pmos pmos1(pmos1_out, pwr, a);
    pmos pmos2(out, pmos1_out, b);

    nmos nmos1(out, gnd, a);
    nmos nmos2(out, gnd, b);
endmodule

module or_gate(a,b,out);
    input wire a, b;
    output out;

    wire nor_out;
    
    nor_gate nor_gate1(a, b, nor_out);
    not_gate not_gate1(nor_out, out);
endmodule