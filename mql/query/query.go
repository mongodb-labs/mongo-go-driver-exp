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
type FieldCondition bson.D

// Filter represents a complete MongoDB query document, e.g. { field: { $gt: v } }.
// Construct via Field or the logical combinators (And, Or, etc.).
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
		merged = append(merged, c...)
	}
	return Filter{{Key: name, Value: merged}}
}

// Option is a functional option that configures the optional parameters of an
// operator with type T.
type Option[T any] func(*T)

// All creates a FieldCondition matching arrays that contain all of the given
// values: { $all: [ ... ] }. Values are usually plain scalars, but may also be
// ElemMatch conditions to match arrays of embedded documents.
func All(values ...any) FieldCondition {
	return FieldCondition{{Key: "$all", Value: bson.A(values)}}
}

// And creates a Filter for logical AND: { $and: [ filter1, filter2, ... ] }.
func And(filters ...Filter) Filter {
	clauses := make(bson.A, 0, len(filters))
	for _, filter := range filters {
		clauses = append(clauses, bson.D(filter))
	}
	return Filter{{Key: "$and", Value: clauses}}
}

// Bitmask is the set of types accepted by the $bits* operators: an integer
// bitmask, a bson.Binary bitmask, or an array of bit positions.
type Bitmask interface {
	~int | ~int32 | ~int64 | bson.Binary | ~[]int | ~[]int32 | ~[]int64
}

// BitsAllClear creates a FieldCondition matching numeric or binary values in
// which all of the given bit positions are 0: { $bitsAllClear: bitmask }. The
// bitmask may be an int, a bson.Binary, or an array of bit positions.
func BitsAllClear[T Bitmask](bitmask T) FieldCondition {
	return FieldCondition{{Key: "$bitsAllClear", Value: bitmask}}
}

// BitsAllSet creates a FieldCondition matching numeric or binary values in which
// all of the given bit positions are 1: { $bitsAllSet: bitmask }. The bitmask
// may be an int, a bson.Binary, or an array of bit positions.
func BitsAllSet[T Bitmask](bitmask T) FieldCondition {
	return FieldCondition{{Key: "$bitsAllSet", Value: bitmask}}
}

// BitsAnyClear creates a FieldCondition matching numeric or binary values in
// which any of the given bit positions are 0: { $bitsAnyClear: bitmask }. The
// bitmask may be an int, a bson.Binary, or an array of bit positions.
func BitsAnyClear[T Bitmask](bitmask T) FieldCondition {
	return FieldCondition{{Key: "$bitsAnyClear", Value: bitmask}}
}

// BitsAnySet creates a FieldCondition matching numeric or binary values in which
// any of the given bit positions are 1: { $bitsAnySet: bitmask }. The bitmask
// may be an int, a bson.Binary, or an array of bit positions.
func BitsAnySet[T Bitmask](bitmask T) FieldCondition {
	return FieldCondition{{Key: "$bitsAnySet", Value: bitmask}}
}

// Box creates a legacy rectangular box geometry ($box) from the bottom-left and
// top-right coordinate pairs, for use with GeoWithin.
func Box(bottomLeft, topRight []float64) Geometry {
	return Geometry{doc: bson.D{{Key: "$box", Value: bson.A{bottomLeft, topRight}}}}
}

// Center creates a legacy circle geometry ($center) from a center coordinate
// pair and a radius (in coordinate units), for use with GeoWithin using planar
// geometry.
func Center(center []float64, radius float64) Geometry {
	return Geometry{doc: bson.D{{Key: "$center", Value: bson.A{center, radius}}}}
}

// CenterSphere creates a legacy spherical circle geometry ($centerSphere) from
// a center coordinate pair and a radius (in radians), for use with GeoWithin
// using spherical geometry.
func CenterSphere(center []float64, radius float64) Geometry {
	return Geometry{doc: bson.D{{Key: "$centerSphere", Value: bson.A{center, radius}}}}
}

// Comment creates a Filter that adds a comment to a query predicate:
// { $comment: comment }.
func Comment(comment string) Filter {
	return Filter{{Key: "$comment", Value: comment}}
}

