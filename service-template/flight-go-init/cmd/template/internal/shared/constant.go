package shared

const (
	FALSE        = "false"
	TRUE         = "true"
	EMPTY_STRING = ""

	DateTimeFormatRFC3339Zoneless = "2006-01-02T15:04:05"
	TimeFormatHourMinute          = "15:04"

	DISTRIBUTION_TYPE = "{{PROVIDER}}"

	SESSION_ID = "sessionId"

	REQUEST           = "request"
	RESPONSE          = "response"
	PROCESS           = "process"
	BODY              = "body"
	ERROR             = "error"
	ERROR_REQUEST_API = "ERROR_CALL_API"
	ERROR_UNMARSHAL   = "ERROR_UNMARSHALLING"
	ERROR_PUBLISH     = "ERROR_PUBLISH"

	MASKED_STRING = "****"

	// - enum log
	REST_IMPL       = "REST_IMPL"
	SERVICE_IMPL    = "SERVICE_IMPL"
	OUTBOUND_IMPL   = "OUTBOUND_IMPL"
	INBOUND_IMPL    = "INBOUND_IMPL"
	REPOSITORY_IMPL = "REPOSITORY_IMPL"
	SCHEDULER_IMPL  = "SCHEDULER_IMPL"

	// - query parameters key
	Q_EXECUTION_ID_KEY = "executionid"
	Q_RELOC_KEY        = "reloc"
	PNR                = "pnr"

	// - cache
	CACHE_PREFIX = "com.tiket.tix.flight-{{PROVIDER}}-integrator"

	// - predefined
	PREDEFINED_CREDENTIAL_CACHE_KEY_PREFIX = "predefined.credential."
	PREDEFINED_SESSION_CACHE_KEY_TEMPLATE  = "predefined.session.%s"

	// - context key
	CTX_ANALYTIC_TOKEN_KEY = "token"
	COOKIE                 = "Cookie"

	// - system parameter key
	SYS_PARAM_PUBLISH_ON_BOOKING_FLAG_KEY      = "publish-on-booking-flag"
	SYS_PARAM_ADDON_BAGGAGE                    = "addons.baggage"
	SYS_PARAM_SEARCH_SSR_CODES                 = "search.ssr.codes"
	SYS_PARAM_CACHE_SSR_FLAG_KEY               = "cache.ssr.flag"
	SYS_PARAM_KAFKA_SEARCH_LISTENER_CONFIG_KEY = "kafka.search.listener.config"
	SYS_PARAM_CHUNK_FLAG_KEY                   = "publish.chunk.flag"
	SYS_PARAM_CHUNK_SIZE_KEY                   = "publish.chunk.size"
	SYS_PARAM_INCLUSIVE_BAGGAGE_CABIN          = "inclusive.baggage.cabin"
	SYS_PARAM_ENABLE_CABIN_CLASS_LIST          = "enabled.cabin.class.list"

	DEPARTURE_INDEX = 0
	RETURN_INDEX    = 1

	PAX_SUPPLIER_TYPE_ADULT  = "ADT"
	PAX_SUPPLIER_TYPE_CHILD  = "CHD"
	PAX_SUPPLIER_TYPE_INFANT = "INF"

	CURRENCY_IDR = "IDR"
	SPACE        = " "

	PAX_TYPE_ADULT  = "adult"
	PAX_TYPE_CHILD  = "child"
	PAX_TYPE_INFANT = "infant"

	SUPPLIER_PAX_TYPE_ADULT  = "ADULT"
	SUPPLIER_PAX_TYPE_CHILD  = "CHD"
	SUPPLIER_PAX_TYPE_INFANT = "INF"

	ROUND_TRIP = "ROUND_TRIP"
	DEPARTURE  = "DEPARTURE"
	RETURN     = "RETURN"

	PIPE_SEPARATOR              = "|"
	FARE_CLASS_SEPARATOR        = PIPE_SEPARATOR
	CONNECTING_FLIGHT_SEPARATOR = PIPE_SEPARATOR
	TRANSIT_SEPARATOR           = PIPE_SEPARATOR
	FlightNumbersSeparator      = PIPE_SEPARATOR
	FareSellKeySeparator        = PIPE_SEPARATOR

	TildeSeparator        = "~"
	ITINERARY_SEPARATOR   = TildeSeparator
	FlightSelectSeparator = TildeSeparator

	FlightNumberSpacelessFormat            = "%s%s"
	FlightNumberFormat                     = "%s%4s"
	FlightReferenceDateFormat              = "%s"
	FlightReferenceOriginDestinationFormat = "%s%s"
	FlightReferenceFormat                  = FlightReferenceDateFormat + " " + FlightNumberFormat + " " + FlightReferenceOriginDestinationFormat

	DOT_SEPARATOR = "."

	SUPPLIER_CODE = "{{SUPPLIER_CODE}}"

	DATE_FORMAT        = "2025-01-20"
	HOUR_MINUTE_FORMAT = "15:30"
	DATE_TIME_FORMAT   = "2006-01-02T15:04:05"

	// Login-related
	ExtendedDataKeyCurrency       = "currency"
	ExtendedDataKeyDomainCode     = "domainCode"
	ExtendedDataKeyAgentName      = "agentName"
	ExtendedDataKeyAgentPassword  = "agentPassword"
	ExtendedDataKeyOrganizationId = "organizationId"
	ExtendedDataKeyPaymentMethod  = "paymentMethod"
	EXTENDED_DATA_AGENT_CODE      = "agentCode"
	EXTENDED_DATA_AGENT_PASSWORD  = "agentPassword"

	HEADER_AUTHORIZATION = "Authorization"
	LOGGER               = "logger"
	SUPPLIER_ID_KEY      = "supplierId"

	QUOTE_CHARGE_TYPE_FARE = "FARE"

	IWJR_CHARGE_TYPE = "IWJR"

	PSC_CHARGE_TYPE = "PSC"
)
