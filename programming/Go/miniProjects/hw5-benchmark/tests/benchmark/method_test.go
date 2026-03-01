package benchmark

import (
	"reflect"
	"testing"
)

type Student struct{ n string }

func (s Student) Name() string { return s.n }

type Namer interface {
	Name() string
}

func BenchmarkStudentName_Direct(b *testing.B) {
	s := Student{n: "alice"}
	var sink string

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sink = s.Name()
	}

	_ = sink
}

func BenchmarkStudentName_Interface(b *testing.B) {
	var s Namer = Student{n: "alice"}
	var sink string

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sink = s.Name()
	}

	_ = sink
}

func BenchmarkStudentName_MethodValue(b *testing.B) {
	s := Student{n: "alice"}
	methodValue := s.Name
	var sink string

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sink = methodValue()
	}

	_ = sink
}

func BenchmarkStudentName_MethodExpression(b *testing.B) {
	s := Student{n: "alice"}
	methodExpr := Student.Name
	var sink string

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sink = methodExpr(s)
	}

	_ = sink
}

func BenchmarkStudentName_Reflection(b *testing.B) {
	s := Student{n: "alice"}
	var sink string

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		val := reflect.ValueOf(s)
		method := val.MethodByName("Name")
		result := method.Call(nil)
		sink = result[0].String()
	}

	_ = sink
}

func BenchmarkStudentName_ReflectionPrepared(b *testing.B) {
	s := Student{n: "alice"}
	val := reflect.ValueOf(s)
	method := val.MethodByName("Name")
	var sink string

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result := method.Call(nil)
		sink = result[0].String()
	}

	_ = sink
}

func BenchmarkStudentName_ReflectionType(b *testing.B) {
	s := Student{n: "alice"}
	t := reflect.TypeOf(s)
	method, _ := t.MethodByName("Name")
	val := reflect.ValueOf(s)
	var sink string

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result := method.Func.Call([]reflect.Value{val})
		sink = result[0].String()
	}

	_ = sink
}

func BenchmarkStudentName_ReflectionInterface(b *testing.B) {
	var s Namer = Student{n: "alice"}
	val := reflect.ValueOf(s)
	method := val.MethodByName("Name")
	var sink string

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result := method.Call(nil)
		sink = result[0].String()
	}

	_ = sink
}

func BenchmarkStudentName_TypeAssertion(b *testing.B) {
	var iface interface{} = Student{n: "alice"}
	var sink string

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s := iface.(Student)
		sink = s.Name()
	}

	_ = sink
}

func BenchmarkStudentName_InterfaceTypeAssertion(b *testing.B) {
	var iface interface{} = Student{n: "alice"}
	var sink string

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s := iface.(Namer)
		sink = s.Name()
	}

	_ = sink
}

type Calculator struct{}

func (c Calculator) Add(a, b int) int {
	return a + b
}

func BenchmarkCalculator_Add_Direct(b *testing.B) {
	c := Calculator{}
	var sink int

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sink = c.Add(10, 20)
	}

	_ = sink
}

func BenchmarkCalculator_Add_Reflection(b *testing.B) {
	c := Calculator{}
	val := reflect.ValueOf(c)
	method := val.MethodByName("Add")
	args := []reflect.Value{
		reflect.ValueOf(10),
		reflect.ValueOf(20),
	}
	var sink int

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result := method.Call(args)
		sink = int(result[0].Int())
	}

	_ = sink
}

type PointerStudent struct{ n string }

func (s *PointerStudent) Name() string { return s.n }

func BenchmarkPointerStudentName_Direct(b *testing.B) {
	s := &PointerStudent{n: "alice"}
	var sink string

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sink = s.Name()
	}

	_ = sink
}

func BenchmarkPointerStudentName_Reflection(b *testing.B) {
	s := &PointerStudent{n: "alice"}
	val := reflect.ValueOf(s)
	method := val.MethodByName("Name")
	var sink string

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result := method.Call(nil)
		sink = result[0].String()
	}

	_ = sink
}

func BenchmarkStudentName_Direct_Allocs(b *testing.B) {
	s := Student{n: "alice"}
	var sink string

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sink = s.Name()
	}

	_ = sink
}

func BenchmarkStudentName_Reflection_Allocs(b *testing.B) {
	s := Student{n: "alice"}
	val := reflect.ValueOf(s)
	method := val.MethodByName("Name")
	var sink string

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result := method.Call(nil)
		sink = result[0].String()
	}

	_ = sink
}

func BenchmarkDifferentStudents(b *testing.B) {
	students := []Student{
		{n: "alice"},
		{n: "bob"},
		{n: "charlie"},
		{n: "diana"},
		{n: "eve"},
	}
	var sink string

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		student := students[i%len(students)]
		sink = student.Name()
	}

	_ = sink
}