// ElemMatch creates a FieldCondition matching arrays with at least one element
// that satisfies all of the given queries: { $elemMatch: { ... } }. Pass Filters
// to match arrays of embedded documents (e.g. query.Field("product", query.Eq("xyz")))
// or FieldConditions to match scalar elements (e.g. query.Gte(80), query.Lt(85)).
func ElemMatch[T Filter | FieldCondition](queries ...T) FieldCondition {
	merged := bson.D{}
	for _, q := range queries {
		merged = append(merged, q...)
	}
	return FieldCondition{{Key: "$elemMatch", Value: merged}}
}

// Eq creates a FieldCondition for equality: { $eq: value }.
func Eq(value any) FieldCondition {
	return FieldCondition{{Key: "$eq", Value: value}}
}

// Exists creates a FieldCondition matching documents that have (or lack) the
// field: { $exists: exists }.
func Exists(exists bool) FieldCondition {
	return FieldCondition{{Key: "$exists", Value: exists}}
}

// Expr creates a Filter that allows an aggregation expression within the query
// language: { $expr: expression }.
func Expr(expr any) Filter {
	return Filter{{Key: "$expr", Value: expr}}
}

// Geometry represents a geometry value supplied to the geospatial query
// operators (GeoWithin, GeoIntersects, Near, NearSphere). Construct via GeoJSON
// or the legacy shape helpers (Box, Center, CenterSphere, Polygon).
type Geometry struct {
	doc bson.D
}

// Coordinates constrains the coordinate values of a GeoJSON geometry to array
// shapes. The nesting depth depends on the geometry type: a Point is a single
// position ([]float64), a Polygon is [][][]float64, and so on. Use bson.A for
// dynamic or mixed shapes.
type Coordinates interface {
	[]float64 | [][]float64 | [][][]float64 | [][][][]float64 | bson.A
}

type geoJSONOptions struct {
	crs bson.D
}

// WithGeoJSONCRS sets the coordinate reference system for a GeoJSON geometry,
// e.g. to request a big (strict CRS84) polygon.
func WithGeoJSONCRS(crs bson.D) Option[geoJSONOptions] {
	return func(opts *geoJSONOptions) {
		opts.crs = crs
	}
}

// GeoJSON creates a GeoJSON geometry ($geometry) of the given type (e.g.
// "Point", "Polygon") and coordinates. The coordinate shape depends on the
// geometry type. Optionally specify a coordinate reference system via
// WithGeoJSONCRS.
func GeoJSON[C Coordinates](geoType string, coordinates C, opts ...Option[geoJSONOptions]) Geometry {
	var o geoJSONOptions
	for _, opt := range opts {
		opt(&o)
	}
	geo := bson.D{
		{Key: "type", Value: geoType},
		{Key: "coordinates", Value: coordinates},
	}
	if o.crs != nil {
		geo = append(geo, bson.E{Key: "crs", Value: o.crs})
	}
	return Geometry{doc: bson.D{{Key: "$geometry", Value: geo}}}
}

// GeoIntersects creates a FieldCondition matching geometries that intersect the
// given GeoJSON geometry: { $geoIntersects: { $geometry: ... } }. Use GeoJSON to
// construct the geometry; $geoIntersects does not support the legacy shapes.
func GeoIntersects(g Geometry) FieldCondition {
	return FieldCondition{{Key: "$geoIntersects", Value: g.doc}}
}

// GeoWithin creates a FieldCondition matching geometries within the given
// bounding geometry: { $geoWithin: { ... } }. Accepts a GeoJSON geometry or any
// of the legacy shapes (Box, Center, CenterSphere, Polygon).
func GeoWithin(g Geometry) FieldCondition {
	return FieldCondition{{Key: "$geoWithin", Value: g.doc}}
}

// Gt creates a FieldCondition for greater than: { $gt: value }.
func Gt(value any) FieldCondition {
	return FieldCondition{{Key: "$gt", Value: value}}
}

// Gte creates a FieldCondition for greater than or equal to: { $gte: value }.
func Gte(value any) FieldCondition {
	return FieldCondition{{Key: "$gte", Value: value}}
}

// In creates a FieldCondition matching any of the given values: { $in: [ ... ] }.
func In(values ...any) FieldCondition {
	return FieldCondition{{Key: "$in", Value: bson.A(values)}}
}

// JSONSchema creates a Filter that validates documents against the given JSON
// Schema: { $jsonSchema: schema }.
func JSONSchema(schema any) Filter {
	return Filter{{Key: "$jsonSchema", Value: schema}}
}

