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

// MatchStage produces a $match stage from a query.Filter.
// Build the filter in the query sub-package, e.g.:
//
//	agg.MatchStage(query.Field("qty", query.Eq(20)))
func MatchStage(f query.Filter) Stage {
	return Stage{{Key: "$match", Value: f}}
}

// --- $set ---

// SetField pairs a field name with the expression to assign to it.
// Construct via Assign.
type SetField struct {
	name string
	expr Expr
}

// Assign creates a SetField that sets the named field to expr.
func Assign(field string, expr Expr) SetField {
	return SetField{name: field, expr: expr}
}

// SetStage produces a $set stage that adds or overwrites the given fields.
func SetStage(fields ...SetField) Stage {
	doc := make(bson.D, len(fields))
	for i, f := range fields {
		doc[i] = bson.E{Key: f.name, Value: f.expr}
	}
	return Stage{{Key: "$set", Value: doc}}
}

// --- $project ---

// ProjectionField specifies what to do with one field in a $project stage.
// Construct via Include, Exclude, or Compute.
type ProjectionField interface{ projectionField() projectionField }

type projectionField struct {
	name string
	// val is int32(1), int32(0), or an Expr. Constrained by constructors.
	val any
}

func (pf projectionField) projectionField() projectionField {
	return pf
}

// Include retains the named field in the output document (sets it to 1).
func Include(field string) ProjectionField {
	return projectionField{name: field, val: int32(1)}
}

// Exclude removes the named field from the output document (sets it to 0).
func Exclude(field string) ProjectionField {
	return projectionField{name: field, val: int32(0)}
}

// Compute adds a new (or replaces an existing) field whose value is the
// result of expr.
func Compute(field string, expr Expr) ProjectionField {
	return projectionField{name: field, val: expr}
}

// ProjectStage produces a $project stage from the given field specs.
func ProjectStage(specs ...ProjectionField) Stage {
	doc := make(bson.D, len(specs))
	for i, s := range specs {
		pf := s.projectionField()
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

// SortSpec pairs a field name with a sort direction.
type SortSpec struct {
	Field string
	Order SortOrder
}

// SortStage produces a $sort stage.
// Field order is preserved, which matters for multi-key sorts.
func SortStage(specs ...SortSpec) Stage {
	doc := make(bson.D, len(specs))
	for i, s := range specs {
		doc[i] = bson.E{Key: s.Field, Value: s.Order.bsonValue()}
	}
	return Stage{{Key: "$sort", Value: doc}}
}
