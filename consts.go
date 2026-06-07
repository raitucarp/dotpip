package dotpip

type Status string

const (
	StatusOK Status = "OK"
)

type ObjectType string

const (
	ObjectTypeNone    ObjectType = "none"
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

type ObjectEncoding string

const (
	// ObjectEncodingJSON represents JSON encoding.
	ObjectEncodingJSON ObjectEncoding = "json"
	// ObjectEncodingYAML represents YAML encoding.
	ObjectEncodingYAML ObjectEncoding = "yaml"
	// ObjectEncodingTOML represents TOML encoding.
	ObjectEncodingTOML ObjectEncoding = "toml"
	// ObjectEncodingRAW represents RAW encoding.
	ObjectEncodingRAW  ObjectEncoding = "raw"
)

type GraphKeyword string

const (
	// GraphKeywordCreate represents the CREATE keyword.
	GraphKeywordCreate               GraphKeyword = "CREATE"
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

type ErrorMessage string

const (
	// ErrMsgBusyKey represents target key already exists error.
	ErrMsgBusyKey                 ErrorMessage = "BUSYKEY Target key name already exists"
	// ErrMsgWrongType represents wrong kind of value error.
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
	ErrMsgGeoSearchBy             ErrorMessage = "ERR either BYRADIUS or BYBOX must be provided"
	ErrMsgInvalidStreamID         ErrorMessage = "ERR Invalid stream ID specified as string"
	ErrMsgXAddIDEqualSmaller      ErrorMessage = "ERR The ID specified in XADD is equal or smaller than the target stream top item"
	ErrMsgXAddIDGreaterZero       ErrorMessage = "ERR The ID specified in XADD must be greater than 0-0"
	ErrMsgXGroupKeyExists         ErrorMessage = "ERR The XGROUP subcommand requires the key to exist"
	ErrMsgBusyGroup               ErrorMessage = "BUSYGROUP Consumer Group name already exists"
	ErrMsgNoGroup                 ErrorMessage = "ERR NOGROUP No such key '%s' or consumer group '%s'"
	ErrMsgUnbalancedXRead         ErrorMessage = "ERR Unbalanced XREAD list of streams and IDs"
	ErrMsgFailedToEncodeArray     ErrorMessage = "failed to encode array"
	ErrMsgUnknownOperation        ErrorMessage = "unknown operation: %s"
	ErrMsgInvalidType             ErrorMessage = "invalid type"
	ErrMsgInvalidTypeFormat       ErrorMessage = "invalid type format"
	ErrMsgInvalidBits             ErrorMessage = "invalid bits"
	ErrMsgInvalidArgumentType     ErrorMessage = "invalid argument type"
	ErrMsgSyntaxError             ErrorMessage = "syntax error"
	ErrMsgInvalidOverflowArg      ErrorMessage = "invalid overflow argument"
	ErrMsgInvalidOverflowType     ErrorMessage = "invalid overflow type"
	ErrMsgUnknownSubcommand       ErrorMessage = "unknown subcommand"
	ErrMsgUnsupportedEncodingType ErrorMessage = "unsupported encoding type: %s"
	ErrMsgRAWStringDecodeExpected ErrorMessage = "RAW stringDecode expected []byte or string, got %T"
	ErrMsgGeospatialDecoderNot    ErrorMessage = "geospatial decoder not configured"
	ErrMsgGeospatialEncoderNot    ErrorMessage = "geospatial encoder not configured"

	ErrMsgReadOnlyQuery ErrorMessage = "read-only query contains write operations"
)
const (
	// ErrWrongTypeVectorSet represents wrong kind of value error for vector set.
	ErrWrongTypeVectorSet ErrorMessage = "WRONGTYPE Operation against a key holding the wrong kind of value"
	// ObjectVectorSet represents a vector set object type.
	ObjectVectorSet ObjectType = "vector_set"
)
