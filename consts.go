package dotpip

// Status represents the response status for commands.
// Status represents the response status for commands.
type Status string

const (
	// StatusOK represents OK.
	// StatusOK represents a successful OK response.
	StatusOK Status = "OK"
)

// ObjectType represents the type of a value.
type ObjectType string

const (
	// ObjectTypeNone represents no object type.
	ObjectTypeNone    ObjectType = "none"
	// ObjectTypeString represents a string object type.
	// ObjectTypeString represents a string object type.
	ObjectTypeString  ObjectType = "string"
	// ObjectTypeHash represents a hash object type.
	ObjectTypeHash    ObjectType = "hash"
	// ObjectTypeList represents a list object type.
	ObjectTypeList    ObjectType = "list"
	// ObjectTypeSet represents a set object type.
	ObjectTypeSet     ObjectType = "set"
	// ObjectTypeZSet represents a zset object type.
	ObjectTypeZSet    ObjectType = "zset"
	// ObjectTypeStream represents a stream object type.
	ObjectTypeStream  ObjectType = "stream"
	// ObjectTypeUnknown represents an unknown object type.
	ObjectTypeUnknown ObjectType = "unknown"
)

// ObjectEncoding represents the encoding of a value.
type ObjectEncoding string

const (
	// ObjectEncodingJSON represents JSON encoding.
	// ObjectEncodingJSON represents JSON encoding.
	ObjectEncodingJSON ObjectEncoding = "json"
	// ObjectEncodingYAML represents YAML encoding.
	// ObjectEncodingYAML represents YAML encoding.
	ObjectEncodingYAML ObjectEncoding = "yaml"
	// ObjectEncodingTOML represents TOML encoding.
	ObjectEncodingTOML ObjectEncoding = "toml"
	// ObjectEncodingRAW represents RAW encoding.
	ObjectEncodingRAW  ObjectEncoding = "raw"
)

// GraphKeyword represents keywords for Graph operations.
type GraphKeyword string

const (
	// GraphKeywordCreate represents the CREATE keyword.
	// GraphKeywordCreate represents the CREATE keyword.
	GraphKeywordCreate               GraphKeyword = "CREATE"
	// GraphKeywordMatch represents the MATCH keyword.
	// GraphKeywordMatch represents the MATCH keyword.
	GraphKeywordMatch                GraphKeyword = "MATCH"
	// GraphKeywordReturn represents the RETURN keyword.
	GraphKeywordReturn               GraphKeyword = "RETURN"
	// GraphKeywordDelete represents the DELETE keyword.
	GraphKeywordDelete               GraphKeyword = "DELETE"
	// GraphKeywordSet represents the SET keyword.
	GraphKeywordSet                  GraphKeyword = "SET"
	// GraphKeywordNodesFound represents the NodesFound keyword.
	GraphKeywordNodesFound           GraphKeyword = "NodesFound"
	// GraphKeywordPathsMatched represents the PathsMatched keyword.
	GraphKeywordPathsMatched         GraphKeyword = "PathsMatched"
	// GraphKeywordNodesCalculated represents the NodesCalculated keyword.
	GraphKeywordNodesCalculated      GraphKeyword = "NodesCalculated"
	// GraphKeywordEdgesCalculated represents the EdgesCalculated keyword.
	GraphKeywordEdgesCalculated      GraphKeyword = "EdgesCalculated"
	// GraphKeywordLabelsAdded represents the LabelsAdded keyword.
	GraphKeywordLabelsAdded          GraphKeyword = "LabelsAdded"
	// GraphKeywordNodesCreated represents the NodesCreated keyword.
	GraphKeywordNodesCreated         GraphKeyword = "NodesCreated"
	// GraphKeywordPropertiesSet represents the PropertiesSet keyword.
	GraphKeywordPropertiesSet        GraphKeyword = "PropertiesSet"
	// GraphKeywordRelationshipsCreated represents the RelationshipsCreated keyword.
	GraphKeywordRelationshipsCreated GraphKeyword = "RelationshipsCreated"
	// GraphKeywordNodesDeleted represents the NodesDeleted keyword.
	GraphKeywordNodesDeleted         GraphKeyword = "NodesDeleted"
)

