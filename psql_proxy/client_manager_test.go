package psql_proxy

import (
	"github.com/jackc/pgx/v5/pgproto3"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestNewClientManager(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		a, _ := net.Pipe()
		config, err := GetConfig()
		require.NoError(t, err)

		manager := NewClientManager(a, config)

		require.NotNil(t, manager)
	})
}

func TestClientManager_Run(t *testing.T) {
	t.Run("bad profile", func(t *testing.T) {
		a, aa := net.Pipe()
		config, err := GetConfig()
		require.NoError(t, err)
		manager := NewClientManager(a, config)
		require.NotNil(t, manager)

		go func() {
			manager.Run()
		}()

		startupMsg := pgproto3.StartupMessage{
			ProtocolVersion: 196608,
			Parameters: map[string]string{
				"user":     "nonexistant",
				"database": "test",
			},
		}
		b, err := startupMsg.Encode(nil)
		require.NoError(t, err)
		_, err = aa.Write(b)
		require.NoError(t, err)
	})

	t.Run("close conn", func(t *testing.T) {
		a, aa := net.Pipe()
		config, err := GetConfig()
		require.NoError(t, err)
		manager := NewClientManager(a, config)
		require.NotNil(t, manager)

		go func() {
			err = manager.Run()
			require.Error(t, err)
			require.Equal(t, "EOF", err.Error())
		}()

		err = aa.Close()
		require.NoError(t, err)

	})

	t.Run("bad query", func(t *testing.T) {
		a, aa := net.Pipe()
		config, err := GetConfig()
		require.NoError(t, err)
		manager := NewClientManager(a, config)
		require.NotNil(t, manager)

		go func() {
			manager.Run()
		}()

		startupMsg := pgproto3.StartupMessage{
			ProtocolVersion: 196608,
			Parameters: map[string]string{
				"user":     "user",
				"database": "test",
			},
		}
		b, err := startupMsg.Encode(nil)
		require.NoError(t, err)
		_, err = aa.Write(b)
		require.NoError(t, err)

		// auth ok
		_, err = aa.Read(b)
		require.NoError(t, err)

		// parameter status
		_, err = aa.Read(b)
		require.NoError(t, err)

		// ready
		_, err = aa.Read(b)
		require.NoError(t, err)

		// no errors when the query is bad, just log it
		q := pgproto3.Query{
			String: "bad query",
		}
		b, err = q.Encode(nil)
		require.NoError(t, err)
		_, err = aa.Write(b)
	})

	t.Run("good query", func(t *testing.T) {
		a, aa := net.Pipe()
		config, err := GetConfig()
		require.NoError(t, err)
		manager := NewClientManager(a, config)
		require.NotNil(t, manager)

		go func() {
			manager.Run()
		}()

		startupMsg := pgproto3.StartupMessage{
			ProtocolVersion: 196608,
			Parameters: map[string]string{
				"user":     "user",
				"database": "test",
			},
		}
		b, err := startupMsg.Encode(nil)
		require.NoError(t, err)
		_, err = aa.Write(b)
		require.NoError(t, err)

		// auth ok
		_, err = aa.Read(b)
		require.NoError(t, err)

		// parameter status
		_, err = aa.Read(b)
		require.NoError(t, err)

		// ready
		_, err = aa.Read(b)
		require.NoError(t, err)

		// no errors when the query is bad, just log it
		q := pgproto3.Query{
			String: "select * from customer",
		}
		b, err = q.Encode(nil)
		require.NoError(t, err)
		_, err = aa.Write(b)
		require.NoError(t, err)
	})

}
