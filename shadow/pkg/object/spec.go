package object

type (
	MySQL struct {
		Uris     []string `yaml:"uris" json:"uris" jsonschema:"required"`
		UserName string   `yaml:"userName" json:"userName" jsonschema:"required"`
		Password string   `yaml:"password" json:"password" jsonschema:"required"`
	}

	Kafka struct {
		Uris []string `yaml:"uris" json:"uris" jsonschema:"required"`
	}

	Redis struct {
		Uris     []string `yaml:"uris" json:"uris" jsonschema:"required"`
		UserName string   `yaml:"userName" json:"userName" jsonschema:"required"`
		Password string   `yaml:"password" json:"password" jsonschema:"required"`
	}

	RabbitMQ struct {
		Uris     []string `yaml:"uris" json:"uris" jsonschema:"required"`
		UserName string   `yaml:"userName" json:"userName" jsonschema:"required"`
		Password string   `yaml:"password" json:"password" jsonschema:"required"`
	}

	ElasticSearch struct {
		Uris     []string `yaml:"uris" json:"uris" jsonschema:"required"`
		UserName string   `yaml:"userName" json:"userName" jsonschema:"required"`
		Password string   `yaml:"password" json:"password" jsonschema:"required"`
	}

	ShadowService struct {
		Name          string         `yaml:"name" jsonschema:"required"`
		ServiceName   string         `yaml:"serviceName" jsonschema:"required"`
		Namespace     string         `yaml:"namespace" jsonschema:"required"`
		MySQL         *MySQL         `yaml:"mysql" jsonschema:"omitempty"`
		Kafka         *Kafka         `yaml:"kafka" jsonschema:"omitempty"`
		Redis         *Redis         `ymal:"redis" jsonschema:"omitempty"`
		RabbitMQ      *RabbitMQ      `ymal:"rabbitMq" jsonschema:"omitempty"`
		ElasticSearch *ElasticSearch `yaml:"elasticSearch" jsonschema:"omitempty"`
	}

	CustomObjectKind struct {
		Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
		// JSONSchema is the json schema to validate a custom object of this kind
		JsonSchema string `protobuf:"bytes,2,opt,name=jsonSchema,proto3" json:"jsonSchema,omitempty"`
	}
)
