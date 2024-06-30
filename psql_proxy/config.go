package psql_proxy

import (
	"errors"
	"fmt"
	"github.com/alecthomas/chroma/v2/quick"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type ActionType string
type LogLevel string

const (
	hideColumn ActionType = "hide column"
	hideRow               = "hide row"

	replaceColumn = "replace column"
)

const (
	debug LogLevel = "debug"
	info           = "info"
)

var singletonConf *Conf

type DataFilter struct {
	Function ActionType `yaml:"function"`
	Table    string     `yaml:"table"`
	Column   string     `yaml:"column,omitempty"`
	Value    string     `yaml:"value,omitempty"`
}

type Profile struct {
	Name        string `yaml:"name"`
	DatabaseEnv string `yaml:"database_env"`

	UsernameEnv string       `yaml:"username_env"`
	PasswordEnv string       `yaml:"password_env"`
	Username    string       `yaml:"-"`
	password    string       `yaml:"-"`
	Filters     []DataFilter `yaml:"filters"`
}

type Conf struct {
	PostgresVersion string    `yaml:"postgres_version"`
	LogLevel        LogLevel  `yaml:"log_level"`
	Profiles        []Profile `yaml:"profiles"`

	logger     *log.Logger
	loggerLock *sync.Mutex
}

func GetConfig() (*Conf, error) {
	if singletonConf != nil {
		return singletonConf, nil
	}

	filepath := os.Getenv("CONFIG_FILEPATH")
	fmt.Printf("filepath: %s\n", filepath)
	if filepath == "" {
		return nil, errors.New("CONFIG_FILEPATH not set")
	}

	// get file content
	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		return nil, errors.New("Unable to read file")
	}

	// unmarshal to struct
	c := &Conf{
		loggerLock: &sync.Mutex{},
		logger:     log.New(os.Stdout, "psql_proxy: ", log.LstdFlags),
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, errors.New("Unable to unmarshal file")
	}

	for i, p := range c.Profiles {
		c.Profiles[i].Username = os.Getenv(p.UsernameEnv)
		c.Profiles[i].password = os.Getenv(p.PasswordEnv)
	}
	singletonConf = c
	return c, nil
}

func (c *Conf) getProfile(username string) (*Profile, error) {
	for _, profile := range c.Profiles {
		if profile.Username == username {
			return &profile, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("configProfile '%s' not found", username))
}

func (c *Conf) debugLog(format string, a ...any) {
	c.loggerLock.Lock()
	str := fmt.Sprintf(format, a...)
	if c.LogLevel == debug {
		c.logger.Println(str)
	}
	c.loggerLock.Unlock()
}

func (c *Conf) infoLogQuery(prefix, query string) {
	shouldLog := false
	validStatuses := []LogLevel{"info", "debug"}
	for _, vs := range validStatuses {
		if c.LogLevel == vs {
			shouldLog = true
		}
	}
	if !shouldLog {
		return
	}

	var builder strings.Builder
	var w io.Writer = &builder
	err := quick.Highlight(w, query, "postgresql", "terminal256", "monokai")
	if err != nil {
		c.logger.Fatal(err)
	}

	c.loggerLock.Lock()
	fmt.Printf("%s %s\n", prefix, builder.String())
	c.loggerLock.Unlock()
}
