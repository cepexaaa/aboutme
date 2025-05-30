`include "implication.v"

module mux_4_1(data0, control, out);
    input [31:0] control;
    input data0, data1, data2, data3;
    output reg out;
    always @(data0, control) begin
        case (control)
            0: out <= data0;
            // 1: out <= data1;
            // 2: out <= data2;
            // 3: out <= data3;
        endcase
    end
endmodule

module tester2();
    reg a, b, oute;			    // values from testvectors
    //wire [1:0] min_out, max_out, any_out, consensus_out;
    wire out, imp_out;
    reg [31:0] vectornum, errors, i;  // bookkeeping variables
    reg [5:0] testvectors[0:3];		// array of testvectors
    reg verdict;
    implication_mos implication_mos1(a, b, imp_out);
    // ternary_min min_to_test(a, b, min_out);
    // ternary_max max_to_test(a, b, max_out);
    // ternary_any any_to_test(a, b, any_out);
    // ternary_consensus consensus_to_test(a, b, consensus_out);
    mux_4_1 mux(imp_out, i, out);
    initial begin
        verdict = 1;
        for (i = 0; i < 1; i = i + 1) begin
            case (i)
                0: begin $display("Test implication"); $readmemb("imp.mem", testvectors); end
                // 0: begin $display("Test min"); $readmemb("min.mem", testvectors); end
                // 1: begin $display("Test max"); $readmemb("max.mem", testvectors); end
                // 2: begin $display("Test any"); $readmemb("any.mem", testvectors); end
                // 3: begin $display("Test consensus"); $readmemb("consensus.mem", testvectors); end
            endcase
            errors = 0;
            for (vectornum = 0; vectornum < 4; vectornum = vectornum + 1) begin
                {a, b, oute} = testvectors[vectornum];
                #1
                if (out !== oute)
                    begin
                        $display("Error: inputs a=0b%b b=0b%b", a, b);
                        $display("  outputs: expected 0b%b, actual 0b%b", oute, out);
                        errors = 1 + errors;
                    end
            end
            $display("%d tests completed with %d errors", vectornum, errors);
            // End simulation:
            if (errors !== 0)
                begin
                    $display("FAIL");
                    verdict = 0;
                end
            else
                begin
                    $display("OK");
                end
            #5;
        end
        if (verdict) $display("Verdict: OK");
        else $display("Verdict: FAIL");
        $finish;
    end
endmodule