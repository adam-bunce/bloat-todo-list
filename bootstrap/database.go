package bootstrap

import (
	"database/sql"
	logr "github.com/adam-bunce/grpc-todo/helpers"
	"github.com/adam-bunce/grpc-todo/variables"
	_ "github.com/jackc/pgx/v5/stdlib" // import for side effect
	"time"
)

func InitDB() {
	logr.Info("Initializing DB from config")
	dsn := variables.GlobalConfig.DbConfig.GetDSN()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Second * 45)

	variables.DB = db
}
