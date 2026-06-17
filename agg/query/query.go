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
// e.g. { $gt: value }. Construct via the operator functions (Eq, Gt, etc.).
type FieldCondition struct{ doc bson.D }

// Filter represents a complete MongoDB query document, e.g. { field: { $gt: v } }.
// Construct via Field or the logical combinators And and Or.
type Filter struct{ doc bson.D }

// MarshalBSON implements bson.Marshaler so Filter can be passed directly as
// a bson.D element value.
func (f Filter) MarshalBSON() ([]byte, error) {
	return bson.Marshal(f.doc)
}

// Field creates a Filter for the named field. cond may be a FieldCondition
// (constructed via Eq, Gt, etc.) for operator matches, or any plain value
// for an implicit equality match: { field: value }.
func Field(name string, cond any) Filter {
	if fc, ok := cond.(FieldCondition); ok {
		return Filter{doc: bson.D{{Key: name, Value: fc.doc}}}
	}
	return Filter{doc: bson.D{{Key: name, Value: cond}}}
}

// Eq creates a FieldCondition for equality: { $eq: value }.
func Eq(value any) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$eq", Value: value}}}
}

// Gt creates a FieldCondition for greater than: { $gt: value }.
func Gt(value any) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$gt", Value: value}}}
}

// Gte creates a FieldCondition for greater than or equal to: { $gte: value }.
func Gte(value any) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$gte", Value: value}}}
}

// Lt creates a FieldCondition for less than: { $lt: value }.
func Lt(value any) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$lt", Value: value}}}
}

// Lte creates a FieldCondition for less than or equal to: { $lte: value }.
func Lte(value any) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$lte", Value: value}}}
}

// And creates a Filter that merges all given filters into a single document,
// producing an implicit AND: { field1: cond1, field2: cond2, ... }.
func And(filters ...Filter) Filter {
	var merged bson.D
	for _, f := range filters {
		merged = append(merged, f.doc...)
	}
	return Filter{doc: merged}
}

// Or creates a Filter for logical OR: { $or: [ filter1, filter2, ... ] }.
func Or(filters ...Filter) Filter {
	clauses := make(bson.A, 0, len(filters))
	for _, f := range filters {
		clauses = append(clauses, f.doc)
	}
	return Filter{doc: bson.D{{Key: "$or", Value: clauses}}}
}
