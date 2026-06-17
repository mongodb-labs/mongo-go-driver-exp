// Package query provides a typed builder for MongoDB query filters,
// intended for use as the argument to agg.MatchStage.
//
// Only a starter set of field conditions is implemented; this package is
// designed to grow independently of the aggregation expression system.
package query

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// FieldCondition represents a condition applied to a single document field,
// e.g. { $eq: value }. Construct via the operator functions (Eq, etc.).
type FieldCondition struct{ doc bson.D }

// Filter represents a complete MongoDB query document, e.g. { field: { $eq: v } }.
// Construct via Field or (in future) logical combinators like And, Or.
type Filter struct{ doc bson.D }

// MarshalBSON implements bson.Marshaler so Filter can be passed directly as
// a bson.D element value.
func (f Filter) MarshalBSON() ([]byte, error) {
	return bson.Marshal(f.doc)
}

// Field creates a Filter that matches documents where the named field
// satisfies cond, e.g. { name: { $eq: value } }.
func Field(name string, cond FieldCondition) Filter {
	return Filter{doc: bson.D{{Key: name, Value: cond.doc}}}
}

// Eq creates a FieldCondition for equality: { $eq: value }.
// value may be any BSON-marshalable Go value.
func Eq(value any) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$eq", Value: value}}}
}

// Gt creates a FieldCondition for greater than: { $gt: value }.
// value may be any BSON-marshalable Go value.
func Gt(value any) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$gt", Value: value}}}
}

// Gte creates a FieldCondition for greater than or equal to: { $gte: value }.
// value may be any BSON-marshalable Go value.
func Gte(value any) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$gte", Value: value}}}
}

// Lt creates a FieldCondition for less than: { $lt: value }.
// value may be any BSON-marshalable Go value.
func Lt(value any) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$lt", Value: value}}}
}

// Lte creates a FieldCondition for less than or equal to: { $lte: value }.
// value may be any BSON-marshalable Go value.
func Lte(value any) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$lte", Value: value}}}
}

// And creates a Filter for logical AND:
// { $and: [ filter1, filter2, ... ] }
func And(filters ...Filter) Filter {
	clauses := make(bson.A, 0, len(filters))
	for _, f := range filters {
		clauses = append(clauses, f.doc)
	}
	return Filter{doc: bson.D{{Key: "$and", Value: clauses}}}
}

// Or creates a Filter for logical OR:
// { $or: [ filter1, filter2, ... ] }
func Or(filters ...Filter) Filter {
	clauses := make(bson.A, 0, len(filters))
	for _, f := range filters {
		clauses = append(clauses, f.doc)
	}
	return Filter{doc: bson.D{{Key: "$or", Value: clauses}}}
}
