package storage

import (
	"context"
	"errors"
	"fmt"
	"strconv"

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
	rows, err := p.conn.Query(context.Background(), "select raw_protocol, domain, port from proxies")
	if err != nil {
		return []string{}, errors.New(fmt.Sprintf("Query failed: %v\n", err))
	}
	for rows.Next() {
		var raw_protocol int8
		var domain string
		var port int64
		err = rows.Scan(&raw_protocol, &domain, &port)
		if err != nil {
			return []string{}, err
		}
		proxies = append(proxies, fmt.Sprintf("%s://%s:%s", decodeProtocol(raw_protocol), domain, strconv.FormatInt(port, 10)))
	}
	return []string{}, nil
}

func decodeProtocol(rawProto int8) string {
	switch rawProto {
	case 0:
		return "http"
	case 1:
		return "https"
	case 2:
		return "socks"
	}
	return ""
}
