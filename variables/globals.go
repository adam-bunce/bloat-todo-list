package variables

import (
	"bytes"
	"database/sql"
	logr "github.com/adam-bunce/grpc-todo/helpers"
	"github.com/spf13/viper"
	"text/template"
	"time"
)

var GlobalConfig *Config
var DB *sql.DB

type Config struct {
	*DbConfig
	AppName    string
	ServerPort int
}

type DbConfig struct {
	Name     string
	User     string
	Password string
	Port     string
	Host     string
}

type db interface {
	GetDSN() string
	SetGlobalDB() error
}

func (d *DbConfig) GetDSN() string {
	dsnTemplate := "host={{.Host}} user={{.User}} password={{.Password}} dbname={{.Name}} port={{.Port}} sslmode=disable"
	t, err := template.New("dsnTemplate").Parse(dsnTemplate)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, d); err != nil {
		panic(err)
	}
	return buf.String()
}

func (d *DbConfig) SetGlobalDB() error {
	dsn := d.GetDSN()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}

	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Second * 45)
	DB = db

	return nil
}

func parseDBConfig(database map[string]interface{}) DbConfig {
	return DbConfig{
		Name:     database["Name"].(string),
		User:     database["User"].(string),
		Password: database["Password"].(string),
		Port:     database["Port"].(string),
		Host:     database["Host"].(string),
	}
}

func InitConfig() {
	logr.Info("Initializing Config")

	viper.SetConfigName("config")
	viper.SetConfigType("hcl")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		// cant function w/o config so panic
		panic(err)
	}
	appName := viper.GetString("app_name")
	serverPort := viper.GetInt("server_port")
	databaseConfig := parseDBConfig(viper.Get("database").([]map[string]interface{})[0])

	GlobalConfig = &Config{
		DbConfig:   &databaseConfig,
		AppName:    appName,
		ServerPort: serverPort,
	}

	GlobalConfig.SetGlobalDB()
}
