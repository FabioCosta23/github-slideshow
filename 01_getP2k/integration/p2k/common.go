package p2k

var config = map[string]string{
	"LOG_LEVEL":       "debug",
	"ENVIRONMENT":     "local",
	"DB_P2K_DRIVER":   "oracle",
	"DB_P2K_USER":     "app_stock_worker",
	"DB_P2K_PASSWORD": "rkrCuEw3j",
	"DB_P2K_HOST":     "172.16.154.105",
	"DB_P2K_NAME":     "p2k",

	"RECEIPT_GET_INTERVAL":     "10",
	"SENDER_ID":                "P2K",
	"KAFKA_HOST":               "localhost:9092",
	"KAFKA_USERNAME":           "",
	"KAFKA_PASSWORD":           "",
	"KAFKA_INVOICE_TOPIC_NAME": "invoice",

	// stock database
	"DB_HOST":              "127.0.0.1",
	"DB_USER":              "postgres",
	"DB_PASSWORD":          "postgres",
	"DB_PORT":              "5432",
	"DB_NAME":              "postgres",
	"DB_MAX_OPEN_CONNS":    "10",
	"DB_MAX_IDLE_CONNS":    "5",
	"DB_CONN_MAX_LIFETIME": "300",
}

const ErrPrefix = "[P2K Service] "

func GetEnv(configKey string) string {
	return config[configKey]
}
