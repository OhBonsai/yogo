package model

const (
	DATABASE_DRIVER_SQLITE   = "sqlite3"
	DATABASE_DRIVER_MYSQL    = "mysql"
	DATABASE_DRIVER_POSTGRES = "postgres"

	SQL_SETTINGS_DEFAULT_DATA_SOURCE = "postgres://bonsai:123456@127.0.0.1:54321/bonsai?sslmode=disable"
	SERVICE_SETTINGS_DEFAULT_SITE_URL           = ""
	SERVICE_SETTINGS_DEFAULT_LISTEN_AND_ADDRESS = ":8080"
)

type Config struct {
	ServiceSettings ServiceSettings
	SqlSettings SqlSettings
	LogSettings LogSettings
	JobSettings JobSettings
}

func (c *Config) SetDefaults() {
	c.SqlSettings.SetDefaults()
	c.SqlSettings.SetDefaults()
	c.LogSettings.SetDefaults()
	c.JobSettings.SetDefaults()
}

func (c *Config) IsValid() *AppError{
	return nil
}

type ServiceSettings struct {
	SiteURL                                           *string
	ListenAddress                                     *string
}

func (s *ServiceSettings) SetDefaults() {
	if s.SiteURL == nil {
		s.SiteURL = NewString(SERVICE_SETTINGS_DEFAULT_SITE_URL)
	}

	if s.ListenAddress == nil {
		s.ListenAddress = NewString(SERVICE_SETTINGS_DEFAULT_LISTEN_AND_ADDRESS)
	}

}


type SqlSettings struct {
	DriverName               *string
	DataSource               *string
	DataSourceReplicas       []string
	DataSourceSearchReplicas []string
	MaxIdleConns             *int
	MaxOpenConns             *int
	Trace                    bool
	AtRestEncryptKey         string
	QueryTimeout             *int
}

func (s *SqlSettings) SetDefaults() *SqlSettings{
	if s.DriverName == nil {
		s.DriverName = NewString(DATABASE_DRIVER_POSTGRES)
	}

	if s.DataSource == nil {
		s.DataSource = NewString(SQL_SETTINGS_DEFAULT_DATA_SOURCE)
	}

	if len(s.AtRestEncryptKey) == 0 {
		s.AtRestEncryptKey = NewRandomString(32)
	}

	if s.MaxIdleConns == nil {
		s.MaxIdleConns = NewInt(20)
	}

	if s.MaxOpenConns == nil {
		s.MaxOpenConns = NewInt(300)
	}

	if s.QueryTimeout == nil {
		s.QueryTimeout = NewInt(30)
	}

	return s
}

type LogSettings struct {
	EnableConsole          bool
	ConsoleLevel           string
	ConsoleJson            *bool
	EnableFile             bool
	FileLevel              string
	FileJson               *bool
	FileLocation           string
	EnableWebhookDebugging bool
	EnableDiagnostics      *bool
}

func (s *LogSettings) SetDefaults() {
	if s.EnableDiagnostics == nil {
		s.EnableDiagnostics = NewBool(true)
	}

	if s.ConsoleJson == nil {
		s.ConsoleJson = NewBool(true)
	}

	if s.FileJson == nil {
		s.FileJson = NewBool(true)
	}
}

type JobSettings struct {
	RunJobs      *bool
	RunScheduler *bool
}

func (s *JobSettings) SetDefaults() {
	if s.RunJobs == nil {
		s.RunJobs = NewBool(true)
	}

	if s.RunScheduler == nil {
		s.RunScheduler = NewBool(true)
	}
}


type ConfigFunc func() *Config





