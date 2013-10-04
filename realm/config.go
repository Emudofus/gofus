package main

import (
	"encoding/json"
	"flag"
	"github.com/Blackrush/gofus/realm/network/backend"
	"github.com/Blackrush/gofus/realm/network/frontend"
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
	workers = flag.Int("workers", invalid_int_flag, "the number of workers to spawn")

	bladdr = flag.String("bladdr", invalid_string_flag, "the address and port the backend service will connect to")
	bpass  = flag.String("bpass", invalid_string_flag, "the password used to secure backend service")

	id         = flag.Int("id", invalid_int_flag, "the id of the realm server")
	addr       = flag.String("addr", invalid_string_flag, "the address of the realm server")
	completion = flag.Int("completion", invalid_int_flag, "the completion of the realm server")
	community  = flag.Int("community", invalid_int_flag, "the community id of the realm server")

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
		ServerId:         1,
		ServerAddr:       "127.0.0.1",
		ServerPort:       5555,
		ServerCompletion: 0,
		Laddr:            ":5554",
		Password:         "",
	}
	cfg.Frontend = frontend.Configuration{
		Port:        5555,
		Workers:     1,
		CommunityId: 1,
		ServerId:    1,
	}
}

func overwrite_config_values(cfg *config) {
	if *data_source_name != invalid_string_flag {
		cfg.Database.DataSourceName = *data_source_name
	}
	if *fport != invalid_int_flag {
		cfg.Frontend.Port = uint16(*fport)
		cfg.Backend.ServerPort = uint16(*fport)
	}
	if *workers != invalid_int_flag {
		cfg.Frontend.Workers = *workers
	}
	if *bladdr != invalid_string_flag {
		cfg.Backend.Laddr = *bladdr
	}
	if *bpass != invalid_string_flag {
		cfg.Backend.Password = *bpass
	}
	if *id != invalid_int_flag {
		cfg.Backend.ServerId = uint(*id)
		cfg.Frontend.ServerId = uint(*id)
	}
	if *addr != invalid_string_flag {
		cfg.Backend.ServerAddr = *addr
	}
	if *completion != invalid_int_flag {
		cfg.Backend.ServerCompletion = *completion
	}
	if *community != invalid_int_flag {
		cfg.Frontend.CommunityId = *community
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
