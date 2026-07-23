package agg

import (
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/mongodb-labs/mongo-go-driver-exp/mql/query"
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

// --- shared field types and helpers ---
//
// These are used by more than one stage, so they live here rather than above a
// single stage. Everything else (types, With* options) sits directly above the
// stage it belongs to.

// setFieldsDoc renders a slice of SetFields as a { name: expr } document,
// shared by the let option of $lookup and $merge. Bindings are built with
// Assign, the same way the $let expression operator takes its variables.
func setFieldsDoc(vars []SetField) bson.D {
	doc := make(bson.D, len(vars))
	for i, v := range vars {
		sf := v.setField()
		doc[i] = bson.E{Key: sf.name, Value: sf.expr}
	}
	return doc
}

// SessionUser identifies a user for the users option of $listSessions and
// $listLocalSessions.
type SessionUser struct {
	User string
	DB   string
}

func sessionUsersArray(users []SessionUser) bson.A {
	arr := make(bson.A, len(users))
	for i, u := range users {
		arr[i] = bson.D{{Key: "user", Value: u.User}, {Key: "db", Value: u.DB}}
	}
	return arr
}

// listSessionsOptions is shared by $listSessions and $listLocalSessions, which
// take identical arguments.
type listSessionsOptions struct {
	users    []SessionUser
	allUsers any
}

func pipelineDocs(stages []Stage) []bson.D {
	docs := make([]bson.D, len(stages))
	for i, s := range stages {
		docs[i] = bson.D(s)
	}
	return docs
}

// groupFieldsDoc renders a slice of GroupFields as an accumulator document,
// shared by the $bucket, $bucketAuto, and $group output specifications.
func groupFieldsDoc(fields []GroupField) bson.D {
	doc := make(bson.D, len(fields))
	for i, f := range fields {
		gf := f.groupField()
		doc[i] = bson.E{Key: gf.name, Value: gf.accumulator}
	}
	return doc
}

// sortFieldsDoc renders a slice of SortFields as a sort document, shared by the
// sortBy option of $fill and $setWindowFields.
func sortFieldsDoc(fields []SortField) bson.D {
	doc := make(bson.D, len(fields))
	for i, f := range fields {
		sf := f.sortField()
		doc[i] = bson.E{Key: sf.name, Value: sf.order.bsonValue()}
	}
	return doc
}

// --- $addFields ---

// AddFieldsStage produces an $addFields stage that adds or overwrites the given fields.
// It is an alias for $set. Construct via Assign.
func AddFieldsStage(fields ...SetField) Stage {
	doc := make(bson.D, len(fields))
	for i, f := range fields {
		sf := f.setField()
		doc[i] = bson.E{Key: sf.name, Value: sf.expr}
	}
	return Stage{{Key: "$addFields", Value: doc}}
}

// --- $bucket ---

type bucketOptions struct {
	def    any
	output []GroupField
}

// WithBucketDefault sets the _id of an additional bucket collecting all
// documents whose groupBy value falls outside the given boundaries. Without it,
// a value outside every bucket causes a runtime error.
func WithBucketDefault(def Expr) Option[bucketOptions] {
	return func(o *bucketOptions) { o.def = def }
}

// WithBucketOutput sets the per-bucket output fields, computed with
// accumulators (construct via Accumulate). Without it, each bucket has only a
// count field.
func WithBucketOutput(output ...GroupField) Option[bucketOptions] {
	return func(o *bucketOptions) { o.output = output }
}

// BucketStage produces a $bucket stage that groups documents into buckets by
// groupBy, using boundaries as the inclusive lower / exclusive upper edges of
// each bucket. boundaries must hold at least two values in ascending order.
func BucketStage(groupBy Expr, boundaries []any, opts ...Option[bucketOptions]) Stage {
	var o bucketOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{
		{Key: "groupBy", Value: groupBy},
		{Key: "boundaries", Value: boundaries},
	}
	if o.def != nil {
		doc = append(doc, bson.E{Key: "default", Value: o.def})
	}
	if o.output != nil {
		doc = append(doc, bson.E{Key: "output", Value: groupFieldsDoc(o.output)})
	}
	return Stage{{Key: "$bucket", Value: doc}}
}

// --- $bucketAuto ---

type bucketAutoOptions struct {
	output      []GroupField
	granularity any
}

// WithBucketAutoOutput sets the per-bucket output fields, computed with
// accumulators (construct via Accumulate). Without it, each bucket has only a
// count field.
func WithBucketAutoOutput(output ...GroupField) Option[bucketAutoOptions] {
	return func(o *bucketAutoOptions) { o.output = output }
}

// WithBucketAutoGranularity sets the preferred number series used to round
// bucket boundary edges (e.g. "R5", "1-2-5", "POWERSOF2"). Valid only when all
// groupBy values are numeric.
func WithBucketAutoGranularity(granularity string) Option[bucketAutoOptions] {
	return func(o *bucketAutoOptions) { o.granularity = granularity }
}

// BucketAutoStage produces a $bucketAuto stage that distributes documents into
// buckets automatically-determined groups by groupBy, attempting to spread them
// evenly across the given number of buckets.
func BucketAutoStage(groupBy Expr, buckets int32, opts ...Option[bucketAutoOptions]) Stage {
	var o bucketAutoOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{
		{Key: "groupBy", Value: groupBy},
		{Key: "buckets", Value: buckets},
	}
	if o.output != nil {
		doc = append(doc, bson.E{Key: "output", Value: groupFieldsDoc(o.output)})
	}
	if o.granularity != nil {
		doc = append(doc, bson.E{Key: "granularity", Value: o.granularity})
	}
	return Stage{{Key: "$bucketAuto", Value: doc}}
}

// --- $changeStream ---

type changeStreamOptions struct {
	allChangesForCluster     any
	fullDocument             any
	fullDocumentBeforeChange any
	resumeAfter              any
	showExpandedEvents       any
	startAfter               any
	startAtOperationTime     any
}

// WithChangeStreamAllChangesForCluster reports all changes across the
// deployment, except on internal databases and collections.
func WithChangeStreamAllChangesForCluster(v bool) Option[changeStreamOptions] {
	return func(o *changeStreamOptions) { o.allChangesForCluster = v }
}

// WithChangeStreamFullDocument controls whether update events include a copy of
// the full modified document: one of "default", "updateLookup", "whenAvailable",
// or "required".
func WithChangeStreamFullDocument(fullDocument string) Option[changeStreamOptions] {
	return func(o *changeStreamOptions) { o.fullDocument = fullDocument }
}

// WithChangeStreamFullDocumentBeforeChange controls whether events include the
// pre-image of the modified document: one of "off", "whenAvailable", or
// "required".
func WithChangeStreamFullDocumentBeforeChange(value string) Option[changeStreamOptions] {
	return func(o *changeStreamOptions) { o.fullDocumentBeforeChange = value }
}

// WithChangeStreamResumeAfter resumes the stream after the given resume token.
// Cannot be combined with WithChangeStreamStartAfter or
// WithChangeStreamStartAtOperationTime.
//
// The spec types this field as an int, but a resume token is a document; the
// option accepts any so the real token document can be passed.
func WithChangeStreamResumeAfter(token any) Option[changeStreamOptions] {
	return func(o *changeStreamOptions) { o.resumeAfter = token }
}

// WithChangeStreamShowExpandedEvents includes additional change events such as
// DDL and index operations.
func WithChangeStreamShowExpandedEvents(v bool) Option[changeStreamOptions] {
	return func(o *changeStreamOptions) { o.showExpandedEvents = v }
}

// WithChangeStreamStartAfter starts the stream after the given resume token,
// including after an invalidate event. Cannot be combined with
// WithChangeStreamResumeAfter or WithChangeStreamStartAtOperationTime.
func WithChangeStreamStartAfter(token any) Option[changeStreamOptions] {
	return func(o *changeStreamOptions) { o.startAfter = token }
}

// WithChangeStreamStartAtOperationTime starts the stream at the given cluster
// time. Cannot be combined with WithChangeStreamResumeAfter or
// WithChangeStreamStartAfter.
func WithChangeStreamStartAtOperationTime(ts bson.Timestamp) Option[changeStreamOptions] {
	return func(o *changeStreamOptions) { o.startAtOperationTime = ts }
}

// ChangeStreamStage produces a $changeStream stage that opens a change stream
// cursor for the collection or database. It must be the first stage in the
// pipeline.
func ChangeStreamStage(opts ...Option[changeStreamOptions]) Stage {
	var o changeStreamOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{}
	if o.allChangesForCluster != nil {
		doc = append(doc, bson.E{Key: "allChangesForCluster", Value: o.allChangesForCluster})
	}
	if o.fullDocument != nil {
		doc = append(doc, bson.E{Key: "fullDocument", Value: o.fullDocument})
	}
	if o.fullDocumentBeforeChange != nil {
		doc = append(doc, bson.E{Key: "fullDocumentBeforeChange", Value: o.fullDocumentBeforeChange})
	}
	if o.resumeAfter != nil {
		doc = append(doc, bson.E{Key: "resumeAfter", Value: o.resumeAfter})
	}
	if o.showExpandedEvents != nil {
		doc = append(doc, bson.E{Key: "showExpandedEvents", Value: o.showExpandedEvents})
	}
	if o.startAfter != nil {
		doc = append(doc, bson.E{Key: "startAfter", Value: o.startAfter})
	}
	if o.startAtOperationTime != nil {
		doc = append(doc, bson.E{Key: "startAtOperationTime", Value: o.startAtOperationTime})
	}
	return Stage{{Key: "$changeStream", Value: doc}}
}

// --- $changeStreamSplitLargeEvent ---

// ChangeStreamSplitLargeEventStage produces a $changeStreamSplitLargeEvent stage
// that splits change stream events exceeding 16 MB into smaller fragments. It
// can only be used as the final stage of a $changeStream pipeline.
func ChangeStreamSplitLargeEventStage() Stage {
	return Stage{{Key: "$changeStreamSplitLargeEvent", Value: bson.D{}}}
}

// --- $collStats ---

type collStatsOptions struct {
	latencyStats   any
	storageStats   any
	count          any
	queryExecStats any
	haveLatency    bool
	haveStorage    bool
	haveCount      bool
	haveQueryExec  bool
}

// WithCollStatsLatencyStats includes latency statistics in the output. Set
// histograms to true to include latency histograms.
func WithCollStatsLatencyStats(histograms bool) Option[collStatsOptions] {
	return func(o *collStatsOptions) {
		o.latencyStats = bson.D{{Key: "histograms", Value: histograms}}
		o.haveLatency = true
	}
}

// WithCollStatsStorageStats includes storage statistics in the output.
//
// The optional scale factor and other storageStats fields are not modeled; the
// builder emits an empty storageStats document.
func WithCollStatsStorageStats() Option[collStatsOptions] {
	return func(o *collStatsOptions) {
		o.storageStats = bson.D{}
		o.haveStorage = true
	}
}

// WithCollStatsCount includes the document count in the output.
func WithCollStatsCount() Option[collStatsOptions] {
	return func(o *collStatsOptions) {
		o.count = bson.D{}
		o.haveCount = true
	}
}

// WithCollStatsQueryExecStats includes query execution statistics in the output.
func WithCollStatsQueryExecStats() Option[collStatsOptions] {
	return func(o *collStatsOptions) {
		o.queryExecStats = bson.D{}
		o.haveQueryExec = true
	}
}

// CollStatsStage produces a $collStats stage that returns statistics about the
// collection or view. Select which statistics to include with the With*
// options. It must be the first stage in the pipeline.
func CollStatsStage(opts ...Option[collStatsOptions]) Stage {
	var o collStatsOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{}
	if o.haveLatency {
		doc = append(doc, bson.E{Key: "latencyStats", Value: o.latencyStats})
	}
	if o.haveStorage {
		doc = append(doc, bson.E{Key: "storageStats", Value: o.storageStats})
	}
	if o.haveCount {
		doc = append(doc, bson.E{Key: "count", Value: o.count})
	}
	if o.haveQueryExec {
		doc = append(doc, bson.E{Key: "queryExecStats", Value: o.queryExecStats})
	}
	return Stage{{Key: "$collStats", Value: doc}}
}

// --- $count ---

// CountStage produces a $count stage that counts the documents reaching this
// point in the pipeline and stores the total in the named output field.
// (Distinct from the $count accumulator.)
func CountStage(field string) Stage {
	return Stage{{Key: "$count", Value: field}}
}

// --- $currentOp ---

type currentOpOptions struct {
	allUsers        any
	idleConnections any
	idleCursors     any
	idleSessions    any
	localOps        any
}

// WithCurrentOpAllUsers reports operations for all users, not just the current
// one.
func WithCurrentOpAllUsers(v bool) Option[currentOpOptions] {
	return func(o *currentOpOptions) { o.allUsers = v }
}

// WithCurrentOpIdleConnections includes idle connections in the output.
func WithCurrentOpIdleConnections(v bool) Option[currentOpOptions] {
	return func(o *currentOpOptions) { o.idleConnections = v }
}

// WithCurrentOpIdleCursors includes idle cursors in the output.
func WithCurrentOpIdleCursors(v bool) Option[currentOpOptions] {
	return func(o *currentOpOptions) { o.idleCursors = v }
}

// WithCurrentOpIdleSessions includes idle sessions in the output.
func WithCurrentOpIdleSessions(v bool) Option[currentOpOptions] {
	return func(o *currentOpOptions) { o.idleSessions = v }
}

// WithCurrentOpLocalOps reports operations local to the targeted mongos rather
// than the shards they run on.
func WithCurrentOpLocalOps(v bool) Option[currentOpOptions] {
	return func(o *currentOpOptions) { o.localOps = v }
}

// CurrentOpStage produces a $currentOp stage that reports active and/or dormant
// operations for the deployment. It must be the first stage and be run against
// the admin database.
func CurrentOpStage(opts ...Option[currentOpOptions]) Stage {
	var o currentOpOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{}
	if o.allUsers != nil {
		doc = append(doc, bson.E{Key: "allUsers", Value: o.allUsers})
	}
	if o.idleConnections != nil {
		doc = append(doc, bson.E{Key: "idleConnections", Value: o.idleConnections})
	}
	if o.idleCursors != nil {
		doc = append(doc, bson.E{Key: "idleCursors", Value: o.idleCursors})
	}
	if o.idleSessions != nil {
		doc = append(doc, bson.E{Key: "idleSessions", Value: o.idleSessions})
	}
	if o.localOps != nil {
		doc = append(doc, bson.E{Key: "localOps", Value: o.localOps})
	}
	return Stage{{Key: "$currentOp", Value: doc}}
}

// --- $densify ---

// DensifyBounds specifies the range over which $densify generates missing
// documents. Construct via DensifyBoundsFull, DensifyBoundsPartition, or
// DensifyBoundsValues.
type DensifyBounds struct{ v any }

// DensifyBoundsFull densifies across the full range of values in the collection.
func DensifyBoundsFull() DensifyBounds { return DensifyBounds{v: "full"} }

// DensifyBoundsPartition densifies across the range of values within each
// partition.
func DensifyBoundsPartition() DensifyBounds { return DensifyBounds{v: "partition"} }

// DensifyBoundsValues densifies between the explicit lower (inclusive) and upper
// (exclusive) bounds. Both must be numbers, or both dates matching the field.
func DensifyBoundsValues(lower, upper any) DensifyBounds {
	return DensifyBounds{v: bson.A{lower, upper}}
}

type densifyOptions struct {
	unit              any
	partitionByFields []string
}

// WithDensifyUnit sets the time unit of step for date densification (e.g.
// "hour"). Omit it for numeric densification.
func WithDensifyUnit(unit string) Option[densifyOptions] {
	return func(o *densifyOptions) { o.unit = unit }
}

// WithDensifyPartitionByFields densifies independently within partitions
// defined by the compound key of the named fields.
func WithDensifyPartitionByFields(fields ...string) Option[densifyOptions] {
	return func(o *densifyOptions) { o.partitionByFields = fields }
}

// DensifyStage produces a $densify stage that fills gaps in field by generating
// documents at the given step over bounds. field values must be all numeric or
// all dates.
func DensifyStage[T Number](field string, step T, bounds DensifyBounds, opts ...Option[densifyOptions]) Stage {
	var o densifyOptions
	for _, opt := range opts {
		opt(&o)
	}
	rng := bson.D{
		{Key: "bounds", Value: bounds.v},
		{Key: "step", Value: step},
	}
	if o.unit != nil {
		rng = append(rng, bson.E{Key: "unit", Value: o.unit})
	}
	doc := bson.D{{Key: "field", Value: field}}
	if o.partitionByFields != nil {
		doc = append(doc, bson.E{Key: "partitionByFields", Value: o.partitionByFields})
	}
	doc = append(doc, bson.E{Key: "range", Value: rng})
	return Stage{{Key: "$densify", Value: doc}}
}

// --- $documents ---

// DocumentsStage produces a $documents stage that emits the documents the given
// expression resolves to. documents may be a literal array of objects (e.g.
// bson.A{bson.D{...}, ...}) or any expression that resolves to an array of
// objects — a system variable such as "$$SEARCH_META", a Let expression, or a
// variable in scope from a $lookup. Expressions that reference the current
// document (e.g. "$myField" or "$$ROOT") resolve to an error at runtime.
func DocumentsStage[T ArrayResolver](documents T) Stage {
	return Stage{{Key: "$documents", Value: documents}}
}

// --- $facet ---

// FacetField pairs a facet name with the sub-pipeline to run for it in a $facet
// stage. Construct via Facet.
type FacetField interface{ facetField() facetField }

type facetField struct {
	name     string
	pipeline Pipeline
}

func (ff facetField) facetField() facetField { return ff }

// Facet creates a FacetField that runs the given sub-pipeline under name.
func Facet(name string, stages ...Stage) FacetField {
	return facetField{name: name, pipeline: stages}
}

// FacetStage produces a $facet stage that runs multiple sub-pipelines over the
// same input documents, collecting each named result into its own field.
func FacetStage(facets ...FacetField) Stage {
	doc := make(bson.D, len(facets))
	for i, f := range facets {
		ff := f.facetField()
		doc[i] = bson.E{Key: ff.name, Value: pipelineDocs(ff.pipeline)}
	}
	return Stage{{Key: "$facet", Value: doc}}
}

// --- $fill ---

// FillOutput specifies how a single field is filled in a $fill stage. Construct
// via FillWithValue or FillWithMethod.
type FillOutput interface{ fillOutput() fillOutput }

type fillOutput struct {
	field string
	spec  bson.D
}

func (fo fillOutput) fillOutput() fillOutput { return fo }

// FillWithValue fills missing values of field with the given constant or
// computed value.
func FillWithValue(field string, value Expr) FillOutput {
	return fillOutput{field: field, spec: bson.D{{Key: "value", Value: value}}}
}

// FillWithMethod fills missing values of field using the given interpolation
// method: "linear" (interpolate from surrounding values and the sortBy field)
// or "locf" (carry the last observed value forward).
func FillWithMethod(field, method string) FillOutput {
	return fillOutput{field: field, spec: bson.D{{Key: "method", Value: method}}}
}

type fillOptions struct {
	partitionBy       any
	partitionByFields []string
	sortBy            []SortField
}

// WithFillPartitionBy groups documents into partitions by the given expression.
// Mutually exclusive with WithFillPartitionByFields.
func WithFillPartitionBy(expr Expr) Option[fillOptions] {
	return func(o *fillOptions) { o.partitionBy = expr }
}

// WithFillPartitionByFields groups documents into partitions by the compound
// key of the named fields. Mutually exclusive with WithFillPartitionBy.
func WithFillPartitionByFields(fields ...string) Option[fillOptions] {
	return func(o *fillOptions) { o.partitionByFields = fields }
}

// WithFillSortBy sorts documents within each partition (construct via Sort),
// as required by the "linear" and "locf" fill methods.
func WithFillSortBy(fields ...SortField) Option[fillOptions] {
	return func(o *fillOptions) { o.sortBy = fields }
}

// FillStage produces a $fill stage that populates null and missing values of
// the given output fields (construct via FillWithValue / FillWithMethod).
func FillStage(output []FillOutput, opts ...Option[fillOptions]) Stage {
	var o fillOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{}
	if o.sortBy != nil {
		doc = append(doc, bson.E{Key: "sortBy", Value: sortFieldsDoc(o.sortBy)})
	}
	if o.partitionBy != nil {
		doc = append(doc, bson.E{Key: "partitionBy", Value: o.partitionBy})
	}
	if o.partitionByFields != nil {
		doc = append(doc, bson.E{Key: "partitionByFields", Value: o.partitionByFields})
	}
	outDoc := make(bson.D, len(output))
	for i, f := range output {
		fo := f.fillOutput()
		outDoc[i] = bson.E{Key: fo.field, Value: fo.spec}
	}
	doc = append(doc, bson.E{Key: "output", Value: outDoc})
	return Stage{{Key: "$fill", Value: doc}}
}

// --- $geoNear ---

type geoNearOptions struct {
	distanceField      any
	distanceMultiplier any
	includeLocs        any
	key                any
	maxDistance        any
	minDistance        any
	query              query.Filter
	spherical          any
}

// WithGeoNearDistanceField names the output field holding each document's
// calculated distance.
func WithGeoNearDistanceField(field string) Option[geoNearOptions] {
	return func(o *geoNearOptions) { o.distanceField = field }
}

// WithGeoNearDistanceMultiplier scales all returned distances by factor.
func WithGeoNearDistanceMultiplier[T Number](factor T) Option[geoNearOptions] {
	return func(o *geoNearOptions) { o.distanceMultiplier = factor }
}

// WithGeoNearIncludeLocs names the output field identifying the location used
// to calculate the distance.
func WithGeoNearIncludeLocs(field string) Option[geoNearOptions] {
	return func(o *geoNearOptions) { o.includeLocs = field }
}

// WithGeoNearKey selects the geospatial indexed field to use for the distance
// calculation.
func WithGeoNearKey(field string) Option[geoNearOptions] {
	return func(o *geoNearOptions) { o.key = field }
}

// WithGeoNearMaxDistance limits results to documents within this distance of
// near (meters for GeoJSON, radians for legacy coordinates).
func WithGeoNearMaxDistance[T Number](d T) Option[geoNearOptions] {
	return func(o *geoNearOptions) { o.maxDistance = d }
}

// WithGeoNearMinDistance limits results to documents beyond this distance from
// near (meters for GeoJSON, radians for legacy coordinates).
func WithGeoNearMinDistance[T Number](d T) Option[geoNearOptions] {
	return func(o *geoNearOptions) { o.minDistance = d }
}

// WithGeoNearQuery limits results to documents matching the given filters
// (merged as an implicit AND), using the same syntax as MatchStage.
func WithGeoNearQuery(filters ...query.Filter) Option[geoNearOptions] {
	return func(o *geoNearOptions) {
		merged := query.Filter{}
		for _, f := range filters {
			merged = append(merged, f...)
		}
		o.query = merged
	}
}

// WithGeoNearSpherical uses spherical geometry ($nearSphere semantics) for the
// distance calculation. Defaults to false when omitted.
func WithGeoNearSpherical(v bool) Option[geoNearOptions] {
	return func(o *geoNearOptions) { o.spherical = v }
}

// GeoNearStage produces a $geoNear stage that returns documents ordered by
// proximity to near. near is a GeoJSON Point document, a legacy coordinate
// pair, or an expression resolving to one. It must be the first stage in the
// pipeline.
func GeoNearStage(near Expr, opts ...Option[geoNearOptions]) Stage {
	var o geoNearOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{{Key: "near", Value: near}}
	if o.distanceField != nil {
		doc = append(doc, bson.E{Key: "distanceField", Value: o.distanceField})
	}
	if o.distanceMultiplier != nil {
		doc = append(doc, bson.E{Key: "distanceMultiplier", Value: o.distanceMultiplier})
	}
	if o.maxDistance != nil {
		doc = append(doc, bson.E{Key: "maxDistance", Value: o.maxDistance})
	}
	if o.minDistance != nil {
		doc = append(doc, bson.E{Key: "minDistance", Value: o.minDistance})
	}
	if o.query != nil {
		doc = append(doc, bson.E{Key: "query", Value: o.query})
	}
	if o.includeLocs != nil {
		doc = append(doc, bson.E{Key: "includeLocs", Value: o.includeLocs})
	}
	if o.key != nil {
		doc = append(doc, bson.E{Key: "key", Value: o.key})
	}
	if o.spherical != nil {
		doc = append(doc, bson.E{Key: "spherical", Value: o.spherical})
	}
	return Stage{{Key: "$geoNear", Value: doc}}
}

// --- $graphLookup ---

type graphLookupOptions struct {
	maxDepth                any
	depthField              any
	restrictSearchWithMatch query.Filter
}

// WithGraphLookupMaxDepth caps the recursion depth of the search.
func WithGraphLookupMaxDepth(maxDepth int32) Option[graphLookupOptions] {
	return func(o *graphLookupOptions) { o.maxDepth = maxDepth }
}

// WithGraphLookupDepthField names a field added to each traversed document
// holding its recursion depth (starting at zero).
func WithGraphLookupDepthField(field string) Option[graphLookupOptions] {
	return func(o *graphLookupOptions) { o.depthField = field }
}

// WithGraphLookupRestrictSearchWithMatch limits the recursive search to
// documents matching the given filters (merged as an implicit AND), using the
// same syntax as MatchStage.
func WithGraphLookupRestrictSearchWithMatch(filters ...query.Filter) Option[graphLookupOptions] {
	return func(o *graphLookupOptions) {
		merged := query.Filter{}
		for _, f := range filters {
			merged = append(merged, f...)
		}
		o.restrictSearchWithMatch = merged
	}
}

// GraphLookupStage produces a $graphLookup stage that recursively searches the
// from collection, matching connectFromField against connectToField starting
// from startWith, and stores the traversal results in an array field named by
// as.
func GraphLookupStage(from string, startWith Expr, connectFromField, connectToField, as string, opts ...Option[graphLookupOptions]) Stage {
	var o graphLookupOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{
		{Key: "from", Value: from},
		{Key: "startWith", Value: startWith},
		{Key: "connectFromField", Value: connectFromField},
		{Key: "connectToField", Value: connectToField},
		{Key: "as", Value: as},
	}
	if o.maxDepth != nil {
		doc = append(doc, bson.E{Key: "maxDepth", Value: o.maxDepth})
	}
	if o.depthField != nil {
		doc = append(doc, bson.E{Key: "depthField", Value: o.depthField})
	}
	if o.restrictSearchWithMatch != nil {
		doc = append(doc, bson.E{Key: "restrictSearchWithMatch", Value: o.restrictSearchWithMatch})
	}
	return Stage{{Key: "$graphLookup", Value: doc}}
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
func GroupStage(_id Expr, fields ...GroupField) Stage {
	doc := make(bson.D, 0, len(fields)+1)
	doc = append(doc, bson.E{Key: "_id", Value: _id})
	for _, f := range fields {
		gf := f.groupField()
		doc = append(doc, bson.E{Key: gf.name, Value: gf.accumulator})
	}
	return Stage{{Key: "$group", Value: doc}}
}

// --- $indexStats ---

// IndexStatsStage produces an $indexStats stage that returns usage statistics
// for each index on the collection.
func IndexStatsStage() Stage {
	return Stage{{Key: "$indexStats", Value: bson.D{}}}
}

// --- $limit ---

// LimitStage produces a $limit stage that restricts the pipeline to the first
// limit documents.
func LimitStage(limit int32) Stage {
	return Stage{{Key: "$limit", Value: limit}}
}

// --- $listLocalSessions ---

// WithListLocalSessionsUsers lists local sessions only for the specified users.
func WithListLocalSessionsUsers(users ...SessionUser) Option[listSessionsOptions] {
	return WithListSessionsUsers(users...)
}

// WithListLocalSessionsAllUsers lists local sessions for all users.
func WithListLocalSessionsAllUsers(v bool) Option[listSessionsOptions] {
	return WithListSessionsAllUsers(v)
}

// ListLocalSessionsStage produces a $listLocalSessions stage that lists sessions
// active on the currently connected instance, which may not yet have propagated
// to the system.sessions collection. It must be the first stage. It shares the
// listSessions options.
func ListLocalSessionsStage(opts ...Option[listSessionsOptions]) Stage {
	var o listSessionsOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{}
	if o.users != nil {
		doc = append(doc, bson.E{Key: "users", Value: sessionUsersArray(o.users)})
	}
	if o.allUsers != nil {
		doc = append(doc, bson.E{Key: "allUsers", Value: o.allUsers})
	}
	return Stage{{Key: "$listLocalSessions", Value: doc}}
}

// --- $listSampledQueries ---

type listSampledQueriesOptions struct {
	namespace any
}

// WithListSampledQueriesNamespace limits results to sampled queries for the
// given namespace (database.collection).
func WithListSampledQueriesNamespace(namespace string) Option[listSampledQueriesOptions] {
	return func(o *listSampledQueriesOptions) { o.namespace = namespace }
}

// ListSampledQueriesStage produces a $listSampledQueries stage that lists sampled
// queries for all collections or, with WithListSampledQueriesNamespace, a
// specific collection.
func ListSampledQueriesStage(opts ...Option[listSampledQueriesOptions]) Stage {
	var o listSampledQueriesOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{}
	if o.namespace != nil {
		doc = append(doc, bson.E{Key: "namespace", Value: o.namespace})
	}
	return Stage{{Key: "$listSampledQueries", Value: doc}}
}

// --- $listSearchIndexes ---

type listSearchIndexesOptions struct {
	id   any
	name any
}

// WithListSearchIndexesID returns information only about the index with the
// given id.
func WithListSearchIndexesID(id string) Option[listSearchIndexesOptions] {
	return func(o *listSearchIndexesOptions) { o.id = id }
}

// WithListSearchIndexesName returns information only about the index with the
// given name.
func WithListSearchIndexesName(name string) Option[listSearchIndexesOptions] {
	return func(o *listSearchIndexesOptions) { o.name = name }
}

// ListSearchIndexesStage produces a $listSearchIndexes stage that returns
// information about Atlas Search indexes on the collection. Without options it
// returns all indexes.
func ListSearchIndexesStage(opts ...Option[listSearchIndexesOptions]) Stage {
	var o listSearchIndexesOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{}
	if o.id != nil {
		doc = append(doc, bson.E{Key: "id", Value: o.id})
	}
	if o.name != nil {
		doc = append(doc, bson.E{Key: "name", Value: o.name})
	}
	return Stage{{Key: "$listSearchIndexes", Value: doc}}
}

// --- $listSessions ---

// WithListSessionsUsers lists sessions only for the specified users.
func WithListSessionsUsers(users ...SessionUser) Option[listSessionsOptions] {
	return func(o *listSessionsOptions) { o.users = users }
}

// WithListSessionsAllUsers lists sessions for all users.
func WithListSessionsAllUsers(v bool) Option[listSessionsOptions] {
	return func(o *listSessionsOptions) { o.allUsers = v }
}

// ListSessionsStage produces a $listSessions stage that lists sessions that have
// propagated to the system.sessions collection. It must be the first stage and
// be run against the system.sessions collection in the config database.
func ListSessionsStage(opts ...Option[listSessionsOptions]) Stage {
	var o listSessionsOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{}
	if o.users != nil {
		doc = append(doc, bson.E{Key: "users", Value: sessionUsersArray(o.users)})
	}
	if o.allUsers != nil {
		doc = append(doc, bson.E{Key: "allUsers", Value: o.allUsers})
	}
	return Stage{{Key: "$listSessions", Value: doc}}
}

// --- $lookup ---

type lookupOptions struct {
	from         any
	localField   any
	foreignField any
	let          []SetField
	pipeline     []Stage
}

// WithLookupFrom sets the collection in the same database to join with. It may
// be omitted when the join uses a $documents stage in WithLookupPipeline.
func WithLookupFrom(coll string) Option[lookupOptions] {
	return func(o *lookupOptions) { o.from = coll }
}

// WithLookupLocalField sets the input-document field matched against the
// foreign field. Used together with WithLookupForeignField for equality joins.
func WithLookupLocalField(field string) Option[lookupOptions] {
	return func(o *lookupOptions) { o.localField = field }
}

// WithLookupForeignField sets the joined-collection field matched against the
// local field.
func WithLookupForeignField(field string) Option[lookupOptions] {
	return func(o *lookupOptions) { o.foreignField = field }
}

// WithLookupLet declares variables (construct via Assign) that expose
// input-document values to the WithLookupPipeline stages.
func WithLookupLet(vars ...SetField) Option[lookupOptions] {
	return func(o *lookupOptions) { o.let = vars }
}

// WithLookupPipeline sets the pipeline to run on the joined collection. It
// cannot contain $out or $merge.
func WithLookupPipeline(stages ...Stage) Option[lookupOptions] {
	return func(o *lookupOptions) { o.pipeline = stages }
}

// LookupStage produces a $lookup stage that performs a left outer join, adding
// the matching joined documents as an array field named by as. Configure the
// join with the With* options: an equality join uses WithLookupFrom,
// WithLookupLocalField, and WithLookupForeignField; a subquery join uses
// WithLookupPipeline (with WithLookupLet).
func LookupStage(as string, opts ...Option[lookupOptions]) Stage {
	var o lookupOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{}
	if o.from != nil {
		doc = append(doc, bson.E{Key: "from", Value: o.from})
	}
	if o.localField != nil {
		doc = append(doc, bson.E{Key: "localField", Value: o.localField})
	}
	if o.foreignField != nil {
		doc = append(doc, bson.E{Key: "foreignField", Value: o.foreignField})
	}
	if o.let != nil {
		doc = append(doc, bson.E{Key: "let", Value: setFieldsDoc(o.let)})
	}
	if o.pipeline != nil {
		doc = append(doc, bson.E{Key: "pipeline", Value: pipelineDocs(o.pipeline)})
	}
	doc = append(doc, bson.E{Key: "as", Value: as})
	return Stage{{Key: "$lookup", Value: doc}}
}

// --- $match ---

// MatchStage produces a $match stage from one or more query.Filters. Multiple
// filters are merged into a single document (implicit AND). Build filters in
// the query sub-package, e.g.:
//
//	agg.MatchStage(query.Field("qty", query.Gt(20)), query.Field("name", query.Eq("Alice")))
func MatchStage(filters ...query.Filter) Stage {
	merged := bson.D{}
	for _, f := range filters {
		merged = append(merged, f...)
	}
	return Stage{{Key: "$match", Value: merged}}
}

// --- $merge ---

type mergeOptions struct {
	intoDB         any
	on             []string
	let            []SetField
	whenMatched    any
	whenNotMatched any
}

// WithMergeIntoDB writes to the output collection in the named database instead
// of the current one.
func WithMergeIntoDB(db string) Option[mergeOptions] {
	return func(o *mergeOptions) { o.intoDB = db }
}

// WithMergeOn sets the field(s) that uniquely identify a document for matching
// against the output collection. A single field is emitted as a string;
// multiple fields as an array.
func WithMergeOn(fields ...string) Option[mergeOptions] {
	return func(o *mergeOptions) { o.on = fields }
}

// WithMergeLet declares variables (construct via Assign) for use in the
// whenMatched pipeline set by WithMergeWhenMatchedPipeline.
func WithMergeLet(vars ...SetField) Option[mergeOptions] {
	return func(o *mergeOptions) { o.let = vars }
}

// WithMergeWhenMatched sets the action taken when a result document matches an
// existing document: one of "replace", "keepExisting", "merge", or "fail". For
// a custom update pipeline, use WithMergeWhenMatchedPipeline instead.
func WithMergeWhenMatched(action string) Option[mergeOptions] {
	return func(o *mergeOptions) { o.whenMatched = action }
}

// WithMergeWhenMatchedPipeline sets a custom update pipeline applied when a
// result document matches an existing document. Mutually exclusive with
// WithMergeWhenMatched.
func WithMergeWhenMatchedPipeline(stages ...Stage) Option[mergeOptions] {
	return func(o *mergeOptions) { o.whenMatched = pipelineDocs(stages) }
}

// WithMergeWhenNotMatched sets the action taken when a result document does not
// match any existing document: one of "insert", "discard", or "fail".
func WithMergeWhenNotMatched(action string) Option[mergeOptions] {
	return func(o *mergeOptions) { o.whenNotMatched = action }
}

// MergeStage produces a $merge stage that writes the pipeline results into coll,
// incorporating them according to the whenMatched / whenNotMatched behavior. It
// must be the last stage in the pipeline. Use WithMergeIntoDB to target a
// different database.
func MergeStage(coll string, opts ...Option[mergeOptions]) Stage {
	var o mergeOptions
	for _, opt := range opts {
		opt(&o)
	}
	var into any = coll
	if o.intoDB != nil {
		into = bson.D{{Key: "db", Value: o.intoDB}, {Key: "coll", Value: coll}}
	}
	doc := bson.D{{Key: "into", Value: into}}
	if o.on != nil {
		if len(o.on) == 1 {
			doc = append(doc, bson.E{Key: "on", Value: o.on[0]})
		} else {
			doc = append(doc, bson.E{Key: "on", Value: o.on})
		}
	}
	if o.let != nil {
		doc = append(doc, bson.E{Key: "let", Value: setFieldsDoc(o.let)})
	}
	if o.whenMatched != nil {
		doc = append(doc, bson.E{Key: "whenMatched", Value: o.whenMatched})
	}
	if o.whenNotMatched != nil {
		doc = append(doc, bson.E{Key: "whenNotMatched", Value: o.whenNotMatched})
	}
	return Stage{{Key: "$merge", Value: doc}}
}

// --- $out ---

// Timeseries configures $out when writing to a time series collection.
// Construct via NewTimeseries.
type Timeseries struct {
	timeField             string
	metaField             any
	granularity           any
	bucketMaxSpanSeconds  any
	bucketRoundingSeconds any
}

type timeseriesOptions struct {
	metaField             any
	granularity           any
	bucketMaxSpanSeconds  any
	bucketRoundingSeconds any
}

// WithTimeseriesMetaField sets the field holding metadata in each document.
func WithTimeseriesMetaField(field string) Option[timeseriesOptions] {
	return func(o *timeseriesOptions) { o.metaField = field }
}

// WithTimeseriesGranularity sets the granularity of time measurements: one of
// "seconds", "minutes", or "hours".
func WithTimeseriesGranularity(granularity string) Option[timeseriesOptions] {
	return func(o *timeseriesOptions) { o.granularity = granularity }
}

// WithTimeseriesBucketMaxSpanSeconds sets the maximum time span between
// measurements in a bucket.
func WithTimeseriesBucketMaxSpanSeconds(seconds int32) Option[timeseriesOptions] {
	return func(o *timeseriesOptions) { o.bucketMaxSpanSeconds = seconds }
}

// WithTimeseriesBucketRoundingSeconds sets the interval that determines the
// starting timestamp for a new bucket.
func WithTimeseriesBucketRoundingSeconds(seconds int32) Option[timeseriesOptions] {
	return func(o *timeseriesOptions) { o.bucketRoundingSeconds = seconds }
}

// NewTimeseries builds a time series configuration for WithOutTimeseries.
// timeField (the date field of each document) is required. It uses the New
// prefix because the exported type name Timeseries is unavailable as a
// constructor name.
func NewTimeseries(timeField string, opts ...Option[timeseriesOptions]) Timeseries {
	var o timeseriesOptions
	for _, opt := range opts {
		opt(&o)
	}
	return Timeseries{
		timeField:             timeField,
		metaField:             o.metaField,
		granularity:           o.granularity,
		bucketMaxSpanSeconds:  o.bucketMaxSpanSeconds,
		bucketRoundingSeconds: o.bucketRoundingSeconds,
	}
}

func (t Timeseries) doc() bson.D {
	doc := bson.D{{Key: "timeField", Value: t.timeField}}
	if t.metaField != nil {
		doc = append(doc, bson.E{Key: "metaField", Value: t.metaField})
	}
	if t.granularity != nil {
		doc = append(doc, bson.E{Key: "granularity", Value: t.granularity})
	}
	if t.bucketMaxSpanSeconds != nil {
		doc = append(doc, bson.E{Key: "bucketMaxSpanSeconds", Value: t.bucketMaxSpanSeconds})
	}
	if t.bucketRoundingSeconds != nil {
		doc = append(doc, bson.E{Key: "bucketRoundingSeconds", Value: t.bucketRoundingSeconds})
	}
	return doc
}

type outOptions struct {
	db         any
	timeseries *Timeseries
}

// WithOutDB writes to coll in the named database instead of the current one.
func WithOutDB(db string) Option[outOptions] {
	return func(o *outOptions) { o.db = db }
}

// WithOutTimeseries writes to a time series collection with the given
// configuration (construct via NewTimeseries).
func WithOutTimeseries(ts Timeseries) Option[outOptions] {
	return func(o *outOptions) { o.timeseries = &ts }
}

// OutStage produces an $out stage that writes the pipeline results to coll. It
// must be the last stage in the pipeline. The builder always emits the verbose
// (object) form.
func OutStage(coll string, opts ...Option[outOptions]) Stage {
	var o outOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{}
	if o.db != nil {
		doc = append(doc, bson.E{Key: "db", Value: o.db})
	}
	doc = append(doc, bson.E{Key: "coll", Value: coll})
	if o.timeseries != nil {
		doc = append(doc, bson.E{Key: "timeseries", Value: o.timeseries.doc()})
	}
	return Stage{{Key: "$out", Value: doc}}
}

// --- $planCacheStats ---

// PlanCacheStatsStage produces a $planCacheStats stage that returns plan cache
// information for the collection.
func PlanCacheStatsStage() Stage {
	return Stage{{Key: "$planCacheStats", Value: bson.D{}}}
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

// --- $redact ---

// RedactStage produces a $redact stage that restricts document content based on
// expr. The expression is evaluated at each document level and must resolve to
// one of the system variables $$DESCEND, $$PRUNE, or $$KEEP.
func RedactStage(expr Expr) Stage {
	return Stage{{Key: "$redact", Value: expr}}
}

// --- $replaceRoot ---

// ReplaceRootStage produces a $replaceRoot stage that promotes newRoot to the
// top level of each document, replacing all existing fields. newRoot must
// resolve to an object.
func ReplaceRootStage[T ObjectResolver](newRoot T) Stage {
	return Stage{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: newRoot}}}}
}

