package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
)

type Postgresql struct {
	conn   *pgx.Conn
	DBURL  string
	DBUser string
	DBPass string
	DBName string
}

func (p *Postgresql) Open() error {
	var err error
	connString := fmt.Sprintf("postgresql://%s:%s@%s/%s", p.DBUser, p.DBPass, p.DBURL, p.DBName)
	p.conn, err = pgx.Connect(context.Background(), connString)
	if err != nil {
		errors.New(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}
	return nil
}

func (p *Postgresql) Load() ([]string, error) {
	var proxies []string
	rows, err := p.conn.Query(context.Background(), "select protocol, ip_addr, port from good_proxy")
	if err != nil {
		return []string{}, errors.New(fmt.Sprintf("Query failed: %v\n", err))
	}
	for rows.Next() {
		var protocol string
		var ip string
		var port string
		err = rows.Scan(&protocol, &ip, &port)
		if err != nil {
			return []string{}, err
		}
		proxies = append(proxies, fmt.Sprintf("%s://%s:%s", protocol, ip, port))
	}
	return []string{}, nil
}
