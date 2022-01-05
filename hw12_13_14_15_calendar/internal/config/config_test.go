package config

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	file, _ := ioutil.TempFile("/tmp", "config.*.yaml")
	_, err := file.Write([]byte(
		`app:
  storage:
    name: sql
    sql:
      dsn: "postgres://postgres:5432/calendar?sslmode=disable"
      maxIdleConns: 2
      maxOpenConns: 2
logger:
  output: stderr
  level: debug
  devMode: true
server:
  addr: ":8080"
`))
	require.NoError(t, err)
	conf, err := NewConfig(file.Name())

	require.NoError(t, err)
	require.Equal(t, "sql", conf.App.Storage.Name)
	require.Equal(t, "postgres://postgres:5432/calendar?sslmode=disable", conf.App.Storage.SQL.DSN)
	require.Equal(t, 2, conf.App.Storage.SQL.MaxIdleConns)
	require.Equal(t, 2, conf.App.Storage.SQL.MaxOpenConns)
	require.Equal(t, []string{"stderr"}, conf.Logger.Output)
	require.Equal(t, "debug", conf.Logger.Level)
	require.Equal(t, true, conf.Logger.DevMode)
	require.Equal(t, ":8080", conf.Server.Addr)

	_ = os.Remove(file.Name())
}
