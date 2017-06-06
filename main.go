package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"runtime"
	"sync"
	"time"
)

type Check struct {
	Command  string   `json:"check"`
	Args     []string `json:"args"`
	Interval int      `json:"interval"`
}

type Checks []Check

type Config struct {
	Checks Checks `json:"checks"`
}

var wg sync.WaitGroup

func init() {

	runtime.GOMAXPROCS(runtime.NumCPU())

}

func main() {

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/ckicka/")
	viper.AddConfigPath("$HOME/.chicka")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	done := make(chan error)

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {

		cfg := Config{}

		err := viper.Unmarshal(&cfg)
		if err != nil {
			done <- err
		}

		for _, check := range cfg.Checks {

			err := check.Validate()
			if err != nil {
				done <- err
			}

			go func(check Check) {

				for true {
					fmt.Println(check.Command)
					time.Sleep(time.Duration(check.Interval) * time.Second)
				}

			}(check)

		}
	})

	panic(<-done)
}

func (c *Check) Validate() error {

	if c.Interval < 5 {
		return errors.New("the interval must be greater than 5")
	}

	return nil
}