// --- $replaceWith ---

// ReplaceWithStage produces a $replaceWith stage that replaces each document
// with expr, which must resolve to an object. It is an alias for $replaceRoot.
func ReplaceWithStage[T ObjectResolver](expr T) Stage {
	return Stage{{Key: "$replaceWith", Value: expr}}
}

// --- $sample ---

// SampleStage produces a $sample stage that randomly selects size documents
// from its input.
func SampleStage(size int32) Stage {
	return Stage{{Key: "$sample", Value: bson.D{{Key: "size", Value: size}}}}
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

// --- $setWindowFields ---

// WindowBound is a single boundary of a window in a $setWindowFields output
// field. Use the WindowUnbounded / WindowCurrent vars, or WindowOffset.
type WindowBound struct{ v any }

var (
	// WindowUnbounded extends the window to the partition boundary.
	WindowUnbounded = WindowBound{v: "unbounded"}
	// WindowCurrent bounds the window at the current document.
	WindowCurrent = WindowBound{v: "current"}
)

// WindowOffset bounds a window at n positions (for a documents window) or n
// units of the sortBy value (for a range window) relative to the current
// document. Negative values look behind, positive ahead.
func WindowOffset[T Number](n T) WindowBound { return WindowBound{v: n} }

// WindowField pairs an output field name with a window accumulator and its
// optional window boundaries. Construct via WindowOutput.
type WindowField interface{ windowField() windowField }

type windowField struct {
	name   string
	acc    Accumulator
	window bson.D
}

func (wf windowField) windowField() windowField { return wf }

type windowOptions struct {
	documents *[2]any
	rangeVals *[2]any
	unit      any
}

// WithWindowDocuments bounds the window by a fixed number of documents relative
// to the current document. Mutually exclusive with WithWindowRange.
func WithWindowDocuments(lower, upper WindowBound) Option[windowOptions] {
	return func(o *windowOptions) {
		o.documents = &[2]any{lower.v, upper.v}
	}
}

// WithWindowRange bounds the window by the range of the sortBy field value
// relative to the current document. Combine with WithWindowRangeUnit for a time
// range. Mutually exclusive with WithWindowDocuments.
func WithWindowRange(lower, upper WindowBound) Option[windowOptions] {
	return func(o *windowOptions) {
		o.rangeVals = &[2]any{lower.v, upper.v}
	}
}

// WithWindowRangeUnit sets the time unit (e.g. "month") for a range window
// whose sortBy field is a date. Use together with WithWindowRange.
func WithWindowRangeUnit(unit string) Option[windowOptions] {
	return func(o *windowOptions) {
		o.unit = unit
	}
}

// WindowOutput creates a WindowField that computes acc (an accumulator or
// window operator) over the window optionally bounded by WithWindowDocuments or
// WithWindowRange, storing the result in the named field.
func WindowOutput(name string, acc Accumulator, opts ...Option[windowOptions]) WindowField {
	var o windowOptions
	for _, opt := range opts {
		opt(&o)
	}
	var window bson.D
	switch {
	case o.documents != nil:
		window = bson.D{{Key: "documents", Value: bson.A{o.documents[0], o.documents[1]}}}
	case o.rangeVals != nil:
		window = bson.D{{Key: "range", Value: bson.A{o.rangeVals[0], o.rangeVals[1]}}}
		if o.unit != nil {
			window = append(window, bson.E{Key: "unit", Value: o.unit})
		}
	}
	return windowField{name: name, acc: acc, window: window}
}

type setWindowFieldsOptions struct {
	partitionBy any
	sortBy      []SortField
}

// WithSetWindowFieldsPartitionBy groups documents into partitions by the given
// expression. Without it, the whole collection is one partition.
func WithSetWindowFieldsPartitionBy(expr Expr) Option[setWindowFieldsOptions] {
	return func(o *setWindowFieldsOptions) { o.partitionBy = expr }
}

// WithSetWindowFieldsSortBy sorts documents within each partition (construct
// via Sort). It is required by many window operators (e.g. those using a
// range window), though the stage itself does not mandate it.
func WithSetWindowFieldsSortBy(fields ...SortField) Option[setWindowFieldsOptions] {
	return func(o *setWindowFieldsOptions) { o.sortBy = fields }
}

// SetWindowFieldsStage produces a $setWindowFields stage that computes the given
// output fields (construct via WindowOutput) over windows of documents within
// each partition.
func SetWindowFieldsStage(output []WindowField, opts ...Option[setWindowFieldsOptions]) Stage {
	var o setWindowFieldsOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{}
	if o.partitionBy != nil {
		doc = append(doc, bson.E{Key: "partitionBy", Value: o.partitionBy})
	}
	if o.sortBy != nil {
		doc = append(doc, bson.E{Key: "sortBy", Value: sortFieldsDoc(o.sortBy)})
	}
	outDoc := make(bson.D, len(output))
	for i, f := range output {
		wf := f.windowField()
		// The accumulator marshals to a single-key document, e.g. {$sum: "$x"};
		// its key/value are merged with an optional window sub-document.
		accBytes, err := wf.acc.MarshalBSON()
		if err != nil {
			panic("agg: marshal window accumulator: " + err.Error())
		}
		var accDoc bson.D
		if err := bson.Unmarshal(accBytes, &accDoc); err != nil {
			panic("agg: unmarshal window accumulator: " + err.Error())
		}
		if wf.window != nil {
			accDoc = append(accDoc, bson.E{Key: "window", Value: wf.window})
		}
		outDoc[i] = bson.E{Key: wf.name, Value: accDoc}
	}
	doc = append(doc, bson.E{Key: "output", Value: outDoc})
	return Stage{{Key: "$setWindowFields", Value: doc}}
}

// --- $shardedDataDistribution ---

// ShardedDataDistributionStage produces a $shardedDataDistribution stage that
// returns data and size distribution information for sharded collections.
func ShardedDataDistributionStage() Stage {
	return Stage{{Key: "$shardedDataDistribution", Value: bson.D{}}}
}

// --- $skip ---

// SkipStage produces a $skip stage that discards the first skip documents and
// passes the remainder to the next stage.
func SkipStage(skip int32) Stage {
	return Stage{{Key: "$skip", Value: skip}}
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

// --- $sortByCount ---

// SortByCountStage produces a $sortByCount stage that groups documents by expr,
// counts each group, and sorts the groups by descending count.
func SortByCountStage(expr Expr) Stage {
	return Stage{{Key: "$sortByCount", Value: expr}}
}

// --- $unionWith ---

type unionWithOptions struct {
	pipeline []Stage
}

// WithUnionWithPipeline sets a pipeline to apply to the unioned collection
// before combining its results. It cannot contain $out or $merge.
func WithUnionWithPipeline(stages ...Stage) Option[unionWithOptions] {
	return func(o *unionWithOptions) { o.pipeline = stages }
}

// UnionWithStage produces a $unionWith stage that combines the results of the
// current pipeline with those of coll (optionally transformed by
// WithUnionWithPipeline). The builder always emits the verbose (object) form.
func UnionWithStage(coll string, opts ...Option[unionWithOptions]) Stage {
	var o unionWithOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{{Key: "coll", Value: coll}}
	if o.pipeline != nil {
		doc = append(doc, bson.E{Key: "pipeline", Value: pipelineDocs(o.pipeline)})
	}
	return Stage{{Key: "$unionWith", Value: doc}}
}

// --- $unset ---

// UnsetStage produces an $unset stage that removes the named fields from each
// document. It is an alias for a $project stage that excludes fields. Fields
// may use dot notation to remove embedded fields. Even for a single field, the
// builder emits the array form.
func UnsetStage(fields ...string) Stage {
	return Stage{{Key: "$unset", Value: fields}}
}

// --- $unwind ---

type unwindOptions struct {
	includeArrayIndex          any
	preserveNullAndEmptyArrays any
}

// WithUnwindIncludeArrayIndex adds a field of the given name to each output
// document holding the array index of the unwound element.
func WithUnwindIncludeArrayIndex(field string) Option[unwindOptions] {
	return func(o *unwindOptions) { o.includeArrayIndex = field }
}

// WithUnwindPreserveNullAndEmptyArrays outputs a document even when path is
// null, missing, or an empty array. Defaults to false when omitted.
func WithUnwindPreserveNullAndEmptyArrays(v bool) Option[unwindOptions] {
	return func(o *unwindOptions) { o.preserveNullAndEmptyArrays = v }
}

// UnwindStage produces an $unwind stage that deconstructs the array at path
// into one output document per element. path is a field path and must be
// prefixed with a dollar sign, e.g. "$sizes". The builder always emits the
// verbose (object) form.
func UnwindStage(path string, opts ...Option[unwindOptions]) Stage {
	var o unwindOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{{Key: "path", Value: path}}
	if o.includeArrayIndex != nil {
		doc = append(doc, bson.E{Key: "includeArrayIndex", Value: o.includeArrayIndex})
	}
	if o.preserveNullAndEmptyArrays != nil {
		doc = append(doc, bson.E{Key: "preserveNullAndEmptyArrays", Value: o.preserveNullAndEmptyArrays})
	}
	return Stage{{Key: "$unwind", Value: doc}}
}

// --- Atlas Search stages (deferred) ---

// The Atlas Search stage family is intentionally not implemented yet. These
// stages are Atlas-only and carry large, rapidly-evolving option surfaces:
//
//   - $search and $searchMeta take a searchOperator, an entire query DSL
//     (compound, text, near, range, autocomplete, ...) that is not yet modeled
//     anywhere in this module.
//   - $vectorSearch, $rankFusion, $scoreFusion, $rerank, and $score have more
//     bounded arguments but still depend on Atlas-specific concepts.
//
// TODO: implement SearchStage, SearchMetaStage, VectorSearchStage,
// RankFusionStage, ScoreFusionStage, RerankStage, and ScoreStage (and the
// underlying searchOperator DSL) as a separate effort.
