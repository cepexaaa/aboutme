package benchmark

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {

	flag.Parse()

	if testing.Short() {
		fmt.Println("Skipping benchmarks in short mode")
		os.Exit(0)
	}

	os.Exit(m.Run())
}

func TestAllMethodsCorrectness(t *testing.T) {
	s := Student{n: "test"}

	if s.Name() != "test" {
		t.Errorf("Direct call failed")
	}

	var namer Namer = s
	if namer.Name() != "test" {
		t.Errorf("Interface call failed")
	}

	methodValue := s.Name
	if methodValue() != "test" {
		t.Errorf("Method value failed")
	}

	methodExpr := Student.Name
	if methodExpr(s) != "test" {
		t.Errorf("Method expression failed")
	}

	val := reflect.ValueOf(s)
	method := val.MethodByName("Name")
	result := method.Call(nil)
	if result[0].String() != "test" {
		t.Errorf("Reflection call failed")
	}
}
