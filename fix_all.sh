#!/bin/bash
sed -i 's/	BitOpOr  BitOp = "OR"/\t\/\/ BitOpOr represents the OR operation.\n\tBitOpOr  BitOp = "OR"/' commands.go
sed -i 's/	BitOpXor BitOp = "XOR"/\t\/\/ BitOpXor represents the XOR operation.\n\tBitOpXor BitOp = "XOR"/' commands.go
sed -i 's/	BitOpNot BitOp = "NOT"/\t\/\/ BitOpNot represents the NOT operation.\n\tBitOpNot BitOp = "NOT"/' commands.go

sed -i 's/	ObjectTypeString  ObjectType = "string"/\t\/\/ ObjectTypeString represents a string object type.\n\tObjectTypeString  ObjectType = "string"/' consts.go
sed -i 's/	ObjectTypeList    ObjectType = "list"/\t\/\/ ObjectTypeList represents a list object type.\n\tObjectTypeList    ObjectType = "list"/' consts.go
sed -i 's/	ObjectTypeSet     ObjectType = "set"/\t\/\/ ObjectTypeSet represents a set object type.\n\tObjectTypeSet     ObjectType = "set"/' consts.go
sed -i 's/	ObjectTypeZSet    ObjectType = "zset"/\t\/\/ ObjectTypeZSet represents a zset object type.\n\tObjectTypeZSet    ObjectType = "zset"/' consts.go
sed -i 's/	ObjectTypeHash    ObjectType = "hash"/\t\/\/ ObjectTypeHash represents a hash object type.\n\tObjectTypeHash    ObjectType = "hash"/' consts.go
sed -i 's/	ObjectTypeStream  ObjectType = "stream"/\t\/\/ ObjectTypeStream represents a stream object type.\n\tObjectTypeStream  ObjectType = "stream"/' consts.go
sed -i 's/	ObjectTypeUnknown ObjectType = "unknown"/\t\/\/ ObjectTypeUnknown represents an unknown object type.\n\tObjectTypeUnknown ObjectType = "unknown"/' consts.go
sed -i 's/	ObjectVectorSet ObjectType = "vector_set"/\t\/\/ ObjectVectorSet represents a vector set object type.\n\tObjectVectorSet ObjectType = "vector_set"/' consts.go

sed -i 's/	ObjectEncodingJSON ObjectEncoding = "json"/\t\/\/ ObjectEncodingJSON represents JSON encoding.\n\tObjectEncodingJSON ObjectEncoding = "json"/' consts.go
sed -i 's/	ObjectEncodingYAML ObjectEncoding = "yaml"/\t\/\/ ObjectEncodingYAML represents YAML encoding.\n\tObjectEncodingYAML ObjectEncoding = "yaml"/' consts.go
sed -i 's/	ObjectEncodingTOML ObjectEncoding = "toml"/\t\/\/ ObjectEncodingTOML represents TOML encoding.\n\tObjectEncodingTOML ObjectEncoding = "toml"/' consts.go
sed -i 's/	ObjectEncodingRAW  ObjectEncoding = "raw"/\t\/\/ ObjectEncodingRAW represents RAW encoding.\n\tObjectEncodingRAW  ObjectEncoding = "raw"/' consts.go

