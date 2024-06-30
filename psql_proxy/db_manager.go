package psql_proxy

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
	"github.com/prusiksoftware/monorepo/obscurio/psql_proxy/modify_queries"
	"github.com/prusiksoftware/monorepo/obscurio/psql_proxy/schema"
	"log"
	"net"
	"os"
)

type dbManager struct {
	name string

	tables []schema.Table

	mutators []modify_queries.ModifierInterface
	username string
	password string
	dbUri    string

	dbConn     *pgconn.PgConn
	clientConn *net.Conn
}

func (sp *dbManager) GetDBConnection() *pgconn.PgConn {
	if sp.dbConn != nil {
		return sp.dbConn
	}

	pgConn, err := pgconn.Connect(context.Background(), sp.dbUri)
	if err != nil {
		log.Fatalln("pgconn failed to connect:", err)
	}
	sp.dbConn = pgConn

	return pgConn
}

func (sp *dbManager) ProxyQuery(query string, clientConnection *net.Conn) error {
	dbConn := sp.GetDBConnection().Conn()
	dbFrontend := pgproto3.NewFrontend(dbConn, dbConn)

	// Send the query to the database
	msg := &pgproto3.Query{String: query}
	dbFrontend.SendQuery(msg)

	// Flush the query to the database
	err := dbFrontend.Flush()
	if err != nil {
		return fmt.Errorf("error sending query: %w", err)
	}

	for {
		// receive a message from the database
		msg, err := dbFrontend.Receive()
		if err != nil {
			return fmt.Errorf("error receiving message: %w", err)
		}

		// convert to bytes
		b, err := msg.Encode(nil)
		if err != nil {
			return fmt.Errorf("error encoding message: %w", err)
		}

		// write the message to the client
		if clientConnection != nil {
			n, err := (*clientConnection).Write(b)
			if err != nil {
				fmt.Printf("error proxying message: %d; %v\n", n, err)
			}
		}

		if _, ok := msg.(*pgproto3.ReadyForQuery); ok {
			return nil
		}
	}
}

func (sp *dbManager) modifyQuery(query string) (string, error) {
	qm, err := modify_queries.NewQueryModifier(query, sp.mutators)
	if err != nil {
		return "", fmt.Errorf("failed to modify query: %w", err)
	}
	err = qm.Modify()
	if err != nil {
		return "", fmt.Errorf("failed to modify query: %w", err)
	}
	return qm.Query()
}

func (sp *dbManager) mutatorNames() []string {
	names := []string{}
	for _, mutator := range sp.mutators {
		names = append(names, mutator.String())
	}
	return names
}

func createDBManager(profile Profile) (*dbManager, error) {
	dbURI, exists := os.LookupEnv(profile.DatabaseEnv)
	if !exists {
		return nil, fmt.Errorf("database environment variable %s not set", profile.DatabaseEnv)
	}

	username, exists := os.LookupEnv(profile.UsernameEnv)
	if !exists {
		return nil, fmt.Errorf("database environment variable %s not set", profile.DatabaseEnv)
	}

	password, exists := os.LookupEnv(profile.PasswordEnv)
	if !exists {
		return nil, fmt.Errorf("database environment variable %s not set", profile.DatabaseEnv)
	}

	dbm := &dbManager{
		name:     profile.Name,
		username: username,
		password: password,
		dbUri:    dbURI,
		tables:   schema.GetTables(dbURI),
	}

	hiddenCols := map[string][]string{}
	for _, filter := range profile.Filters {
		if filter.Function == hideColumn {
			if _, ok := hiddenCols[filter.Table]; !ok {
				hiddenCols[filter.Table] = []string{}
			}
			hiddenCols[filter.Table] = append(hiddenCols[filter.Table], filter.Column)
		}
	}

	mutators := []modify_queries.ModifierInterface{
		modify_queries.NewWildcardExpander(dbm.tables, hiddenCols),
		modify_queries.NewColumnHider(hiddenCols),
	}

	for _, filter := range profile.Filters {
		if filter.Function == hideRow {
			mutators = append(mutators, modify_queries.NewRowHider(filter.Table, filter.Column, modify_queries.NotEqual, filter.Value))

		}
	}

	replacedCols := map[string]map[string]string{}
	for _, filter := range profile.Filters {
		if filter.Function == replaceColumn {
			if _, ok := replacedCols[filter.Table]; !ok {
				replacedCols[filter.Table] = map[string]string{}
			}
			replacedCols[filter.Table][filter.Column] = filter.Value
		}
	}
	for table, columns := range replacedCols {
		mutators = append(mutators, modify_queries.NewColumnReplacer(table, columns))
	}

	dbm.mutators = mutators
	return dbm, nil
}
