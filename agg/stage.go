package agg

import (
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/mongodb-labs/mongo-go-driver-exp/agg/query"
)

// Stage is a single aggregation pipeline stage, e.g. { $match: ... }.
type Stage bson.D

// Pipeline is an ordered sequence of stages.
type Pipeline []Stage

// MarshalBSON implements bson.Marshaler for Pipeline.
func (p Pipeline) MarshalBSON() ([]byte, error) {
	stages := make([]bson.D, len(p))
	for i, s := range p {
		stages[i] = bson.D(s)
	}
	return bson.Marshal(bson.D{{Key: "pipeline", Value: stages}})
}

// --- $match ---

// MatchStage produces a $match stage from one or more query.Filters.
// Multiple filters are merged into a single document (implicit AND).
// Build filters in the query sub-package, e.g.:
//
//	agg.MatchStage(query.Field("qty", query.Gt(20)), query.Field("name", "Alice"))
func MatchStage(filters ...query.Filter) Stage {
	merged := query.And(filters...)
	return Stage{{Key: "$match", Value: merged}}
}

// --- $set ---

// SetField pairs a field name with the expression to assign to it.
// Construct via Assign.
type SetField interface{ setField() setField }

type setField struct {
	name string
	expr Expr
}

func (sf setField) setField() setField {
	return sf
}

// Assign creates a SetField that sets the named field to expr.
func Assign(field string, expr Expr) SetField {
	return setField{name: field, expr: expr}
}

// SetStage produces a $set stage that adds or overwrites the given fields.
func SetStage(fields ...SetField) Stage {
	doc := make(bson.D, len(fields))
	for i, f := range fields {
		sf := f.setField()
		doc[i] = bson.E{Key: sf.name, Value: sf.expr}
	}
	return Stage{{Key: "$set", Value: doc}}
}

// --- $project ---

// ProjectionField specifies what to do with one field in a $project stage.
// Construct via Include, Exclude, or Compute.
type ProjectionField interface{ projectField() projectField }

type projectField struct {
	name string
	// val is int32(1), int32(0), or an Expr. Constrained by constructors.
	val any
}

func (pf projectField) projectField() projectField {
	return pf
}

// Include retains the named field in the output document (sets it to 1).
func Include(field string) ProjectionField {
	return projectField{name: field, val: int32(1)}
}

// Exclude removes the named field from the output document (sets it to 0).
func Exclude(field string) ProjectionField {
	return projectField{name: field, val: int32(0)}
}

// Compute adds a new (or replaces an existing) field whose value is the
// result of expr.
func Compute(field string, expr Expr) ProjectionField {
	return projectField{name: field, val: expr}
}

// ProjectStage produces a $project stage from the given field specs.
func ProjectStage(specs ...ProjectionField) Stage {
	doc := make(bson.D, len(specs))
	for i, s := range specs {
		pf := s.projectField()
		doc[i] = bson.E{Key: pf.name, Value: pf.val}
	}
	return Stage{{Key: "$project", Value: doc}}
}

// --- $sort ---

// sortOrderKind is the underlying enum type for SortOrder.
type sortOrderKind uint8

const (
	sortKindAsc sortOrderKind = iota
	sortKindDesc
	sortKindTextScore
)

// SortOrder represents a sort direction. Use the package-level vars Asc,
// Desc, and TextScore — do not construct directly.
type SortOrder struct{ kind sortOrderKind }

var (
	// Asc sorts a field in ascending order.
	Asc = SortOrder{sortKindAsc}
	// Desc sorts a field in descending order.
	Desc = SortOrder{sortKindDesc}
	// TextScore sorts by the text search score of the document.
	TextScore = SortOrder{sortKindTextScore}
)

func (s SortOrder) bsonValue() any {
	switch s.kind {
	case sortKindAsc:
		return int32(1)
	case sortKindDesc:
		return int32(-1)
	case sortKindTextScore:
		return bson.D{{Key: "$meta", Value: "textScore"}}
	default:
		panic("agg: invalid SortOrder")
	}
}

// SortField pairs a field name with a sort direction for use in a $sort stage.
// Construct via Sort.
type SortField interface{ sortField() sortField }

type sortField struct {
	name  string
	order SortOrder
}

func (sf sortField) sortField() sortField { return sf }

// Sort creates a SortField that sorts the named field in the given direction.
func Sort(field string, order SortOrder) SortField {
	return sortField{name: field, order: order}
}

// SortStage produces a $sort stage.
// Field order is preserved, which matters for multi-key sorts.
func SortStage(fields ...SortField) Stage {
	doc := make(bson.D, len(fields))
	for i, f := range fields {
		sf := f.sortField()
		doc[i] = bson.E{Key: sf.name, Value: sf.order.bsonValue()}
	}
	return Stage{{Key: "$sort", Value: doc}}
}

// --- $group ---

// GroupField pairs a field name with an accumulator expression for use in a
// $group stage. Construct via Accumulate.
type GroupField interface{ groupField() groupField }

type groupField struct {
	name        string
	accumulator Accumulator
}

func (gf groupField) groupField() groupField {
	return gf
}

// Accumulate creates a GroupField that computes acc for each group and stores
// the result in the named field.
func Accumulate(field string, acc Accumulator) GroupField {
	return groupField{name: field, accumulator: acc}
}

// GroupStage produces a $group stage that groups documents by _id and computes
// the given accumulator fields for each group.
func GroupStage[T AnyExpr | ObjectExpr | string](_id T, fields ...GroupField) Stage {
	doc := make(bson.D, 0, len(fields)+1)
	doc = append(doc, bson.E{Key: "_id", Value: _id})
	for _, f := range fields {
		gf := f.groupField()
		doc = append(doc, bson.E{Key: gf.name, Value: gf.accumulator})
	}
	return Stage{{Key: "$group", Value: doc}}
}
