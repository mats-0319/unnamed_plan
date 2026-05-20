package mconfig

type Internal struct {
	HTTP *ConfigHTTP `json:"http"`
	DB   *ConfigDB   `json:"db"`
}

type ConfigHTTP struct {
	Port int `json:"port"`
}

type ConfigDB struct {
	DSN string `json:"dsn"`
}

func GetHTTPConfig() *ConfigHTTP {
	return conf.Internal.HTTP
}

func GetDBConfig() *ConfigDB {
	return conf.Internal.DB
}
