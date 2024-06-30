package psql_proxy

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetConfig(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		config, err := GetConfig()

		require.NoError(t, err)
		require.NotNil(t, config)
	})

	t.Run("singleton", func(t *testing.T) {
		config1, err := GetConfig()
		config2, err := GetConfig()

		require.NoError(t, err)
		require.Equal(t, config1, config2)
	})

	t.Run("no env var", func(t *testing.T) {
		t.Setenv("CONFIG_FILEPATH", "")

		singletonConf = nil
		config, err := GetConfig()

		require.Nil(t, config)
		require.Error(t, err)
	})

	t.Run("no file", func(t *testing.T) {
		t.Setenv("CONFIG_FILEPATH", "/app/nonexistent.yaml")

		singletonConf = nil
		config, err := GetConfig()

		require.Nil(t, config)
		require.Error(t, err)
	})

	t.Run("unmarshal error", func(t *testing.T) {
		t.Setenv("CONFIG_FILEPATH", "/app/main.go")

		singletonConf = nil
		config, err := GetConfig()

		require.Nil(t, config)
		require.Error(t, err)
	})

}

func TestConf_getProfile(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		config, err := GetConfig()
		require.NoError(t, err)

		profile, err := config.getProfile("user")
		require.NoError(t, err)
		require.NotNil(t, profile)
	})

	t.Run("not found", func(t *testing.T) {
		config, err := GetConfig()
		require.NoError(t, err)

		profile, err := config.getProfile("nonexistent")
		require.Error(t, err)
		require.Nil(t, profile)
	})
}

func TestConf_debugLog(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		config, err := GetConfig()
		require.NoError(t, err)

		config.debugLog("test")
	})
}

func TestConf_infoLogQuery(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		config, err := GetConfig()
		require.NoError(t, err)

		config.infoLogQuery("user", "select * from customer;")
	})
}
