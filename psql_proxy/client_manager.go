package psql_proxy

import (
	"errors"
	"fmt"
	"github.com/jackc/pgproto3/v2"
	"io"
	"log"
	"net"
	"strings"
)

type QueryHandlerFunc func(query string, conn *net.Conn) error
type ClientManager struct {
	config        *Conf
	clientConn    net.Conn
	clientBackend *pgproto3.Backend
	configProfile *Profile
	serverProfile *dbManager
}

func NewClientManager(conn net.Conn, config *Conf) *ClientManager {
	backend := pgproto3.NewBackend(
		pgproto3.NewChunkReader(conn),
		conn,
	)
	return &ClientManager{
		config:        config,
		clientConn:    conn,
		clientBackend: backend,
		configProfile: nil,
		serverProfile: nil,
	}
}

func (p *ClientManager) Run() error {

	profile, err := p.handleStartup()
	if err != nil {
		return err
	}
	p.configProfile = profile

	dbManager, err := createDBManager(*profile)
	if err != nil {
		return err
	}
	p.serverProfile = dbManager

	for {
		msg, err := p.clientBackend.Receive()
		if errors.Is(err, io.ErrUnexpectedEOF) {
			log.Printf("connection closed by client")
			return err
		}

		queryMsg, ok := msg.(*pgproto3.Query)
		if ok {
			originalQuery := queryMsg.String

			modifiedQuery, err := dbManager.modifyQuery(originalQuery)
			if err != nil {
				log.Printf("failed to modify query: %v", err)
				return err
			}
			p.config.infoLogQuery("client query:", originalQuery)
			p.config.debugLog("with mutators: %v", strings.Join(dbManager.mutatorNames(), ", "))
			p.config.infoLogQuery("server query:", modifiedQuery)

			err = p.serverProfile.ProxyQuery(modifiedQuery, &p.clientConn)
			if err != nil {
				log.Printf("failed to proxy query: %v", err)
				return err
			}
		}

	}
}

func (p *ClientManager) handleStartup() (*Profile, error) {
	for {
		msg, err := p.clientBackend.ReceiveStartupMessage()
		if err != nil {
			log.Printf("failed to receive startup message: %v", err)
			return nil, err
		}

		switch startupMsg := msg.(type) {

		case *pgproto3.SSLRequest:
			_, err = p.clientConn.Write([]byte("N"))
			if err != nil {
				return nil, fmt.Errorf("error sending deny SSL request: %w", err)
			}

		case *pgproto3.StartupMessage:
			username := startupMsg.Parameters["user"]
			profile, err := p.config.getProfile(username)
			if err != nil {
				log.Printf("failed to get configProfile: %v", err)
				return nil, err
			}

			authOk := &pgproto3.AuthenticationOk{}
			err = p.clientBackend.Send(authOk)
			if err != nil {
				log.Printf("failed to send AuthenticationOk: %v", err)
				return nil, err
			}

			err = p.clientBackend.Send(&pgproto3.ParameterStatus{
				Name:  "server_version",
				Value: "16.3",
			})
			if err != nil {
				log.Printf("failed to send ParameterStatus: %v", err)
			}

			// Send ReadyForQuery
			readyForQuery := &pgproto3.ReadyForQuery{TxStatus: 'I'}
			err = p.clientBackend.Send(readyForQuery)
			if err != nil {
				log.Printf("failed to send ReadyForQuery: %v", err)
				return nil, err
			}

			return profile, nil

		default:
			log.Printf("received unexpected message: %#v", startupMsg)
			log.Printf("%v\n", msg)
			return nil, errors.New("unexpected message")
		}
	}
}

func (p *ClientManager) Close() error {

	return nil
	//return p.clientConn.Close()
}
