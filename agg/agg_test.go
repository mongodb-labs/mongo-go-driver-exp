package agg_test

import (
	"bytes"
	"testing"

	"github.com/mongodb-labs/mongo-go-driver-exp/agg"
	"github.com/mongodb-labs/mongo-go-driver-exp/agg/query"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// --- $sort ---

func TestSortStage_AscendingDescendingSort(t *testing.T) {
	pipeline := agg.Pipeline{
		agg.SortStage(
			agg.SortSpec{
				Field: "age",
				Order: agg.Desc},
			agg.SortSpec{
				Field: "posts",
				Order: agg.Asc,
			},
		),
	}

	got, err := bson.Marshal(pipeline)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	want, err := bson.Marshal(bson.D{{Key: "pipeline", Value: bson.A{
		bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "age", Value: -1},
				{Key: "posts", Value: 1},
			}},
		},
	}}})
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if !bytes.Equal(want, got) {
		t.Errorf(
			"Pipelines don't match.\nWant: %s\nGot:  %s",
			bson.Raw(want).String(),
			bson.Raw(got).String())
	}
}

// --- $set ---

func TestSetStage_AddingFieldsToEmbeddedDoc(t *testing.T) {
	pipeline := agg.Pipeline{
		agg.SetStage(
			agg.Assign("specs.fuel_type", "unleaded"),
		),
	}

	got, err := bson.Marshal(pipeline)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	want, err := bson.Marshal(bson.D{{Key: "pipeline", Value: bson.A{
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "specs.fuel_type", Value: "unleaded"}},
			}},
	}}})
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if !bytes.Equal(want, got) {
		t.Errorf(
			"Pipelines don't match.\nWant: %s\nGot:  %s",
			bson.Raw(want).String(),
			bson.Raw(got).String())
	}
}

func TestSetStage_OverwriteExistingField(t *testing.T) {
	pipeline := agg.Pipeline{
		agg.SetStage(
			agg.Assign("cats", 20),
		),
	}

	got, err := bson.Marshal(pipeline)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	want, err := bson.Marshal(bson.D{{Key: "pipeline", Value: bson.A{
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "cats", Value: 20}},
			}},
	}}})
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if !bytes.Equal(want, got) {
		t.Errorf(
			"Pipelines don't match.\nWant: %s\nGot:  %s",
			bson.Raw(want).String(),
			bson.Raw(got).String())
	}
}

// --- $project ---

// TODO: Add limit of 1 when limit stage is implemented
func TestProjectStage_IncludeSpecificFieldsInOutputDocs(t *testing.T) {
	pipeline := agg.Pipeline{
		agg.MatchStage(query.Field("title", query.Eq("The Great Train Robbery"))),
		agg.ProjectStage(
			agg.Include("title"),
			agg.Include("rated"),
		),
	}

	got, err := bson.Marshal(pipeline)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	want, err := bson.Marshal(bson.D{{
		Key: "pipeline",
		Value: bson.A{
			bson.D{{
				Key: "$match",
				Value: bson.D{{
					Key: "title",
					Value: bson.D{{
						Key:   "$eq",
						Value: "The Great Train Robbery",
					}},
				}},
			}},
			bson.D{{
				Key: "$project",
				Value: bson.D{
					{Key: "title", Value: 1},
					{Key: "rated", Value: 1},
				},
			}},
		},
	}})
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if !bytes.Equal(want, got) {
		t.Errorf(
			"Pipelines don't match.\nWant: %s\nGot:  %s",
			bson.Raw(want).String(),
			bson.Raw(got).String())
	}
}

// TODO: Add limit of 1 when limit stage is implemented
func TestProjectStage_SuppressIdFieldInOutputDocs(t *testing.T) {
	pipeline := agg.Pipeline{
		agg.MatchStage(query.Field("title", query.Eq("The Great Train Robbery"))),
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Include("title"),
			agg.Include("rated"),
		),
	}

	got, err := bson.Marshal(pipeline)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	want, err := bson.Marshal(bson.D{{
		Key: "pipeline",
		Value: bson.A{
			bson.D{{
				Key: "$match",
				Value: bson.D{{
					Key: "title",
					Value: bson.D{{
						Key:   "$eq",
						Value: "The Great Train Robbery",
					}},
				}},
			}},
			bson.D{{
				Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "title", Value: 1},
					{Key: "rated", Value: 1},
				},
			}},
		},
	}})
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if !bytes.Equal(want, got) {
		t.Errorf(
			"Pipelines don't match.\nWant: %s\nGot:  %s",
			bson.Raw(want).String(),
			bson.Raw(got).String())
	}
}

