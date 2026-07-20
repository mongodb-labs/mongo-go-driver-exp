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
//	query.Field("qty", query.Gt(5), query.Lt(15))
//	// { qty: { $gt: 5, $lt: 15 } }
func Field(name string, conds ...FieldCondition) Filter {
	merged := bson.D{}
	for _, c := range conds {
		merged = append(merged, c.doc...)
	}
	return Filter{{Key: name, Value: merged}}
}

// All creates a FieldCondition matching arrays that contain all of the given
// values: { $all: [ ... ] }. Values are usually plain scalars, but may also be
// ElemMatch conditions to match arrays of embedded documents.
func All(values ...any) FieldCondition {
	arr := make(bson.A, len(values))
	for i, v := range values {
		if fc, ok := v.(FieldCondition); ok {
			arr[i] = fc.doc
		} else {
			arr[i] = v
		}
	}
	return FieldCondition{doc: bson.D{{Key: "$all", Value: arr}}}
}

// And creates a Filter for logical AND: { $and: [ filter1, filter2, ... ] }.
func And(filters ...Filter) Filter {
	clauses := make(bson.A, 0, len(filters))
	for _, f := range filters {
		clauses = append(clauses, bson.D(f))
	}
	return Filter{{Key: "$and", Value: clauses}}
}

// ElemMatch creates a FieldCondition matching arrays with at least one element
// that satisfies all of the given queries: { $elemMatch: { ... } }. Pass Filters
// to match arrays of embedded documents (e.g. query.Field("product", query.Eq("xyz")))
// or FieldConditions to match scalar elements (e.g. query.Gte(80), query.Lt(85)).
func ElemMatch[T Filter | FieldCondition](queries ...T) FieldCondition {
	inner := bson.D{}
	for _, q := range queries {
		switch v := any(q).(type) {
		case Filter:
			inner = append(inner, v...)
		case FieldCondition:
			inner = append(inner, v.doc...)
		}
	}
	return FieldCondition{doc: bson.D{{Key: "$elemMatch", Value: inner}}}
}

// Eq creates a FieldCondition for equality: { $eq: value }.
func Eq(value any) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$eq", Value: value}}}
}

// Exists creates a FieldCondition matching documents that have (or lack) the
// field: { $exists: exists }.
func Exists(exists bool) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$exists", Value: exists}}}
}

// Gt creates a FieldCondition for greater than: { $gt: value }.
func Gt(value any) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$gt", Value: value}}}
}

// Gte creates a FieldCondition for greater than or equal to: { $gte: value }.
func Gte(value any) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$gte", Value: value}}}
}

// In creates a FieldCondition matching any of the given values: { $in: [ ... ] }.
func In(values ...any) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$in", Value: bson.A(values)}}}
}

// Lt creates a FieldCondition for less than: { $lt: value }.
func Lt(value any) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$lt", Value: value}}}
}

// Lte creates a FieldCondition for less than or equal to: { $lte: value }.
func Lte(value any) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$lte", Value: value}}}
}

// Ne creates a FieldCondition matching values not equal to value: { $ne: value }.
func Ne(value any) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$ne", Value: value}}}
}

// Nin creates a FieldCondition matching none of the given values: { $nin: [ ... ] }.
func Nin(values ...any) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$nin", Value: bson.A(values)}}}
}

// Nor creates a Filter for logical NOR, matching documents that fail every
// clause: { $nor: [ filter1, filter2, ... ] }.
func Nor(filters ...Filter) Filter {
	clauses := make(bson.A, 0, len(filters))
	for _, f := range filters {
		clauses = append(clauses, bson.D(f))
	}
	return Filter{{Key: "$nor", Value: clauses}}
}

// Not creates a FieldCondition that inverts another field-level condition:
// { $not: <arg> }. The argument may be a FieldCondition, e.g.
// query.Not(query.Gt(1.99)) yields { $not: { $gt: 1.99 } }, or a bson.Regex,
// e.g. query.Not(bson.Regex{Pattern: "^p.*"}) yields { $not: /^p.*/ }. MongoDB's
// $not requires an operator expression or regex; it does not accept a plain
// scalar value.
func Not[T FieldCondition | bson.Regex](arg T) FieldCondition {
	var value any
	switch a := any(arg).(type) {
	case FieldCondition:
		value = a.doc
	case bson.Regex:
		value = a
	}
	return FieldCondition{doc: bson.D{{Key: "$not", Value: value}}}
}

// Size creates a FieldCondition matching arrays with the given number of
// elements: { $size: value }.
func Size(value int) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$size", Value: value}}}
}

// Type creates a FieldCondition matching documents where the field is one of the
// specified BSON types: { $type: [ ... ] }. Each type may be an alias string or
// numeric code. The verbose array form is always emitted.
func Type(types ...any) FieldCondition {
	return FieldCondition{doc: bson.D{{Key: "$type", Value: bson.A(types)}}}
}

// Or creates a Filter for logical OR: { $or: [ filter1, filter2, ... ] }.
func Or(filters ...Filter) Filter {
	clauses := make(bson.A, 0, len(filters))
	for _, f := range filters {
		clauses = append(clauses, bson.D(f))
	}
	return Filter{{Key: "$or", Value: clauses}}
}
