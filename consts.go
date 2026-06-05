package dotpip

type Status string

const (
	StatusOK Status = "OK"
)

type ObjectType string

const (
	ObjectTypeNone    ObjectType = "none"
	ObjectTypeString  ObjectType = "string"
	ObjectTypeHash    ObjectType = "hash"
	ObjectTypeList    ObjectType = "list"
	ObjectTypeSet     ObjectType = "set"
	ObjectTypeZSet    ObjectType = "zset"
	ObjectTypeStream  ObjectType = "stream"
	ObjectTypeUnknown ObjectType = "unknown"
)

type ObjectEncoding string

const (
	ObjectEncodingJSON ObjectEncoding = "json"
	ObjectEncodingYAML ObjectEncoding = "yaml"
	ObjectEncodingTOML ObjectEncoding = "toml"
	ObjectEncodingRAW  ObjectEncoding = "raw"
)

type GraphKeyword string

const (
	GraphKeywordCreate GraphKeyword = "CREATE"
	GraphKeywordMatch GraphKeyword = "MATCH"
	GraphKeywordReturn GraphKeyword = "RETURN"
	GraphKeywordDelete GraphKeyword = "DELETE"
	GraphKeywordSet GraphKeyword = "SET"
	GraphKeywordNodesFound GraphKeyword = "NodesFound"
	GraphKeywordPathsMatched GraphKeyword = "PathsMatched"
	GraphKeywordNodesCalculated GraphKeyword = "NodesCalculated"
	GraphKeywordEdgesCalculated GraphKeyword = "EdgesCalculated"
	GraphKeywordLabelsAdded GraphKeyword = "LabelsAdded"
	GraphKeywordNodesCreated GraphKeyword = "NodesCreated"
	GraphKeywordPropertiesSet GraphKeyword = "PropertiesSet"
	GraphKeywordRelationshipsCreated GraphKeyword = "RelationshipsCreated"
	GraphKeywordNodesDeleted GraphKeyword = "NodesDeleted"
)

type ErrorMessage string

const (
	ErrMsgBusyKey                 ErrorMessage = "BUSYKEY Target key name already exists"
	ErrMsgWrongType               ErrorMessage = "WRONGTYPE Operation against a key holding the wrong kind of value"
	ErrMsgLFUEviction             ErrorMessage = "ERR LFU eviction not supported"
	ErrMsgMigrateNotSupported     ErrorMessage = "MIGRATE is not supported in fs mode over network"
	ErrMsgValueNotInt             ErrorMessage = "ERR value is not an integer or out of range"
	ErrMsgValueNotFloat           ErrorMessage = "ERR value is not a valid float"
	ErrMsgOffsetOutOfRange        ErrorMessage = "ERR offset is out of range"
	ErrMsgRankZero                ErrorMessage = "ERR RANK can't be zero"
	ErrMsgNoSuchKey               ErrorMessage = "ERR no such key"
	ErrMsgIndexOutOfRange         ErrorMessage = "ERR index out of range"
	ErrMsgIndexOutOfBounds        ErrorMessage = "ERR index out of bounds"
	ErrMsgNewObjectsRoot          ErrorMessage = "ERR new objects must be created at the root"
	ErrMsgInvalidLonLat           ErrorMessage = "ERR invalid longitude,latitude pair %f,%f"
	ErrMsgMembersNotExist         ErrorMessage = "ERR one or both members do not exist"
	ErrMsgUnsupportedUnit         ErrorMessage = "ERR unsupported unit provided"
	ErrMsgCouldNotDecodeZSet      ErrorMessage = "ERR could not decode requested zset member"
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

	ErrMsgReadOnlyQuery           ErrorMessage = "read-only query contains write operations"
)