sed -i 's/	GraphKeywordCreate               GraphKeyword = "CREATE"/\t\/\/ GraphKeywordCreate represents the CREATE keyword.\n\tGraphKeywordCreate               GraphKeyword = "CREATE"/' consts.go
sed -i 's/	GraphKeywordMatch                GraphKeyword = "MATCH"/\t\/\/ GraphKeywordMatch represents the MATCH keyword.\n\tGraphKeywordMatch                GraphKeyword = "MATCH"/' consts.go
sed -i 's/	GraphKeywordReturn               GraphKeyword = "RETURN"/\t\/\/ GraphKeywordReturn represents the RETURN keyword.\n\tGraphKeywordReturn               GraphKeyword = "RETURN"/' consts.go
sed -i 's/	GraphKeywordDelete               GraphKeyword = "DELETE"/\t\/\/ GraphKeywordDelete represents the DELETE keyword.\n\tGraphKeywordDelete               GraphKeyword = "DELETE"/' consts.go
sed -i 's/	GraphKeywordSet                  GraphKeyword = "SET"/\t\/\/ GraphKeywordSet represents the SET keyword.\n\tGraphKeywordSet                  GraphKeyword = "SET"/' consts.go
sed -i 's/	GraphKeywordNodesFound           GraphKeyword = "NodesFound"/\t\/\/ GraphKeywordNodesFound represents the NodesFound keyword.\n\tGraphKeywordNodesFound           GraphKeyword = "NodesFound"/' consts.go
sed -i 's/	GraphKeywordPathsMatched         GraphKeyword = "PathsMatched"/\t\/\/ GraphKeywordPathsMatched represents the PathsMatched keyword.\n\tGraphKeywordPathsMatched         GraphKeyword = "PathsMatched"/' consts.go
sed -i 's/	GraphKeywordNodesCalculated      GraphKeyword = "NodesCalculated"/\t\/\/ GraphKeywordNodesCalculated represents the NodesCalculated keyword.\n\tGraphKeywordNodesCalculated      GraphKeyword = "NodesCalculated"/' consts.go
sed -i 's/	GraphKeywordEdgesCalculated      GraphKeyword = "EdgesCalculated"/\t\/\/ GraphKeywordEdgesCalculated represents the EdgesCalculated keyword.\n\tGraphKeywordEdgesCalculated      GraphKeyword = "EdgesCalculated"/' consts.go
sed -i 's/	GraphKeywordNodesCreated         GraphKeyword = "NodesCreated"/\t\/\/ GraphKeywordNodesCreated represents the NodesCreated keyword.\n\tGraphKeywordNodesCreated         GraphKeyword = "NodesCreated"/' consts.go
sed -i 's/	GraphKeywordPropertiesSet        GraphKeyword = "PropertiesSet"/\t\/\/ GraphKeywordPropertiesSet represents the PropertiesSet keyword.\n\tGraphKeywordPropertiesSet        GraphKeyword = "PropertiesSet"/' consts.go
sed -i 's/	GraphKeywordRelationshipsCreated GraphKeyword = "RelationshipsCreated"/\t\/\/ GraphKeywordRelationshipsCreated represents the RelationshipsCreated keyword.\n\tGraphKeywordRelationshipsCreated GraphKeyword = "RelationshipsCreated"/' consts.go
sed -i 's/	GraphKeywordNodesDeleted         GraphKeyword = "NodesDeleted"/\t\/\/ GraphKeywordNodesDeleted represents the NodesDeleted keyword.\n\tGraphKeywordNodesDeleted         GraphKeyword = "NodesDeleted"/' consts.go
sed -i 's/	GraphKeywordLabelsAdded          GraphKeyword = "LabelsAdded"/\t\/\/ GraphKeywordLabelsAdded represents the LabelsAdded keyword.\n\tGraphKeywordLabelsAdded          GraphKeyword = "LabelsAdded"/' consts.go

