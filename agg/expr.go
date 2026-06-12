package agg

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Expr is the sealed base interface for all aggregation expressions.
// The unexported method prevents external implementations; use the
// constructor functions and field-ref helpers in this package.
type Expr interface{ expr() }

// NumericExpr is an Expr that resolves to a number.
type NumericExpr interface {
	Expr
	numericExpr()
}

// BoolExpr is an Expr that resolves to a boolean.
type BoolExpr interface {
	Expr
	boolExpr()
}

// StringExpr is an Expr that resolves to a string.
type StringExpr interface {
	Expr
	stringExpr()
}

// ArrayExpr is an Expr that resolves to an array.
type ArrayExpr interface {
	Expr
	arrayExpr()
}

// DateExpr is an Expr that resolves to a date.
type DateExpr interface {
	Expr
	dateExpr()
}

// ObjectExpr is an Expr that resolves to an object/document.
type ObjectExpr interface {
	Expr
	objectExpr()
}

// --- concrete expression types ---
// Each holds an `any` internally so it can carry either a field-path string
// or an operator document (bson.D), both of which marshal correctly.

type numericExprVal struct{ v any }

func (e numericExprVal) expr()        {}
func (e numericExprVal) numericExpr() {}
func (e numericExprVal) MarshalBSONValue() (byte, []byte, error) {
	return marshalExprValue(e.v)
}

type boolExprVal struct{ v any }

func (e boolExprVal) expr()     {}
func (e boolExprVal) boolExpr() {}
func (e boolExprVal) MarshalBSONValue() (byte, []byte, error) {
	return marshalExprValue(e.v)
}

type stringExprVal struct{ v any }

func (e stringExprVal) expr()       {}
func (e stringExprVal) stringExpr() {}
func (e stringExprVal) MarshalBSONValue() (byte, []byte, error) {
	return marshalExprValue(e.v)
}

type arrayExprVal struct{ v any }

func (e arrayExprVal) expr()      {}
func (e arrayExprVal) arrayExpr() {}
func (e arrayExprVal) MarshalBSONValue() (byte, []byte, error) {
	return marshalExprValue(e.v)
}

type dateExprVal struct{ v any }

func (e dateExprVal) expr()     {}
func (e dateExprVal) dateExpr() {}
func (e dateExprVal) MarshalBSONValue() (byte, []byte, error) {
	return marshalExprValue(e.v)
}

type objectExprVal struct{ v any }

func (e objectExprVal) expr()       {}
func (e objectExprVal) objectExpr() {}
func (e objectExprVal) MarshalBSONValue() (byte, []byte, error) {
	return marshalExprValue(e.v)
}

// genericExprVal backs Expr (base interface) return types.
type genericExprVal struct{ v any }

func (e genericExprVal) expr() {}
func (e genericExprVal) MarshalBSONValue() (byte, []byte, error) {
	return marshalExprValue(e.v)
}

// --- field reference helpers ---
// Pass the bare field name (no leading "$"); the helper adds the prefix.

// Field returns a field reference as a base Expr.
// Use typed variants (NumericField, BoolField, etc.) when the operator
// requires a specific expression category.
func Field(name string) Expr {
	return genericExprVal{v: "$" + name}
}

// NumericField returns a field reference that resolves to a number.
func NumericField(name string) NumericExpr {
	return numericExprVal{v: "$" + name}
}

// BoolField returns a field reference that resolves to a boolean.
func BoolField(name string) BoolExpr {
	return boolExprVal{v: "$" + name}
}

// StringField returns a field reference that resolves to a string.
func StringField(name string) StringExpr {
	return stringExprVal{v: "$" + name}
}

// ArrayField returns a field reference that resolves to an array.
func ArrayField(name string) ArrayExpr {
	return arrayExprVal{v: "$" + name}
}

// DateField returns a field reference that resolves to a date.
func DateField(name string) DateExpr {
	return dateExprVal{v: "$" + name}
}

// Literal wraps v in a $literal expression so the aggregation engine never
// interprets it as a field path or operator.
func Literal(v any) Expr {
	return genericExprVal{v: bson.D{{Key: "$literal", Value: v}}}
}

// marshalExprValue is the shared marshal implementation for all expression types.
func marshalExprValue(v any) (byte, []byte, error) {
	typ, b, err := bson.MarshalValue(v)
	return byte(typ), b, err
}
