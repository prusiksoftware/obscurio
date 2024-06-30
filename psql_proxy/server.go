package psql_proxy

import (
	"context"
	"fmt"
	wire "github.com/jeroenrinzema/psql-wire"
	_ "github.com/lib/pq"
	"github.com/prusiksoftware/monorepo/obscurio/analytics"
	"log"
	"net"
)

type Status string

const (
	Booting Status = "booting"
	Running        = "running"
	Done           = "done"
)

type Server struct {
	config           *Conf
	StatusChan       chan Status
	dbManagers       []*dbManager
	analyticsTracker *analytics.Analytics
}

func NewServer(config *Conf, a *analytics.Analytics) (*Server, error) {
	var serverProfiles []*dbManager
	for _, configProfile := range config.Profiles {
		p, err := createDBManager(configProfile)
		if err != nil {
			return nil, err
		}
		serverProfiles = append(serverProfiles, p)
	}

	s := Server{
		config:           config,
		StatusChan:       make(chan Status),
		dbManagers:       serverProfiles,
		analyticsTracker: a,
	}

	return &s, nil
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", "0.0.0.0:5432")
	if err != nil {
		return err
	}
	defer listener.Close()

	s.StatusChan <- Running
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go s.handleConnection(conn, s.config)
	}
}

func (s *Server) handleConnection(conn net.Conn, config *Conf) {
	log.Println("Accepted connection from", conn.RemoteAddr())
	b := NewClientManager(conn, config)
	err := b.Run()
	if err != nil {
		log.Println(err)
	}
	log.Println("Closed connection from", conn.RemoteAddr())
}

func (s *Server) getDBManager(ctx context.Context) (*dbManager, error) {
	connectionUsername := wire.AuthenticatedUsername(ctx)
	for _, profile := range s.dbManagers {
		if profile.username == connectionUsername {
			return profile, nil
		}
	}

	return nil, fmt.Errorf("no configProfile found for user '%s'", connectionUsername)
}
