
module control_unit(opcode, funct, MemtoReg, MemWrite, BranchN, BranchE, ALUSrc, RegDst, RegWrite, ALUControl, Jump, Jal, Jr);

//We need to accept the data and return flags

input wire [5:0] opcode, funct;
output wire MemtoReg, MemWrite, BranchN, BranchE, ALUSrc, RegDst, RegWrite, Jump, Jal, Jr;
output reg [2:0] ALUControl;

reg tempMemtoReg, tempMemWrite, tempBranchN, tempBranchE, tempALUSrc, tempRegDst, tempRegWrite, tempJump, tempJal, tempJr;

reg [1:0] ALUop;


always @* begin

case (opcode)
  // R-type
  0: begin
    tempRegWrite = 1;
    tempRegDst = 1;
    tempALUSrc = 0;
    tempMemtoReg = 0;
    tempJal = 0;
    tempBranchN = 0;
    tempBranchE = 0;
    tempMemWrite = 0;
    ALUop = 2'b10;
    tempJump = 0;
    tempJr = 0;
  end
  
  // lw
  35: begin
    tempRegWrite = 1;
    tempRegDst = 0;
    tempALUSrc = 1;
    tempMemtoReg = 1;
    tempJal = 0;
    tempBranchN = 0;
    tempBranchE = 0;
    tempMemWrite = 0;
    ALUop = 2'b00;
    tempJump = 0;
    tempJr = 0;
  end

  // sw
  43: begin
    tempRegWrite = 0;
    tempRegDst = 0;
    tempALUSrc = 1;
    tempMemtoReg = 0;
    tempJal = 0;
    tempBranchN = 0;
    tempBranchE = 0;
    tempMemWrite = 1;
    ALUop = 2'b00;
    tempJump = 0;
    tempJr = 0;
  end

  // beq
  4: begin
    tempRegWrite = 0;
    tempRegDst = 0;
    tempALUSrc = 0;
    tempMemtoReg = 0;
    tempJal = 0;
    tempBranchN = 0;
    tempBranchE = 1;
    tempMemWrite = 0;
    ALUop = 2'b01;
    tempJump = 0;
    tempJr = 0;
  end

  // bne
  5: begin
    tempRegWrite = 0;
    tempRegDst = 0;
    tempALUSrc = 0;
    tempMemtoReg = 0;
    tempJal = 0;
    tempBranchN = 1;
    tempBranchE = 0;
    tempMemWrite = 0;
    ALUop = 2'b01;
    tempJump = 0;
    tempJr = 0;
  end

  // addi
  8: begin
    tempRegWrite = 1;
    tempRegDst = 0;
    tempALUSrc = 1;
    tempMemtoReg = 0;
    tempJal = 0;
    tempBranchN = 0;
    tempBranchE = 0;
    tempMemWrite = 0;
    ALUop = 2'b00;
    tempJump = 0;
    tempJr = 0;
  end

  // andi
  12: begin
    tempRegWrite = 1;
    tempRegDst = 0;
    tempALUSrc = 1;
    tempMemtoReg = 0;
    tempJal = 0;
    tempBranchN = 0;
    tempBranchE = 0;
    tempMemWrite = 0;
    ALUop = 2'b11;
    tempJump = 0;
    tempJr = 0;
  end

  // j
  2: begin
    tempRegWrite = 0;
    tempRegDst = 0;
    tempALUSrc = 0;
    tempMemtoReg = 0;
    tempJal = 0;
    tempBranchN = 0;
    tempBranchE = 0;
    tempMemWrite = 0;
    ALUop = 2'b10;
    tempJump = 1;
    tempJr = 0;
  end

  // jal
  3: begin
    tempRegWrite = 1;
    tempRegDst = 0;
    tempALUSrc = 0;
    tempMemtoReg = 0;
    tempJal = 1;
    tempBranchN = 0;
    tempBranchE = 0;
    tempMemWrite = 0;
    ALUop = 2'b10;
    tempJump = 1;
    tempJr = 0;
  end

  // if we catch exception, we return 000000000000.....
  default: begin
    tempRegWrite = 0;
    tempRegDst = 0;
    tempALUSrc = 0;
    tempMemtoReg = 0;
    tempJal = 0;
    tempBranchN = 0;
    tempBranchE = 0;
    tempMemWrite = 0;
    ALUop = 0;
    tempJump = 0;
    tempJr = 0;
  end
endcase


if (ALUop == 2'b00)
begin
  ALUControl = 3'b010;
end else if (ALUop == 2'b11) 
begin
  ALUControl = 3'b000;
end else if (ALUop == 2'b01)
begin
  ALUControl = 3'b110;
end else begin
  case (funct)
  6'b100000: ALUControl = 3'b010;  // +
  6'b100010: ALUControl = 3'b110;  // -
  6'b100100: ALUControl = 3'b000;  // &
  6'b100101: ALUControl = 3'b001;  // |
  6'b101010: ALUControl = 3'b111;  // slt
  6'b001000: begin 
    tempJr = 1;
    tempRegWrite = 0;
  end
  endcase
end
end

assign MemtoReg = tempMemtoReg;
assign RegDst = tempRegDst;
assign ALUSrc = tempALUSrc;
assign RegWrite = tempRegWrite;
assign BranchE = tempBranchE;
assign MemWrite = tempMemWrite;
assign BranchN = tempBranchN;
assign Jump = tempJump;
assign Jal = tempJal;
assign Jr = tempJr; 

endmodule