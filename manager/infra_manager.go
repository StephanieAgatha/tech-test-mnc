package manager

import (
	"database/sql"
	"fmt"
	"github.com/gookit/slog"
	_ "github.com/lib/pq"
	"mnc-test/config"
)

type InfraManager interface {
	Connect() *sql.DB
}

type infraManager struct {
	db  *sql.DB
	cfg *config.Config
}

func (i *infraManager) Connect() *sql.DB {
	return i.db
}

func (i *infraManager) initdb() error {
	//init dsn disini
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		i.cfg.DbConfig.Host,
		i.cfg.DbConfig.Port,
		i.cfg.DbConfig.Username,
		i.cfg.DbConfig.Password,
		i.cfg.DbConfig.DBName,
	)

	//sql open
	db, err := sql.Open(i.cfg.DBDriver, dsn)
	if err != nil {
		slog.Errorf("Failed to open db %v", err.Error())
		return err
	}
	i.db = db
	return nil
}

func NewInfraManager(cfg *config.Config) (InfraManager, error) {
	conn := &infraManager{
		cfg: cfg,
	}
	if err := conn.initdb(); err != nil {
		return nil, fmt.Errorf("Failed on infra manager %v", err.Error())
	}

	return conn, nil
}
