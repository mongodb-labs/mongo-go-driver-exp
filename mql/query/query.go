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
type Filter bson.D

// Field creates a Filter for the named field from one or more FieldConditions
// (constructed via Eq, Gt, etc.). Multiple conditions are merged into a single
// document to apply several conditions to the same field, e.g.
//
//	query.Field("qty", query.Exists(true), query.Nin(5, 15))
//	// { qty: { $exists: true, $nin: [ 5, 15 ] } }
func Field(name string, conds ...FieldCondition) Filter {
	merged := bson.D{}
	for _, c := range conds {
		merged = append(merged, c.doc...)
	}
	return Filter{{Key: name, Value: merged}}
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

// And creates a Filter for logical AND: { $and: [ filter1, filter2, ... ] }.
func And(filters ...Filter) Filter {
	clauses := make(bson.A, 0, len(filters))
	for _, f := range filters {
		clauses = append(clauses, bson.D(f))
	}
	return Filter{{Key: "$and", Value: clauses}}
}

// Or creates a Filter for logical OR: { $or: [ filter1, filter2, ... ] }.
func Or(filters ...Filter) Filter {
	clauses := make(bson.A, 0, len(filters))
	for _, f := range filters {
		clauses = append(clauses, bson.D(f))
	}
	return Filter{{Key: "$or", Value: clauses}}
}
