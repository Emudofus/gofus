package main

import (
	"encoding/json"
	"flag"
	"github.com/Blackrush/gofus/login/network/backend"
	"github.com/Blackrush/gofus/login/network/frontend"
	"github.com/Blackrush/gofus/shared/db"
	"os"
)

const (
	invalid_string_flag = "this flag will be ignored if not set"
	invalid_int_flag    = -1
)

var (
	cfg_file = flag.String("cfg", "./config.json", "the location where to load the configuration")

	fport   = flag.Int("fport", invalid_int_flag, "the port the frontend server will listen on")
	bport   = flag.Int("bport", invalid_int_flag, "the port the backend server will listen on")
	bpass   = flag.String("bpass", invalid_string_flag, "the password used to secure the backend server")
	workers = flag.Int("workers", invalid_int_flag, "the number of workers to start")

	data_source_name = flag.String("dsn", invalid_string_flag, "the source parameters used to connect to the PostgreSQL database")
)

type config struct {
	Database db.Configuration
	Backend  backend.Configuration
	Frontend frontend.Configuration
}

func set_default_config_values(cfg *config) {
	cfg.Database = db.Configuration{
		Driver:         "postgres",
		DataSourceName: "user=postgres dbname=gofus password='' sslmode=disable",
	}
	cfg.Backend = backend.Configuration{
		Port:     5554,
		Password: "",
	}
	cfg.Frontend = frontend.Configuration{
		Port:    5555,
		Workers: 1,
	}
}

func overwrite_config_values(cfg *config) {
	if *data_source_name != invalid_string_flag {
		cfg.Database.DataSourceName = *data_source_name
	}
	if *bport != -1 {
		cfg.Backend.Port = uint16(*bport)
	}
	if *bpass != invalid_string_flag {
		cfg.Backend.Password = *bpass
	}
	if *fport != invalid_int_flag {
		cfg.Frontend.Port = uint16(*fport)
	}
	if *workers != invalid_int_flag {
		cfg.Frontend.Workers = *workers
	}
}

func load_config() *config {
	cfg := config{}

	if file, err := os.Open(*cfg_file); err == nil {
		defer file.Close()

		json.NewDecoder(file).Decode(&cfg)
		overwrite_config_values(&cfg)
	} else if file, err := os.Create(*cfg_file); err == nil {
		defer file.Close()

		set_default_config_values(&cfg)
		if data, err := json.MarshalIndent(cfg, "", "  "); err == nil {
			if _, err := file.Write(data); err != nil {
				panic(err.Error())
			}
		} else {
			panic(err.Error())
		}

		println("default configuration file created")
	} else {
		panic(err.Error())
	}

	return &cfg
}
