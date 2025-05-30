module alu(a, b, control, res);
  input [3:0] a, b;
  input [2:0] control; // Выбор операции. Стоит поменять в зависимости от того, как нам подают инф-у

  output reg [3:0] res; 

  // Необходимо добавить выход "zero" в вашем описании модуля
  output reg zero;

  always @*
    case (control)
      3'b000: res = a & b;
      3'b001: res = ~(a & b);
      3'b010: res = a | b;
      3'b011: res = ~(a | b);
      3'b100: res = a + b;
      3'b101: res = a - b;
      3'b110: begin
        // Исправленная логика для операции сравнения
        res = ((a[3] == b[3]) & (a[2:0] < b[2:0]) | (a[3] == 1'b1 & b[3] == 1'b0)) ? 4'b0001 : 4'b0000;
      end
      default: res = 4'b0; // Обработка некорректного значения control
    endcase

  // Установка выхода "zero" в 1, если результат равен 0
  //assign zero = (res == 32'b0);

endmodule




// module alu(a, b, control, res);
//   input [3:0] a, b; // Операнды
//   input [2:0] control; // Управляющие сигналы для выбора операции

//   output [3:0] res; // Результат
//   // TODO: implementation

// wire [3:0] and_res, nand_res, or_res, nor_res, add_res,inv_b, add_res_1, sub_res, slt_res;

//   assign and_res = a & b;
//   assign nand_res = ~(a & b);
//   assign or_res = a | b;
//   assign nor_res = ~(a | b);
//   assign add_res = a + b;
//   assign sub_res = a - b;
//   assign slt_res = ((a[3] == b[3]) & (a[2:0] < b[2:0]) | (a[3] == 1'b1 & b[3] == 1'b0)) ? 4'b0001 : 4'b0000;

// //   and_gate_for and_gate_for1(a, b, and_res); 
// //   two_comp two_comp1(and_res, nand_res); 
// //   or_gate_for or_gate_for1(a, b, or_res);
// //   two_comp two_comp2(or_res, nor_res);
// //   full_adder_4b full_adder_4b_1(add_res, a, b);
// //   two_comp two_comp3(inv_b, b);
// //   full_adder_4b full_adder_4b_2(sub_res, a, inv_b);
// //   assign slt_res = a < b ? 4'b0001 : 4'b0000;

//   wire sel0, sel1, sel2, sel3, sel4, sel5, sel6;
//   assign sel0 = ~control[2] & ~control[1] & ~control[0];
//   assign sel1 = ~control[2] & ~control[1] & control[0];
//   assign sel2 = ~control[2] & control[1] & ~control[0];
//   assign sel3 = ~control[2] & control[1] & control[0];
//   assign sel4 = control[2] & ~control[1] & ~control[0];
//   assign sel5 = control[2] & ~control[1] & control[0];
//   assign sel6 = control[2] & control[1] & ~control[0];
  
//   assign res[0] = (and_res[0] & sel0) | (nand_res[0] & sel1) | (or_res[0] & sel2) | (nor_res[0] & sel3) | (add_res[0] & sel4) | (sub_res[0] & sel5) | (slt_res[0] & sel6);
//   assign res[1] = (and_res[1] & sel0) | (nand_res[1] & sel1) | (or_res[1] & sel2) | (nor_res[1] & sel3) | (add_res[1] & sel4) | (sub_res[1] & sel5) | (slt_res[1] & sel6);
//   assign res[2] = (and_res[2] & sel0) | (nand_res[2] & sel1) | (or_res[2] & sel2) | (nor_res[2] & sel3) | (add_res[2] & sel4) | (sub_res[2] & sel5) | (slt_res[2] & sel6);
//   assign res[3] = (and_res[3] & sel0) | (nand_res[3] & sel1) | (or_res[3] & sel2) | (nor_res[3] & sel3) | (add_res[3] & sel4) | (sub_res[3] & sel5) | (slt_res[3] & sel6);

// endmodule

module d_latch(clk, d, we, q);
  input clk; // Сигнал синхронизации
  input d; // Бит для записи в ячейку
  input we; // Необходимо ли перезаписать содержимое ячейки

  output reg q; // Сама ячейка
  // Изначально в ячейке хранится 0
  initial begin
    q <= 0;
  end
  // Значение изменяется на переданное на спаде сигнала синхронизации
  always @ (negedge clk) begin
    // Запись происходит при we = 1
    if (we) begin
      q <= d;
    end
  end
