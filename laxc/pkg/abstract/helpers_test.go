package abstract

import (
	"laxc/internal/env"
	"laxc/pkg/attributed"
	"testing"
)

func Test_coerceBothToOneOf(t *testing.T) {
	a := attributed.Assignment{
		Var: &env.Variable{
			Name: "a",
			Type: env.RefType{Underlying: env.IntegerType{}},
		},
	}
	b := attributed.Assignment{
		Var: &env.Variable{
			Name: "b",
			Type: env.RefType{Underlying: env.RealType{}},
		},
	}

	aOut, bOut, err := coerceBothToOneOf(a, b, env.IntegerType{}, env.RealType{})
	if err != nil {
		t.Fatal(err)
	}

	if !aOut.Type().Equals(bOut.Type()) {
		t.Fatalf("%s should equal %s", aOut.Type(), bOut.Type())
	}

	if !aOut.Type().Equals(env.RealType{}) {
		t.Fatalf("should have been coerced to %s but was %s", env.RealType{}, aOut.Type())
	}
}