// Lt creates a FieldCondition for less than: { $lt: value }.
func Lt(value any) FieldCondition {
	return FieldCondition{{Key: "$lt", Value: value}}
}

// Lte creates a FieldCondition for less than or equal to: { $lte: value }.
func Lte(value any) FieldCondition {
	return FieldCondition{{Key: "$lte", Value: value}}
}

// MaxDistance limits results to within d of the query point — radians for a
// legacy 2d index, meters for GeoJSON on a 2dsphere index. d must be
// non-negative.
//
// The server treats the distance as a double and clamps it to the maximum
// distance on the sphere (~2.0037e7 m, i.e. an antipodal great-circle
// distance; a few radians for legacy indexes). Any larger value matches the
// whole sphere, so the precision limits of float64 above 2^53 and the range
// of int64 are never reachable here — even int32 far exceeds any meaningful
// distance. The Number constraint (int/float kinds, no unsigned) is chosen
// for call-site ergonomics, not range.

// MaxDistance creates a FieldCondition limiting Near and NearSphere results to
// at most the given distance from the center point: { $maxDistance: value }.
func MaxDistance[T int | float64](value T) FieldCondition {
	return FieldCondition{{Key: "$maxDistance", Value: value}}
}

// MinDistance creates a FieldCondition limiting Near and NearSphere results to
// at least the given distance from the center point: { $minDistance: value }.
func MinDistance[T int | float64](value T) FieldCondition {
	return FieldCondition{{Key: "$minDistance", Value: value}}
}

// Mod creates a FieldCondition that performs a modulo operation on the field
// and matches the specified result: { $mod: [ divisor, remainder ] }.
func Mod(divisor int, remainder int) FieldCondition {
	return FieldCondition{{Key: "$mod", Value: bson.A{divisor, remainder}}}
}

type nearOptions struct {
	minDistance any
	maxDistance any
}

// WithNearMinDistance limits Near/NearSphere results to at least the given
// distance (in meters) from the center point.
func WithNearMinDistance[T int | float64](d T) Option[nearOptions] {
	return func(opts *nearOptions) { opts.minDistance = d }
}

// WithNearMaxDistance limits Near/NearSphere results to at most the given
// distance (in meters) from the center point.
func WithNearMaxDistance[T int | float64](d T) Option[nearOptions] {
	return func(opts *nearOptions) { opts.maxDistance = d }
}

// nearDoc merges the geometry with the optional distance bounds into the value
// document shared by $near and $nearSphere.
func nearDoc(g Geometry, opts []Option[nearOptions]) bson.D {
	var o nearOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := append(bson.D(nil), g.doc...)
	if o.minDistance != nil {
		doc = append(doc, bson.E{Key: "$minDistance", Value: o.minDistance})
	}
	if o.maxDistance != nil {
		doc = append(doc, bson.E{Key: "$maxDistance", Value: o.maxDistance})
	}
	return doc
}

// Ne creates a FieldCondition matching values not equal to value: { $ne: value }.
func Ne(value any) FieldCondition {
	return FieldCondition{{Key: "$ne", Value: value}}
}

// Near creates a FieldCondition matching geospatial objects in proximity to the
// given point, sorted by distance: { $near: { $geometry: ..., $minDistance,
// $maxDistance } }. Bounds are set via WithNearMinDistance and
// WithNearMaxDistance.
func Near(g Geometry, opts ...Option[nearOptions]) FieldCondition {
	return FieldCondition{{Key: "$near", Value: nearDoc(g, opts)}}
}

// NearSphere creates a FieldCondition matching geospatial objects in proximity
// to the given point on a sphere, sorted by distance: { $nearSphere: {
// $geometry: ..., $minDistance, $maxDistance } }. Bounds are set via
// WithNearMinDistance and WithNearMaxDistance.
func NearSphere(g Geometry, opts ...Option[nearOptions]) FieldCondition {
	return FieldCondition{{Key: "$nearSphere", Value: nearDoc(g, opts)}}
}

// Nin creates a FieldCondition matching none of the given values: { $nin: [ ... ] }.
func Nin(values ...any) FieldCondition {
	return FieldCondition{{Key: "$nin", Value: bson.A(values)}}
}

