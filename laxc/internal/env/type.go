package env

import "fmt"

type Type interface {
	fmt.Stringer
	Name() string
	Equals(Type) bool
	Dereference() (Type, error)
	Coerce() Type
	IsSubtypeOf(Type) bool
}

type IntegerType struct{}

func (t IntegerType) Name() string {
	return "integer"
}

func (t IntegerType) String() string {
	return t.Name()
}

func (t IntegerType) Equals(other Type) bool {
	return t.String() == other.String()
}

func (t IntegerType) Dereference() (Type, error) {
	return nil, fmt.Errorf("%s cannot be dereferenced", t.Name())
}

func (t IntegerType) Coerce() Type {
	return RealType{}
}

func (t IntegerType) IsSubtypeOf(t2 Type) bool {
	return t.Equals(t2) || t.Coerce().IsSubtypeOf(t2)
}

type RealType struct{}

func (t RealType) Name() string {
	return "real"
}

func (t RealType) String() string {
	return t.Name()
}

func (t RealType) Equals(other Type) bool {
	return t.String() == other.String()
}

func (t RealType) Dereference() (Type, error) {
	return nil, fmt.Errorf("%s cannot be dereferenced", t.Name())
}

func (t RealType) Coerce() Type {
	return VoidType{}
}

func (t RealType) IsSubtypeOf(t2 Type) bool {
	return t.Equals(t2) || t.Coerce().IsSubtypeOf(t2)
}

type BooleanType struct {
}

func (t BooleanType) Name() string {
	return "boolean"
}

func (t BooleanType) String() string {
	return t.Name()
}

func (t BooleanType) Equals(other Type) bool {
	return t.String() == other.String()
}

func (t BooleanType) Dereference() (Type, error) {
	return nil, fmt.Errorf("%s cannot be dereferenced", t.Name())
}

func (t BooleanType) Coerce() Type {
	return VoidType{}
}

func (t BooleanType) IsSubtypeOf(t2 Type) bool {
	return t.Equals(t2) || t.Coerce().IsSubtypeOf(t2)
}

type VoidType struct{}

func (t VoidType) Name() string {
	return "void"
}

func (t VoidType) String() string {
	return t.Name()
}

func (t VoidType) Equals(other Type) bool {
	return t.String() == other.String()
}

func (t VoidType) Dereference() (Type, error) {
	return nil, fmt.Errorf("%s cannot be dereferenced", t.Name())
}

func (t VoidType) Coerce() Type {
	return VoidType{}
}

func (t VoidType) IsSubtypeOf(t2 Type) bool {
	return t.Equals(t2)
}

type RefType struct {
	Underlying Type
}

func (t RefType) Name() string {
	return "ref " + t.Underlying.Name()
}

func (t RefType) String() string {
	return t.Name()
}

func (t RefType) Equals(other Type) bool {
	return t.String() == other.String()
}

func (t RefType) Dereference() (Type, error) {
	return t.Underlying, nil
}

func (t RefType) Coerce() Type {
	return t.Underlying
}

func (t RefType) IsSubtypeOf(t2 Type) bool {
	return t.Equals(t2) || t.Coerce().IsSubtypeOf(t2)
}

func GreatestCommonDenominator(t1, t2 Type) Type {
	if t1.IsSubtypeOf(t2) {
		return t2
	} else if t2.IsSubtypeOf(t1) {
		return t1
	}

	for !t1.IsSubtypeOf(t2) {
		t2 = t2.Coerce()
	}

	return t2
}
