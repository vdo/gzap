package gzap

import (
	"testing"

	graylog "github.com/Devatoria/go-graylog"
)

func TestInitLogger(t *testing.T) {
	type args struct {
		graylogAppName     string
		isTestEnv          bool
		graylogHost        string
		graylogHandlerType graylog.Transport
		graylogLogEnvName  string
		jsonformatter      bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     string
	}{
		{
			"InitLogger should return a noop logger when running a test",
			args{
				isTestEnv: true,
			},
			false,
			"",
		},
		{
			"InitLogger should return a dev logger when no GRAYLOG_HOST is set",
			args{},
			false,
			"",
		},
		{
			"InitLogger should return a json logger when ENABLE_DATADOG_JSON_FORMATTER is set",
			args{
				jsonformatter: true,
			},
			false,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Instaniate new MockEnvConfig.
			cfg := MockEnvConfig{}
			cfg.On("enableJSONFormatter").Return(tt.args.jsonformatter)
			cfg.On("getGraylogAppName").Return(tt.args.graylogAppName)
			cfg.On("getIsTestEnv").Return(tt.args.isTestEnv)
			cfg.On("getGraylogHost").Return(tt.args.graylogHost)
			cfg.On("getGraylogHandlerType").Return(tt.args.graylogHandlerType)
			cfg.On("getGraylogLogEnvName").Return(tt.args.graylogLogEnvName)
			cfg.On("useColoredConsolelogs").Return(true)

			err := initLogger(&cfg, false)

			if tt.wantErr && err == nil {
				t.Errorf("initLogger() expected error = \"%v\"; got \"nil\"", tt.err)
			}

			if err != nil && err.Error() != tt.err {
				t.Errorf("initLogger() expected error = \"%v\";  got \"%v\"", tt.err, err.Error())
			}
		})
	}
}

func ExampleInitLogger() {
	if err := InitLogger(); err != nil {
		panic(err)
	}

	defer Logger.Sync()

	Logger.Info("this is a test info log")
}
