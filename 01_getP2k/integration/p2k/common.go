package p2k

var config = map[string]string{
	"LOG_LEVEL":            "debug",
	"ENVIRONMENT":          "local",
	"DB_P2K_DRIVER":        "oracle",
	"DB_P2K_USER":          "app_stock_worker",
	"DB_P2K_PASSWORD":      "rkrCuEw3j",
	"DB_P2K_HOST":          "172.16.154.105",
	"DB_P2K_NAME":          "p2k",
	"RECEIPT_GET_INTERVAL": "10",
}

func GetEnv(configKey string) string {
	return config[configKey]
}
