// module alu(a, b, control, res);
//   input [3:0] a, b;
//   input [2:0] control; // Выбор операции. Стоит поменять в зависимости от того, как нам подают инф-у

//   output reg [3:0] res; 

//   
//   output reg zero;

//   always @*
//     case (control)
//       3'b000: res = a & b;
//       3'b001: res = ~(a & b);
//       3'b010: res = a | b;
//       3'b011: res = ~(a | b);
//       3'b100: res = a + b;
//       3'b101: res = a - b;
//       3'b110: begin
//         
//         res = ((a[3] == b[3]) & (a[2:0] < b[2:0]) | (a[3] == 1'b1 & b[3] == 1'b0)) ? 4'b0001 : 4'b0000;
//       end
//       default: res = 4'b0; 
//     endcase

//   
//   //assign zero = (res == 32'b0);

// endmodule




// module alu(a, b, control, res, zero);
//   input signed [31:0] a, b;
//   input [2:0] control; // Выбор операции. Стоит поменять в зависимости от того, как нам подают инф-у

//   output reg [31:0] res; 

//   output reg zero;

//   always @*
//     case (control[1:0])
//       3'b000: res = a & b;
//       3'b001: res = a | b;
//       3'b010: res = a + b;
//       3'b011: res = a - b;
//       3'b110: begin
//       
//         res = (a > b) ? 32'b0000_0000_0000_0000_0000_0000_0000_0001 : 32'b0000_0000_0000_0000_0000_0000_0000_0000;
//       end
//       default: res = 32'b0; 
//     endcase

//  
//   assign zero = (res == 32'b0);

// endmodule


module alu(a, b, control, res, zero);
  input signed [31:0] a, b;
  input [2:0] control;
  output reg [31:0] res;
  output reg zero;

  reg [31:0] tempb;

  always @(control or a or b) begin
    if (control[2] == 0) 
    begin 
      tempb = b;
    end else begin
      tempb = ~b;
    end

    case (control[1:0])
      0:
         res = a & tempb;
      1:
         res = a | tempb;
      2:
         res = a + tempb + control[2];
      3:
        if (control[2] == 1) begin
          if (a < b) begin
             res = 1;
          end else begin
             res = 0;
          end
        end 
      endcase

      if (res == 0) begin
        zero = 1;
      end else begin
        zero = 0;
      end
  end

endmodule




// module alu(a, b, control, res, zero);
//   input signed [31:0] a, b;
//   input [2:0] control;
//   output reg [31:0] res;
//   output reg zero;

//   reg [31:0] tempb;

//   always @(control or a or b) begin

//     case (control[1:0])
//       3'b010:
//         res = a + b;
//       3'b000:
//          res = a & tempb;
//       3'b001:
//          res = a | tempb;
//       3'b110:
//          res = a - b;
//       3'b111:
//         res = ((a[31] == b[31]) & (a[30:0] < b[30:0]) | (a[31] == 1'b1 & b[31] == 1'b0)) ? 32'b0000_0000_0000_0000_0000_0000_0000_0001 : 32'b0000_0000_0000_0000_0000_0000_0000_0000;
//       endcase

//       if (res == 0) begin
//         zero = 1;
//       end else begin
//         zero = 0;
//       end
//   end

// endmodule






// module alu(a, b, control, res, zero);
//   input [31:0] a, b;
//   input [2:0] control; //выбор операции. Стоит поменять в зависимости от того, как нам подают инф-у

//   output reg [31:0] res; 


// // wire [31:0] and_res, nand_res, or_res, nor_res, add_res,inv_b, add_res_1, sub_res, slt_res;

// //   assign and_res = a & b;
// //   assign nand_res = ~(a & b);
// //   assign or_res = a | b;
// //   assign nor_res = ~(a | b);
// //   assign add_res = a + b;
// //   assign sub_res = a - b;
// //   assign slt_res = ((a[31] == b[31]) & (a[30:0] < b[30:0]) | (a[31] == 1'b1 & b[31] == 1'b0)) ? 4'b00000000000000000000000000000001 : 4'b00000000000000000000000000000000;

// integer index;
//  case (control)
//     3'b000: begin
//         res = a & b;
//     end
//     3'b001: begin
//         res = ~(a & b);
//     end
//     3'b010: begin
//         res = a | b;
//     end
//     3'b011: begin
//         res = ~(a | b);
//     end
//     3'b100: begin
//         res = a + b;
//     end
//     3'b101: begin
//         res = a - b;
//     end
//     3'b110: begin
//         res = ((a[31] == b[31]) & (a[30:0] < b[30:0]) | (a[31] == 1'b1 & b[31] == 1'b0)) ? 4'b00000000000000000000000000000001 : 4'b00000000000000000000000000000000;
//     end

//  endcase


// //   wire sel0, sel1, sel2, sel3, sel4, sel5, sel6;
// //   assign sel0 = ~control[2] & ~control[1] & ~control[0];
// //   assign sel1 = ~control[2] & ~control[1] & control[0];
// //   assign sel2 = ~control[2] & control[1] & ~control[0];
// //   assign sel3 = ~control[2] & control[1] & control[0];
// //   assign sel4 = control[2] & ~control[1] & ~control[0];
// //   assign sel5 = control[2] & ~control[1] & control[0];
// //   assign sel6 = control[2] & control[1] & ~control[0];
  
// //   assign res[0] = (and_res[0] & sel0) | (nand_res[0] & sel1) | (or_res[0] & sel2) | (nor_res[0] & sel3) | (add_res[0] & sel4) | (sub_res[0] & sel5) | (slt_res[0] & sel6);
// //   assign res[1] = (and_res[1] & sel0) | (nand_res[1] & sel1) | (or_res[1] & sel2) | (nor_res[1] & sel3) | (add_res[1] & sel4) | (sub_res[1] & sel5) | (slt_res[1] & sel6);
// //   assign res[2] = (and_res[2] & sel0) | (nand_res[2] & sel1) | (or_res[2] & sel2) | (nor_res[2] & sel3) | (add_res[2] & sel4) | (sub_res[2] & sel5) | (slt_res[2] & sel6);
// //   assign res[3] = (and_res[3] & sel0) | (nand_res[3] & sel1) | (or_res[3] & sel2) | (nor_res[3] & sel3) | (add_res[3] & sel4) | (sub_res[3] & sel5) | (slt_res[3] & sel6);

// endmodule


