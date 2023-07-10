package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/robfig/cron/v3"
	"google.golang.org/grpc/grpclog"
)

var log = grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)

type conf struct {
	Targets []string `yaml:"targets"`
}

func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Info("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func main() {
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	var conf conf
	conf.getConf()
	log.Info(conf.Targets)
	c := cron.New(cron.WithSeconds())

	healthCheckJob, err := c.AddFunc("*/10 * * * * *", func() {
		HealthCheckJob(conf.Targets)
	})

	if err != nil {
		log.Error(err.Error())
		c.Remove(healthCheckJob)
	}

	log.Info("CRON JOB STARTED")
	c.Run()
}

func HealthCheckJob(targets []string) {
	var wg = new(sync.WaitGroup)
	for _, uri := range targets {
		wg.Add(1)

		go func(uri string) {
			defer wg.Done()

			requestURL := fmt.Sprintf("http://%s", uri)
			res, err := http.Get(requestURL)
			if err != nil {
				log.Errorf("%s | FAILED | with error %s\n", uri, err.Error())
				return
			}

			if res.StatusCode >= 200 && res.StatusCode < 400 {
				log.Infof("%s | SUCCESS | %s", uri, res.Status)
			} else {
				log.Errorf("%s | FAILED | %s", uri, res.Status)
			}
		}(uri)
	}

	wg.Wait()
}
