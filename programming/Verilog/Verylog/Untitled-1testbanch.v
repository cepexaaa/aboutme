// `include "implication.v"

// module testbench();

//   reg a = 0;
//   reg b = 0;
//   wire and_result;
//   wire nand_result;
//   wire s;
//   wire c_out;
//   wire c_out1, c_out2, c_out3, c_out4;

//   //not_gate not_gate1(a, result);
//   implication_mos implication_mos1(a, b, c_out1)
//   // nand_gate nand_gate1(a, b, nand_result);
//   // and_gate and_gate1(a, b, and_result);
//   // half_adder half_adder1(a, b, c_out, s);
//   initial begin

//     $display("a = %b, b = %b, out = %b", a, b, c_out1);
//     #5 a = 0; b = 1;
//     implication_mos implication_mos2(a, b, c_out2)
//     $display("a = %b, b = %b, c_out = %b", a, b, c_out, s);
//     #5 a = 1; b = 0;
//     implication_mos implication_mos3(a, b, c_out3)
//     $display("a = %b, b = %b, c_out = %b", a, b, c_out, s);
//     #5 a = 1; b = 1;
//     implication_mos implication_mos4(a, b, c_out4)
//     $display("a = %b, b = %b, c_out = %b", a, b, c_out, s);
//   end
// endmodule

`include "implication.v"

module testbench();

  reg a = 0;
  reg b = 0;
  wire c_out1, c_out2, c_out3, c_out4;

  implication_mos implication_mos1(a, b, c_out1);
  implication_mos implication_mos2(a, b, c_out2);
  implication_mos implication_mos3(a, b, c_out3);
  implication_mos implication_mos4(a, b, c_out4);

  // Uncomment the following lines if needed
  // wire and_result;
  // wire nand_result;
  // wire s;
  // wire c_out;

  initial begin
    $display("a = %b, b = %b, c_out = %b", a, b, c_out1);
    #5 a = 0; b = 1;
    $display("a = %b, b = %b, c_out = %b", a, b, c_out2);
    #5 a = 1; b = 0;
    $display("a = %b, b = %b, c_out = %b", a, b, c_out3);
    #5 a = 1; b = 1;
    $display("a = %b, b = %b, c_out = %b", a, b, c_out4);
  end
endmodule