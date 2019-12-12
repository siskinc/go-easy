package class

type MysqlConfig struct {
	Hostname string `json:"hostname"`
	Password string `json:"password"`
	Username string `json:"username"`
	Database string `json:"database"`
	Port     string `json:"port"`
}
