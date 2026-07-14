package agg

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Number interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int |
		~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint | ~uintptr |
		~float32 | ~float64
}

type ArrayResolver interface {
	AnyExpr | ArrayExpr | string | []any | bson.A | []string | []int | []int32 | []int64 | []float64
}

type NumberResolver interface {
	AnyExpr | NumberExpr | Number | string
}

type StringResolver interface {
	AnyExpr | StringExpr | string
}

type BoolResolver interface {
	AnyExpr | BoolExpr | bool | string
}

type ObjectResolver interface {
	AnyExpr | ObjectExpr | string | bson.D | bson.M
}

type DateResolver interface {
	AnyExpr | DateExpr | time.Time | string
}

type TimestampResolver interface {
	AnyExpr | bson.Timestamp | string
}

type ObjectIDResolver interface {
	AnyExpr | bson.ObjectID | string
}

type Expr any

type Option[T any] func(*T)

type AnyExpr struct {
	expr Expr
}

func (ae AnyExpr) MarshalBSONValue() (byte, []byte, error) {
	if ae.expr == nil {
		return 0x0A, nil, nil // BSON null
	}
	typ, b, err := bson.MarshalValue(ae.expr)
	return byte(typ), b, err
}

// Null is an AnyExpr that evaluates to BSON null. Use as the _id in
// GroupStage to accumulate all documents into a single group.
var Null = AnyExpr{expr: nil}

type NumberExpr struct {
	expr Expr
}

func (ne NumberExpr) MarshalBSONValue() (byte, []byte, error) {
	typ, b, err := bson.MarshalValue(ne.expr)
	return byte(typ), b, err
}

type ArrayExpr struct {
	expr Expr
}

func (ae ArrayExpr) MarshalBSONValue() (byte, []byte, error) {
	typ, b, err := bson.MarshalValue(ae.expr)
	return byte(typ), b, err
}

func Array[T any](values []T) ArrayExpr {
	return ArrayExpr{
		expr: values,
	}
}

type StringExpr struct {
	expr Expr
}

func (se StringExpr) MarshalBSONValue() (byte, []byte, error) {
	typ, b, err := bson.MarshalValue(se.expr)
	return byte(typ), b, err
}

type ObjectExpr struct {
	expr Expr
}

func (oe ObjectExpr) MarshalBSONValue() (byte, []byte, error) {
	typ, b, err := bson.MarshalValue(oe.expr)
	return byte(typ), b, err
}

func RootObject() ObjectExpr {
	return ObjectExpr{
		expr: "$$ROOT",
	}
}

type BoolExpr struct {
	expr Expr
}

func (be BoolExpr) MarshalBSONValue() (byte, []byte, error) {
	typ, b, err := bson.MarshalValue(be.expr)
	return byte(typ), b, err
}

type DateExpr struct {
	expr Expr
}

func (de DateExpr) MarshalBSONValue() (byte, []byte, error) {
	typ, b, err := bson.MarshalValue(de.expr)
	return byte(typ), b, err
}
