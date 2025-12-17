package db

import (
	"context"
	"database/sql"
	"fgw_web_aforms/internal/config"
	"fgw_web_aforms/pkg/common"
	msg2 "fgw_web_aforms/pkg/common/msg"
	"fmt"
	"log"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

func NewConnMSSQL(ctx context.Context, configDB *config.MSSQLCfg, logger *common.Logger) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("%s://%s:%s@%s?database=%s&charset=%s",
		configDB.MSSQL.Driver,
		configDB.MSSQL.User,
		configDB.MSSQL.Passwd,
		configDB.MSSQL.Server,
		configDB.MSSQL.Name,
		configDB.MSSQL.Charset)
	db, err := sql.Open("mssql", dataSourceName)
	if err != nil {
		logger.LogE(msg2.E3200, err)

		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err = db.PingContext(pingCtx); err != nil {
		Close(db)
		log.Printf("%s: %v", msg2.E3201, err)

		return nil, err
	}

	return db, nil
}

func Close(db *sql.DB) {
	if db == nil {
		return
	}

	if err := db.Close(); err != nil {
		log.Printf("%s: %v", msg2.E3201, err)

		return
	}
	log.Printf(msg2.I2200)
}

func RowsClose(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Printf("%s: %v", msg2.E3203, err)
	}
	log.Printf(msg2.I2201)
}