// ErrorMessage represents standard error messages.
type ErrorMessage string

const (
	// ErrMsgBusyKey represents target key already exists error.
	// ErrMsgBusyKey represents an error when target key already exists.
	ErrMsgBusyKey                 ErrorMessage = "BUSYKEY Target key name already exists"
	// ErrMsgWrongType represents wrong kind of value error.
	// ErrMsgWrongType represents a WRONGTYPE error.
	ErrMsgWrongType               ErrorMessage = "WRONGTYPE Operation against a key holding the wrong kind of value"
	// ErrMsgLFUEviction represents LFU eviction not supported error.
	ErrMsgLFUEviction             ErrorMessage = "ERR LFU eviction not supported"
	// ErrMsgMigrateNotSupported represents MIGRATE not supported error.
	ErrMsgMigrateNotSupported     ErrorMessage = "MIGRATE is not supported in fs mode over network"
	// ErrMsgValueNotInt represents value is not an integer error.
	ErrMsgValueNotInt             ErrorMessage = "ERR value is not an integer or out of range"
	// ErrMsgValueNotFloat represents value is not a valid float error.
	ErrMsgValueNotFloat           ErrorMessage = "ERR value is not a valid float"
	// ErrMsgOffsetOutOfRange represents offset out of range error.
	ErrMsgOffsetOutOfRange        ErrorMessage = "ERR offset is out of range"
	// ErrMsgRankZero represents RANK cannot be zero error.
	ErrMsgRankZero                ErrorMessage = "ERR RANK can't be zero"
	// ErrMsgNoSuchKey represents no such key error.
	ErrMsgNoSuchKey               ErrorMessage = "ERR no such key"
	// ErrMsgIndexOutOfRange represents index out of range error.
	ErrMsgIndexOutOfRange         ErrorMessage = "ERR index out of range"
	// ErrMsgIndexOutOfBounds represents index out of bounds error.
	ErrMsgIndexOutOfBounds        ErrorMessage = "ERR index out of bounds"
	// ErrMsgNewObjectsRoot represents new objects root error.
	ErrMsgNewObjectsRoot          ErrorMessage = "ERR new objects must be created at the root"
	// ErrMsgInvalidLonLat represents invalid longitude latitude pair error.
	ErrMsgInvalidLonLat           ErrorMessage = "ERR invalid longitude,latitude pair %f,%f"
	// ErrMsgMembersNotExist represents members do not exist error.
	ErrMsgMembersNotExist         ErrorMessage = "ERR one or both members do not exist"
	// ErrMsgUnsupportedUnit represents unsupported unit error.
	ErrMsgUnsupportedUnit         ErrorMessage = "ERR unsupported unit provided"
	// ErrMsgCouldNotDecodeZSet represents could not decode zset member error.
	ErrMsgCouldNotDecodeZSet      ErrorMessage = "ERR could not decode requested zset member"
	// ErrMsgGeoSearchFrom represents geo search from error.
	ErrMsgGeoSearchFrom           ErrorMessage = "ERR either FROMMEMBER or FROMLONLAT must be provided"
	// ErrMsgGeoSearchBy indicates an error when GeoSearch is called without BYRADIUS or BYBOX
	ErrMsgGeoSearchBy             ErrorMessage = "ERR either BYRADIUS or BYBOX must be provided"
	// ErrMsgInvalidStreamID indicates an error when a stream ID is invalid
	ErrMsgInvalidStreamID         ErrorMessage = "ERR Invalid stream ID specified as string"
	// ErrMsgXAddIDEqualSmaller indicates an error when an XADD ID is equal or smaller
	ErrMsgXAddIDEqualSmaller      ErrorMessage = "ERR The ID specified in XADD is equal or smaller than the target stream top item"
	// ErrMsgXAddIDGreaterZero indicates an error when an XADD ID is not greater than 0-0
	ErrMsgXAddIDGreaterZero       ErrorMessage = "ERR The ID specified in XADD must be greater than 0-0"
	// ErrMsgXGroupKeyExists indicates an error when XGROUP key does not exist
	ErrMsgXGroupKeyExists         ErrorMessage = "ERR The XGROUP subcommand requires the key to exist"
	// ErrMsgBusyGroup indicates an error when a consumer group already exists
	ErrMsgBusyGroup               ErrorMessage = "BUSYGROUP Consumer Group name already exists"
	// ErrMsgNoGroup indicates an error when a consumer group or key does not exist
	ErrMsgNoGroup                 ErrorMessage = "ERR NOGROUP No such key '%s' or consumer group '%s'"
	// ErrMsgUnbalancedXRead indicates an error when XREAD stream and ID lists are unbalanced
	ErrMsgUnbalancedXRead         ErrorMessage = "ERR Unbalanced XREAD list of streams and IDs"
	// ErrMsgFailedToEncodeArray indicates an error when encoding an array fails
	ErrMsgFailedToEncodeArray     ErrorMessage = "failed to encode array"
	// ErrMsgUnknownOperation indicates an error for an unknown operation
	ErrMsgUnknownOperation        ErrorMessage = "unknown operation: %s"
	// ErrMsgInvalidType indicates an error for an invalid type
	ErrMsgInvalidType             ErrorMessage = "invalid type"
	// ErrMsgInvalidTypeFormat indicates an error.
	ErrMsgInvalidTypeFormat       ErrorMessage = "invalid type format"
	// ErrMsgInvalidBits indicates an error.
	ErrMsgInvalidBits             ErrorMessage = "invalid bits"
	// ErrMsgInvalidArgumentType indicates an error.
	ErrMsgInvalidArgumentType     ErrorMessage = "invalid argument type"
	// ErrMsgSyntaxError indicates an error.
	ErrMsgSyntaxError             ErrorMessage = "syntax error"
	// ErrMsgInvalidOverflowArg indicates an error.
	ErrMsgInvalidOverflowArg      ErrorMessage = "invalid overflow argument"
	// ErrMsgInvalidOverflowType indicates an error.
	ErrMsgInvalidOverflowType     ErrorMessage = "invalid overflow type"
	// ErrMsgUnknownSubcommand indicates an error.
	ErrMsgUnknownSubcommand       ErrorMessage = "unknown subcommand"
	// ErrMsgUnsupportedEncodingType indicates an error.
	ErrMsgUnsupportedEncodingType ErrorMessage = "unsupported encoding type: %s"
	// ErrMsgRAWStringDecodeExpected indicates an error.
	ErrMsgRAWStringDecodeExpected ErrorMessage = "RAW stringDecode expected []byte or string, got %T"
	// ErrMsgGeospatialDecoderNot indicates an error.
	ErrMsgGeospatialDecoderNot    ErrorMessage = "geospatial decoder not configured"
	// ErrMsgGeospatialEncoderNot indicates an error.
	ErrMsgGeospatialEncoderNot    ErrorMessage = "geospatial encoder not configured"

	// ErrMsgReadOnlyQuery indicates an error.
	ErrMsgReadOnlyQuery ErrorMessage = "read-only query contains write operations"
)
const (
	// ErrWrongTypeVectorSet represents wrong kind of value error for vector set.
	ErrWrongTypeVectorSet ErrorMessage = "WRONGTYPE Operation against a key holding the wrong kind of value"
	// ObjectVectorSet represents a vector set object type.
	ObjectVectorSet ObjectType = "vector_set"
)