endmodule

module register_4bit(clk, d, we, q);
  input [3:0] d;
  input clk, we;
  output [3:0] q;

  d_latch d0(clk, d[0], we, q[0]);
  d_latch d1(clk, d[1], we, q[1]);
  d_latch d2(clk, d[2], we, q[2]);
  d_latch d3(clk, d[3], we, q[3]);

endmodule

module register_file(clk, rd_addr, we_addr, we_data, rd_data, we);
  input clk; // Сигнал синхронизации
  input [1:0] rd_addr, we_addr; // Номера регистров для чтения и записи
  input [3:0] we_data; // Данные для записи в регистровый файл
  input we; // Необходимо ли перезаписать содержимое регистра
  
  output [3:0] rd_data; // Данные, полученные в результате чтения из регистрового файла

  wire we0, we1, we2, we3, we_addr_00, we_addr_11, we_addr_01, we_addr_10, we_addr_n1, we_addr_n0;
   not_gate not_gate1(we_addr[0], we_addr_n0);
   not_gate not_gate2(we_addr[1], we_addr_n1);

   and_gate and_gate0(we_addr_n0, we_addr_n1, we_addr_00);
   and_gate and_gate1(we_addr_n1, we_addr[0], we_addr_01);
   and_gate and_gate2(we_addr[1], we_addr_n0, we_addr_10);
   and_gate and_gate3(we_addr[0], we_addr[1], we_addr_11);

   and_gate and_gate4(we, we_addr_00, we0);
  and_gate and_gate5(we, we_addr_00, we1);
  and_gate and_gate6(we, we_addr_00, we2);
  and_gate and_gate7(we, we_addr_00, we3);


   wire [3:0] reg0, reg1, reg2, reg3;
  register_4bit r0(clk, we_data, we && (we_addr == 2'b00), reg0);//we && (we_addr == 2'b00)
  register_4bit r1(clk, we_data, we && (we_addr == 2'b01), reg1);//we && (we_addr == 2'b01)
  register_4bit r2(clk, we_data, we && (we_addr == 2'b10), reg2);
  register_4bit r3(clk, we_data, we && (we_addr == 2'b11), reg3);

  // Мультиплексоры для чтения и записи данных
  assign rd_data = rd_addr == 2'b00 ? reg0 :
                   rd_addr == 2'b01 ? reg1 :
                   rd_addr == 2'b10 ? reg2 :
                   reg3;

endmodule

// module d_trigger();
// endmodule


// module mem_cell(output [3:0] cell_, input [3:0] d, input wr, input clk, input reset);
//     wire not_reset;
//     not n(not_reset, reset);         // одна ячейка памяти
//     wire w3, w2, w1, w0;
//     and a3(w3, not_reset, wr, d[3]);
//     and a2(w2, not_reset, wr, d[2]);
//     and a1(w1, not_reset, wr, d[1]);
//     and a0(w0, not_reset, wr, d[0]);
//     wire w_clk;
//     and a(w_clk, wr, clk);
//     wire tr_clk;
//     or o(tr_clk, w_clk, reset);
//     d_trigger d3(cell_[3], w3, tr_clk);
//     d_trigger d2(cell_[2], w2, tr_clk);
//     d_trigger d1(cell_[1], w1, tr_clk);
//     d_trigger d0(cell_[0], w0, tr_clk);
// endmodule

// module mem(output [3:0] c0, output [3:0] c1, output [3:0] c2, output [3:0] c3, output [3:0] c4, 
//            input [3:0] d, input write, input clk, input in_c0, input in_c1, input in_c2, input in_c3, input in_c4, input reset);
//     wire w0, w1, w2, w3, w4, w5; // память из 5-ти ячеек
//     and a0(w0, clk, in_c0);
//     and a1(w1, clk, in_c1);
//     and a2(w2, clk, in_c2);
//     and a3(w3, clk, in_c3);
//     and a4(w4, clk, in_c4);
//     mem_cell m0(c0, d, write, w0, reset);
//     mem_cell m1(c1, d, write, w1, reset);
//     mem_cell m2(c2, d, write, w2, reset);
//     mem_cell m3(c3, d, write, w3, reset);
//     mem_cell m4(c4, d, write, w4, reset);
// endmodule

// module counter(clk, addr, control, immediate, data);
//   input clk; // Сигнал синхронизации
//   input [1:0] addr; // Номер значения счетчика которое читается или изменяется
//   input [3:0] immediate; // Целочисленная константа, на которую увеличивается/уменьшается значение счетчика
//   input control; // 0 - операция инкремента, 1 - операция декремента

//   output [3:0] data; // Данные из значения под номером addr, подающиеся на выход
//   // TODO: implementation
// endmodule

module counter(clk, addr, control, immediate, data);
  input clk; // Сигнал синхронизации
  input [1:0] addr; // Номер значения счетчика которое читается или изменяется
  input [3:0] immediate; // Целочисленная константа, на которую увеличивается/уменьшается значение счетчика
  input control; // 0 - операция инкремента, 1 - операция декремента

  output reg [3:0] data; // Данные из значения под номером addr, подающиеся на выход

  reg [3:0] counters [0:3]; // Регистры для хранения значений счетчиков

endmodule











// module two_comp(output [3:0] res, input [3:0] num);
//     wire [3:0] n;
//     not n_0(n[0], num[0]);
//     not n_1(n[1], num[1]);
//     not n_2(n[2], num[2]);
//     not n_3(n[3], num[3]);
//     full_adder_4b fa4(res, n, 4'b0001);
// endmodule

// module hsum(output c, output s, input a, input b);
//     and_gate a_(a, b, c);
//     xor_gate x_(a, b, s);
// endmodule

// module sum(output c_out, output s, input a, input b, input c_in);
//     wire hc, hs;
//     hsum h_sum(hc, hs, a, b);
//     wire hc2;
//     hsum h_sum_2(hc2, s, hs, c_in);
//     or_gate o_(c_out, hc, hc2);
// endmodule

// module full_adder_3b(output [3:0] res, input [2:0] a, input [2:0] b);
//     wire c_in_1, c_in_2;
//     hsum h(c_in_1, res[0], a[0], b[0]);
//     sum s1(c_in_2, res[1], a[1], b[1], c_in_1);
//     sum s2(res[3], res[2], a[2], b[2], c_in_2);
// endmodule

// module full_adder_4b(output [3:0] res, input [3:0] a, input [3:0] b);
//     wire [3:0] sum_3b;
//     full_adder_3b f3(sum_3b, a[2:0], b[2:0]);
//     assign res[2:0] = sum_3b[2:0];
//     wire w;
//     sum s(w, res[3], a[3], b[3], sum_3b[3]);
// endmodule

module and_gate_for(a, b, out_alu);
input [3:0] a, b;
output [3:0] out_alu;
wire a1, a2, a3, a0;
and_gate and_gate0(a[0], b[0], a0);
and_gate and_gate1(a[1], b[1], a1);
and_gate and_gate2(a[2], b[2], a2);
and_gate and_gate3(a[3], b[3], a3);
assign out_alu[0] = a0;
assign out_alu[1] = a1;
assign out_alu[2] = a2;
assign out_alu[3] = a3;
    
endmodule

module or_gate_for(a, b, out_alu);
input [3:0] a, b;
output [3:0] out_alu;
wire a0, a1, a2, a3;
or_gate or_gate0(a[0], b[0], a0);
or_gate or_gate1(a[1], b[1], a1);
or_gate or_gate2(a[2], b[2], a2);
or_gate or_gate3(a[3], b[3], a3);
assign out_alu[0] = a0;
assign out_alu[1] = a1;
assign out_alu[2] = a2;
assign out_alu[3] = a3;
endmodule

module xor_gate(a, b, xor_out);
input wire a, b;
output wire xor_out;
wire not_a, not_b, and_ab_1, and_ab_2;
//assign xor_out = a&~b|~a&b
not_gate not_gate1(a, not_a);
not_gate not_gate2(b, not_b);
and_gate and_gate3(a, not_b, and_ab_1);
and_gate and_gate4(b, not_a, and_ab_2);
or_gate or_gate1(and_ab_1, and_ab_2, xor_out);
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