// TODO: Add limit of 1 when limit stage is implemented
func TestProjectStage_ExcludeFieldsFromOutputDocs(t *testing.T) {
	pipeline := agg.Pipeline{
		agg.MatchStage(query.Field("title", query.Eq("The Great Train Robbery"))),
		agg.ProjectStage(
			agg.Exclude("rated"),
		),
	}

	got, err := bson.Marshal(pipeline)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	want, err := bson.Marshal(bson.D{{
		Key: "pipeline",
		Value: bson.A{
			bson.D{{
				Key: "$match",
				Value: bson.D{{
					Key: "title",
					Value: bson.D{{
						Key:   "$eq",
						Value: "The Great Train Robbery",
					}},
				}},
			}},
			bson.D{{
				Key: "$project",
				Value: bson.D{
					{Key: "rated", Value: 0},
				},
			}},
		},
	}})
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if !bytes.Equal(want, got) {
		t.Errorf(
			"Pipelines don't match.\nWant: %s\nGot:  %s",
			bson.Raw(want).String(),
			bson.Raw(got).String())
	}
}

// TODO: Add limit of 1 when limit stage is implemented
func TestProjectStage_ExcludeFieldsFromEmbeddedDocs(t *testing.T) {
	pipeline := agg.Pipeline{
		agg.MatchStage(query.Field("title", query.Eq("The Great Train Robbery"))),
		agg.ProjectStage(
			agg.Exclude("imdb.id"),
			agg.Exclude("type"),
		),
	}

	got, err := bson.Marshal(pipeline)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	want, err := bson.Marshal(bson.D{{
		Key: "pipeline",
		Value: bson.A{
			bson.D{{
				Key: "$match",
				Value: bson.D{{
					Key: "title",
					Value: bson.D{{
						Key:   "$eq",
						Value: "The Great Train Robbery",
					}},
				}},
			}},
			bson.D{{
				Key: "$project",
				Value: bson.D{
					{Key: "imdb.id", Value: 0},
					{Key: "type", Value: 0},
				},
			}},
		},
	}})
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if !bytes.Equal(want, got) {
		t.Errorf(
			"Pipelines don't match.\nWant: %s\nGot:  %s",
			bson.Raw(want).String(),
			bson.Raw(got).String())
	}
}

// TODO: Add limit of 1 when limit stage is implemented
func TestProjectStage_IncludeSpecificFieldsFromEmbeddedDocs(t *testing.T) {
	pipeline := agg.Pipeline{
		agg.MatchStage(query.Field("title", query.Eq("The Great Train Robbery"))),
		agg.ProjectStage(
			agg.Include("imdb.rating"),
		),
	}

	got, err := bson.Marshal(pipeline)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	want, err := bson.Marshal(bson.D{{
		Key: "pipeline",
		Value: bson.A{
			bson.D{{
				Key: "$match",
				Value: bson.D{{
					Key: "title",
					Value: bson.D{{
						Key:   "$eq",
						Value: "The Great Train Robbery",
					}},
				}},
			}},
			bson.D{{
				Key: "$project",
				Value: bson.D{
					{Key: "imdb.rating", Value: 1},
				},
			}},
		},
	}})
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if !bytes.Equal(want, got) {
		t.Errorf(
			"Pipelines don't match.\nWant: %s\nGot:  %s",
			bson.Raw(want).String(),
			bson.Raw(got).String())
	}
}

// TODO: Add limit of 1 when limit stage is implemented
// TODO: Add computed field leadActor: { $arrayElemAt: [ "$cast", 0 ] }
// when arrayElemAt operator is implemented
func TestProjectStage_IncludeComputedFields(t *testing.T) {
	pipeline := agg.Pipeline{
		agg.MatchStage(query.Field("title", query.Eq("The Great Train Robbery"))),
		agg.ProjectStage(
			agg.Include("title"),
			agg.Compute("releaseYear", "$year"),
		),
	}

	got, err := bson.Marshal(pipeline)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	want, err := bson.Marshal(bson.D{{
		Key: "pipeline",
		Value: bson.A{
			bson.D{{
				Key: "$match",
				Value: bson.D{{
					Key:   "title",
					Value: "The Great Train Robbery",
				}},
			}},
			bson.D{{
				Key: "$project",
				Value: bson.D{
					{Key: "title", Value: 1},
					{Key: "releaseYear", Value: "$year"},
				},
			}},
		},
	}})
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if !bytes.Equal(want, got) {
		t.Errorf(
			"Pipelines don't match.\nWant: %s\nGot:  %s",
			bson.Raw(want).String(),
			bson.Raw(got).String())
	}
}

// --- $match ---

func TestMatchStage_EqualityMatch(t *testing.T) {
	pipeline := agg.Pipeline{
		agg.MatchStage(query.And(query.Field("rated", query.Eq("TV-PG")), query.Field("runtime", query.Gt(1000)))),
	}

	got, err := bson.Marshal(pipeline)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	want, err := bson.Marshal(bson.D{{
		Key: "pipeline",
		Value: bson.A{
			bson.D{{
				Key: "$match",
				Value: bson.D{
					{Key: "rated", Value: "TV-PG"},
					{Key: "runtime", Value: bson.D{{Key: "$gt", Value: 1000}}},
				},
			}},
		},
	}})

	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if !bytes.Equal(want, got) {
		t.Errorf(
			"Pipelines don't match.\nWant: %s\nGot:  %s",
			bson.Raw(want).String(),
			bson.Raw(got).String())
	}
}
