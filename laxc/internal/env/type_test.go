package env

import "testing"

func assert(t *testing.T, name string, actual bool) {
	t.Run(name, func(t *testing.T) {
		if !actual {
			t.Fail()
		}
	})
}

func TestCoercion(t *testing.T) {
	assert(t, "ref boolean < boolean", RefType{BooleanType{}}.IsSubtypeOf(BooleanType{}))
	assert(t, "ref integer < integer", RefType{IntegerType{}}.IsSubtypeOf(IntegerType{}))
	assert(t, "ref integer < real", RefType{IntegerType{}}.IsSubtypeOf(RealType{}))
	assert(t, "ref integer < void", RefType{IntegerType{}}.IsSubtypeOf(VoidType{}))
	assert(t, "!(ref integer < ref boolean)", !RefType{IntegerType{}}.IsSubtypeOf(RefType{BooleanType{}}))
	assert(t, "!(ref integer < ref real)", !RefType{IntegerType{}}.IsSubtypeOf(RefType{RealType{}}))
	assert(t, "!(ref real < ref integer)", !RefType{RealType{}}.IsSubtypeOf(RefType{IntegerType{}}))
}

func TestGreatestCommonDenominator(t *testing.T) {
	assert(
		t, "GCD(ref integer < ref real) = real",
		GreatestCommonDenominator(
			RefType{RealType{}},
			RefType{IntegerType{}},
		).Equals(RealType{}),
	)
	assert(
		t, "GCD(ref integer < ref boolean) = void",
		GreatestCommonDenominator(
			RefType{RealType{}},
			RefType{BooleanType{}},
		).Equals(VoidType{}),
	)
}