sed -i 's/	ErrMsgBusyKey                 ErrorMessage = "BUSYKEY Target key name already exists"/\t\/\/ ErrMsgBusyKey represents target key already exists error.\n\tErrMsgBusyKey                 ErrorMessage = "BUSYKEY Target key name already exists"/' consts.go
sed -i 's/	ErrMsgWrongType               ErrorMessage = "WRONGTYPE Operation against a key holding the wrong kind of value"/\t\/\/ ErrMsgWrongType represents wrong kind of value error.\n\tErrMsgWrongType               ErrorMessage = "WRONGTYPE Operation against a key holding the wrong kind of value"/' consts.go
sed -i 's/	ErrMsgLFUEviction             ErrorMessage = "ERR LFU eviction not supported"/\t\/\/ ErrMsgLFUEviction represents LFU eviction not supported error.\n\tErrMsgLFUEviction             ErrorMessage = "ERR LFU eviction not supported"/' consts.go
sed -i 's/	ErrMsgMigrateNotSupported     ErrorMessage = "MIGRATE is not supported in fs mode over network"/\t\/\/ ErrMsgMigrateNotSupported represents MIGRATE not supported error.\n\tErrMsgMigrateNotSupported     ErrorMessage = "MIGRATE is not supported in fs mode over network"/' consts.go
sed -i 's/	ErrMsgValueNotInt             ErrorMessage = "ERR value is not an integer or out of range"/\t\/\/ ErrMsgValueNotInt represents value is not an integer error.\n\tErrMsgValueNotInt             ErrorMessage = "ERR value is not an integer or out of range"/' consts.go
sed -i 's/	ErrMsgValueNotFloat           ErrorMessage = "ERR value is not a valid float"/\t\/\/ ErrMsgValueNotFloat represents value is not a valid float error.\n\tErrMsgValueNotFloat           ErrorMessage = "ERR value is not a valid float"/' consts.go
sed -i 's/	ErrMsgOffsetOutOfRange        ErrorMessage = "ERR offset is out of range"/\t\/\/ ErrMsgOffsetOutOfRange represents offset out of range error.\n\tErrMsgOffsetOutOfRange        ErrorMessage = "ERR offset is out of range"/' consts.go
sed -i 's/	ErrMsgRankZero                ErrorMessage = "ERR RANK can'\''t be zero"/\t\/\/ ErrMsgRankZero represents RANK cannot be zero error.\n\tErrMsgRankZero                ErrorMessage = "ERR RANK can'\''t be zero"/' consts.go
sed -i 's/	ErrMsgNoSuchKey               ErrorMessage = "ERR no such key"/\t\/\/ ErrMsgNoSuchKey represents no such key error.\n\tErrMsgNoSuchKey               ErrorMessage = "ERR no such key"/' consts.go
sed -i 's/	ErrWrongTypeVectorSet ErrorMessage = "WRONGTYPE Operation against a key holding the wrong kind of value"/\t\/\/ ErrWrongTypeVectorSet represents wrong kind of value error for vector set.\n\tErrWrongTypeVectorSet ErrorMessage = "WRONGTYPE Operation against a key holding the wrong kind of value"/' consts.go
sed -i 's/	ErrMsgIndexOutOfRange         ErrorMessage = "ERR index out of range"/\t\/\/ ErrMsgIndexOutOfRange represents index out of range error.\n\tErrMsgIndexOutOfRange         ErrorMessage = "ERR index out of range"/' consts.go
sed -i 's/	ErrMsgIndexOutOfBounds        ErrorMessage = "ERR index out of bounds"/\t\/\/ ErrMsgIndexOutOfBounds represents index out of bounds error.\n\tErrMsgIndexOutOfBounds        ErrorMessage = "ERR index out of bounds"/' consts.go
sed -i 's/	ErrMsgNewObjectsRoot          ErrorMessage = "ERR new objects must be created at the root"/\t\/\/ ErrMsgNewObjectsRoot represents new objects root error.\n\tErrMsgNewObjectsRoot          ErrorMessage = "ERR new objects must be created at the root"/' consts.go
sed -i 's/	ErrMsgInvalidLonLat           ErrorMessage = "ERR invalid longitude,latitude pair %f,%f"/\t\/\/ ErrMsgInvalidLonLat represents invalid longitude latitude pair error.\n\tErrMsgInvalidLonLat           ErrorMessage = "ERR invalid longitude,latitude pair %f,%f"/' consts.go
sed -i 's/	ErrMsgMembersNotExist         ErrorMessage = "ERR one or both members do not exist"/\t\/\/ ErrMsgMembersNotExist represents members do not exist error.\n\tErrMsgMembersNotExist         ErrorMessage = "ERR one or both members do not exist"/' consts.go
sed -i 's/	ErrMsgUnsupportedUnit         ErrorMessage = "ERR unsupported unit provided"/\t\/\/ ErrMsgUnsupportedUnit represents unsupported unit error.\n\tErrMsgUnsupportedUnit         ErrorMessage = "ERR unsupported unit provided"/' consts.go
sed -i 's/	ErrMsgCouldNotDecodeZSet      ErrorMessage = "ERR could not decode requested zset member"/\t\/\/ ErrMsgCouldNotDecodeZSet represents could not decode zset member error.\n\tErrMsgCouldNotDecodeZSet      ErrorMessage = "ERR could not decode requested zset member"/' consts.go
sed -i 's/	ErrMsgGeoSearchFrom           ErrorMessage = "ERR either FROMMEMBER or FROMLONLAT must be provided"/\t\/\/ ErrMsgGeoSearchFrom represents geo search from error.\n\tErrMsgGeoSearchFrom           ErrorMessage = "ERR either FROMMEMBER or FROMLONLAT must be provided"/' consts.go

