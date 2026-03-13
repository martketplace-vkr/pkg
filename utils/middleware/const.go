package middleware

const (

	// FaaSInvokedNameKey is the attribute Key conforming to the
	// "faas.invoked_name" semantic conventions. It represents the name of the
	// invoked function.
	//
	// Type: string
	// RequirementLevel: Required
	// Stability: experimental
	// Examples: 'my-function'
	// Note: SHOULD be equal to the `faas.name` resource attribute of the
	// invoked function.
	FaaSInvokedNameKey = "faas.invoked_name"

	// FaaSInvokedProviderKey is the attribute Key conforming to the
	// "faas.invoked_provider" semantic conventions. It represents the cloud
	// provider of the invoked function.
	//
	// Type: Enum
	// RequirementLevel: Required
	// Stability: experimental
	// Note: SHOULD be equal to the `cloud.provider` resource attribute of the
	// invoked function.
	FaaSInvokedProviderKey = "faas.invoked_provider"

	// FaaSInvokedRegionKey is the attribute Key conforming to the
	// "faas.invoked_region" semantic conventions. It represents the cloud
	// region of the invoked function.
	//
	// Type: string
	// RequirementLevel: ConditionallyRequired (For some cloud providers, like
	// AWS or GCP, the region in which a function is hosted is essential to
	// uniquely identify the function and also part of its endpoint. Since it's
	// part of the endpoint being called, the region is always known to
	// clients. In these cases, `faas.invoked_region` MUST be set accordingly.
	// If the region is unknown to the client or not required for identifying
	// the invoked function, setting `faas.invoked_region` is optional.)
	// Stability: experimental
	// Examples: 'eu-central-1'
	// Note: SHOULD be equal to the `cloud.region` resource attribute of the
	// invoked function.
	FaaSInvokedRegionKey = "faas.invoked_region"

	// FaaSTriggerKey is the attribute Key conforming to the "faas.trigger"
	// semantic conventions. It represents the type of the trigger which caused
	// this function invocation.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	FaaSTriggerKey = "faas.trigger"

	// EventNameKey is the attribute Key conforming to the "event.name"
	// semantic conventions. It represents the identifies the class / type of
	// event.
	//
	// Type: string
	// RequirementLevel: Required
	// Stability: experimental
	// Examples: 'browser.mouse.click', 'device.app.lifecycle'
	// Note: Event names are subject to the same rules as [attribute
	// names](https://github.com/open-telemetry/opentelemetry-specification/tree/v1.26.0/specification/common/attribute-naming.md).
	// Notably, event names are namespaced to avoid collisions and provide a
	// clean separation of semantics for events in separate domains like
	// browser, mobile, and kubernetes.
	EventNameKey = "event.name"

	// LogRecordUIDKey is the attribute Key conforming to the "log.record.uid"
	// semantic conventions. It represents a unique identifier for the Log
	// Record.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '01ARZ3NDEKTSV4RRFFQ69G5FAV'
	// Note: If an id is provided, other log records with the same id will be
	// considered duplicates and can be removed safely. This means, that two
	// distinguishable log records MUST have different values.
	// The id MAY be an [Universally Unique Lexicographically Sortable
	// Identifier (ULID)](https://github.com/ulid/spec), but other identifiers
	// (e.g. UUID) may be used as needed.
	LogRecordUIDKey = "log.record.uid"

	// LogIostreamKey is the attribute Key conforming to the "log.iostream"
	// semantic conventions. It represents the stream associated with the log.
	// See below for a list of well-known values.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	LogIostreamKey = "log.iostream"

	// LogFileNameKey is the attribute Key conforming to the "log.file.name"
	// semantic conventions. It represents the basename of the file.
	//
	// Type: string
	// RequirementLevel: Recommended
	// Stability: experimental
	// Examples: 'audit.log'
	LogFileNameKey = "log.file.name"

	// LogFileNameResolvedKey is the attribute Key conforming to the
	// "log.file.name_resolved" semantic conventions. It represents the
	// basename of the file, with symlinks resolved.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'uuid.log'
	LogFileNameResolvedKey = "log.file.name_resolved"

	// LogFilePathKey is the attribute Key conforming to the "log.file.path"
	// semantic conventions. It represents the full path to the file.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '/var/log/mysql/audit.log'
	LogFilePathKey = "log.file.path"

	// LogFilePathResolvedKey is the attribute Key conforming to the
	// "log.file.path_resolved" semantic conventions. It represents the full
	// path to the file, with symlinks resolved.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '/var/lib/docker/uuid.log'
	LogFilePathResolvedKey = "log.file.path_resolved"

	// PoolNameKey is the attribute Key conforming to the "pool.name" semantic
	// conventions. It represents the name of the connection pool; unique
	// within the instrumented application. In case the connection pool
	// implementation doesn't provide a name, then the
	// [db.connection_string](/docs/database/database-spans.md#connection-level-attributes)
	// should be used
	//
	// Type: string
	// RequirementLevel: Required
	// Stability: experimental
	// Examples: 'myDataSource'
	PoolNameKey = "pool.name"

	// StateKey is the attribute Key conforming to the "state" semantic
	// conventions. It represents the state of a connection in the pool
	//
	// Type: Enum
	// RequirementLevel: Required
	// Stability: experimental
	// Examples: 'idle'
	StateKey = "state"

	// AspnetcoreDiagnosticsHandlerTypeKey is the attribute Key conforming to
	// the "aspnetcore.diagnostics.handler.type" semantic conventions. It
	// represents the full type name of the
	// [`IExceptionHandler`](https://learn.microsoft.com/dotnet/api/microsoft.aspnetcore.diagnostics.iexceptionhandler)
	// implementation that handled the exception.
	//
	// Type: string
	// RequirementLevel: ConditionallyRequired (if and only if the exception
	// was handled by this handler.)
	// Stability: experimental
	// Examples: 'Contoso.MyHandler'
	AspnetcoreDiagnosticsHandlerTypeKey = "aspnetcore.diagnostics.handler.type"

	// AspnetcoreRateLimitingPolicyKey is the attribute Key conforming to the
	// "aspnetcore.rate_limiting.policy" semantic conventions. It represents
	// the rate limiting policy name.
	//
	// Type: string
	// RequirementLevel: ConditionallyRequired (if the matched endpoint for the
	// request had a rate-limiting policy.)
	// Stability: experimental
	// Examples: 'fixed', 'sliding', 'token'
	AspnetcoreRateLimitingPolicyKey = "aspnetcore.rate_limiting.policy"

	// AspnetcoreRateLimitingResultKey is the attribute Key conforming to the
	// "aspnetcore.rate_limiting.result" semantic conventions. It represents
	// the rate-limiting result, shows whether the lease was acquired or
	// contains a rejection reason
	//
	// Type: Enum
	// RequirementLevel: Required
	// Stability: experimental
	// Examples: 'acquired', 'request_canceled'
	AspnetcoreRateLimitingResultKey = "aspnetcore.rate_limiting.result"

	// AspnetcoreRequestIsUnhandledKey is the attribute Key conforming to the
	// "aspnetcore.request.is_unhandled" semantic conventions. It represents
	// the flag indicating if request was handled by the application pipeline.
	//
	// Type: boolean
	// RequirementLevel: ConditionallyRequired (if and only if the request was
	// not handled.)
	// Stability: experimental
	// Examples: True
	AspnetcoreRequestIsUnhandledKey = "aspnetcore.request.is_unhandled"

	// AspnetcoreRoutingIsFallbackKey is the attribute Key conforming to the
	// "aspnetcore.routing.is_fallback" semantic conventions. It represents a
	// value that indicates whether the matched route is a fallback route.
	//
	// Type: boolean
	// RequirementLevel: ConditionallyRequired (If and only if a route was
	// successfully matched.)
	// Stability: experimental
	// Examples: True
	AspnetcoreRoutingIsFallbackKey = "aspnetcore.routing.is_fallback"

	// SignalrConnectionStatusKey is the attribute Key conforming to the
	// "signalr.connection.status" semantic conventions. It represents the
	// signalR HTTP connection closure status.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'app_shutdown', 'timeout'
	SignalrConnectionStatusKey = "signalr.connection.status"

	// SignalrTransportKey is the attribute Key conforming to the
	// "signalr.transport" semantic conventions. It represents the [SignalR
	// transport
	// type](https://github.com/dotnet/aspnetcore/blob/main/src/SignalR/docs/specs/TransportProtocols.md)
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'web_sockets', 'long_polling'
	SignalrTransportKey = "signalr.transport"

	// JvmBufferPoolNameKey is the attribute Key conforming to the
	// "jvm.buffer.pool.name" semantic conventions. It represents the name of
	// the buffer pool.
	//
	// Type: string
	// RequirementLevel: Recommended
	// Stability: experimental
	// Examples: 'mapped', 'direct'
	// Note: Pool names are generally obtained via
	// [BufferPoolMXBean#getName()](https://docs.oracle.com/en/java/javase/11/docs/api/java.management/java/lang/management/BufferPoolMXBean.html#getName()).
	JvmBufferPoolNameKey = "jvm.buffer.pool.name"

	// JvmMemoryPoolNameKey is the attribute Key conforming to the
	// "jvm.memory.pool.name" semantic conventions. It represents the name of
	// the memory pool.
	//
	// Type: string
	// RequirementLevel: Recommended
	// Stability: stable
	// Examples: 'G1 Old Gen', 'G1 Eden space', 'G1 Survivor Space'
	// Note: Pool names are generally obtained via
	// [MemoryPoolMXBean#getName()](https://docs.oracle.com/en/java/javase/11/docs/api/java.management/java/lang/management/MemoryPoolMXBean.html#getName()).
	JvmMemoryPoolNameKey = "jvm.memory.pool.name"

	// JvmMemoryTypeKey is the attribute Key conforming to the
	// "jvm.memory.type" semantic conventions. It represents the type of
	// memory.
	//
	// Type: Enum
	// RequirementLevel: Recommended
	// Stability: stable
	// Examples: 'heap', 'non_heap'
	JvmMemoryTypeKey = "jvm.memory.type"

	// SystemDeviceKey is the attribute Key conforming to the "system.device"
	// semantic conventions. It represents the device identifier
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '(identifier)'
	SystemDeviceKey = "system.device"

	// SystemCPULogicalNumberKey is the attribute Key conforming to the
	// "system.cpu.logical_number" semantic conventions. It represents the
	// logical CPU number [0..n-1]
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 1
	SystemCPULogicalNumberKey = "system.cpu.logical_number"

	// SystemCPUStateKey is the attribute Key conforming to the
	// "system.cpu.state" semantic conventions. It represents the state of the
	// CPU
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'idle', 'interrupt'
	SystemCPUStateKey = "system.cpu.state"

	// SystemMemoryStateKey is the attribute Key conforming to the
	// "system.memory.state" semantic conventions. It represents the memory
	// state
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'free', 'cached'
	SystemMemoryStateKey = "system.memory.state"

	// SystemPagingDirectionKey is the attribute Key conforming to the
	// "system.paging.direction" semantic conventions. It represents the paging
	// access direction
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'in'
	SystemPagingDirectionKey = "system.paging.direction"

	// SystemPagingStateKey is the attribute Key conforming to the
	// "system.paging.state" semantic conventions. It represents the memory
	// paging state
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'free'
	SystemPagingStateKey = "system.paging.state"

	// SystemPagingTypeKey is the attribute Key conforming to the
	// "system.paging.type" semantic conventions. It represents the memory
	// paging type
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'minor'
	SystemPagingTypeKey = "system.paging.type"

	// SystemFilesystemModeKey is the attribute Key conforming to the
	// "system.filesystem.mode" semantic conventions. It represents the
	// filesystem mode
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'rw, ro'
	SystemFilesystemModeKey = "system.filesystem.mode"

	// SystemFilesystemMountpointKey is the attribute Key conforming to the
	// "system.filesystem.mountpoint" semantic conventions. It represents the
	// filesystem mount path
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '/mnt/data'
	SystemFilesystemMountpointKey = "system.filesystem.mountpoint"

	// SystemFilesystemStateKey is the attribute Key conforming to the
	// "system.filesystem.state" semantic conventions. It represents the
	// filesystem state
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'used'
	SystemFilesystemStateKey = "system.filesystem.state"

	// SystemFilesystemTypeKey is the attribute Key conforming to the
	// "system.filesystem.type" semantic conventions. It represents the
	// filesystem type
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'ext4'
	SystemFilesystemTypeKey = "system.filesystem.type"

	// SystemNetworkStateKey is the attribute Key conforming to the
	// "system.network.state" semantic conventions. It represents a stateless
	// protocol MUST NOT set this attribute
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'close_wait'
	SystemNetworkStateKey = "system.network.state"

	// SystemProcessesStatusKey is the attribute Key conforming to the
	// "system.processes.status" semantic conventions. It represents the
	// process state, e.g., [Linux Process State
	// Codes](https://man7.org/linux/man-pages/man1/ps.1.html#PROCESS_STATE_CODES)
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'running'
	SystemProcessesStatusKey = "system.processes.status"

	// ClientAddressKey is the attribute Key conforming to the "client.address"
	// semantic conventions. It represents the client address - domain name if
	// available without reverse DNS lookup; otherwise, IP address or Unix
	// domain socket name.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'client.example.com', '10.1.2.80', '/tmp/my.sock'
	// Note: When observed from the server side, and when communicating through
	// an intermediary, `client.address` SHOULD represent the client address
	// behind any intermediaries,  for example proxies, if it's available.
	ClientAddressKey = "client.address"

	// ClientPortKey is the attribute Key conforming to the "client.port"
	// semantic conventions. It represents the client port number.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 65123
	// Note: When observed from the server side, and when communicating through
	// an intermediary, `client.port` SHOULD represent the client port behind
	// any intermediaries,  for example proxies, if it's available.
	ClientPortKey = "client.port"

	// DBCassandraConsistencyLevelKey is the attribute Key conforming to the
	// "db.cassandra.consistency_level" semantic conventions. It represents the
	// consistency level of the query. Based on consistency values from
	// [CQL](https://docs.datastax.com/en/cassandra-oss/3.0/cassandra/dml/dmlConfigConsistency.html).
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	DBCassandraConsistencyLevelKey = "db.cassandra.consistency_level"

	// DBCassandraCoordinatorDCKey is the attribute Key conforming to the
	// "db.cassandra.coordinator.dc" semantic conventions. It represents the
	// data center of the coordinating node for a query.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'us-west-2'
	DBCassandraCoordinatorDCKey = "db.cassandra.coordinator.dc"

	// DBCassandraCoordinatorIDKey is the attribute Key conforming to the
	// "db.cassandra.coordinator.id" semantic conventions. It represents the ID
	// of the coordinating node for a query.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'be13faa2-8574-4d71-926d-27f16cf8a7af'
	DBCassandraCoordinatorIDKey = "db.cassandra.coordinator.id"

	// DBCassandraIdempotenceKey is the attribute Key conforming to the
	// "db.cassandra.idempotence" semantic conventions. It represents the
	// whether or not the query is idempotent.
	//
	// Type: boolean
	// RequirementLevel: Optional
	// Stability: experimental
	DBCassandraIdempotenceKey = "db.cassandra.idempotence"

	// DBCassandraPageSizeKey is the attribute Key conforming to the
	// "db.cassandra.page_size" semantic conventions. It represents the fetch
	// size used for paging, i.e. how many rows will be returned at once.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 5000
	DBCassandraPageSizeKey = "db.cassandra.page_size"

	// DBCassandraSpeculativeExecutionCountKey is the attribute Key conforming
	// to the "db.cassandra.speculative_execution_count" semantic conventions.
	// It represents the number of times a query was speculatively executed.
	// Not set or `0` if the query was not executed speculatively.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 0, 2
	DBCassandraSpeculativeExecutionCountKey = "db.cassandra.speculative_execution_count"

	// DBCassandraTableKey is the attribute Key conforming to the
	// "db.cassandra.table" semantic conventions. It represents the name of the
	// primary Cassandra table that the operation is acting upon, including the
	// keyspace name (if applicable).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'mytable'
	// Note: This mirrors the db.sql.table attribute but references cassandra
	// rather than sql. It is not recommended to attempt any client-side
	// parsing of `db.statement` just to get this property, but it should be
	// set if it is provided by the library being instrumented. If the
	// operation is acting upon an anonymous table, or more than one table,
	// this value MUST NOT be set.
	DBCassandraTableKey = "db.cassandra.table"

	// DBConnectionStringKey is the attribute Key conforming to the
	// "db.connection_string" semantic conventions. It represents the
	// connection string used to connect to the database. It is recommended to
	// remove embedded credentials.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'Server=(localdb)\\v11.0;Integrated Security=true;'
	DBConnectionStringKey = "db.connection_string"

	// DBCosmosDBClientIDKey is the attribute Key conforming to the
	// "db.cosmosdb.client_id" semantic conventions. It represents the unique
	// Cosmos client instance id.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '3ba4827d-4422-483f-b59f-85b74211c11d'
	DBCosmosDBClientIDKey = "db.cosmosdb.client_id"

	// DBCosmosDBConnectionModeKey is the attribute Key conforming to the
	// "db.cosmosdb.connection_mode" semantic conventions. It represents the
	// cosmos client connection mode.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	DBCosmosDBConnectionModeKey = "db.cosmosdb.connection_mode"

	// DBCosmosDBContainerKey is the attribute Key conforming to the
	// "db.cosmosdb.container" semantic conventions. It represents the cosmos
	// DB container name.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'anystring'
	DBCosmosDBContainerKey = "db.cosmosdb.container"

	// DBCosmosDBOperationTypeKey is the attribute Key conforming to the
	// "db.cosmosdb.operation_type" semantic conventions. It represents the
	// cosmosDB Operation Type.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	DBCosmosDBOperationTypeKey = "db.cosmosdb.operation_type"

	// DBCosmosDBRequestChargeKey is the attribute Key conforming to the
	// "db.cosmosdb.request_charge" semantic conventions. It represents the rU
	// consumed for that operation
	//
	// Type: double
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 46.18, 1.0
	DBCosmosDBRequestChargeKey = "db.cosmosdb.request_charge"

	// DBCosmosDBRequestContentLengthKey is the attribute Key conforming to the
	// "db.cosmosdb.request_content_length" semantic conventions. It represents
	// the request payload size in bytes
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	DBCosmosDBRequestContentLengthKey = "db.cosmosdb.request_content_length"

	// DBCosmosDBStatusCodeKey is the attribute Key conforming to the
	// "db.cosmosdb.status_code" semantic conventions. It represents the cosmos
	// DB status code.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 200, 201
	DBCosmosDBStatusCodeKey = "db.cosmosdb.status_code"

	// DBCosmosDBSubStatusCodeKey is the attribute Key conforming to the
	// "db.cosmosdb.sub_status_code" semantic conventions. It represents the
	// cosmos DB sub status code.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 1000, 1002
	DBCosmosDBSubStatusCodeKey = "db.cosmosdb.sub_status_code"

	// DBElasticsearchClusterNameKey is the attribute Key conforming to the
	// "db.elasticsearch.cluster.name" semantic conventions. It represents the
	// represents the identifier of an Elasticsearch cluster.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'e9106fc68e3044f0b1475b04bf4ffd5f'
	DBElasticsearchClusterNameKey = "db.elasticsearch.cluster.name"

	// DBElasticsearchNodeNameKey is the attribute Key conforming to the
	// "db.elasticsearch.node.name" semantic conventions. It represents the
	// represents the human-readable identifier of the node/instance to which a
	// request was routed.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'instance-0000000001'
	DBElasticsearchNodeNameKey = "db.elasticsearch.node.name"

	// DBInstanceIDKey is the attribute Key conforming to the "db.instance.id"
	// semantic conventions. It represents an identifier (address, unique name,
	// or any other identifier) of the database instance that is executing
	// queries or mutations on the current connection. This is useful in cases
	// where the database is running in a clustered environment and the
	// instrumentation is able to record the node executing the query. The
	// client may obtain this value in databases like MySQL using queries like
	// `select @@hostname`.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'mysql-e26b99z.example.com'
	DBInstanceIDKey = "db.instance.id"

	// DBJDBCDriverClassnameKey is the attribute Key conforming to the
	// "db.jdbc.driver_classname" semantic conventions. It represents the
	// fully-qualified class name of the [Java Database Connectivity
	// (JDBC)](https://docs.oracle.com/javase/8/docs/technotes/guides/jdbc/)
	// driver used to connect.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'org.postgresql.Driver',
	// 'com.microsoft.sqlserver.jdbc.SQLServerDriver'
	DBJDBCDriverClassnameKey = "db.jdbc.driver_classname"

	// DBMongoDBCollectionKey is the attribute Key conforming to the
	// "db.mongodb.collection" semantic conventions. It represents the MongoDB
	// collection being accessed within the database stated in `db.name`.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'customers', 'products'
	DBMongoDBCollectionKey = "db.mongodb.collection"

	// DBMSSQLInstanceNameKey is the attribute Key conforming to the
	// "db.mssql.instance_name" semantic conventions. It represents the
	// Microsoft SQL Server [instance
	// name](https://docs.microsoft.com/sql/connect/jdbc/building-the-connection-url?view=sql-server-ver15)
	// connecting to. This name is used to determine the port of a named
	// instance.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'MSSQLSERVER'
	// Note: If setting a `db.mssql.instance_name`, `server.port` is no longer
	// required (but still recommended if non-standard).
	DBMSSQLInstanceNameKey = "db.mssql.instance_name"

	// DBNameKey is the attribute Key conforming to the "db.name" semantic
	// conventions. It represents the this attribute is used to report the name
	// of the database being accessed. For commands that switch the database,
	// this should be set to the target database (even if the command fails).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'customers', 'main'
	// Note: In some SQL databases, the database name to be used is called
	// "schema name". In case there are multiple layers that could be
	// considered for database name (e.g. Oracle instance name and schema
	// name), the database name to be used is the more specific layer (e.g.
	// Oracle schema name).
	DBNameKey = "db.name"

	// DBOperationKey is the attribute Key conforming to the "db.operation"
	// semantic conventions. It represents the name of the operation being
	// executed, e.g. the [MongoDB command
	// name](https://docs.mongodb.com/manual/reference/command/#database-operations)
	// such as `findAndModify`, or the SQL keyword.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'findAndModify', 'HMSET', 'SELECT'
	// Note: When setting this to an SQL keyword, it is not recommended to
	// attempt any client-side parsing of `db.statement` just to get this
	// property, but it should be set if the operation name is provided by the
	// library being instrumented. If the SQL statement has an ambiguous
	// operation, or performs more than one operation, this value may be
	// omitted.
	DBOperationKey = "db.operation"

	// DBRedisDBIndexKey is the attribute Key conforming to the
	// "db.redis.database_index" semantic conventions. It represents the index
	// of the database being accessed as used in the [`SELECT`
	// command](https://redis.io/commands/select), provided as an integer. To
	// be used instead of the generic `db.name` attribute.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 0, 1, 15
	DBRedisDBIndexKey = "db.redis.database_index"

	// DBSQLTableKey is the attribute Key conforming to the "db.sql.table"
	// semantic conventions. It represents the name of the primary table that
	// the operation is acting upon, including the database name (if
	// applicable).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'public.users', 'customers'
	// Note: It is not recommended to attempt any client-side parsing of
	// `db.statement` just to get this property, but it should be set if it is
	// provided by the library being instrumented. If the operation is acting
	// upon an anonymous table, or more than one table, this value MUST NOT be
	// set.
	DBSQLTableKey = "db.sql.table"

	// DBStatementKey is the attribute Key conforming to the "db.statement"
	// semantic conventions. It represents the database statement being
	// executed.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'SELECT * FROM wuser_table', 'SET mykey "WuValue"'
	DBStatementKey = "db.statement"

	// DBSystemKey is the attribute Key conforming to the "db.system" semantic
	// conventions. It represents an identifier for the database management
	// system (DBMS) product being used. See below for a list of well-known
	// identifiers.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	DBSystemKey = "db.system"

	// DBUserKey is the attribute Key conforming to the "db.user" semantic
	// conventions. It represents the username for accessing the database.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'readonly_user', 'reporting_user'
	DBUserKey = "db.user"

	// HTTPFlavorKey is the attribute Key conforming to the "http.flavor"
	// semantic conventions.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: deprecated
	// Deprecated: use `network.protocol.name` instead.
	HTTPFlavorKey = "http.flavor"

	// HTTPMethodKey is the attribute Key conforming to the "http.method"
	// semantic conventions.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: 'GET', 'POST', 'HEAD'
	// Deprecated: use `http.request.method` instead.
	HTTPMethodKey = "http.method"

	// HTTPRequestContentLengthKey is the attribute Key conforming to the
	// "http.request_content_length" semantic conventions.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: 3495
	// Deprecated: use `http.request.header.content-length` instead.
	HTTPRequestContentLengthKey = "http.request_content_length"

	// HTTPResponseContentLengthKey is the attribute Key conforming to the
	// "http.response_content_length" semantic conventions.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: 3495
	// Deprecated: use `http.response.header.content-length` instead.
	HTTPResponseContentLengthKey = "http.response_content_length"

	// HTTPSchemeKey is the attribute Key conforming to the "http.scheme"
	// semantic conventions.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: 'http', 'https'
	// Deprecated: use `url.scheme` instead.
	HTTPSchemeKey = "http.scheme"

	// HTTPStatusCodeKey is the attribute Key conforming to the
	// "http.status_code" semantic conventions.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: 200
	// Deprecated: use `http.response.status_code` instead.
	HTTPStatusCodeKey = "http.status_code"

	// HTTPTargetKey is the attribute Key conforming to the "http.target"
	// semantic conventions.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: '/search?q=OpenTelemetry#SemConv'
	// Deprecated: use `url.path` and `url.query` instead.
	HTTPTargetKey = "http.target"

	// HTTPURLKey is the attribute Key conforming to the "http.url" semantic
	// conventions.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: 'https://www.foo.bar/search?q=OpenTelemetry#SemConv'
	// Deprecated: use `url.full` instead.
	HTTPURLKey = "http.url"

	// HTTPUserAgentKey is the attribute Key conforming to the
	// "http.user_agent" semantic conventions.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: 'CERN-LineMode/2.15 libwww/2.17b3', 'Mozilla/5.0 (iPhone; CPU
	// iPhone OS 14_7_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko)
	// Version/14.1.2 Mobile/15E148 Safari/604.1'
	// Deprecated: use `user_agent.original` instead.
	HTTPUserAgentKey = "http.user_agent"

	// NetHostNameKey is the attribute Key conforming to the "net.host.name"
	// semantic conventions.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: 'example.com'
	// Deprecated: use `server.address`.
	NetHostNameKey = "net.host.name"

	// NetHostPortKey is the attribute Key conforming to the "net.host.port"
	// semantic conventions.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: 8080
	// Deprecated: use `server.port`.
	NetHostPortKey = "net.host.port"

	// NetPeerNameKey is the attribute Key conforming to the "net.peer.name"
	// semantic conventions.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: 'example.com'
	// Deprecated: use `server.address` on client spans and `client.address` on
	// server spans.
	NetPeerNameKey = "net.peer.name"

	// NetPeerPortKey is the attribute Key conforming to the "net.peer.port"
	// semantic conventions.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: 8080
	// Deprecated: use `server.port` on client spans and `client.port` on
	// server spans.
	NetPeerPortKey = "net.peer.port"

	// NetProtocolNameKey is the attribute Key conforming to the
	// "net.protocol.name" semantic conventions.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: 'amqp', 'http', 'mqtt'
	// Deprecated: use `network.protocol.name`.
	NetProtocolNameKey = "net.protocol.name"

	// NetProtocolVersionKey is the attribute Key conforming to the
	// "net.protocol.version" semantic conventions.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: '3.1.1'
	// Deprecated: use `network.protocol.version`.
	NetProtocolVersionKey = "net.protocol.version"

	// NetSockFamilyKey is the attribute Key conforming to the
	// "net.sock.family" semantic conventions.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: deprecated
	// Deprecated: use `network.transport` and `network.type`.
	NetSockFamilyKey = "net.sock.family"

	// NetSockHostAddrKey is the attribute Key conforming to the
	// "net.sock.host.addr" semantic conventions.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: '/var/my.sock'
	// Deprecated: use `network.local.address`.
	NetSockHostAddrKey = "net.sock.host.addr"

	// NetSockHostPortKey is the attribute Key conforming to the
	// "net.sock.host.port" semantic conventions.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: 8080
	// Deprecated: use `network.local.port`.
	NetSockHostPortKey = "net.sock.host.port"

	// NetSockPeerAddrKey is the attribute Key conforming to the
	// "net.sock.peer.addr" semantic conventions.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: '192.168.0.1'
	// Deprecated: use `network.peer.address`.
	NetSockPeerAddrKey = "net.sock.peer.addr"

	// NetSockPeerNameKey is the attribute Key conforming to the
	// "net.sock.peer.name" semantic conventions.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: '/var/my.sock'
	// Deprecated: no replacement at this time.
	NetSockPeerNameKey = "net.sock.peer.name"

	// NetSockPeerPortKey is the attribute Key conforming to the
	// "net.sock.peer.port" semantic conventions.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: 65531
	// Deprecated: use `network.peer.port`.
	NetSockPeerPortKey = "net.sock.peer.port"

	// NetTransportKey is the attribute Key conforming to the "net.transport"
	// semantic conventions.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: deprecated
	// Deprecated: use `network.transport`.
	NetTransportKey = "net.transport"

	// DestinationAddressKey is the attribute Key conforming to the
	// "destination.address" semantic conventions. It represents the
	// destination address - domain name if available without reverse DNS
	// lookup; otherwise, IP address or Unix domain socket name.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'destination.example.com', '10.1.2.80', '/tmp/my.sock'
	// Note: When observed from the source side, and when communicating through
	// an intermediary, `destination.address` SHOULD represent the destination
	// address behind any intermediaries, for example proxies, if it's
	// available.
	DestinationAddressKey = "destination.address"

	// DestinationPortKey is the attribute Key conforming to the
	// "destination.port" semantic conventions. It represents the destination
	// port number
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 3389, 2888
	DestinationPortKey = "destination.port"

	// DiskIoDirectionKey is the attribute Key conforming to the
	// "disk.io.direction" semantic conventions. It represents the disk IO
	// operation direction.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'read'
	DiskIoDirectionKey = "disk.io.direction"

	// ErrorTypeKey is the attribute Key conforming to the "error.type"
	// semantic conventions. It represents the describes a class of error the
	// operation ended with.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'timeout', 'java.net.UnknownHostException',
	// 'server_certificate_invalid', '500'
	// Note: The `error.type` SHOULD be predictable and SHOULD have low
	// cardinality.
	// Instrumentations SHOULD document the list of errors they report.
	//
	// The cardinality of `error.type` within one instrumentation library
	// SHOULD be low.
	// Telemetry consumers that aggregate data from multiple instrumentation
	// libraries and applications
	// should be prepared for `error.type` to have high cardinality at query
	// time when no
	// additional filters are applied.
	//
	// If the operation has completed successfully, instrumentations SHOULD NOT
	// set `error.type`.
	//
	// If a specific domain defines its own set of error identifiers (such as
	// HTTP or gRPC status codes),
	// it's RECOMMENDED to:
	//
	// * Use a domain-specific attribute
	// * Set `error.type` to capture all errors, regardless of whether they are
	// defined within the domain-specific set or not.
	ErrorTypeKey = "error.type"

	// ExceptionEscapedKey is the attribute Key conforming to the
	// "exception.escaped" semantic conventions. It represents the sHOULD be
	// set to true if the exception event is recorded at a point where it is
	// known that the exception is escaping the scope of the span.
	//
	// Type: boolean
	// RequirementLevel: Optional
	// Stability: experimental
	// Note: An exception is considered to have escaped (or left) the scope of
	// a span,
	// if that span is ended while the exception is still logically "in
	// flight".
	// This may be actually "in flight" in some languages (e.g. if the
	// exception
	// is passed to a Context manager's `__exit__` method in Python) but will
	// usually be caught at the point of recording the exception in most
	// languages.
	//
	// It is usually not possible to determine at the point where an exception
	// is thrown
	// whether it will escape the scope of a span.
	// However, it is trivial to know that an exception
	// will escape, if one checks for an active exception just before ending
	// the span,
	// as done in the [example for recording span
	// exceptions](#recording-an-exception).
	//
	// It follows that an exception may still escape the scope of the span
	// even if the `exception.escaped` attribute was not set or set to false,
	// since the event might have been recorded at a time where it was not
	// clear whether the exception will escape.
	ExceptionEscapedKey = "exception.escaped"

	// ExceptionMessageKey is the attribute Key conforming to the
	// "exception.message" semantic conventions. It represents the exception
	// message.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'Division by zero', "Can't convert 'int' object to str
	// implicitly"
	ExceptionMessageKey = "exception.message"

	// ExceptionStacktraceKey is the attribute Key conforming to the
	// "exception.stacktrace" semantic conventions. It represents a stacktrace
	// as a string in the natural representation for the language runtime. The
	// representation is to be determined and documented by each language SIG.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'Exception in thread "main" java.lang.RuntimeException: Test
	// exception\\n at '
	//  'com.example.GenerateTrace.methodB(GenerateTrace.java:13)\\n at '
	//  'com.example.GenerateTrace.methodA(GenerateTrace.java:9)\\n at '
	//  'com.example.GenerateTrace.main(GenerateTrace.java:5)'
	ExceptionStacktraceKey = "exception.stacktrace"

	// ExceptionTypeKey is the attribute Key conforming to the "exception.type"
	// semantic conventions. It represents the type of the exception (its
	// fully-qualified class name, if applicable). The dynamic type of the
	// exception should be preferred over the static type in languages that
	// support it.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'java.net.ConnectException', 'OSError'
	ExceptionTypeKey = "exception.type"

	// HTTPRequestBodySizeKey is the attribute Key conforming to the
	// "http.request.body.size" semantic conventions. It represents the size of
	// the request payload body in bytes. This is the number of bytes
	// transferred excluding headers and is often, but not always, present as
	// the
	// [Content-Length](https://www.rfc-editor.org/rfc/rfc9110.html#field.content-length)
	// header. For requests using transport encoding, this should be the
	// compressed size.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 3495
	HTTPRequestBodySizeKey = "http.request.body.size"

	// HTTPRequestMethodKey is the attribute Key conforming to the
	// "http.request.method" semantic conventions. It represents the hTTP
	// request method.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'GET', 'POST', 'HEAD'
	// Note: HTTP request method value SHOULD be "known" to the
	// instrumentation.
	// By default, this convention defines "known" methods as the ones listed
	// in [RFC9110](https://www.rfc-editor.org/rfc/rfc9110.html#name-methods)
	// and the PATCH method defined in
	// [RFC5789](https://www.rfc-editor.org/rfc/rfc5789.html).
	//
	// If the HTTP request method is not known to instrumentation, it MUST set
	// the `http.request.method` attribute to `_OTHER`.
	//
	// If the HTTP instrumentation could end up converting valid HTTP request
	// methods to `_OTHER`, then it MUST provide a way to override
	// the list of known HTTP methods. If this override is done via environment
	// variable, then the environment variable MUST be named
	// OTEL_INSTRUMENTATION_HTTP_KNOWN_METHODS and support a comma-separated
	// list of case-sensitive known HTTP methods
	// (this list MUST be a full override of the default known method, it is
	// not a list of known methods in addition to the defaults).
	//
	// HTTP method names are case-sensitive and `http.request.method` attribute
	// value MUST match a known HTTP method name exactly.
	// Instrumentations for specific web frameworks that consider HTTP methods
	// to be case insensitive, SHOULD populate a canonical equivalent.
	// Tracing instrumentations that do so, MUST also set
	// `http.request.method_original` to the original value.
	HTTPRequestMethodKey = "http.request.method"

	// HTTPRequestMethodOriginalKey is the attribute Key conforming to the
	// "http.request.method_original" semantic conventions. It represents the
	// original HTTP method sent by the client in the request line.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'GeT', 'ACL', 'foo'
	HTTPRequestMethodOriginalKey = "http.request.method_original"

	// HTTPRequestResendCountKey is the attribute Key conforming to the
	// "http.request.resend_count" semantic conventions. It represents the
	// ordinal number of request resending attempt (for any reason, including
	// redirects).
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 3
	// Note: The resend count SHOULD be updated each time an HTTP request gets
	// resent by the client, regardless of what was the cause of the resending
	// (e.g. redirection, authorization failure, 503 Server Unavailable,
	// network issues, or any other).
	HTTPRequestResendCountKey = "http.request.resend_count"

	// HTTPResponseBodySizeKey is the attribute Key conforming to the
	// "http.response.body.size" semantic conventions. It represents the size
	// of the response payload body in bytes. This is the number of bytes
	// transferred excluding headers and is often, but not always, present as
	// the
	// [Content-Length](https://www.rfc-editor.org/rfc/rfc9110.html#field.content-length)
	// header. For requests using transport encoding, this should be the
	// compressed size.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 3495
	HTTPResponseBodySizeKey = "http.response.body.size"

	// HTTPResponseStatusCodeKey is the attribute Key conforming to the
	// "http.response.status_code" semantic conventions. It represents the
	// [HTTP response status
	// code](https://tools.ietf.org/html/rfc7231#section-6).
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 200
	HTTPResponseStatusCodeKey = "http.response.status_code"

	// HTTPRouteKey is the attribute Key conforming to the "http.route"
	// semantic conventions. It represents the matched route, that is, the path
	// template in the format used by the respective server framework.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: '/users/:userID?', '{controller}/{action}/{id?}'
	// Note: MUST NOT be populated when this is not supported by the HTTP
	// server framework as the route attribute should have low-cardinality and
	// the URI path can NOT substitute it.
	// SHOULD include the [application
	// root](/docs/http/http-spans.md#http-server-definitions) if there is one.
	HTTPRouteKey = "http.route"

	// MessagingBatchMessageCountKey is the attribute Key conforming to the
	// "messaging.batch.message_count" semantic conventions. It represents the
	// number of messages sent, received, or processed in the scope of the
	// batching operation.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 0, 1, 2
	// Note: Instrumentations SHOULD NOT set `messaging.batch.message_count` on
	// spans that operate with a single message. When a messaging client
	// library supports both batch and single-message API for the same
	// operation, instrumentations SHOULD use `messaging.batch.message_count`
	// for batching APIs and SHOULD NOT use it for single-message APIs.
	MessagingBatchMessageCountKey = "messaging.batch.message_count"

	// MessagingClientIDKey is the attribute Key conforming to the
	// "messaging.client_id" semantic conventions. It represents a unique
	// identifier for the client that consumes or produces a message.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'client-5', 'myhost@8742@s8083jm'
	MessagingClientIDKey = "messaging.client_id"

	// MessagingDestinationAnonymousKey is the attribute Key conforming to the
	// "messaging.destination.anonymous" semantic conventions. It represents a
	// boolean that is true if the message destination is anonymous (could be
	// unnamed or have auto-generated name).
	//
	// Type: boolean
	// RequirementLevel: Optional
	// Stability: experimental
	MessagingDestinationAnonymousKey = "messaging.destination.anonymous"

	// MessagingDestinationNameKey is the attribute Key conforming to the
	// "messaging.destination.name" semantic conventions. It represents the
	// message destination name
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'MyQueue', 'MyTopic'
	// Note: Destination name SHOULD uniquely identify a specific queue, topic
	// or other entity within the broker. If
	// the broker doesn't have such notion, the destination name SHOULD
	// uniquely identify the broker.
	MessagingDestinationNameKey = "messaging.destination.name"

	// MessagingDestinationTemplateKey is the attribute Key conforming to the
	// "messaging.destination.template" semantic conventions. It represents the
	// low cardinality representation of the messaging destination name
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '/customers/{customerID}'
	// Note: Destination names could be constructed from templates. An example
	// would be a destination name involving a user name or product id.
	// Although the destination name in this case is of high cardinality, the
	// underlying template is of low cardinality and can be effectively used
	// for grouping and aggregation.
	MessagingDestinationTemplateKey = "messaging.destination.template"

	// MessagingDestinationTemporaryKey is the attribute Key conforming to the
	// "messaging.destination.temporary" semantic conventions. It represents a
	// boolean that is true if the message destination is temporary and might
	// not exist anymore after messages are processed.
	//
	// Type: boolean
	// RequirementLevel: Optional
	// Stability: experimental
	MessagingDestinationTemporaryKey = "messaging.destination.temporary"

	// MessagingDestinationPublishAnonymousKey is the attribute Key conforming
	// to the "messaging.destination_publish.anonymous" semantic conventions.
	// It represents a boolean that is true if the publish message destination
	// is anonymous (could be unnamed or have auto-generated name).
	//
	// Type: boolean
	// RequirementLevel: Optional
	// Stability: experimental
	MessagingDestinationPublishAnonymousKey = "messaging.destination_publish.anonymous"

	// MessagingDestinationPublishNameKey is the attribute Key conforming to
	// the "messaging.destination_publish.name" semantic conventions. It
	// represents the name of the original destination the message was
	// published to
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'MyQueue', 'MyTopic'
	// Note: The name SHOULD uniquely identify a specific queue, topic, or
	// other entity within the broker. If
	// the broker doesn't have such notion, the original destination name
	// SHOULD uniquely identify the broker.
	MessagingDestinationPublishNameKey = "messaging.destination_publish.name"

	// MessagingGCPPubsubMessageOrderingKeyKey is the attribute Key conforming
	// to the "messaging.gcp_pubsub.message.ordering_key" semantic conventions.
	// It represents the ordering key for a given message. If the attribute is
	// not present, the message does not have an ordering key.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'ordering_key'
	MessagingGCPPubsubMessageOrderingKeyKey = "messaging.gcp_pubsub.message.ordering_key"

	// MessagingKafkaConsumerGroupKey is the attribute Key conforming to the
	// "messaging.kafka.consumer.group" semantic conventions. It represents the
	// name of the Kafka Consumer Group that is handling the message. Only
	// applies to consumers, not producers.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'my-group'
	MessagingKafkaConsumerGroupKey = "messaging.kafka.consumer.group"

	// MessagingKafkaDestinationPartitionKey is the attribute Key conforming to
	// the "messaging.kafka.destination.partition" semantic conventions. It
	// represents the partition the message is sent to.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 2
	MessagingKafkaDestinationPartitionKey = "messaging.kafka.destination.partition"

	// MessagingKafkaMessageKeyKey is the attribute Key conforming to the
	// "messaging.kafka.message.key" semantic conventions. It represents the
	// message keys in Kafka are used for grouping alike messages to ensure
	// they're processed on the same partition. They differ from
	// `messaging.message.id` in that they're not unique. If the key is `null`,
	// the attribute MUST NOT be set.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'myKey'
	// Note: If the key type is not string, it's string representation has to
	// be supplied for the attribute. If the key has no unambiguous, canonical
	// string form, don't include its value.
	MessagingKafkaMessageKeyKey = "messaging.kafka.message.key"

	// MessagingKafkaMessageOffsetKey is the attribute Key conforming to the
	// "messaging.kafka.message.offset" semantic conventions. It represents the
	// offset of a record in the corresponding Kafka partition.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 42
	MessagingKafkaMessageOffsetKey = "messaging.kafka.message.offset"

	// MessagingKafkaMessageTombstoneKey is the attribute Key conforming to the
	// "messaging.kafka.message.tombstone" semantic conventions. It represents
	// a boolean that is true if the message is a tombstone.
	//
	// Type: boolean
	// RequirementLevel: Optional
	// Stability: experimental
	MessagingKafkaMessageTombstoneKey = "messaging.kafka.message.tombstone"

	// MessagingMessageBodySizeKey is the attribute Key conforming to the
	// "messaging.message.body.size" semantic conventions. It represents the
	// size of the message body in bytes.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 1439
	// Note: This can refer to both the compressed or uncompressed body size.
	// If both sizes are known, the uncompressed
	// body size should be used.
	MessagingMessageBodySizeKey = "messaging.message.body.size"

	// MessagingMessageConversationIDKey is the attribute Key conforming to the
	// "messaging.message.conversation_id" semantic conventions. It represents
	// the conversation ID identifying the conversation to which the message
	// belongs, represented as a string. Sometimes called "Correlation ID".
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'MyConversationID'
	MessagingMessageConversationIDKey = "messaging.message.conversation_id"

	// MessagingMessageEnvelopeSizeKey is the attribute Key conforming to the
	// "messaging.message.envelope.size" semantic conventions. It represents
	// the size of the message body and metadata in bytes.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 2738
	// Note: This can refer to both the compressed or uncompressed size. If
	// both sizes are known, the uncompressed
	// size should be used.
	MessagingMessageEnvelopeSizeKey = "messaging.message.envelope.size"

	// MessagingMessageIDKey is the attribute Key conforming to the
	// "messaging.message.id" semantic conventions. It represents a value used
	// by the messaging system as an identifier for the message, represented as
	// a string.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '452a7c7c7c7048c2f887f61572b18fc2'
	MessagingMessageIDKey = "messaging.message.id"

	// MessagingOperationKey is the attribute Key conforming to the
	// "messaging.operation" semantic conventions. It represents a string
	// identifying the kind of messaging operation.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	// Note: If a custom value is used, it MUST be of low cardinality.
	MessagingOperationKey = "messaging.operation"

	// MessagingRabbitmqDestinationRoutingKeyKey is the attribute Key
	// conforming to the "messaging.rabbitmq.destination.routing_key" semantic
	// conventions. It represents the rabbitMQ message routing key.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'myKey'
	MessagingRabbitmqDestinationRoutingKeyKey = "messaging.rabbitmq.destination.routing_key"

	// MessagingRocketmqClientGroupKey is the attribute Key conforming to the
	// "messaging.rocketmq.client_group" semantic conventions. It represents
	// the name of the RocketMQ producer/consumer group that is handling the
	// message. The client type is identified by the SpanKind.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'myConsumerGroup'
	MessagingRocketmqClientGroupKey = "messaging.rocketmq.client_group"

	// MessagingRocketmqConsumptionModelKey is the attribute Key conforming to
	// the "messaging.rocketmq.consumption_model" semantic conventions. It
	// represents the model of message consumption. This only applies to
	// consumer spans.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	MessagingRocketmqConsumptionModelKey = "messaging.rocketmq.consumption_model"

	// MessagingRocketmqMessageDelayTimeLevelKey is the attribute Key
	// conforming to the "messaging.rocketmq.message.delay_time_level" semantic
	// conventions. It represents the delay time level for delay message, which
	// determines the message delay time.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 3
	MessagingRocketmqMessageDelayTimeLevelKey = "messaging.rocketmq.message.delay_time_level"

	// MessagingRocketmqMessageDeliveryTimestampKey is the attribute Key
	// conforming to the "messaging.rocketmq.message.delivery_timestamp"
	// semantic conventions. It represents the timestamp in milliseconds that
	// the delay message is expected to be delivered to consumer.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 1665987217045
	MessagingRocketmqMessageDeliveryTimestampKey = "messaging.rocketmq.message.delivery_timestamp"

	// MessagingRocketmqMessageGroupKey is the attribute Key conforming to the
	// "messaging.rocketmq.message.group" semantic conventions. It represents
	// the it is essential for FIFO message. Messages that belong to the same
	// message group are always processed one by one within the same consumer
	// group.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'myMessageGroup'
	MessagingRocketmqMessageGroupKey = "messaging.rocketmq.message.group"

	// MessagingRocketmqMessageKeysKey is the attribute Key conforming to the
	// "messaging.rocketmq.message.keys" semantic conventions. It represents
	// the key(s) of message, another way to mark message besides message id.
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'keyA', 'keyB'
	MessagingRocketmqMessageKeysKey = "messaging.rocketmq.message.keys"

	// MessagingRocketmqMessageTagKey is the attribute Key conforming to the
	// "messaging.rocketmq.message.tag" semantic conventions. It represents the
	// secondary classifier of message besides topic.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'tagA'
	MessagingRocketmqMessageTagKey = "messaging.rocketmq.message.tag"

	// MessagingRocketmqMessageTypeKey is the attribute Key conforming to the
	// "messaging.rocketmq.message.type" semantic conventions. It represents
	// the type of message.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	MessagingRocketmqMessageTypeKey = "messaging.rocketmq.message.type"

	// MessagingRocketmqNamespaceKey is the attribute Key conforming to the
	// "messaging.rocketmq.namespace" semantic conventions. It represents the
	// namespace of RocketMQ resources, resources in different namespaces are
	// individual.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'myNamespace'
	MessagingRocketmqNamespaceKey = "messaging.rocketmq.namespace"

	// MessagingSystemKey is the attribute Key conforming to the
	// "messaging.system" semantic conventions. It represents an identifier for
	// the messaging system being used. See below for a list of well-known
	// identifiers.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	MessagingSystemKey = "messaging.system"

	// NetworkCarrierIccKey is the attribute Key conforming to the
	// "network.carrier.icc" semantic conventions. It represents the ISO 3166-1
	// alpha-2 2-character country code associated with the mobile carrier
	// network.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'DE'
	NetworkCarrierIccKey = "network.carrier.icc"

	// NetworkCarrierMccKey is the attribute Key conforming to the
	// "network.carrier.mcc" semantic conventions. It represents the mobile
	// carrier country code.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '310'
	NetworkCarrierMccKey = "network.carrier.mcc"

	// NetworkCarrierMncKey is the attribute Key conforming to the
	// "network.carrier.mnc" semantic conventions. It represents the mobile
	// carrier network code.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '001'
	NetworkCarrierMncKey = "network.carrier.mnc"

	// NetworkCarrierNameKey is the attribute Key conforming to the
	// "network.carrier.name" semantic conventions. It represents the name of
	// the mobile carrier.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'sprint'
	NetworkCarrierNameKey = "network.carrier.name"

	// NetworkConnectionSubtypeKey is the attribute Key conforming to the
	// "network.connection.subtype" semantic conventions. It represents the
	// this describes more details regarding the connection.type. It may be the
	// type of cell technology connection, but it could be used for describing
	// details about a wifi connection.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'LTE'
	NetworkConnectionSubtypeKey = "network.connection.subtype"

	// NetworkConnectionTypeKey is the attribute Key conforming to the
	// "network.connection.type" semantic conventions. It represents the
	// internet connection type.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'wifi'
	NetworkConnectionTypeKey = "network.connection.type"

	// NetworkIoDirectionKey is the attribute Key conforming to the
	// "network.io.direction" semantic conventions. It represents the network
	// IO operation direction.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'transmit'
	NetworkIoDirectionKey = "network.io.direction"

	// NetworkLocalAddressKey is the attribute Key conforming to the
	// "network.local.address" semantic conventions. It represents the local
	// address of the network connection - IP address or Unix domain socket
	// name.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: '10.1.2.80', '/tmp/my.sock'
	NetworkLocalAddressKey = "network.local.address"

	// NetworkLocalPortKey is the attribute Key conforming to the
	// "network.local.port" semantic conventions. It represents the local port
	// number of the network connection.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 65123
	NetworkLocalPortKey = "network.local.port"

	// NetworkPeerAddressKey is the attribute Key conforming to the
	// "network.peer.address" semantic conventions. It represents the peer
	// address of the network connection - IP address or Unix domain socket
	// name.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: '10.1.2.80', '/tmp/my.sock'
	NetworkPeerAddressKey = "network.peer.address"

	// NetworkPeerPortKey is the attribute Key conforming to the
	// "network.peer.port" semantic conventions. It represents the peer port
	// number of the network connection.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 65123
	NetworkPeerPortKey = "network.peer.port"

	// NetworkProtocolNameKey is the attribute Key conforming to the
	// "network.protocol.name" semantic conventions. It represents the [OSI
	// application layer](https://osi-model.com/application-layer/) or non-OSI
	// equivalent.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'amqp', 'http', 'mqtt'
	// Note: The value SHOULD be normalized to lowercase.
	NetworkProtocolNameKey = "network.protocol.name"

	// NetworkProtocolVersionKey is the attribute Key conforming to the
	// "network.protocol.version" semantic conventions. It represents the
	// version of the protocol specified in `network.protocol.name`.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: '3.1.1'
	// Note: `network.protocol.version` refers to the version of the protocol
	// used and might be different from the protocol client's version. If the
	// HTTP client has a version of `0.27.2`, but sends HTTP version `1.1`,
	// this attribute should be set to `1.1`.
	NetworkProtocolVersionKey = "network.protocol.version"

	// NetworkTransportKey is the attribute Key conforming to the
	// "network.transport" semantic conventions. It represents the [OSI
	// transport layer](https://osi-model.com/transport-layer/) or
	// [inter-process communication
	// method](https://wikipedia.org/wiki/Inter-process_communication).
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'tcp', 'udp'
	// Note: The value SHOULD be normalized to lowercase.
	//
	// Consider always setting the transport when setting a port number, since
	// a port number is ambiguous without knowing the transport. For example
	// different processes could be listening on TCP port 12345 and UDP port
	// 12345.
	NetworkTransportKey = "network.transport"

	// NetworkTypeKey is the attribute Key conforming to the "network.type"
	// semantic conventions. It represents the [OSI network
	// layer](https://osi-model.com/network-layer/) or non-OSI equivalent.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'ipv4', 'ipv6'
	// Note: The value SHOULD be normalized to lowercase.
	NetworkTypeKey = "network.type"

	// RPCConnectRPCErrorCodeKey is the attribute Key conforming to the
	// "rpc.connect_rpc.error_code" semantic conventions. It represents the
	// [error codes](https://connect.build/docs/protocol/#error-codes) of the
	// Connect request. Error codes are always string values.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	RPCConnectRPCErrorCodeKey = "rpc.connect_rpc.error_code"

	// RPCGRPCStatusCodeKey is the attribute Key conforming to the
	// "rpc.transport.status_code" semantic conventions. It represents the [numeric
	// status
	// code](https://github.com/grpc/grpc/blob/v1.33.2/doc/statuscodes.md) of
	// the gRPC request.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	RPCGRPCStatusCodeKey = "rpc.transport.status_code"

	// RPCJsonrpcErrorCodeKey is the attribute Key conforming to the
	// "rpc.jsonrpc.error_code" semantic conventions. It represents the
	// `error.code` property of response if it is an error response.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: -32700, 100
	RPCJsonrpcErrorCodeKey = "rpc.jsonrpc.error_code"

	// RPCJsonrpcErrorMessageKey is the attribute Key conforming to the
	// "rpc.jsonrpc.error_message" semantic conventions. It represents the
	// `error.message` property of response if it is an error response.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'Parse error', 'User already exists'
	RPCJsonrpcErrorMessageKey = "rpc.jsonrpc.error_message"

	// RPCJsonrpcRequestIDKey is the attribute Key conforming to the
	// "rpc.jsonrpc.request_id" semantic conventions. It represents the `id`
	// property of request or response. Since protocol allows id to be int,
	// string, `null` or missing (for notifications), value is expected to be
	// cast to string for simplicity. Use empty string in case of `null` value.
	// Omit entirely if this is a notification.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '10', 'request-7', ''
	RPCJsonrpcRequestIDKey = "rpc.jsonrpc.request_id"

	// RPCJsonrpcVersionKey is the attribute Key conforming to the
	// "rpc.jsonrpc.version" semantic conventions. It represents the protocol
	// version as in `jsonrpc` property of request/response. Since JSON-RPC 1.0
	// doesn't specify this, the value can be omitted.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '2.0', '1.0'
	RPCJsonrpcVersionKey = "rpc.jsonrpc.version"

	// RPCMethodKey is the attribute Key conforming to the "rpc.method"
	// semantic conventions. It represents the name of the (logical) method
	// being called, must be equal to the $method part in the span name.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'exampleMethod'
	// Note: This is the logical name of the method from the RPC interface
	// perspective, which can be different from the name of any implementing
	// method/function. The `code.function` attribute may be used to store the
	// latter (e.g., method actually executing the call on the server side, RPC
	// client stub method on the client side).
	RPCMethodKey = "rpc.method"

	// RPCServiceKey is the attribute Key conforming to the "rpc.service"
	// semantic conventions. It represents the full (logical) name of the
	// service being called, including its package name, if applicable.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'myservice.EchoService'
	// Note: This is the logical name of the service from the RPC interface
	// perspective, which can be different from the name of any implementing
	// class. The `code.namespace` attribute may be used to store the latter
	// (despite the attribute name, it may include a class name; e.g., class
	// with method actually executing the call on the server side, RPC client
	// stub class on the client side).
	RPCServiceKey = "rpc.service"

	// RPCSystemKey is the attribute Key conforming to the "rpc.system"
	// semantic conventions. It represents a string identifying the remoting
	// system. See below for a list of well-known identifiers.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	RPCSystemKey = "rpc.system"

	// ServerAddressKey is the attribute Key conforming to the "server.address"
	// semantic conventions. It represents the server domain name if available
	// without reverse DNS lookup; otherwise, IP address or Unix domain socket
	// name.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'example.com', '10.1.2.80', '/tmp/my.sock'
	// Note: When observed from the client side, and when communicating through
	// an intermediary, `server.address` SHOULD represent the server address
	// behind any intermediaries, for example proxies, if it's available.
	ServerAddressKey = "server.address"

	// ServerPortKey is the attribute Key conforming to the "server.port"
	// semantic conventions. It represents the server port number.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 80, 8080, 443
	// Note: When observed from the client side, and when communicating through
	// an intermediary, `server.port` SHOULD represent the server port behind
	// any intermediaries, for example proxies, if it's available.
	ServerPortKey = "server.port"

	// SourceAddressKey is the attribute Key conforming to the "source.address"
	// semantic conventions. It represents the source address - domain name if
	// available without reverse DNS lookup; otherwise, IP address or Unix
	// domain socket name.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'source.example.com', '10.1.2.80', '/tmp/my.sock'
	// Note: When observed from the destination side, and when communicating
	// through an intermediary, `source.address` SHOULD represent the source
	// address behind any intermediaries, for example proxies, if it's
	// available.
	SourceAddressKey = "source.address"

	// SourcePortKey is the attribute Key conforming to the "source.port"
	// semantic conventions. It represents the source port number
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 3389, 2888
	SourcePortKey = "source.port"

	// TLSCipherKey is the attribute Key conforming to the "tls.cipher"
	// semantic conventions. It represents the string indicating the
	// [cipher](https://datatracker.ietf.org/doc/html/rfc5246#appendix-A.5)
	// used during the current connection.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'TLS_RSA_WITH_3DES_EDE_CBC_SHA',
	// 'TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256'
	// Note: The values allowed for `tls.cipher` MUST be one of the
	// `Descriptions` of the [registered TLS Cipher
	// Suits](https://www.iana.org/assignments/tls-parameters/tls-parameters.xhtml#table-tls-parameters-4).
	TLSCipherKey = "tls.cipher"

	// TLSClientCertificateKey is the attribute Key conforming to the
	// "tls.client.certificate" semantic conventions. It represents the
	// pEM-encoded stand-alone certificate offered by the client. This is
	// usually mutually-exclusive of `client.certificate_chain` since this
	// value also exists in that list.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'MII...'
	TLSClientCertificateKey = "tls.client.certificate"

	// TLSClientCertificateChainKey is the attribute Key conforming to the
	// "tls.client.certificate_chain" semantic conventions. It represents the
	// array of PEM-encoded certificates that make up the certificate chain
	// offered by the client. This is usually mutually-exclusive of
	// `client.certificate` since that value should be the first certificate in
	// the chain.
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'MII...', 'MI...'
	TLSClientCertificateChainKey = "tls.client.certificate_chain"

	// TLSClientHashMd5Key is the attribute Key conforming to the
	// "tls.client.hash.md5" semantic conventions. It represents the
	// certificate fingerprint using the MD5 digest of DER-encoded version of
	// certificate offered by the client. For consistency with other hash
	// values, this value should be formatted as an uppercase hash.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '0F76C7F2C55BFD7D8E8B8F4BFBF0C9EC'
	TLSClientHashMd5Key = "tls.client.hash.md5"

	// TLSClientHashSha1Key is the attribute Key conforming to the
	// "tls.client.hash.sha1" semantic conventions. It represents the
	// certificate fingerprint using the SHA1 digest of DER-encoded version of
	// certificate offered by the client. For consistency with other hash
	// values, this value should be formatted as an uppercase hash.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '9E393D93138888D288266C2D915214D1D1CCEB2A'
	TLSClientHashSha1Key = "tls.client.hash.sha1"

	// TLSClientHashSha256Key is the attribute Key conforming to the
	// "tls.client.hash.sha256" semantic conventions. It represents the
	// certificate fingerprint using the SHA256 digest of DER-encoded version
	// of certificate offered by the client. For consistency with other hash
	// values, this value should be formatted as an uppercase hash.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples:
	// '0687F666A054EF17A08E2F2162EAB4CBC0D265E1D7875BE74BF3C712CA92DAF0'
	TLSClientHashSha256Key = "tls.client.hash.sha256"

	// TLSClientIssuerKey is the attribute Key conforming to the
	// "tls.client.issuer" semantic conventions. It represents the
	// distinguished name of
	// [subject](https://datatracker.ietf.org/doc/html/rfc5280#section-4.1.2.6)
	// of the issuer of the x.509 certificate presented by the client.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'CN=Example Root CA, OU=Infrastructure Team, DC=example,
	// DC=com'
	TLSClientIssuerKey = "tls.client.issuer"

	// TLSClientJa3Key is the attribute Key conforming to the "tls.client.ja3"
	// semantic conventions. It represents a hash that identifies clients based
	// on how they perform an SSL/TLS handshake.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'd4e5b18d6b55c71272893221c96ba240'
	TLSClientJa3Key = "tls.client.ja3"

	// TLSClientNotAfterKey is the attribute Key conforming to the
	// "tls.client.not_after" semantic conventions. It represents the date/Time
	// indicating when client certificate is no longer considered valid.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '2021-01-01T00:00:00.000Z'
	TLSClientNotAfterKey = "tls.client.not_after"

	// TLSClientNotBeforeKey is the attribute Key conforming to the
	// "tls.client.not_before" semantic conventions. It represents the
	// date/Time indicating when client certificate is first considered valid.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '1970-01-01T00:00:00.000Z'
	TLSClientNotBeforeKey = "tls.client.not_before"

	// TLSClientServerNameKey is the attribute Key conforming to the
	// "tls.client.server_name" semantic conventions. It represents the also
	// called an SNI, this tells the server which hostname to which the client
	// is attempting to connect to.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'opentelemetry.io'
	TLSClientServerNameKey = "tls.client.server_name"

	// TLSClientSubjectKey is the attribute Key conforming to the
	// "tls.client.subject" semantic conventions. It represents the
	// distinguished name of subject of the x.509 certificate presented by the
	// client.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'CN=myclient, OU=Documentation Team, DC=example, DC=com'
	TLSClientSubjectKey = "tls.client.subject"

	// TLSClientSupportedCiphersKey is the attribute Key conforming to the
	// "tls.client.supported_ciphers" semantic conventions. It represents the
	// array of ciphers offered by the client during the client hello.
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
	// "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384", "..."'
	TLSClientSupportedCiphersKey = "tls.client.supported_ciphers"

	// TLSCurveKey is the attribute Key conforming to the "tls.curve" semantic
	// conventions. It represents the string indicating the curve used for the
	// given cipher, when applicable
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'secp256r1'
	TLSCurveKey = "tls.curve"

	// TLSEstablishedKey is the attribute Key conforming to the
	// "tls.established" semantic conventions. It represents the boolean flag
	// indicating if the TLS negotiation was successful and transitioned to an
	// encrypted tunnel.
	//
	// Type: boolean
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: True
	TLSEstablishedKey = "tls.established"

	// TLSNextProtocolKey is the attribute Key conforming to the
	// "tls.next_protocol" semantic conventions. It represents the string
	// indicating the protocol being tunneled. Per the values in the [IANA
	// registry](https://www.iana.org/assignments/tls-extensiontype-values/tls-extensiontype-values.xhtml#alpn-protocol-ids),
	// this string should be lower case.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'http/1.1'
	TLSNextProtocolKey = "tls.next_protocol"

	// TLSProtocolNameKey is the attribute Key conforming to the
	// "tls.protocol.name" semantic conventions. It represents the normalized
	// lowercase protocol name parsed from original string of the negotiated
	// [SSL/TLS protocol
	// version](https://www.openssl.org/docs/man1.1.1/man3/SSL_get_version.html#RETURN-VALUES)
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	TLSProtocolNameKey = "tls.protocol.name"

	// TLSProtocolVersionKey is the attribute Key conforming to the
	// "tls.protocol.version" semantic conventions. It represents the numeric
	// part of the version parsed from the original string of the negotiated
	// [SSL/TLS protocol
	// version](https://www.openssl.org/docs/man1.1.1/man3/SSL_get_version.html#RETURN-VALUES)
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '1.2', '3'
	TLSProtocolVersionKey = "tls.protocol.version"

	// TLSResumedKey is the attribute Key conforming to the "tls.resumed"
	// semantic conventions. It represents the boolean flag indicating if this
	// TLS connection was resumed from an existing TLS negotiation.
	//
	// Type: boolean
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: True
	TLSResumedKey = "tls.resumed"

	// TLSServerCertificateKey is the attribute Key conforming to the
	// "tls.server.certificate" semantic conventions. It represents the
	// pEM-encoded stand-alone certificate offered by the server. This is
	// usually mutually-exclusive of `server.certificate_chain` since this
	// value also exists in that list.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'MII...'
	TLSServerCertificateKey = "tls.server.certificate"

	// TLSServerCertificateChainKey is the attribute Key conforming to the
	// "tls.server.certificate_chain" semantic conventions. It represents the
	// array of PEM-encoded certificates that make up the certificate chain
	// offered by the server. This is usually mutually-exclusive of
	// `server.certificate` since that value should be the first certificate in
	// the chain.
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'MII...', 'MI...'
	TLSServerCertificateChainKey = "tls.server.certificate_chain"

	// TLSServerHashMd5Key is the attribute Key conforming to the
	// "tls.server.hash.md5" semantic conventions. It represents the
	// certificate fingerprint using the MD5 digest of DER-encoded version of
	// certificate offered by the server. For consistency with other hash
	// values, this value should be formatted as an uppercase hash.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '0F76C7F2C55BFD7D8E8B8F4BFBF0C9EC'
	TLSServerHashMd5Key = "tls.server.hash.md5"

	// TLSServerHashSha1Key is the attribute Key conforming to the
	// "tls.server.hash.sha1" semantic conventions. It represents the
	// certificate fingerprint using the SHA1 digest of DER-encoded version of
	// certificate offered by the server. For consistency with other hash
	// values, this value should be formatted as an uppercase hash.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '9E393D93138888D288266C2D915214D1D1CCEB2A'
	TLSServerHashSha1Key = "tls.server.hash.sha1"

	// TLSServerHashSha256Key is the attribute Key conforming to the
	// "tls.server.hash.sha256" semantic conventions. It represents the
	// certificate fingerprint using the SHA256 digest of DER-encoded version
	// of certificate offered by the server. For consistency with other hash
	// values, this value should be formatted as an uppercase hash.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples:
	// '0687F666A054EF17A08E2F2162EAB4CBC0D265E1D7875BE74BF3C712CA92DAF0'
	TLSServerHashSha256Key = "tls.server.hash.sha256"

	// TLSServerIssuerKey is the attribute Key conforming to the
	// "tls.server.issuer" semantic conventions. It represents the
	// distinguished name of
	// [subject](https://datatracker.ietf.org/doc/html/rfc5280#section-4.1.2.6)
	// of the issuer of the x.509 certificate presented by the client.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'CN=Example Root CA, OU=Infrastructure Team, DC=example,
	// DC=com'
	TLSServerIssuerKey = "tls.server.issuer"

	// TLSServerJa3sKey is the attribute Key conforming to the
	// "tls.server.ja3s" semantic conventions. It represents a hash that
	// identifies servers based on how they perform an SSL/TLS handshake.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'd4e5b18d6b55c71272893221c96ba240'
	TLSServerJa3sKey = "tls.server.ja3s"

	// TLSServerNotAfterKey is the attribute Key conforming to the
	// "tls.server.not_after" semantic conventions. It represents the date/Time
	// indicating when server certificate is no longer considered valid.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '2021-01-01T00:00:00.000Z'
	TLSServerNotAfterKey = "tls.server.not_after"

	// TLSServerNotBeforeKey is the attribute Key conforming to the
	// "tls.server.not_before" semantic conventions. It represents the
	// date/Time indicating when server certificate is first considered valid.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '1970-01-01T00:00:00.000Z'
	TLSServerNotBeforeKey = "tls.server.not_before"

	// TLSServerSubjectKey is the attribute Key conforming to the
	// "tls.server.subject" semantic conventions. It represents the
	// distinguished name of subject of the x.509 certificate presented by the
	// server.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'CN=myserver, OU=Documentation Team, DC=example, DC=com'
	TLSServerSubjectKey = "tls.server.subject"

	// URLFragmentKey is the attribute Key conforming to the "url.fragment"
	// semantic conventions. It represents the [URI
	// fragment](https://www.rfc-editor.org/rfc/rfc3986#section-3.5) component
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'SemConv'
	URLFragmentKey = "url.fragment"

	// URLFullKey is the attribute Key conforming to the "url.full" semantic
	// conventions. It represents the absolute URL describing a network
	// resource according to [RFC3986](https://www.rfc-editor.org/rfc/rfc3986)
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'https://www.foo.bar/search?q=OpenTelemetry#SemConv',
	// '//localhost'
	// Note: For network calls, URL usually has
	// `scheme://host[:port][path][?query][#fragment]` format, where the
	// fragment is not transmitted over HTTP, but if it is known, it SHOULD be
	// included nevertheless.
	// `url.full` MUST NOT contain credentials passed via URL in form of
	// `https://username:password@www.example.com/`. In such case username and
	// password SHOULD be redacted and attribute's value SHOULD be
	// `https://REDACTED:REDACTED@www.example.com/`.
	// `url.full` SHOULD capture the absolute URL when it is available (or can
	// be reconstructed) and SHOULD NOT be validated or modified except for
	// sanitizing purposes.
	URLFullKey = "url.full"

	// URLPathKey is the attribute Key conforming to the "url.path" semantic
	// conventions. It represents the [URI
	// path](https://www.rfc-editor.org/rfc/rfc3986#section-3.3) component
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: '/search'
	URLPathKey = "url.path"

	// URLQueryKey is the attribute Key conforming to the "url.query" semantic
	// conventions. It represents the [URI
	// query](https://www.rfc-editor.org/rfc/rfc3986#section-3.4) component
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'q=OpenTelemetry'
	// Note: Sensitive content provided in query string SHOULD be scrubbed when
	// instrumentations can identify it.
	URLQueryKey = "url.query"

	// URLSchemeKey is the attribute Key conforming to the "url.scheme"
	// semantic conventions. It represents the [URI
	// scheme](https://www.rfc-editor.org/rfc/rfc3986#section-3.1) component
	// identifying the used protocol.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'https', 'ftp', 'telnet'
	URLSchemeKey = "url.scheme"

	// UserAgentOriginalKey is the attribute Key conforming to the
	// "user_agent.original" semantic conventions. It represents the value of
	// the [HTTP
	// User-Agent](https://www.rfc-editor.org/rfc/rfc9110.html#field.user-agent)
	// header sent by the client.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'CERN-LineMode/2.15 libwww/2.17b3', 'Mozilla/5.0 (iPhone; CPU
	// iPhone OS 14_7_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko)
	// Version/14.1.2 Mobile/15E148 Safari/604.1'
	UserAgentOriginalKey = "user_agent.original"

	// SessionIDKey is the attribute Key conforming to the "session.id"
	// semantic conventions. It represents a unique id to identify a session.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '00112233-4455-6677-8899-aabbccddeeff'
	SessionIDKey = "session.id"

	// SessionPreviousIDKey is the attribute Key conforming to the
	// "session.previous_id" semantic conventions. It represents the previous
	// `session.id` for this user, when known.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '00112233-4455-6677-8899-aabbccddeeff'
	SessionPreviousIDKey = "session.previous_id"

	// IosStateKey is the attribute Key conforming to the "ios.state" semantic
	// conventions. It represents the this attribute represents the state the
	// application has transitioned into at the occurrence of the event.
	//
	// Type: Enum
	// RequirementLevel: Required
	// Stability: experimental
	// Note: The iOS lifecycle states are defined in the [UIApplicationDelegate
	// documentation](https://developer.apple.com/documentation/uikit/uiapplicationdelegate#1656902),
	// and from which the `OS terminology` column values are derived.
	IosStateKey = "ios.state"

	// AndroidStateKey is the attribute Key conforming to the "android.state"
	// semantic conventions. It represents the this attribute represents the
	// state the application has transitioned into at the occurrence of the
	// event.
	//
	// Type: Enum
	// RequirementLevel: Required
	// Stability: experimental
	// Note: The Android lifecycle states are defined in [Activity lifecycle
	// callbacks](https://developer.android.com/guide/components/activities/activity-lifecycle#lc),
	// and from which the `OS identifiers` are derived.
	AndroidStateKey = "android.state"

	// FeatureFlagKeyKey is the attribute Key conforming to the
	// "feature_flag.key" semantic conventions. It represents the unique
	// identifier of the feature flag.
	//
	// Type: string
	// RequirementLevel: Required
	// Stability: experimental
	// Examples: 'logo-color'
	FeatureFlagKeyKey = "feature_flag.key"

	// FeatureFlagProviderNameKey is the attribute Key conforming to the
	// "feature_flag.provider_name" semantic conventions. It represents the
	// name of the service provider that performs the flag evaluation.
	//
	// Type: string
	// RequirementLevel: Recommended
	// Stability: experimental
	// Examples: 'Flag Manager'
	FeatureFlagProviderNameKey = "feature_flag.provider_name"

	// FeatureFlagVariantKey is the attribute Key conforming to the
	// "feature_flag.variant" semantic conventions. It represents the sHOULD be
	// a semantic identifier for a value. If one is unavailable, a stringified
	// version of the value can be used.
	//
	// Type: string
	// RequirementLevel: Recommended
	// Stability: experimental
	// Examples: 'red', 'true', 'on'
	// Note: A semantic identifier, commonly referred to as a variant, provides
	// a means
	// for referring to a value without including the value itself. This can
	// provide additional context for understanding the meaning behind a value.
	// For example, the variant `red` maybe be used for the value `#c05543`.
	//
	// A stringified version of the value can be used in situations where a
	// semantic identifier is unavailable. String representation of the value
	// should be determined by the implementer.
	FeatureFlagVariantKey = "feature_flag.variant"

	// MessageCompressedSizeKey is the attribute Key conforming to the
	// "message.compressed_size" semantic conventions. It represents the
	// compressed size of the message in bytes.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	MessageCompressedSizeKey = "message.compressed_size"

	// MessageIDKey is the attribute Key conforming to the "message.id"
	// semantic conventions. It represents the mUST be calculated as two
	// different counters starting from `1` one for sent messages and one for
	// received message.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Note: This way we guarantee that the values will be consistent between
	// different implementations.
	MessageIDKey = "message.id"

	// MessageTypeKey is the attribute Key conforming to the "message.type"
	// semantic conventions. It represents the whether this is a received or
	// sent message.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	MessageTypeKey = "message.type"

	// MessageUncompressedSizeKey is the attribute Key conforming to the
	// "message.uncompressed_size" semantic conventions. It represents the
	// uncompressed size of the message in bytes.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	MessageUncompressedSizeKey = "message.uncompressed_size"

	// CloudAccountIDKey is the attribute Key conforming to the
	// "cloud.account.id" semantic conventions. It represents the cloud account
	// ID the resource is assigned to.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '111111111111', 'opentelemetry'
	CloudAccountIDKey = "cloud.account.id"

	// CloudAvailabilityZoneKey is the attribute Key conforming to the
	// "cloud.availability_zone" semantic conventions. It represents the cloud
	// regions often have multiple, isolated locations known as zones to
	// increase availability. Availability zone represents the zone where the
	// resource is running.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'us-east-1c'
	// Note: Availability zones are called "zones" on Alibaba Cloud and Google
	// Cloud.
	CloudAvailabilityZoneKey = "cloud.availability_zone"

	// CloudPlatformKey is the attribute Key conforming to the "cloud.platform"
	// semantic conventions. It represents the cloud platform in use.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	// Note: The prefix of the service SHOULD match the one specified in
	// `cloud.provider`.
	CloudPlatformKey = "cloud.platform"

	// CloudProviderKey is the attribute Key conforming to the "cloud.provider"
	// semantic conventions. It represents the name of the cloud provider.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	CloudProviderKey = "cloud.provider"

	// CloudRegionKey is the attribute Key conforming to the "cloud.region"
	// semantic conventions. It represents the geographical region the resource
	// is running.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'us-central1', 'us-east-1'
	// Note: Refer to your provider's docs to see the available regions, for
	// example [Alibaba Cloud
	// regions](https://www.alibabacloud.com/help/doc-detail/40654.htm), [AWS
	// regions](https://aws.amazon.com/about-aws/global-infrastructure/regions_az/),
	// [Azure
	// regions](https://azure.microsoft.com/global-infrastructure/geographies/),
	// [Google Cloud regions](https://cloud.google.com/about/locations), or
	// [Tencent Cloud
	// regions](https://www.tencentcloud.com/document/product/213/6091).
	CloudRegionKey = "cloud.region"

	// CloudResourceIDKey is the attribute Key conforming to the
	// "cloud.resource_id" semantic conventions. It represents the cloud
	// provider-specific native identifier of the monitored cloud resource
	// (e.g. an
	// [ARN](https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html)
	// on AWS, a [fully qualified resource
	// ID](https://learn.microsoft.com/rest/api/resources/resources/get-by-id)
	// on Azure, a [full resource
	// name](https://cloud.google.com/apis/design/resource_names#full_resource_name)
	// on GCP)
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'arn:aws:lambda:REGION:ACCOUNT_ID:function:my-function',
	// '//run.googleapis.com/projects/PROJECT_ID/locations/LOCATION_ID/services/SERVICE_ID',
	// '/subscriptions/<SUBSCIPTION_GUID>/resourceGroups/<RG>/providers/Microsoft.Web/sites/<FUNCAPP>/functions/<FUNC>'
	// Note: On some cloud providers, it may not be possible to determine the
	// full ID at startup,
	// so it may be necessary to set `cloud.resource_id` as a span attribute
	// instead.
	//
	// The exact value to use for `cloud.resource_id` depends on the cloud
	// provider.
	// The following well-known definitions MUST be used if you set this
	// attribute and they apply:
	//
	// * **AWS Lambda:** The function
	// [ARN](https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html).
	//   Take care not to use the "invoked ARN" directly but replace any
	//   [alias
	// suffix](https://docs.aws.amazon.com/lambda/latest/dg/configuration-aliases.html)
	//   with the resolved function version, as the same runtime instance may
	// be invokable with
	//   multiple different aliases.
	// * **GCP:** The [URI of the
	// resource](https://cloud.google.com/iam/docs/full-resource-names)
	// * **Azure:** The [Fully Qualified Resource
	// ID](https://docs.microsoft.com/rest/api/resources/resources/get-by-id)
	// of the invoked function,
	//   *not* the function app, having the form
	// `/subscriptions/<SUBSCIPTION_GUID>/resourceGroups/<RG>/providers/Microsoft.Web/sites/<FUNCAPP>/functions/<FUNC>`.
	//   This means that a span attribute MUST be used, as an Azure function
	// app can host multiple functions that would usually share
	//   a TracerProvider.
	CloudResourceIDKey = "cloud.resource_id"

	// ContainerCommandKey is the attribute Key conforming to the
	// "container.command" semantic conventions. It represents the command used
	// to run the container (i.e. the command name).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'otelcontribcol'
	// Note: If using embedded credentials or sensitive data, it is recommended
	// to remove them to prevent potential leakage.
	ContainerCommandKey = "container.command"

	// ContainerCommandArgsKey is the attribute Key conforming to the
	// "container.command_args" semantic conventions. It represents the all the
	// command arguments (including the command/executable itself) run by the
	// container. [2]
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'otelcontribcol, --config, config.yaml'
	ContainerCommandArgsKey = "container.command_args"

	// ContainerCommandLineKey is the attribute Key conforming to the
	// "container.command_line" semantic conventions. It represents the full
	// command run by the container as a single string representing the full
	// command. [2]
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'otelcontribcol --config config.yaml'
	ContainerCommandLineKey = "container.command_line"

	// ContainerIDKey is the attribute Key conforming to the "container.id"
	// semantic conventions. It represents the container ID. Usually a UUID, as
	// for example used to [identify Docker
	// containers](https://docs.docker.com/engine/reference/run/#container-identification).
	// The UUID might be abbreviated.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'a3bf90e006b2'
	ContainerIDKey = "container.id"

	// ContainerImageIDKey is the attribute Key conforming to the
	// "container.image.id" semantic conventions. It represents the runtime
	// specific image identifier. Usually a hash algorithm followed by a UUID.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples:
	// 'sha256:19c92d0a00d1b66d897bceaa7319bee0dd38a10a851c60bcec9474aa3f01e50f'
	// Note: Docker defines a sha256 of the image id; `container.image.id`
	// corresponds to the `Image` field from the Docker container inspect
	// [API](https://docs.docker.com/engine/api/v1.43/#tag/Container/operation/ContainerInspect)
	// endpoint.
	// K8S defines a link to the container registry repository with digest
	// `"imageID": "registry.azurecr.io
	// /namespace/service/dockerfile@sha256:bdeabd40c3a8a492eaf9e8e44d0ebbb84bac7ee25ac0cf8a7159d25f62555625"`.
	// The ID is assinged by the container runtime and can vary in different
	// environments. Consider using `oci.manifest.digest` if it is important to
	// identify the same image in different environments/runtimes.
	ContainerImageIDKey = "container.image.id"

	// ContainerImageNameKey is the attribute Key conforming to the
	// "container.image.name" semantic conventions. It represents the name of
	// the image the container was built on.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'gcr.io/opentelemetry/operator'
	ContainerImageNameKey = "container.image.name"

	// ContainerImageRepoDigestsKey is the attribute Key conforming to the
	// "container.image.repo_digests" semantic conventions. It represents the
	// repo digests of the container image as provided by the container
	// runtime.
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples:
	// 'example@sha256:afcc7f1ac1b49db317a7196c902e61c6c3c4607d63599ee1a82d702d249a0ccb',
	// 'internal.registry.example.com:5000/example@sha256:b69959407d21e8a062e0416bf13405bb2b71ed7a84dde4158ebafacfa06f5578'
	// Note:
	// [Docker](https://docs.docker.com/engine/api/v1.43/#tag/Image/operation/ImageInspect)
	// and
	// [CRI](https://github.com/kubernetes/cri-api/blob/c75ef5b473bbe2d0a4fc92f82235efd665ea8e9f/pkg/apis/runtime/v1/api.proto#L1237-L1238)
	// report those under the `RepoDigests` field.
	ContainerImageRepoDigestsKey = "container.image.repo_digests"

	// ContainerImageTagsKey is the attribute Key conforming to the
	// "container.image.tags" semantic conventions. It represents the container
	// image tags. An example can be found in [Docker Image
	// Inspect](https://docs.docker.com/engine/api/v1.43/#tag/Image/operation/ImageInspect).
	// Should be only the `<tag>` section of the full name for example from
	// `registry.example.com/my-org/my-image:<tag>`.
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'v1.27.1', '3.5.7-0'
	ContainerImageTagsKey = "container.image.tags"

	// ContainerNameKey is the attribute Key conforming to the "container.name"
	// semantic conventions. It represents the container name used by container
	// runtime.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'opentelemetry-autoconf'
	ContainerNameKey = "container.name"

	// ContainerRuntimeKey is the attribute Key conforming to the
	// "container.runtime" semantic conventions. It represents the container
	// runtime managing this container.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'docker', 'containerd', 'rkt'
	ContainerRuntimeKey = "container.runtime"

	// DeviceIDKey is the attribute Key conforming to the "device.id" semantic
	// conventions. It represents a unique identifier representing the device
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '2ab2916d-a51f-4ac8-80ee-45ac31a28092'
	// Note: The device identifier MUST only be defined using the values
	// outlined below. This value is not an advertising identifier and MUST NOT
	// be used as such. On iOS (Swift or Objective-C), this value MUST be equal
	// to the [vendor
	// identifier](https://developer.apple.com/documentation/uikit/uidevice/1620059-identifierforvendor).
	// On Android (Java or Kotlin), this value MUST be equal to the Firebase
	// Installation ID or a globally unique UUID which is persisted across
	// sessions in your application. More information can be found
	// [here](https://developer.android.com/training/articles/user-data-ids) on
	// best practices and exact implementation details. Caution should be taken
	// when storing personal data or anything which can identify a user. GDPR
	// and data protection laws may apply, ensure you do your own due
	// diligence.
	DeviceIDKey = "device.id"

	// DeviceManufacturerKey is the attribute Key conforming to the
	// "device.manufacturer" semantic conventions. It represents the name of
	// the device manufacturer
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'Apple', 'Samsung'
	// Note: The Android OS provides this field via
	// [Build](https://developer.android.com/reference/android/os/Build#MANUFACTURER).
	// iOS apps SHOULD hardcode the value `Apple`.
	DeviceManufacturerKey = "device.manufacturer"

	// DeviceModelIdentifierKey is the attribute Key conforming to the
	// "device.model.identifier" semantic conventions. It represents the model
	// identifier for the device
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'iPhone3,4', 'SM-G920F'
	// Note: It's recommended this value represents a machine-readable version
	// of the model identifier rather than the market or consumer-friendly name
	// of the device.
	DeviceModelIdentifierKey = "device.model.identifier"

	// DeviceModelNameKey is the attribute Key conforming to the
	// "device.model.name" semantic conventions. It represents the marketing
	// name for the device model
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'iPhone 6s Plus', 'Samsung Galaxy S6'
	// Note: It's recommended this value represents a human-readable version of
	// the device model rather than a machine-readable alternative.
	DeviceModelNameKey = "device.model.name"

	// HostArchKey is the attribute Key conforming to the "host.arch" semantic
	// conventions. It represents the CPU architecture the host system is
	// running on.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	HostArchKey = "host.arch"

	// HostCPUCacheL2SizeKey is the attribute Key conforming to the
	// "host.cpu.cache.l2.size" semantic conventions. It represents the amount
	// of level 2 memory cache available to the processor (in Bytes).
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 12288000
	HostCPUCacheL2SizeKey = "host.cpu.cache.l2.size"

	// HostCPUFamilyKey is the attribute Key conforming to the
	// "host.cpu.family" semantic conventions. It represents the family or
	// generation of the CPU.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '6', 'PA-RISC 1.1e'
	HostCPUFamilyKey = "host.cpu.family"

	// HostCPUModelIDKey is the attribute Key conforming to the
	// "host.cpu.model.id" semantic conventions. It represents the model
	// identifier. It provides more granular information about the CPU,
	// distinguishing it from other CPUs within the same family.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '6', '9000/778/B180L'
	HostCPUModelIDKey = "host.cpu.model.id"

	// HostCPUModelNameKey is the attribute Key conforming to the
	// "host.cpu.model.name" semantic conventions. It represents the model
	// designation of the processor.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '11th Gen Intel(R) Core(TM) i7-1185G7 @ 3.00GHz'
	HostCPUModelNameKey = "host.cpu.model.name"

	// HostCPUSteppingKey is the attribute Key conforming to the
	// "host.cpu.stepping" semantic conventions. It represents the stepping or
	// core revisions.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 1
	HostCPUSteppingKey = "host.cpu.stepping"

	// HostCPUVendorIDKey is the attribute Key conforming to the
	// "host.cpu.vendor.id" semantic conventions. It represents the processor
	// manufacturer identifier. A maximum 12-character string.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'GenuineIntel'
	// Note: [CPUID](https://wiki.osdev.org/CPUID) command returns the vendor
	// ID string in EBX, EDX and ECX registers. Writing these to memory in this
	// order results in a 12-character string.
	HostCPUVendorIDKey = "host.cpu.vendor.id"

	// HostIDKey is the attribute Key conforming to the "host.id" semantic
	// conventions. It represents the unique host ID. For Cloud, this must be
	// the instance_id assigned by the cloud provider. For non-containerized
	// systems, this should be the `machine-id`. See the table below for the
	// sources to use to determine the `machine-id` based on operating system.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'fdbf79e8af94cb7f9e8df36789187052'
	HostIDKey = "host.id"

	// HostImageIDKey is the attribute Key conforming to the "host.image.id"
	// semantic conventions. It represents the vM image ID or host OS image ID.
	// For Cloud, this value is from the provider.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'ami-07b06b442921831e5'
	HostImageIDKey = "host.image.id"

	// HostImageNameKey is the attribute Key conforming to the
	// "host.image.name" semantic conventions. It represents the name of the VM
	// image or OS install the host was instantiated from.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'infra-ami-eks-worker-node-7d4ec78312', 'CentOS-8-x86_64-1905'
	HostImageNameKey = "host.image.name"

	// HostImageVersionKey is the attribute Key conforming to the
	// "host.image.version" semantic conventions. It represents the version
	// string of the VM image or host OS as defined in [Version
	// Attributes](/docs/resource/README.md#version-attributes).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '0.1'
	HostImageVersionKey = "host.image.version"

	// HostIPKey is the attribute Key conforming to the "host.ip" semantic
	// conventions. It represents the available IP addresses of the host,
	// excluding loopback interfaces.
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '192.168.1.140', 'fe80::abc2:4a28:737a:609e'
	// Note: IPv4 Addresses MUST be specified in dotted-quad notation. IPv6
	// addresses MUST be specified in the [RFC
	// 5952](https://www.rfc-editor.org/rfc/rfc5952.html) format.
	HostIPKey = "host.ip"

	// HostMacKey is the attribute Key conforming to the "host.mac" semantic
	// conventions. It represents the available MAC addresses of the host,
	// excluding loopback interfaces.
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'AC-DE-48-23-45-67', 'AC-DE-48-23-45-67-01-9F'
	// Note: MAC Addresses MUST be represented in [IEEE RA hexadecimal
	// form](https://standards.ieee.org/wp-content/uploads/import/documents/tutorials/eui.pdf):
	// as hyphen-separated octets in uppercase hexadecimal form from most to
	// least significant.
	HostMacKey = "host.mac"

	// HostNameKey is the attribute Key conforming to the "host.name" semantic
	// conventions. It represents the name of the host. On Unix systems, it may
	// contain what the hostname command returns, or the fully qualified
	// hostname, or another name specified by the user.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'opentelemetry-test'
	HostNameKey = "host.name"

	// HostTypeKey is the attribute Key conforming to the "host.type" semantic
	// conventions. It represents the type of host. For Cloud, this must be the
	// machine type.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'n1-standard-1'
	HostTypeKey = "host.type"

	// K8SClusterNameKey is the attribute Key conforming to the
	// "k8s.cluster.name" semantic conventions. It represents the name of the
	// cluster.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'opentelemetry-cluster'
	K8SClusterNameKey = "k8s.cluster.name"

	// K8SClusterUIDKey is the attribute Key conforming to the
	// "k8s.cluster.uid" semantic conventions. It represents a pseudo-ID for
	// the cluster, set to the UID of the `kube-system` namespace.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '218fc5a9-a5f1-4b54-aa05-46717d0ab26d'
	// Note: K8S doesn't have support for obtaining a cluster ID. If this is
	// ever
	// added, we will recommend collecting the `k8s.cluster.uid` through the
	// official APIs. In the meantime, we are able to use the `uid` of the
	// `kube-system` namespace as a proxy for cluster ID. Read on for the
	// rationale.
	//
	// Every object created in a K8S cluster is assigned a distinct UID. The
	// `kube-system` namespace is used by Kubernetes itself and will exist
	// for the lifetime of the cluster. Using the `uid` of the `kube-system`
	// namespace is a reasonable proxy for the K8S ClusterID as it will only
	// change if the cluster is rebuilt. Furthermore, Kubernetes UIDs are
	// UUIDs as standardized by
	// [ISO/IEC 9834-8 and ITU-T
	// X.667](https://www.itu.int/ITU-T/studygroups/com17/oid.html).
	// Which states:
	//
	// > If generated according to one of the mechanisms defined in Rec.
	//   ITU-T X.667 | ISO/IEC 9834-8, a UUID is either guaranteed to be
	//   different from all other UUIDs generated before 3603 A.D., or is
	//   extremely likely to be different (depending on the mechanism chosen).
	//
	// Therefore, UIDs between clusters should be extremely unlikely to
	// conflict.
	K8SClusterUIDKey = "k8s.cluster.uid"

	// K8SContainerNameKey is the attribute Key conforming to the
	// "k8s.container.name" semantic conventions. It represents the name of the
	// Container from Pod specification, must be unique within a Pod. Container
	// runtime usually uses different globally unique name (`container.name`).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'redis'
	K8SContainerNameKey = "k8s.container.name"

	// K8SContainerRestartCountKey is the attribute Key conforming to the
	// "k8s.container.restart_count" semantic conventions. It represents the
	// number of times the container was restarted. This attribute can be used
	// to identify a particular container (running or stopped) within a
	// container spec.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 0, 2
	K8SContainerRestartCountKey = "k8s.container.restart_count"

	// K8SCronJobNameKey is the attribute Key conforming to the
	// "k8s.cronjob.name" semantic conventions. It represents the name of the
	// CronJob.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'opentelemetry'
	K8SCronJobNameKey = "k8s.cronjob.name"

	// K8SCronJobUIDKey is the attribute Key conforming to the
	// "k8s.cronjob.uid" semantic conventions. It represents the UID of the
	// CronJob.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '275ecb36-5aa8-4c2a-9c47-d8bb681b9aff'
	K8SCronJobUIDKey = "k8s.cronjob.uid"

	// K8SDaemonSetNameKey is the attribute Key conforming to the
	// "k8s.daemonset.name" semantic conventions. It represents the name of the
	// DaemonSet.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'opentelemetry'
	K8SDaemonSetNameKey = "k8s.daemonset.name"

	// K8SDaemonSetUIDKey is the attribute Key conforming to the
	// "k8s.daemonset.uid" semantic conventions. It represents the UID of the
	// DaemonSet.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '275ecb36-5aa8-4c2a-9c47-d8bb681b9aff'
	K8SDaemonSetUIDKey = "k8s.daemonset.uid"

	// K8SDeploymentNameKey is the attribute Key conforming to the
	// "k8s.deployment.name" semantic conventions. It represents the name of
	// the Deployment.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'opentelemetry'
	K8SDeploymentNameKey = "k8s.deployment.name"

	// K8SDeploymentUIDKey is the attribute Key conforming to the
	// "k8s.deployment.uid" semantic conventions. It represents the UID of the
	// Deployment.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '275ecb36-5aa8-4c2a-9c47-d8bb681b9aff'
	K8SDeploymentUIDKey = "k8s.deployment.uid"

	// K8SJobNameKey is the attribute Key conforming to the "k8s.job.name"
	// semantic conventions. It represents the name of the Job.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'opentelemetry'
	K8SJobNameKey = "k8s.job.name"

	// K8SJobUIDKey is the attribute Key conforming to the "k8s.job.uid"
	// semantic conventions. It represents the UID of the Job.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '275ecb36-5aa8-4c2a-9c47-d8bb681b9aff'
	K8SJobUIDKey = "k8s.job.uid"

	// K8SNamespaceNameKey is the attribute Key conforming to the
	// "k8s.namespace.name" semantic conventions. It represents the name of the
	// namespace that the pod is running in.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'default'
	K8SNamespaceNameKey = "k8s.namespace.name"

	// K8SNodeNameKey is the attribute Key conforming to the "k8s.node.name"
	// semantic conventions. It represents the name of the Node.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'node-1'
	K8SNodeNameKey = "k8s.node.name"

	// K8SNodeUIDKey is the attribute Key conforming to the "k8s.node.uid"
	// semantic conventions. It represents the UID of the Node.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '1eb3a0c6-0477-4080-a9cb-0cb7db65c6a2'
	K8SNodeUIDKey = "k8s.node.uid"

	// K8SPodNameKey is the attribute Key conforming to the "k8s.pod.name"
	// semantic conventions. It represents the name of the Pod.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'opentelemetry-pod-autoconf'
	K8SPodNameKey = "k8s.pod.name"

	// K8SPodUIDKey is the attribute Key conforming to the "k8s.pod.uid"
	// semantic conventions. It represents the UID of the Pod.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '275ecb36-5aa8-4c2a-9c47-d8bb681b9aff'
	K8SPodUIDKey = "k8s.pod.uid"

	// K8SReplicaSetNameKey is the attribute Key conforming to the
	// "k8s.replicaset.name" semantic conventions. It represents the name of
	// the ReplicaSet.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'opentelemetry'
	K8SReplicaSetNameKey = "k8s.replicaset.name"

	// K8SReplicaSetUIDKey is the attribute Key conforming to the
	// "k8s.replicaset.uid" semantic conventions. It represents the UID of the
	// ReplicaSet.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '275ecb36-5aa8-4c2a-9c47-d8bb681b9aff'
	K8SReplicaSetUIDKey = "k8s.replicaset.uid"

	// K8SStatefulSetNameKey is the attribute Key conforming to the
	// "k8s.statefulset.name" semantic conventions. It represents the name of
	// the StatefulSet.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'opentelemetry'
	K8SStatefulSetNameKey = "k8s.statefulset.name"

	// K8SStatefulSetUIDKey is the attribute Key conforming to the
	// "k8s.statefulset.uid" semantic conventions. It represents the UID of the
	// StatefulSet.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '275ecb36-5aa8-4c2a-9c47-d8bb681b9aff'
	K8SStatefulSetUIDKey = "k8s.statefulset.uid"

	// OciManifestDigestKey is the attribute Key conforming to the
	// "oci.manifest.digest" semantic conventions. It represents the digest of
	// the OCI image manifest. For container images specifically is the digest
	// by which the container image is known.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples:
	// 'sha256:e4ca62c0d62f3e886e684806dfe9d4e0cda60d54986898173c1083856cfda0f4'
	// Note: Follows [OCI Image Manifest
	// Specification](https://github.com/opencontainers/image-spec/blob/main/manifest.md),
	// and specifically the [Digest
	// property](https://github.com/opencontainers/image-spec/blob/main/descriptor.md#digests).
	// An example can be found in [Example Image
	// Manifest](https://docs.docker.com/registry/spec/manifest-v2-2/#example-image-manifest).
	OciManifestDigestKey = "oci.manifest.digest"

	// OSBuildIDKey is the attribute Key conforming to the "os.build_id"
	// semantic conventions. It represents the unique identifier for a
	// particular build or compilation of the operating system.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'TQ3C.230805.001.B2', '20E247', '22621'
	OSBuildIDKey = "os.build_id"

	// OSDescriptionKey is the attribute Key conforming to the "os.description"
	// semantic conventions. It represents the human readable (not intended to
	// be parsed) OS version information, like e.g. reported by `ver` or
	// `lsb_release -a` commands.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'Microsoft Windows [Version 10.0.18363.778]', 'Ubuntu 18.04.1
	// LTS'
	OSDescriptionKey = "os.description"

	// OSNameKey is the attribute Key conforming to the "os.name" semantic
	// conventions. It represents the human readable operating system name.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'iOS', 'Android', 'Ubuntu'
	OSNameKey = "os.name"

	// OSTypeKey is the attribute Key conforming to the "os.type" semantic
	// conventions. It represents the operating system type.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	OSTypeKey = "os.type"

	// OSVersionKey is the attribute Key conforming to the "os.version"
	// semantic conventions. It represents the version string of the operating
	// system as defined in [Version
	// Attributes](/docs/resource/README.md#version-attributes).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '14.2.1', '18.04.1'
	OSVersionKey = "os.version"

	// ProcessCommandKey is the attribute Key conforming to the
	// "process.command" semantic conventions. It represents the command used
	// to launch the process (i.e. the command name). On Linux based systems,
	// can be set to the zeroth string in `proc/[pid]/cmdline`. On Windows, can
	// be set to the first parameter extracted from `GetCommandLineW`.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'cmd/otelcol'
	ProcessCommandKey = "process.command"

	// ProcessCommandArgsKey is the attribute Key conforming to the
	// "process.command_args" semantic conventions. It represents the all the
	// command arguments (including the command/executable itself) as received
	// by the process. On Linux-based systems (and some other Unixoid systems
	// supporting procfs), can be set according to the list of null-delimited
	// strings extracted from `proc/[pid]/cmdline`. For libc-based executables,
	// this would be the full argv vector passed to `main`.
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'cmd/otecol', '--config=config.yaml'
	ProcessCommandArgsKey = "process.command_args"

	// ProcessCommandLineKey is the attribute Key conforming to the
	// "process.command_line" semantic conventions. It represents the full
	// command used to launch the process as a single string representing the
	// full command. On Windows, can be set to the result of `GetCommandLineW`.
	// Do not set this if you have to assemble it just for monitoring; use
	// `process.command_args` instead.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'C:\\cmd\\otecol --config="my directory\\config.yaml"'
	ProcessCommandLineKey = "process.command_line"

	// ProcessExecutableNameKey is the attribute Key conforming to the
	// "process.executable.name" semantic conventions. It represents the name
	// of the process executable. On Linux based systems, can be set to the
	// `Name` in `proc/[pid]/status`. On Windows, can be set to the base name
	// of `GetProcessImageFileNameW`.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'otelcol'
	ProcessExecutableNameKey = "process.executable.name"

	// ProcessExecutablePathKey is the attribute Key conforming to the
	// "process.executable.path" semantic conventions. It represents the full
	// path to the process executable. On Linux based systems, can be set to
	// the target of `proc/[pid]/exe`. On Windows, can be set to the result of
	// `GetProcessImageFileNameW`.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '/usr/bin/cmd/otelcol'
	ProcessExecutablePathKey = "process.executable.path"

	// ProcessOwnerKey is the attribute Key conforming to the "process.owner"
	// semantic conventions. It represents the username of the user that owns
	// the process.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'root'
	ProcessOwnerKey = "process.owner"

	// ProcessParentPIDKey is the attribute Key conforming to the
	// "process.parent_pid" semantic conventions. It represents the parent
	// Process identifier (PPID).
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 111
	ProcessParentPIDKey = "process.parent_pid"

	// ProcessPIDKey is the attribute Key conforming to the "process.pid"
	// semantic conventions. It represents the process identifier (PID).
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 1234
	ProcessPIDKey = "process.pid"

	// ProcessRuntimeDescriptionKey is the attribute Key conforming to the
	// "process.runtime.description" semantic conventions. It represents an
	// additional description about the runtime of the process, for example a
	// specific vendor customization of the runtime environment.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'Eclipse OpenJ9 Eclipse OpenJ9 VM openj9-0.21.0'
	ProcessRuntimeDescriptionKey = "process.runtime.description"

	// ProcessRuntimeNameKey is the attribute Key conforming to the
	// "process.runtime.name" semantic conventions. It represents the name of
	// the runtime of this process. For compiled native binaries, this SHOULD
	// be the name of the compiler.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'OpenJDK Runtime Environment'
	ProcessRuntimeNameKey = "process.runtime.name"

	// ProcessRuntimeVersionKey is the attribute Key conforming to the
	// "process.runtime.version" semantic conventions. It represents the
	// version of the runtime of this process, as returned by the runtime
	// without modification.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '14.0.2'
	ProcessRuntimeVersionKey = "process.runtime.version"

	// AndroidOSAPILevelKey is the attribute Key conforming to the
	// "android.os.api_level" semantic conventions. It represents the uniquely
	// identifies the framework API revision offered by a version
	// (`os.version`) of the android operating system. More information can be
	// found
	// [here](https://developer.android.com/guide/topics/manifest/uses-sdk-element#APILevels).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '33', '32'
	AndroidOSAPILevelKey = "android.os.api_level"

	// BrowserBrandsKey is the attribute Key conforming to the "browser.brands"
	// semantic conventions. It represents the array of brand name and version
	// separated by a space
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: ' Not A;Brand 99', 'Chromium 99', 'Chrome 99'
	// Note: This value is intended to be taken from the [UA client hints
	// API](https://wicg.github.io/ua-client-hints/#interface)
	// (`navigator.userAgentData.brands`).
	BrowserBrandsKey = "browser.brands"

	// BrowserLanguageKey is the attribute Key conforming to the
	// "browser.language" semantic conventions. It represents the preferred
	// language of the user using the browser
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'en', 'en-US', 'fr', 'fr-FR'
	// Note: This value is intended to be taken from the Navigator API
	// `navigator.language`.
	BrowserLanguageKey = "browser.language"

	// BrowserMobileKey is the attribute Key conforming to the "browser.mobile"
	// semantic conventions. It represents a boolean that is true if the
	// browser is running on a mobile device
	//
	// Type: boolean
	// RequirementLevel: Optional
	// Stability: experimental
	// Note: This value is intended to be taken from the [UA client hints
	// API](https://wicg.github.io/ua-client-hints/#interface)
	// (`navigator.userAgentData.mobile`). If unavailable, this attribute
	// SHOULD be left unset.
	BrowserMobileKey = "browser.mobile"

	// BrowserPlatformKey is the attribute Key conforming to the
	// "browser.platform" semantic conventions. It represents the platform on
	// which the browser is running
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'Windows', 'macOS', 'Android'
	// Note: This value is intended to be taken from the [UA client hints
	// API](https://wicg.github.io/ua-client-hints/#interface)
	// (`navigator.userAgentData.platform`). If unavailable, the legacy
	// `navigator.platform` API SHOULD NOT be used instead and this attribute
	// SHOULD be left unset in order for the values to be consistent.
	// The list of possible values is defined in the [W3C User-Agent Client
	// Hints
	// specification](https://wicg.github.io/ua-client-hints/#sec-ch-ua-platform).
	// Note that some (but not all) of these values can overlap with values in
	// the [`os.type` and `os.name` attributes](./os.md). However, for
	// consistency, the values in the `browser.platform` attribute should
	// capture the exact value that the user agent provides.
	BrowserPlatformKey = "browser.platform"

	// AWSECSClusterARNKey is the attribute Key conforming to the
	// "aws.ecs.cluster.arn" semantic conventions. It represents the ARN of an
	// [ECS
	// cluster](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/clusters.html).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'arn:aws:ecs:us-west-2:123456789123:cluster/my-cluster'
	AWSECSClusterARNKey = "aws.ecs.cluster.arn"

	// AWSECSContainerARNKey is the attribute Key conforming to the
	// "aws.ecs.container.arn" semantic conventions. It represents the Amazon
	// Resource Name (ARN) of an [ECS container
	// instance](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ECS_instances.html).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples:
	// 'arn:aws:ecs:us-west-1:123456789123:container/32624152-9086-4f0e-acae-1a75b14fe4d9'
	AWSECSContainerARNKey = "aws.ecs.container.arn"

	// AWSECSLaunchtypeKey is the attribute Key conforming to the
	// "aws.ecs.launchtype" semantic conventions. It represents the [launch
	// type](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/launch_types.html)
	// for an ECS task.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	AWSECSLaunchtypeKey = "aws.ecs.launchtype"

	// AWSECSTaskARNKey is the attribute Key conforming to the
	// "aws.ecs.task.arn" semantic conventions. It represents the ARN of an
	// [ECS task
	// definition](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definitions.html).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples:
	// 'arn:aws:ecs:us-west-1:123456789123:task/10838bed-421f-43ef-870a-f43feacbbb5b'
	AWSECSTaskARNKey = "aws.ecs.task.arn"

	// AWSECSTaskFamilyKey is the attribute Key conforming to the
	// "aws.ecs.task.family" semantic conventions. It represents the task
	// definition family this task definition is a member of.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'opentelemetry-family'
	AWSECSTaskFamilyKey = "aws.ecs.task.family"

	// AWSECSTaskRevisionKey is the attribute Key conforming to the
	// "aws.ecs.task.revision" semantic conventions. It represents the revision
	// for this task definition.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '8', '26'
	AWSECSTaskRevisionKey = "aws.ecs.task.revision"

	// AWSEKSClusterARNKey is the attribute Key conforming to the
	// "aws.eks.cluster.arn" semantic conventions. It represents the ARN of an
	// EKS cluster.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'arn:aws:ecs:us-west-2:123456789123:cluster/my-cluster'
	AWSEKSClusterARNKey = "aws.eks.cluster.arn"

	// AWSLogGroupARNsKey is the attribute Key conforming to the
	// "aws.log.group.arns" semantic conventions. It represents the Amazon
	// Resource Name(s) (ARN) of the AWS log group(s).
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples:
	// 'arn:aws:logs:us-west-1:123456789012:log-group:/aws/my/group:*'
	// Note: See the [log group ARN format
	// documentation](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/iam-access-control-overview-cwl.html#CWL_ARN_Format).
	AWSLogGroupARNsKey = "aws.log.group.arns"

	// AWSLogGroupNamesKey is the attribute Key conforming to the
	// "aws.log.group.names" semantic conventions. It represents the name(s) of
	// the AWS log group(s) an application is writing to.
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '/aws/lambda/my-function', 'opentelemetry-service'
	// Note: Multiple log groups must be supported for cases like
	// multi-container applications, where a single application has sidecar
	// containers, and each write to their own log group.
	AWSLogGroupNamesKey = "aws.log.group.names"

	// AWSLogStreamARNsKey is the attribute Key conforming to the
	// "aws.log.stream.arns" semantic conventions. It represents the ARN(s) of
	// the AWS log stream(s).
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples:
	// 'arn:aws:logs:us-west-1:123456789012:log-group:/aws/my/group:log-stream:logs/main/10838bed-421f-43ef-870a-f43feacbbb5b'
	// Note: See the [log stream ARN format
	// documentation](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/iam-access-control-overview-cwl.html#CWL_ARN_Format).
	// One log group can contain several log streams, so these ARNs necessarily
	// identify both a log group and a log stream.
	AWSLogStreamARNsKey = "aws.log.stream.arns"

	// AWSLogStreamNamesKey is the attribute Key conforming to the
	// "aws.log.stream.names" semantic conventions. It represents the name(s)
	// of the AWS log stream(s) an application is writing to.
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'logs/main/10838bed-421f-43ef-870a-f43feacbbb5b'
	AWSLogStreamNamesKey = "aws.log.stream.names"

	// GCPCloudRunJobExecutionKey is the attribute Key conforming to the
	// "gcp.cloud_run.job.execution" semantic conventions. It represents the
	// name of the Cloud Run
	// [execution](https://cloud.google.com/run/docs/managing/job-executions)
	// being run for the Job, as set by the
	// [`CLOUD_RUN_EXECUTION`](https://cloud.google.com/run/docs/container-contract#jobs-env-vars)
	// environment variable.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'job-name-xxxx', 'sample-job-mdw84'
	GCPCloudRunJobExecutionKey = "gcp.cloud_run.job.execution"

	// GCPCloudRunJobTaskIndexKey is the attribute Key conforming to the
	// "gcp.cloud_run.job.task_index" semantic conventions. It represents the
	// index for a task within an execution as provided by the
	// [`CLOUD_RUN_TASK_INDEX`](https://cloud.google.com/run/docs/container-contract#jobs-env-vars)
	// environment variable.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 0, 1
	GCPCloudRunJobTaskIndexKey = "gcp.cloud_run.job.task_index"

	// GCPGceInstanceHostnameKey is the attribute Key conforming to the
	// "gcp.gce.instance.hostname" semantic conventions. It represents the
	// hostname of a GCE instance. This is the full value of the default or
	// [custom
	// hostname](https://cloud.google.com/compute/docs/instances/custom-hostname-vm).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'my-host1234.example.com',
	// 'sample-vm.us-west1-b.c.my-project.internal'
	GCPGceInstanceHostnameKey = "gcp.gce.instance.hostname"

	// GCPGceInstanceNameKey is the attribute Key conforming to the
	// "gcp.gce.instance.name" semantic conventions. It represents the instance
	// name of a GCE instance. This is the value provided by `host.name`, the
	// visible name of the instance in the Cloud Console UI, and the prefix for
	// the default hostname of the instance as defined by the [default internal
	// DNS
	// name](https://cloud.google.com/compute/docs/internal-dns#instance-fully-qualified-domain-names).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'instance-1', 'my-vm-name'
	GCPGceInstanceNameKey = "gcp.gce.instance.name"

	// HerokuAppIDKey is the attribute Key conforming to the "heroku.app.id"
	// semantic conventions. It represents the unique identifier for the
	// application
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '2daa2797-e42b-4624-9322-ec3f968df4da'
	HerokuAppIDKey = "heroku.app.id"

	// HerokuReleaseCommitKey is the attribute Key conforming to the
	// "heroku.release.commit" semantic conventions. It represents the commit
	// hash for the current release
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'e6134959463efd8966b20e75b913cafe3f5ec'
	HerokuReleaseCommitKey = "heroku.release.commit"

	// HerokuReleaseCreationTimestampKey is the attribute Key conforming to the
	// "heroku.release.creation_timestamp" semantic conventions. It represents
	// the time and date the release was created
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '2022-10-23T18:00:42Z'
	HerokuReleaseCreationTimestampKey = "heroku.release.creation_timestamp"

	// DeploymentEnvironmentKey is the attribute Key conforming to the
	// "deployment.environment" semantic conventions. It represents the name of
	// the [deployment
	// environment](https://wikipedia.org/wiki/Deployment_environment) (aka
	// deployment tier).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'staging', 'production'
	// Note: `deployment.environment` does not affect the uniqueness
	// constraints defined through
	// the `service.namespace`, `service.name` and `service.instance.id`
	// resource attributes.
	// This implies that resources carrying the following attribute
	// combinations MUST be
	// considered to be identifying the same service:
	//
	// * `service.name=frontend`, `deployment.environment=production`
	// * `service.name=frontend`, `deployment.environment=staging`.
	DeploymentEnvironmentKey = "deployment.environment"

	// FaaSInstanceKey is the attribute Key conforming to the "faas.instance"
	// semantic conventions. It represents the execution environment ID as a
	// string, that will be potentially reused for other invocations to the
	// same function/function version.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '2021/06/28/[$LATEST]2f399eb14537447da05ab2a2e39309de'
	// Note: * **AWS Lambda:** Use the (full) log stream name.
	FaaSInstanceKey = "faas.instance"

	// FaaSMaxMemoryKey is the attribute Key conforming to the
	// "faas.max_memory" semantic conventions. It represents the amount of
	// memory available to the serverless function converted to Bytes.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 134217728
	// Note: It's recommended to set this attribute since e.g. too little
	// memory can easily stop a Java AWS Lambda function from working
	// correctly. On AWS Lambda, the environment variable
	// `AWS_LAMBDA_FUNCTION_MEMORY_SIZE` provides this information (which must
	// be multiplied by 1,048,576).
	FaaSMaxMemoryKey = "faas.max_memory"

	// FaaSNameKey is the attribute Key conforming to the "faas.name" semantic
	// conventions. It represents the name of the single function that this
	// runtime instance executes.
	//
	// Type: string
	// RequirementLevel: Required
	// Stability: experimental
	// Examples: 'my-function', 'myazurefunctionapp/some-function-name'
	// Note: This is the name of the function as configured/deployed on the
	// FaaS
	// platform and is usually different from the name of the callback
	// function (which may be stored in the
	// [`code.namespace`/`code.function`](/docs/general/attributes.md#source-code-attributes)
	// span attributes).
	//
	// For some cloud providers, the above definition is ambiguous. The
	// following
	// definition of function name MUST be used for this attribute
	// (and consequently the span name) for the listed cloud
	// providers/products:
	//
	// * **Azure:**  The full name `<FUNCAPP>/<FUNC>`, i.e., function app name
	//   followed by a forward slash followed by the function name (this form
	//   can also be seen in the resource JSON for the function).
	//   This means that a span attribute MUST be used, as an Azure function
	//   app can host multiple functions that would usually share
	//   a TracerProvider (see also the `cloud.resource_id` attribute).
	FaaSNameKey = "faas.name"

	// FaaSVersionKey is the attribute Key conforming to the "faas.version"
	// semantic conventions. It represents the immutable version of the
	// function being executed.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '26', 'pinkfroid-00002'
	// Note: Depending on the cloud provider and platform, use:
	//
	// * **AWS Lambda:** The [function
	// version](https://docs.aws.amazon.com/lambda/latest/dg/configuration-versions.html)
	//   (an integer represented as a decimal string).
	// * **Google Cloud Run (Services):** The
	// [revision](https://cloud.google.com/run/docs/managing/revisions)
	//   (i.e., the function name plus the revision suffix).
	// * **Google Cloud Functions:** The value of the
	//   [`K_REVISION` environment
	// variable](https://cloud.google.com/functions/docs/env-var#runtime_environment_variables_set_automatically).
	// * **Azure Functions:** Not applicable. Do not set this attribute.
	FaaSVersionKey = "faas.version"

	// ServiceNameKey is the attribute Key conforming to the "service.name"
	// semantic conventions. It represents the logical name of the service.
	//
	// Type: string
	// RequirementLevel: Required
	// Stability: experimental
	// Examples: 'shoppingcart'
	// Note: MUST be the same for all instances of horizontally scaled
	// services. If the value was not specified, SDKs MUST fallback to
	// `unknown_service:` concatenated with
	// [`process.executable.name`](process.md#process), e.g.
	// `unknown_service:bash`. If `process.executable.name` is not available,
	// the value MUST be set to `unknown_service`.
	ServiceNameKey = "service.name"

	// ServiceVersionKey is the attribute Key conforming to the
	// "service.version" semantic conventions. It represents the version string
	// of the service API or implementation. The format is not defined by these
	// conventions.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '2.0.0', 'a01dbef8a'
	ServiceVersionKey = "service.version"

	// ServiceInstanceIDKey is the attribute Key conforming to the
	// "service.instance.id" semantic conventions. It represents the string ID
	// of the service instance.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'my-k8s-pod-deployment-1',
	// '627cc493-f310-47de-96bd-71410b7dec09'
	// Note: MUST be unique for each instance of the same
	// `service.namespace,service.name` pair (in other words
	// `service.namespace,service.name,service.instance.id` triplet MUST be
	// globally unique). The ID helps to distinguish instances of the same
	// service that exist at the same time (e.g. instances of a horizontally
	// scaled service). It is preferable for the ID to be persistent and stay
	// the same for the lifetime of the service instance, however it is
	// acceptable that the ID is ephemeral and changes during important
	// lifetime events for the service (e.g. service restarts). If the service
	// has no inherent unique ID that can be used as the value of this
	// attribute it is recommended to generate a random Version 1 or Version 4
	// RFC 4122 UUID (services aiming for reproducible UUIDs may also use
	// Version 5, see RFC 4122 for more recommendations).
	ServiceInstanceIDKey = "service.instance.id"

	// ServiceNamespaceKey is the attribute Key conforming to the
	// "service.namespace" semantic conventions. It represents a namespace for
	// `service.name`.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'Shop'
	// Note: A string value having a meaning that helps to distinguish a group
	// of services, for example the team name that owns a group of services.
	// `service.name` is expected to be unique within the same namespace. If
	// `service.namespace` is not specified in the Resource then `service.name`
	// is expected to be unique for all services that have no explicit
	// namespace defined (so the empty/unspecified namespace is simply one more
	// valid namespace). Zero-length namespace string is assumed equal to
	// unspecified namespace.
	ServiceNamespaceKey = "service.namespace"

	// TelemetrySDKLanguageKey is the attribute Key conforming to the
	// "telemetry.sdk.language" semantic conventions. It represents the
	// language of the telemetry SDK.
	//
	// Type: Enum
	// RequirementLevel: Required
	// Stability: experimental
	TelemetrySDKLanguageKey = "telemetry.sdk.language"

	// TelemetrySDKNameKey is the attribute Key conforming to the
	// "telemetry.sdk.name" semantic conventions. It represents the name of the
	// telemetry SDK as defined above.
	//
	// Type: string
	// RequirementLevel: Required
	// Stability: experimental
	// Examples: 'opentelemetry'
	// Note: The OpenTelemetry SDK MUST set the `telemetry.sdk.name` attribute
	// to `opentelemetry`.
	// If another SDK, like a fork or a vendor-provided implementation, is
	// used, this SDK MUST set the
	// `telemetry.sdk.name` attribute to the fully-qualified class or module
	// name of this SDK's main entry point
	// or another suitable identifier depending on the language.
	// The identifier `opentelemetry` is reserved and MUST NOT be used in this
	// case.
	// All custom identifiers SHOULD be stable across different versions of an
	// implementation.
	TelemetrySDKNameKey = "telemetry.sdk.name"

	// TelemetrySDKVersionKey is the attribute Key conforming to the
	// "telemetry.sdk.version" semantic conventions. It represents the version
	// string of the telemetry SDK.
	//
	// Type: string
	// RequirementLevel: Required
	// Stability: experimental
	// Examples: '1.2.3'
	TelemetrySDKVersionKey = "telemetry.sdk.version"

	// TelemetryDistroNameKey is the attribute Key conforming to the
	// "telemetry.distro.name" semantic conventions. It represents the name of
	// the auto instrumentation agent or distribution, if used.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'parts-unlimited-java'
	// Note: Official auto instrumentation agents and distributions SHOULD set
	// the `telemetry.distro.name` attribute to
	// a string starting with `opentelemetry-`, e.g.
	// `opentelemetry-java-instrumentation`.
	TelemetryDistroNameKey = "telemetry.distro.name"

	// TelemetryDistroVersionKey is the attribute Key conforming to the
	// "telemetry.distro.version" semantic conventions. It represents the
	// version string of the auto instrumentation agent or distribution, if
	// used.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '1.2.3'
	TelemetryDistroVersionKey = "telemetry.distro.version"

	// WebEngineDescriptionKey is the attribute Key conforming to the
	// "webengine.description" semantic conventions. It represents the
	// additional description of the web engine (e.g. detailed version and
	// edition information).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'WildFly Full 21.0.0.Final (WildFly Core 13.0.1.Final) -
	// 2.2.2.Final'
	WebEngineDescriptionKey = "webengine.description"

	// WebEngineNameKey is the attribute Key conforming to the "webengine.name"
	// semantic conventions. It represents the name of the web engine.
	//
	// Type: string
	// RequirementLevel: Required
	// Stability: experimental
	// Examples: 'WildFly'
	WebEngineNameKey = "webengine.name"

	// WebEngineVersionKey is the attribute Key conforming to the
	// "webengine.version" semantic conventions. It represents the version of
	// the web engine.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '21.0.0'
	WebEngineVersionKey = "webengine.version"

	// OTelScopeNameKey is the attribute Key conforming to the
	// "otel.scope.name" semantic conventions. It represents the name of the
	// instrumentation scope - (`InstrumentationScope.Name` in OTLP).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'io.opentelemetry.contrib.mongodb'
	OTelScopeNameKey = "otel.scope.name"

	// OTelScopeVersionKey is the attribute Key conforming to the
	// "otel.scope.version" semantic conventions. It represents the version of
	// the instrumentation scope - (`InstrumentationScope.Version` in OTLP).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '1.0.0'
	OTelScopeVersionKey = "otel.scope.version"

	// OTelLibraryNameKey is the attribute Key conforming to the
	// "otel.library.name" semantic conventions.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: 'io.opentelemetry.contrib.mongodb'
	// Deprecated: use the `otel.scope.name` attribute.
	OTelLibraryNameKey = "otel.library.name"

	// OTelLibraryVersionKey is the attribute Key conforming to the
	// "otel.library.version" semantic conventions.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: deprecated
	// Examples: '1.0.0'
	// Deprecated: use the `otel.scope.version` attribute.
	OTelLibraryVersionKey = "otel.library.version"

	// PeerServiceKey is the attribute Key conforming to the "peer.service"
	// semantic conventions. It represents the
	// [`service.name`](/docs/resource/README.md#service) of the remote
	// service. SHOULD be equal to the actual `service.name` resource attribute
	// of the remote service if any.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'AuthTokenCache'
	PeerServiceKey = "peer.service"

	// EnduserIDKey is the attribute Key conforming to the "enduser.id"
	// semantic conventions. It represents the username or client_id extracted
	// from the access token or
	// [Authorization](https://tools.ietf.org/html/rfc7235#section-4.2) header
	// in the inbound request from outside the system.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'username'
	EnduserIDKey = "enduser.id"

	// EnduserRoleKey is the attribute Key conforming to the "enduser.role"
	// semantic conventions. It represents the actual/assumed role the client
	// is making the request under extracted from token or application security
	// context.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'admin'
	EnduserRoleKey = "enduser.role"

	// EnduserScopeKey is the attribute Key conforming to the "enduser.scope"
	// semantic conventions. It represents the scopes or granted authorities
	// the client currently possesses extracted from token or application
	// security context. The value would come from the scope associated with an
	// [OAuth 2.0 Access
	// Token](https://tools.ietf.org/html/rfc6749#section-3.3) or an attribute
	// value in a [SAML 2.0
	// Assertion](http://docs.oasis-open.org/security/saml/Post2.0/sstc-saml-tech-overview-2.0.html).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'read:message, write:files'
	EnduserScopeKey = "enduser.scope"

	// CodeColumnKey is the attribute Key conforming to the "code.column"
	// semantic conventions. It represents the column number in `code.filepath`
	// best representing the operation. It SHOULD point within the code unit
	// named in `code.function`.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 16
	CodeColumnKey = "code.column"

	// CodeFilepathKey is the attribute Key conforming to the "code.filepath"
	// semantic conventions. It represents the source code file name that
	// identifies the code unit as uniquely as possible (preferably an absolute
	// file path).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '/usr/local/MyApplication/content_root/app/index.php'
	CodeFilepathKey = "code.filepath"

	// CodeFunctionKey is the attribute Key conforming to the "code.function"
	// semantic conventions. It represents the method or function name, or
	// equivalent (usually rightmost part of the code unit's name).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'serveRequest'
	CodeFunctionKey = "code.function"

	// CodeLineNumberKey is the attribute Key conforming to the "code.lineno"
	// semantic conventions. It represents the line number in `code.filepath`
	// best representing the operation. It SHOULD point within the code unit
	// named in `code.function`.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 42
	CodeLineNumberKey = "code.lineno"

	// CodeNamespaceKey is the attribute Key conforming to the "code.namespace"
	// semantic conventions. It represents the "namespace" within which
	// `code.function` is defined. Usually the qualified class or module name,
	// such that `code.namespace` + some separator + `code.function` form a
	// unique identifier for the code unit.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'com.example.MyHTTPService'
	CodeNamespaceKey = "code.namespace"

	// CodeStacktraceKey is the attribute Key conforming to the
	// "code.stacktrace" semantic conventions. It represents a stacktrace as a
	// string in the natural representation for the language runtime. The
	// representation is to be determined and documented by each language SIG.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'at
	// com.example.GenerateTrace.methodB(GenerateTrace.java:13)\\n at '
	//  'com.example.GenerateTrace.methodA(GenerateTrace.java:9)\\n at '
	//  'com.example.GenerateTrace.main(GenerateTrace.java:5)'
	CodeStacktraceKey = "code.stacktrace"

	// ThreadIDKey is the attribute Key conforming to the "thread.id" semantic
	// conventions. It represents the current "managed" thread ID (as opposed
	// to OS thread ID).
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 42
	ThreadIDKey = "thread.id"

	// ThreadNameKey is the attribute Key conforming to the "thread.name"
	// semantic conventions. It represents the current thread name.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'main'
	ThreadNameKey = "thread.name"

	// AWSLambdaInvokedARNKey is the attribute Key conforming to the
	// "aws.lambda.invoked_arn" semantic conventions. It represents the full
	// invoked ARN as provided on the `Context` passed to the function
	// (`Lambda-Runtime-Invoked-Function-ARN` header on the
	// `/runtime/invocation/next` applicable).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'arn:aws:lambda:us-east-1:123456:function:myfunction:myalias'
	// Note: This may be different from `cloud.resource_id` if an alias is
	// involved.
	AWSLambdaInvokedARNKey = "aws.lambda.invoked_arn"

	// CloudeventsEventIDKey is the attribute Key conforming to the
	// "cloudevents.event_id" semantic conventions. It represents the
	// [event_id](https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md#id)
	// uniquely identifies the event.
	//
	// Type: string
	// RequirementLevel: Required
	// Stability: experimental
	// Examples: '123e4567-e89b-12d3-a456-426614174000', '0001'
	CloudeventsEventIDKey = "cloudevents.event_id"

	// CloudeventsEventSourceKey is the attribute Key conforming to the
	// "cloudevents.event_source" semantic conventions. It represents the
	// [source](https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md#source-1)
	// identifies the context in which an event happened.
	//
	// Type: string
	// RequirementLevel: Required
	// Stability: experimental
	// Examples: 'https://github.com/cloudevents',
	// '/cloudevents/spec/pull/123', 'my-service'
	CloudeventsEventSourceKey = "cloudevents.event_source"

	// CloudeventsEventSpecVersionKey is the attribute Key conforming to the
	// "cloudevents.event_spec_version" semantic conventions. It represents the
	// [version of the CloudEvents
	// specification](https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md#specversion)
	// which the event uses.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '1.0'
	CloudeventsEventSpecVersionKey = "cloudevents.event_spec_version"

	// CloudeventsEventSubjectKey is the attribute Key conforming to the
	// "cloudevents.event_subject" semantic conventions. It represents the
	// [subject](https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md#subject)
	// of the event in the context of the event producer (identified by
	// source).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'mynewfile.jpg'
	CloudeventsEventSubjectKey = "cloudevents.event_subject"

	// CloudeventsEventTypeKey is the attribute Key conforming to the
	// "cloudevents.event_type" semantic conventions. It represents the
	// [event_type](https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md#type)
	// contains a value describing the type of event related to the originating
	// occurrence.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'com.github.pull_request.opened',
	// 'com.example.object.deleted.v2'
	CloudeventsEventTypeKey = "cloudevents.event_type"

	// OpentracingRefTypeKey is the attribute Key conforming to the
	// "opentracing.ref_type" semantic conventions. It represents the
	// parent-child Reference type
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	// Note: The causal relationship between a child Span and a parent Span.
	OpentracingRefTypeKey = "opentracing.ref_type"

	// OTelStatusCodeKey is the attribute Key conforming to the
	// "otel.status_code" semantic conventions. It represents the name of the
	// code, either "OK" or "ERROR". MUST NOT be set if the status code is
	// UNSET.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	OTelStatusCodeKey = "otel.status_code"

	// OTelStatusDescriptionKey is the attribute Key conforming to the
	// "otel.status_description" semantic conventions. It represents the
	// description of the Status if it has a value, otherwise not set.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'resource not found'
	OTelStatusDescriptionKey = "otel.status_description"

	// FaaSInvocationIDKey is the attribute Key conforming to the
	// "faas.invocation_id" semantic conventions. It represents the invocation
	// ID of the current function invocation.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'af9d5aa4-a685-4c5f-a22b-444f80b3cc28'
	FaaSInvocationIDKey = "faas.invocation_id"

	// FaaSDocumentCollectionKey is the attribute Key conforming to the
	// "faas.document.collection" semantic conventions. It represents the name
	// of the source on which the triggering operation was performed. For
	// example, in Cloud Storage or S3 corresponds to the bucket name, and in
	// Cosmos DB to the database name.
	//
	// Type: string
	// RequirementLevel: Required
	// Stability: experimental
	// Examples: 'myBucketName', 'myDBName'
	FaaSDocumentCollectionKey = "faas.document.collection"

	// FaaSDocumentNameKey is the attribute Key conforming to the
	// "faas.document.name" semantic conventions. It represents the document
	// name/table subjected to the operation. For example, in Cloud Storage or
	// S3 is the name of the file, and in Cosmos DB the table name.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'myFile.txt', 'myTableName'
	FaaSDocumentNameKey = "faas.document.name"

	// FaaSDocumentOperationKey is the attribute Key conforming to the
	// "faas.document.operation" semantic conventions. It represents the
	// describes the type of the operation that was performed on the data.
	//
	// Type: Enum
	// RequirementLevel: Required
	// Stability: experimental
	FaaSDocumentOperationKey = "faas.document.operation"

	// FaaSDocumentTimeKey is the attribute Key conforming to the
	// "faas.document.time" semantic conventions. It represents a string
	// containing the time when the data was accessed in the [ISO
	// 8601](https://www.iso.org/iso-8601-date-and-time-format.html) format
	// expressed in [UTC](https://www.w3.org/TR/NOTE-datetime).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '2020-01-23T13:47:06Z'
	FaaSDocumentTimeKey = "faas.document.time"

	// FaaSCronKey is the attribute Key conforming to the "faas.cron" semantic
	// conventions. It represents a string containing the schedule period as
	// [Cron
	// Expression](https://docs.oracle.com/cd/E12058_01/doc/doc.1014/e12030/cron_expressions.htm).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '0/5 * * * ? *'
	FaaSCronKey = "faas.cron"

	// FaaSTimeKey is the attribute Key conforming to the "faas.time" semantic
	// conventions. It represents a string containing the function invocation
	// time in the [ISO
	// 8601](https://www.iso.org/iso-8601-date-and-time-format.html) format
	// expressed in [UTC](https://www.w3.org/TR/NOTE-datetime).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '2020-01-23T13:47:06Z'
	FaaSTimeKey = "faas.time"

	// FaaSColdstartKey is the attribute Key conforming to the "faas.coldstart"
	// semantic conventions. It represents a boolean that is true if the
	// serverless function is executed for the first time (aka cold-start).
	//
	// Type: boolean
	// RequirementLevel: Optional
	// Stability: experimental
	FaaSColdstartKey = "faas.coldstart"

	// AWSRequestIDKey is the attribute Key conforming to the "aws.request_id"
	// semantic conventions. It represents the AWS request ID as returned in
	// the response headers `x-amz-request-id` or `x-amz-requestid`.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '79b9da39-b7ae-508a-a6bc-864b2829c622', 'C9ER4AJX75574TDJ'
	AWSRequestIDKey = "aws.request_id"

	// AWSDynamoDBAttributesToGetKey is the attribute Key conforming to the
	// "aws.dynamodb.attributes_to_get" semantic conventions. It represents the
	// value of the `AttributesToGet` request parameter.
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'lives', 'id'
	AWSDynamoDBAttributesToGetKey = "aws.dynamodb.attributes_to_get"

	// AWSDynamoDBConsistentReadKey is the attribute Key conforming to the
	// "aws.dynamodb.consistent_read" semantic conventions. It represents the
	// value of the `ConsistentRead` request parameter.
	//
	// Type: boolean
	// RequirementLevel: Optional
	// Stability: experimental
	AWSDynamoDBConsistentReadKey = "aws.dynamodb.consistent_read"

	// AWSDynamoDBConsumedCapacityKey is the attribute Key conforming to the
	// "aws.dynamodb.consumed_capacity" semantic conventions. It represents the
	// JSON-serialized value of each item in the `ConsumedCapacity` response
	// field.
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '{ "CapacityUnits": number, "GlobalSecondaryIndexes": {
	// "string" : { "CapacityUnits": number, "ReadCapacityUnits": number,
	// "WriteCapacityUnits": number } }, "LocalSecondaryIndexes": { "string" :
	// { "CapacityUnits": number, "ReadCapacityUnits": number,
	// "WriteCapacityUnits": number } }, "ReadCapacityUnits": number, "Table":
	// { "CapacityUnits": number, "ReadCapacityUnits": number,
	// "WriteCapacityUnits": number }, "TableName": "string",
	// "WriteCapacityUnits": number }'
	AWSDynamoDBConsumedCapacityKey = "aws.dynamodb.consumed_capacity"

	// AWSDynamoDBIndexNameKey is the attribute Key conforming to the
	// "aws.dynamodb.index_name" semantic conventions. It represents the value
	// of the `IndexName` request parameter.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'name_to_group'
	AWSDynamoDBIndexNameKey = "aws.dynamodb.index_name"

	// AWSDynamoDBItemCollectionMetricsKey is the attribute Key conforming to
	// the "aws.dynamodb.item_collection_metrics" semantic conventions. It
	// represents the JSON-serialized value of the `ItemCollectionMetrics`
	// response field.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '{ "string" : [ { "ItemCollectionKey": { "string" : { "B":
	// blob, "BOOL": boolean, "BS": [ blob ], "L": [ "AttributeValue" ], "M": {
	// "string" : "AttributeValue" }, "N": "string", "NS": [ "string" ],
	// "NULL": boolean, "S": "string", "SS": [ "string" ] } },
	// "SizeEstimateRangeGB": [ number ] } ] }'
	AWSDynamoDBItemCollectionMetricsKey = "aws.dynamodb.item_collection_metrics"

	// AWSDynamoDBLimitKey is the attribute Key conforming to the
	// "aws.dynamodb.limit" semantic conventions. It represents the value of
	// the `Limit` request parameter.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 10
	AWSDynamoDBLimitKey = "aws.dynamodb.limit"

	// AWSDynamoDBProjectionKey is the attribute Key conforming to the
	// "aws.dynamodb.projection" semantic conventions. It represents the value
	// of the `ProjectionExpression` request parameter.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'Title', 'Title, Price, Color', 'Title, Body,
	// RelatedItems, ProductReviews'
	AWSDynamoDBProjectionKey = "aws.dynamodb.projection"

	// AWSDynamoDBProvisionedReadCapacityKey is the attribute Key conforming to
	// the "aws.dynamodb.provisioned_read_capacity" semantic conventions. It
	// represents the value of the `ProvisionedThroughput.ReadCapacityUnits`
	// request parameter.
	//
	// Type: double
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 1.0, 2.0
	AWSDynamoDBProvisionedReadCapacityKey = "aws.dynamodb.provisioned_read_capacity"

	// AWSDynamoDBProvisionedWriteCapacityKey is the attribute Key conforming
	// to the "aws.dynamodb.provisioned_write_capacity" semantic conventions.
	// It represents the value of the
	// `ProvisionedThroughput.WriteCapacityUnits` request parameter.
	//
	// Type: double
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 1.0, 2.0
	AWSDynamoDBProvisionedWriteCapacityKey = "aws.dynamodb.provisioned_write_capacity"

	// AWSDynamoDBSelectKey is the attribute Key conforming to the
	// "aws.dynamodb.select" semantic conventions. It represents the value of
	// the `Select` request parameter.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'ALL_ATTRIBUTES', 'COUNT'
	AWSDynamoDBSelectKey = "aws.dynamodb.select"

	// AWSDynamoDBTableNamesKey is the attribute Key conforming to the
	// "aws.dynamodb.table_names" semantic conventions. It represents the keys
	// in the `RequestItems` object field.
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'Users', 'Cats'
	AWSDynamoDBTableNamesKey = "aws.dynamodb.table_names"

	// AWSDynamoDBGlobalSecondaryIndexesKey is the attribute Key conforming to
	// the "aws.dynamodb.global_secondary_indexes" semantic conventions. It
	// represents the JSON-serialized value of each item of the
	// `GlobalSecondaryIndexes` request field
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '{ "IndexName": "string", "KeySchema": [ { "AttributeName":
	// "string", "KeyType": "string" } ], "Projection": { "NonKeyAttributes": [
	// "string" ], "ProjectionType": "string" }, "ProvisionedThroughput": {
	// "ReadCapacityUnits": number, "WriteCapacityUnits": number } }'
	AWSDynamoDBGlobalSecondaryIndexesKey = "aws.dynamodb.global_secondary_indexes"

	// AWSDynamoDBLocalSecondaryIndexesKey is the attribute Key conforming to
	// the "aws.dynamodb.local_secondary_indexes" semantic conventions. It
	// represents the JSON-serialized value of each item of the
	// `LocalSecondaryIndexes` request field.
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '{ "IndexARN": "string", "IndexName": "string",
	// "IndexSizeBytes": number, "ItemCount": number, "KeySchema": [ {
	// "AttributeName": "string", "KeyType": "string" } ], "Projection": {
	// "NonKeyAttributes": [ "string" ], "ProjectionType": "string" } }'
	AWSDynamoDBLocalSecondaryIndexesKey = "aws.dynamodb.local_secondary_indexes"

	// AWSDynamoDBExclusiveStartTableKey is the attribute Key conforming to the
	// "aws.dynamodb.exclusive_start_table" semantic conventions. It represents
	// the value of the `ExclusiveStartTableName` request parameter.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'Users', 'CatsTable'
	AWSDynamoDBExclusiveStartTableKey = "aws.dynamodb.exclusive_start_table"

	// AWSDynamoDBTableCountKey is the attribute Key conforming to the
	// "aws.dynamodb.table_count" semantic conventions. It represents the the
	// number of items in the `TableNames` response parameter.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 20
	AWSDynamoDBTableCountKey = "aws.dynamodb.table_count"

	// AWSDynamoDBScanForwardKey is the attribute Key conforming to the
	// "aws.dynamodb.scan_forward" semantic conventions. It represents the
	// value of the `ScanIndexForward` request parameter.
	//
	// Type: boolean
	// RequirementLevel: Optional
	// Stability: experimental
	AWSDynamoDBScanForwardKey = "aws.dynamodb.scan_forward"

	// AWSDynamoDBCountKey is the attribute Key conforming to the
	// "aws.dynamodb.count" semantic conventions. It represents the value of
	// the `Count` response parameter.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 10
	AWSDynamoDBCountKey = "aws.dynamodb.count"

	// AWSDynamoDBScannedCountKey is the attribute Key conforming to the
	// "aws.dynamodb.scanned_count" semantic conventions. It represents the
	// value of the `ScannedCount` response parameter.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 50
	AWSDynamoDBScannedCountKey = "aws.dynamodb.scanned_count"

	// AWSDynamoDBSegmentKey is the attribute Key conforming to the
	// "aws.dynamodb.segment" semantic conventions. It represents the value of
	// the `Segment` request parameter.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 10
	AWSDynamoDBSegmentKey = "aws.dynamodb.segment"

	// AWSDynamoDBTotalSegmentsKey is the attribute Key conforming to the
	// "aws.dynamodb.total_segments" semantic conventions. It represents the
	// value of the `TotalSegments` request parameter.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 100
	AWSDynamoDBTotalSegmentsKey = "aws.dynamodb.total_segments"

	// AWSDynamoDBAttributeDefinitionsKey is the attribute Key conforming to
	// the "aws.dynamodb.attribute_definitions" semantic conventions. It
	// represents the JSON-serialized value of each item in the
	// `AttributeDefinitions` request field.
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '{ "AttributeName": "string", "AttributeType": "string" }'
	AWSDynamoDBAttributeDefinitionsKey = "aws.dynamodb.attribute_definitions"

	// AWSDynamoDBGlobalSecondaryIndexUpdatesKey is the attribute Key
	// conforming to the "aws.dynamodb.global_secondary_index_updates" semantic
	// conventions. It represents the JSON-serialized value of each item in the
	// the `GlobalSecondaryIndexUpdates` request field.
	//
	// Type: string[]
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '{ "Create": { "IndexName": "string", "KeySchema": [ {
	// "AttributeName": "string", "KeyType": "string" } ], "Projection": {
	// "NonKeyAttributes": [ "string" ], "ProjectionType": "string" },
	// "ProvisionedThroughput": { "ReadCapacityUnits": number,
	// "WriteCapacityUnits": number } }'
	AWSDynamoDBGlobalSecondaryIndexUpdatesKey = "aws.dynamodb.global_secondary_index_updates"

	// AWSS3BucketKey is the attribute Key conforming to the "aws.s3.bucket"
	// semantic conventions. It represents the S3 bucket name the request
	// refers to. Corresponds to the `--bucket` parameter of the [S3
	// API](https://docs.aws.amazon.com/cli/latest/reference/s3api/index.html)
	// operations.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'some-bucket-name'
	// Note: The `bucket` attribute is applicable to all S3 operations that
	// reference a bucket, i.e. that require the bucket name as a mandatory
	// parameter.
	// This applies to almost all S3 operations except `list-buckets`.
	AWSS3BucketKey = "aws.s3.bucket"

	// AWSS3CopySourceKey is the attribute Key conforming to the
	// "aws.s3.copy_source" semantic conventions. It represents the source
	// object (in the form `bucket`/`key`) for the copy operation.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'someFile.yml'
	// Note: The `copy_source` attribute applies to S3 copy operations and
	// corresponds to the `--copy-source` parameter
	// of the [copy-object operation within the S3
	// API](https://docs.aws.amazon.com/cli/latest/reference/s3api/copy-object.html).
	// This applies in particular to the following operations:
	//
	// -
	// [copy-object](https://docs.aws.amazon.com/cli/latest/reference/s3api/copy-object.html)
	// -
	// [upload-part-copy](https://docs.aws.amazon.com/cli/latest/reference/s3api/upload-part-copy.html)
	AWSS3CopySourceKey = "aws.s3.copy_source"

	// AWSS3DeleteKey is the attribute Key conforming to the "aws.s3.delete"
	// semantic conventions. It represents the delete request container that
	// specifies the objects to be deleted.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples:
	// 'Objects=[{Key=string,VersionID=string},{Key=string,VersionID=string}],Quiet=boolean'
	// Note: The `delete` attribute is only applicable to the
	// [delete-object](https://docs.aws.amazon.com/cli/latest/reference/s3api/delete-object.html)
	// operation.
	// The `delete` attribute corresponds to the `--delete` parameter of the
	// [delete-objects operation within the S3
	// API](https://docs.aws.amazon.com/cli/latest/reference/s3api/delete-objects.html).
	AWSS3DeleteKey = "aws.s3.delete"

	// AWSS3KeyKey is the attribute Key conforming to the "aws.s3.key" semantic
	// conventions. It represents the S3 object key the request refers to.
	// Corresponds to the `--key` parameter of the [S3
	// API](https://docs.aws.amazon.com/cli/latest/reference/s3api/index.html)
	// operations.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'someFile.yml'
	// Note: The `key` attribute is applicable to all object-related S3
	// operations, i.e. that require the object key as a mandatory parameter.
	// This applies in particular to the following operations:
	//
	// -
	// [copy-object](https://docs.aws.amazon.com/cli/latest/reference/s3api/copy-object.html)
	// -
	// [delete-object](https://docs.aws.amazon.com/cli/latest/reference/s3api/delete-object.html)
	// -
	// [get-object](https://docs.aws.amazon.com/cli/latest/reference/s3api/get-object.html)
	// -
	// [head-object](https://docs.aws.amazon.com/cli/latest/reference/s3api/head-object.html)
	// -
	// [put-object](https://docs.aws.amazon.com/cli/latest/reference/s3api/put-object.html)
	// -
	// [restore-object](https://docs.aws.amazon.com/cli/latest/reference/s3api/restore-object.html)
	// -
	// [select-object-content](https://docs.aws.amazon.com/cli/latest/reference/s3api/select-object-content.html)
	// -
	// [abort-multipart-upload](https://docs.aws.amazon.com/cli/latest/reference/s3api/abort-multipart-upload.html)
	// -
	// [complete-multipart-upload](https://docs.aws.amazon.com/cli/latest/reference/s3api/complete-multipart-upload.html)
	// -
	// [create-multipart-upload](https://docs.aws.amazon.com/cli/latest/reference/s3api/create-multipart-upload.html)
	// -
	// [list-parts](https://docs.aws.amazon.com/cli/latest/reference/s3api/list-parts.html)
	// -
	// [upload-part](https://docs.aws.amazon.com/cli/latest/reference/s3api/upload-part.html)
	// -
	// [upload-part-copy](https://docs.aws.amazon.com/cli/latest/reference/s3api/upload-part-copy.html)
	AWSS3KeyKey = "aws.s3.key"

	// AWSS3PartNumberKey is the attribute Key conforming to the
	// "aws.s3.part_number" semantic conventions. It represents the part number
	// of the part being uploaded in a multipart-upload operation. This is a
	// positive integer between 1 and 10,000.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 3456
	// Note: The `part_number` attribute is only applicable to the
	// [upload-part](https://docs.aws.amazon.com/cli/latest/reference/s3api/upload-part.html)
	// and
	// [upload-part-copy](https://docs.aws.amazon.com/cli/latest/reference/s3api/upload-part-copy.html)
	// operations.
	// The `part_number` attribute corresponds to the `--part-number` parameter
	// of the
	// [upload-part operation within the S3
	// API](https://docs.aws.amazon.com/cli/latest/reference/s3api/upload-part.html).
	AWSS3PartNumberKey = "aws.s3.part_number"

	// AWSS3UploadIDKey is the attribute Key conforming to the
	// "aws.s3.upload_id" semantic conventions. It represents the upload ID
	// that identifies the multipart upload.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'dfRtDYWFbkRONycy.Yxwh66Yjlx.cph0gtNBtJ'
	// Note: The `upload_id` attribute applies to S3 multipart-upload
	// operations and corresponds to the `--upload-id` parameter
	// of the [S3
	// API](https://docs.aws.amazon.com/cli/latest/reference/s3api/index.html)
	// multipart operations.
	// This applies in particular to the following operations:
	//
	// -
	// [abort-multipart-upload](https://docs.aws.amazon.com/cli/latest/reference/s3api/abort-multipart-upload.html)
	// -
	// [complete-multipart-upload](https://docs.aws.amazon.com/cli/latest/reference/s3api/complete-multipart-upload.html)
	// -
	// [list-parts](https://docs.aws.amazon.com/cli/latest/reference/s3api/list-parts.html)
	// -
	// [upload-part](https://docs.aws.amazon.com/cli/latest/reference/s3api/upload-part.html)
	// -
	// [upload-part-copy](https://docs.aws.amazon.com/cli/latest/reference/s3api/upload-part-copy.html)
	AWSS3UploadIDKey = "aws.s3.upload_id"

	// GraphqlDocumentKey is the attribute Key conforming to the
	// "graphql.document" semantic conventions. It represents the GraphQL
	// document being executed.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'query findBookByID { bookByID(id: ?) { name } }'
	// Note: The value may be sanitized to exclude sensitive information.
	GraphqlDocumentKey = "graphql.document"

	// GraphqlOperationNameKey is the attribute Key conforming to the
	// "graphql.operation.name" semantic conventions. It represents the name of
	// the operation being executed.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'findBookByID'
	GraphqlOperationNameKey = "graphql.operation.name"

	// GraphqlOperationTypeKey is the attribute Key conforming to the
	// "graphql.operation.type" semantic conventions. It represents the type of
	// the operation being executed.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'query', 'mutation', 'subscription'
	GraphqlOperationTypeKey = "graphql.operation.type"
)
