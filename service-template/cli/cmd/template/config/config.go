package config

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

type Commission struct {
	Type  CommissionType `yaml:"type"`
	Value float64        `yaml:"value"`
}

type CommissionType string

const (
	CommissionTypePercentage CommissionType = "PERCENTAGE"
	CommissionTypeNominal    CommissionType = "NOMINAL"
	CommissionTypeNone       CommissionType = "NONE"
)

type (
	StatsDConfig struct {
		Addr string `yaml:"address"`
		Port string `yaml:"port"`
	}

	LogConfig struct {
		Level     string `yaml:"level"`
		Formatter string `yaml:"formatter"`
	}

	MongoDBConfig struct {
		Uri               string        `yaml:"uri"`
		DBName            string        `yaml:"db_name"`
		ConnectionTimeout time.Duration `yaml:"connection_timeout"`
		MaxPool           int           `yaml:"max_pool"`
		MinPool           int           `yaml:"min_pool"`
	}

	MongoConfig struct {
		IntegratorMongoDBConfig MongoDBConfig `yaml:"integrator"`
	}

	IssuedConfig struct {
		DummyIssued bool `yaml:"dummy_issued"`
	}

	KafkaTopics struct {
		IntegratorSearchResponse      string `yaml:"integrator_search_response"`
		IntegratorSearchRequest       string `yaml:"integrator_search_request"`
		IntegratorSearchChunkResponse string `yaml:"integrator_search_chunk_response"`
		AnalyticBAUBooking            string `yaml:"analytic_bau_booking"`
		AnalyticBAUSearch             string `yaml:"analytic_bau_search"`
		AnalyticBAUIssuance           string `yaml:"analytic_bau_issued"`
		AnalyticBAULog                string `yaml:"analytic_bau_log"`
	}

	KafkaServerConfig struct {
		Version       string        `yaml:"version"`
		Host          []string      `yaml:"host"`
		ConsumerGroup string        `yaml:"consumer_group"`
		CaptureMetric bool          `yaml:"capture_metric"`
		ReadInterval  string        `yaml:"read_interval"`
		ReadTimeout   time.Duration `yaml:"read_timeout"`
		CorePool      int           `yaml:"core_pool"`
	}

	KafkaConfig struct {
		FlightConfig  KafkaServerConfig `yaml:"flight"`
		GeneralConfig KafkaServerConfig `yaml:"general"`
		Topics        KafkaTopics       `yaml:"topic"`
	}

	RedisServerConfig struct {
		Database      int           `yaml:"database"`
		Address       []string      `yaml:"address"`
		Password      string        `yaml:"password"`
		PoolSize      int           `yaml:"pool_size"`
		MinIdleConn   int           `yaml:"min_idle_conn"`
		CaptureMetric bool          `yaml:"capture_metric"`
		DialTimeout   time.Duration `yaml:"dial_timeout"`
		PoolTimeout   time.Duration `yaml:"pool_timeout"`
		ReadTimeout   time.Duration `yaml:"read_timeout"`
		WriteTimeout  time.Duration `yaml:"write_timeout"`
		MaxConnAge    time.Duration `yaml:"max_conn_age"`
	}

	RedisConfig struct {
		IntegratorRedisConfig RedisServerConfig `yaml:"integrator"`
	}

	ClientHttpConfig struct {
		Url     string        `yaml:"url"`
		Timeout time.Duration `yaml:"timeout"`
		Retries int           `yaml:"retries"`
		OpName  string        `yaml:"op_name"`
	}

	HttpConfig struct {
		BaseUrl                           string           `yaml:"base_url"`
		FlightAvailClientConfig           ClientHttpConfig `yaml:"flight_avail_client"`
		FlightLowestFareAvailClientConfig ClientHttpConfig `yaml:"flight_lowest_fare_avail_client"`
		GetQuoteSummaryClientConfig       ClientHttpConfig `yaml:"get_quote_summary_client"`
		ServiceInitializeClientConfig     ClientHttpConfig `yaml:"service_initialize_client"`
		FlightAddClientConfig             ClientHttpConfig `yaml:"flight_add_client"`
		BookingGetSessionClientConfig     ClientHttpConfig `yaml:"booking_get_session_client"`
		BookingSaveClientConfig           ClientHttpConfig `yaml:"booking_save_client"`
		BookingGetItineraryClientConfig   ClientHttpConfig `yaml:"booking_get_itinerary_client"`
		AddPaymentClientConfig            ClientHttpConfig `yaml:"add_payment_client"`
	}

	OutboundConfig struct {
		SupplierOutboundHttpConfig HttpConfig `yaml:"supplier"`
	}

	CacheConfig struct {
		DefaultTtlDuration             time.Duration `yaml:"default_ttl_duration"`
		DefaultCleanupIntervalDuration time.Duration `yaml:"default_cleanup_interval_duration"`
	}

	SchedulerConfig struct {
		CredentialIntervalSecs      string `yaml:"credential_interval_secs"`
		SystemParameterIntervalSecs string `yaml:"system_parameter_interval_secs"`
	}

	Config struct {
		HttpServerPort          int                     `yaml:"http_server_port" json:"HttpServerPort,omitempty"`
		ServiceName             string                  `yaml:"service_name" json:"ServiceName,omitempty"`
		IsStaging               bool                    `yaml:"is_staging" json:"IsStaging,omitempty"`
		StatsDConfig            StatsDConfig            `yaml:"statsd" json:"StatsDConfig"`
		LogConfig               LogConfig               `yaml:"log" json:"LogConfig"`
		IssuedConfig            IssuedConfig            `yaml:"issued_config" json:"IssuedConfig"`
		PromotionCodes          string                  `yaml:"promotion_codes" json:"PromotionCodes,omitempty"`
		KafkaConfig             KafkaConfig             `yaml:"kafka" json:"KafkaConfig"`
		RedisConfig             RedisConfig             `yaml:"redis" json:"RedisConfig"`
		MongoConfig             MongoConfig             `yaml:"mongo" json:"MongoConfig"`
		CacheConfig             CacheConfig             `yaml:"cache" json:"CacheConfig"`
		SchedulerConfig         SchedulerConfig         `yaml:"scheduler" json:"SchedulerConfig"`
		CurrencyV1Config        CurrencyV1Config        `yaml:"currency_v1" json:"CurrencyV1Config"`
		OutboundConfig          OutboundConfig          `yaml:"outbound" json:"OutboundConfig"`
	{{CONFIG_PROPERTY_STRUCT}} {{CONFIG_PROPERTY_STRUCT}} `yaml:"{{PROVIDER}}_config_properties"`
		DefaultCurrency         string                  `yaml:"default_currency" json:"DefaultCurrency,omitempty"`
		ContextTimeout          int                     `yaml:"context_timeout" json:"ContextTimeout"`
		DefaultTimeZone         string                  `yaml:"default_timezone" json:"DefaultTimezone,omitempty"`
		BookingExpiredTimeLimit int                     `yaml:"booking_expired_time_limit"`
	}

   {{CONFIG_PROPERTY_STRUCT}} struct {
		Commission Commission `yaml:"commission"`
		PromoCode  string     `yaml:"promo_code"`
	}

	CurrencyV1Config struct {
		CacheKeyCurrencyPrefixIntegrator string        `yaml:"cache_key_currency_prefix_integrator"`
		CacheCurrencyTTLV1               time.Duration `yaml:"cache_currency_ttl_v1"`
	}
)

func New() (*Config, error) {
	var cfg Config

	file, err := os.ReadFile("application.yml")
	if err != nil {
		return nil, err
	}

	file = []byte(os.ExpandEnv(string(file)))
	decoder := yaml.NewDecoder(bytes.NewReader(file))
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	if cfg.ContextTimeout == 0 {
		cfg.ContextTimeout = 5
	}

	return &cfg, nil
}