sed -i 's/type FileEncodeType string/\/\/ FileEncodeType represents a file encoding type.\ntype FileEncodeType string/' fs/encoding.go
sed -i 's/	JSON FileEncodeType = "json"/\t\/\/ JSON represents JSON encoding type.\n\tJSON FileEncodeType = "json"/' fs/encoding.go
sed -i 's/	YAML FileEncodeType = "yaml"/\t\/\/ YAML represents YAML encoding type.\n\tYAML FileEncodeType = "yaml"/' fs/encoding.go
sed -i 's/	TOML FileEncodeType = "toml"/\t\/\/ TOML represents TOML encoding type.\n\tTOML FileEncodeType = "toml"/' fs/encoding.go
sed -i 's/	RAW  FileEncodeType = "raw"/\t\/\/ RAW represents RAW encoding type.\n\tRAW  FileEncodeType = "raw"/' fs/encoding.go

sed -i 's/func (f \*FileSystem) EncodeType/\/\/ EncodeType configures the encoding type.\nfunc (f \*FileSystem) EncodeType/' fs/encoding.go
sed -i 's/func (f \*FileSystem) JSONEncode/\/\/ JSONEncode encodes a value to JSON.\nfunc (f \*FileSystem) JSONEncode/' fs/encoding.go
sed -i 's/func (f \*FileSystem) JSONDecode/\/\/ JSONDecode decodes a JSON value.\nfunc (f \*FileSystem) JSONDecode/' fs/encoding.go

sed -i 's/func (f \*FileSystem) Expire(/\/\/ Expire sets the expiration time.\nfunc (f \*FileSystem) Expire(/' fs/expire.go
sed -i 's/func (f \*FileSystem) PExpire(/\/\/ PExpire sets the expiration time in milliseconds.\nfunc (f \*FileSystem) PExpire(/' fs/expire.go
sed -i 's/func (f \*FileSystem) ExpireAt/\/\/ ExpireAt sets the expiration to a specific unix time.\nfunc (f \*FileSystem) ExpireAt/' fs/expire.go
sed -i 's/func (f \*FileSystem) PExpireAt/\/\/ PExpireAt sets the expiration to a specific unix time in milliseconds.\nfunc (f \*FileSystem) PExpireAt/' fs/expire.go
sed -i 's/func (f \*FileSystem) ExpireTime/\/\/ ExpireTime returns the absolute Unix timestamp in seconds at which the given key will expire.\nfunc (f \*FileSystem) ExpireTime/' fs/expire.go
sed -i 's/func (f \*FileSystem) PExpireTime/\/\/ PExpireTime returns the absolute Unix timestamp in milliseconds at which the given key will expire.\nfunc (f \*FileSystem) PExpireTime/' fs/expire.go
sed -i 's/func (f \*FileSystem) TTL/\/\/ TTL returns the remaining time to live of a key that has a timeout.\nfunc (f \*FileSystem) TTL/' fs/expire.go
sed -i 's/func (f \*FileSystem) PTTL/\/\/ PTTL returns the remaining time to live of a key that has a timeout, in milliseconds.\nfunc (f \*FileSystem) PTTL/' fs/expire.go
sed -i 's/func (f \*FileSystem) Persist/\/\/ Persist removes the expiration from a key.\nfunc (f \*FileSystem) Persist/' fs/expire.go

sed -i 's/type FileSystem struct {/\/\/ FileSystem represents a file-system backed DotPip implementation.\ntype FileSystem struct {/' fs/fs.go
sed -i 's/func NewFileSystem(pathRoot string) \*FileSystem {/\/\/ NewFileSystem creates a new FileSystem.\nfunc NewFileSystem(pathRoot string) \*FileSystem {/' fs/fs.go
sed -i 's/func (f \*FileSystem) Close() {/\/\/ Close closes the FileSystem and cleans up resources.\nfunc (f \*FileSystem) Close() {/' fs/fs.go
sed -i 's/func (f \*FileSystem) ConfigSet(parameter string, value string) error {/\/\/ ConfigSet sets a configuration parameter.\nfunc (f \*FileSystem) ConfigSet(parameter string, value string) error {/' fs/fs.go
sed -i 's/func (f \*FileSystem) ConfigGet(parameter string) (map\[string\]string, error) {/\/\/ ConfigGet gets a configuration parameter.\nfunc (f \*FileSystem) ConfigGet(parameter string) (map\[string\]string, error) {/' fs/fs.go
