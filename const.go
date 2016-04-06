package dynago

// CapacityDetail describes what kind of ConsumedCapacity response to get.
type CapacityDetail string

// Define all the ways you can request consumed capacity.
const (
	CapacityIndexes CapacityDetail = "INDEXES"
	CapacityTotal   CapacityDetail = "TOTAL"
	CapacityNone    CapacityDetail = "NONE"
)

// Select describes which attributes to return in Scan and Query operations.
type Select string

// Define all the ways you can ask for attributes
const (
	SelectAllAttributes      Select = "ALL_ATTRIBUTES"
	SelectAllProjected       Select = "ALL_PROJECTED_ATTRIBUTES"
	SelectCount              Select = "COUNT"
	SelectSpecificAttributes Select = "SPECIFIC_ATTRIBUTES"
)

// ReturnValues enable returning some or all changed data from put operations.
type ReturnValues string

// Define all the ways you can ask for return values.
// Note not all values are valid in all contexts.
const (
	ReturnNone       ReturnValues = "NONE"
	ReturnAllOld     ReturnValues = "ALL_OLD"
	ReturnUpdatedOld ReturnValues = "UPDATED_OLD"
	ReturnAllNew     ReturnValues = "ALL_NEW"
	ReturnUpdatedNew ReturnValues = "UPDATED_NEW"
)

// DebugFlags controls Dynago debugging
type DebugFlags uint

// All the available debug bit-flags.
const (
	DebugRequests DebugFlags = 1 << iota
	DebugResponses
	DebugAuth
)

// Time format for low-level storage
const iso8601compact = "2006-01-02T15:04:05Z"
