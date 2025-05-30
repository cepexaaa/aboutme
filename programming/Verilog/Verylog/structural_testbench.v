`include "rabota1.v"

module testbench();

  reg a = 0;
  reg b = 0;
  wire and_result;
  wire nand_result;
  wire s;
  wire c_out;

  //not_gate not_gate1(a, result);
  nand_gate nand_gate1(a, b, nand_result);
  and_gate and_gate1(a, b, and_result);
  half_adder half_adder1(a, b, c_out, s);
  initial begin
    $display("a = %b, b = %b, c_out = %b, s = %b", a, b, c_out, s);
    #5 a = 0; b = 1;
    $display("a = %b, b = %b, c_out = %b, s = %b", a, b, c_out, s);
    #5 a = 1; b = 0;
    $display("a = %b, b = %b, c_out = %b, s = %b", a, b, c_out, s);
    #5 a = 1; b = 1;
    $display("a = %b, b = %b, c_out = %b, s = %b", a, b, c_out, s);
  end
endmodule