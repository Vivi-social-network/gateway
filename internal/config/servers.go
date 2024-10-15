package config

import "time"

type HTTPServer struct {
	UnescapePath      bool          `yaml:"unescape_path"`
	BodyLimit         int           `yaml:"body_limit"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
	AppName           string        `yaml:"app_name"`
	EnablePprof       bool          `yaml:"pprof_enabled"`
	Address           string        `yaml:"address"`
	EnablePrintRoutes bool          `yaml:"enable_print_routes"`
}

type Servers struct {
	HTTP HTTPServer `yaml:"http"`
}