// Nor creates a Filter for logical NOR, matching documents that fail every
// clause: { $nor: [ filter1, filter2, ... ] }.
func Nor(filters ...Filter) Filter {
	clauses := make(bson.A, 0, len(filters))
	for _, filter := range filters {
		clauses = append(clauses, bson.D(filter))
	}
	return Filter{{Key: "$nor", Value: clauses}}
}

// Not creates a FieldCondition that inverts another field-level condition:
// { $not: <arg> }. The argument may be a FieldCondition, e.g.
// query.Not(query.Gt(1.99)) yields { $not: { $gt: 1.99 } }, or a bson.Regex,
// e.g. query.Not(bson.Regex{Pattern: "^p.*"}) yields { $not: /^p.*/ }. MongoDB's
// $not requires an operator expression or regex; it does not accept a plain
// scalar value.
func Not[T FieldCondition | bson.Regex](value T) FieldCondition {
	return FieldCondition{{Key: "$not", Value: value}}
}

// Or creates a Filter for logical OR: { $or: [ filter1, filter2, ... ] }.
func Or(filters ...Filter) Filter {
	clauses := make(bson.A, 0, len(filters))
	for _, filter := range filters {
		clauses = append(clauses, bson.D(filter))
	}
	return Filter{{Key: "$or", Value: clauses}}
}

// Polygon creates a legacy polygon geometry ($polygon) from a series of
// coordinate pairs defining the polygon's vertices, for use with GeoWithin.
func Polygon(points ...[]float64) Geometry {
	return Geometry{doc: bson.D{{Key: "$polygon", Value: [][]float64(points)}}}
}

// Regex creates a FieldCondition matching values against a regular expression:
// { $regex: regex }. Use the Options field of bson.Regex for flags like "i".
func Regex(regex bson.Regex) FieldCondition {
	return FieldCondition{{Key: "$regex", Value: regex}}
}

// SampleRate creates a Filter that randomly selects documents at the given rate
// (a probability between 0 and 1): { $sampleRate: rate }.
func SampleRate(rate float64) Filter {
	return Filter{{Key: "$sampleRate", Value: rate}}
}

// Size creates a FieldCondition matching arrays with the given number of
// elements: { $size: value }.
func Size(value int) FieldCondition {
	return FieldCondition{{Key: "$size", Value: value}}
}

type textOptions struct {
	language           any
	caseSensitive      any
	diacriticSensitive any
}

// WithTextLanguage sets the $language option for a Text search.
func WithTextLanguage(language string) Option[textOptions] {
	return func(opts *textOptions) {
		opts.language = language
	}
}

// WithTextCaseSensitive sets the $caseSensitive option for a Text search.
func WithTextCaseSensitive(caseSensitive bool) Option[textOptions] {
	return func(opts *textOptions) {
		opts.caseSensitive = caseSensitive
	}
}

// WithTextDiacriticSensitive sets the $diacriticSensitive option for a Text search.
func WithTextDiacriticSensitive(diacriticSensitive bool) Option[textOptions] {
	return func(opts *textOptions) {
		opts.diacriticSensitive = diacriticSensitive
	}
}

// Text creates a Filter that performs a text search: { $text: { $search: ... } }.
// Optional behavior is set via the WithText* options.
func Text(search string, opts ...Option[textOptions]) Filter {
	var o textOptions
	for _, opt := range opts {
		opt(&o)
	}
	args := bson.D{bson.E{Key: "$search", Value: search}}
	if o.language != nil {
		args = append(args, bson.E{Key: "$language", Value: o.language})
	}
	if o.caseSensitive != nil {
		args = append(args, bson.E{Key: "$caseSensitive", Value: o.caseSensitive})
	}
	if o.diacriticSensitive != nil {
		args = append(args, bson.E{Key: "$diacriticSensitive", Value: o.diacriticSensitive})
	}
	return Filter{{Key: "$text", Value: args}}
}

// Type creates a FieldCondition matching documents where the field is one of the
// specified BSON types: { $type: [ ... ] }. Each type may be an alias string or
// numeric code. The verbose array form is always emitted.
func Type(types ...any) FieldCondition {
	return FieldCondition{{Key: "$type", Value: bson.A(types)}}
}

// Where creates a Filter matching documents that satisfy a JavaScript
// expression: { $where: function }.
func Where(function bson.JavaScript) Filter {
	return Filter{{Key: "$where", Value: function}}
}
