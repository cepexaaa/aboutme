#!/bin/bash

echo "Run Benchmark-tests..."
echo "=========================="

mkdir -p tests/benchmark_results

echo "1. All tests (5 sec for test):"
go test -bench=. -benchmem -benchtime=5s ./tests/benchmark/... -count=3 | tee benchmark_results/all_benchmarks.txt

echo -e "\n2. Only main tests (Student):"
go test -bench="BenchmarkStudentName" -benchmem -benchtime=5s ./tests/benchmark/... -count=3 | tee benchmark_results/student_benchmarks.txt

echo -e "\n3. With allocations:"
go test -bench=".*Allocs" -benchmem -benchtime=5s ./tests/benchmark/... -count=3 | tee benchmark_results/allocation_benchmarks.txt

echo -e "\n4.reflection vs direct:"
go test -bench=".*Direct|.*Reflection" -benchmem -benchtime=5s ./tests/benchmark/... -count=3 | tee benchmark_results/reflection_vs_direct.txt

echo -e "\nResult in benchmark_results/"