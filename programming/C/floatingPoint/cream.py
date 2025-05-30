import subprocess

def test():
    type = input("Enter test-type\nTests supported: mult, add, divide\n->")
    if type == "mult":
        file_name = "source\\true_gen_float_mult_tests.tsv"
    elif type == "add":
        file_name = "source\\true_gen_float_+-_tests.tsv"
    elif type == "divide":
        file_name = "source\\true_gen_float_div_tests.tsv"
    else:
        print("chleni?")
        return
    
    path = input("ENTER PATH TO EXE FILE:\n->") #change to your path
    test_f = open(file_name, "r")
    tests = test_f.readlines()

    half = 0
    single = 0
    half_pass = 0
    single_pass = 0
    count = 0
    passed = 0
    for test_id in range(0, len(tests)):
        arg = tests[test_id][:-1]
        arg = arg.split(',')
        if (len(arg) == 0):
            continue
        res = arg[-1]
        arg = arg[0].split()
        arg = [path] + arg
        output = subprocess.check_output(arg).decode("utf-8")
        count += 1
        if (arg[1] == 'h'):
            half += 1
        elif (arg[1] == 'f'):
            single += 1
        if res.strip('\r\n ') != output.strip('\r\n '):
            print("FAILED test: " + str(test_id))
            for i in range(1, len(arg)):
                print(arg[i], end=" ")
            print("Expected:", res, end=" ")
            print("Actual:", output)
        else:
            if (arg[1] == 'h'):
                half_pass += 1
            elif (arg[1] == 'f'):
                single_pass += 1
            passed += 1
            #print("PASSED")

        if (count % 1000 == 0):
            print('----------------------------------')
            print("STATS:" + str(passed) + '/' + str(count) )
            print("HALF PASSED: " + str(half_pass) + "/" + str(half))
            if single > 0:
                print("SINGLE PASSED: " + str(single_pass) + "/" + str(single))
            print('----------------------------------')
    print(str(passed) + "/" + str(count) + " PASSED")
test()