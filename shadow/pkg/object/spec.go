package object

type (
	MySQL struct {
		Hosts    []string
		username string
		password string
	}

	Kafka struct {
		Hosts []string
	}

	Redis struct {
		Hosts    []*string
		username string
		password string
	}

	RabbitMQ struct {
		Hosts    []string
		username string
		password string
	}

	ElasticSearch struct {
		Hosts    []string
		username string
		password string
	}

	ShadowService struct {
		Name          string
		ServiceName   string
		NameSpace     string
		MySQL         *MySQL
		Redis         *Redis
		Kafka         *Kafka
		RabbitMQ      *RabbitMQ
		ElasticSearch *ElasticSearch
	}
)
