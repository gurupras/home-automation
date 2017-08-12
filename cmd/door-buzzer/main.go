package main

import (
	"io/ioutil"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/gurupras/go-easyfiles"
	"github.com/gurupras/home-automation/door-buzzer"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	app     = kingpin.New("ThermaBox", "Temperature-controller")
	conf    = app.Arg("conf", "Configuration file (YAML)").Required().String()
	verbose = app.Flag("verbose", "Verbose logging").Short('v').Default("false").Bool()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	if !easyfiles.Exists(*conf) {
		log.Fatalf("Configuration file '%v' does not exist")
	}

	db := &door_buzzer.DoorBuzzer{}
	data, err := ioutil.ReadFile(*conf)
	if err != nil {
		log.Fatalf("Failed to read conf file: %v", err)
	}

	if err := yaml.Unmarshal(data, db); err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}
	db.Initialize()
	// At this point, the webserver should've started
	// So just wait forever
	c := make(chan struct{})
	_ = <-c
}
