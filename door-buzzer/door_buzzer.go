package door_buzzer

import (
	"fmt"
	"net/http"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/gurupras/home-automation/web-api"
	"github.com/gurupras/thermabox"
)

type DoorBuzzer struct {
	Relay           thermabox.RelayInterface `yaml:"relay"`
	buzzPeriod      time.Duration            `yaml:"buzz_period"`
	*web_api.WebAPI `yaml:"web_api"`
}

func (d *DoorBuzzer) UnmarshalYAML(unmarshal func(i interface{}) error) error {
	m := make(map[string]interface{})
	if err := unmarshal(&m); err != nil {
		return err
	}
	relayUnmarshaler := func(i interface{}) error {
		b, _ := yaml.Marshal(m["relay"])
		return yaml.Unmarshal(b, i)
	}
	if d.Relay == nil {
		d.Relay = &thermabox.Relay{}
	}
	if err := d.Relay.UnmarshalYAML(relayUnmarshaler); err != nil {
		return err
	}
	if _, ok := m["buzz_period"]; !ok {
		d.buzzPeriod = 1 * time.Second
	} else {
		d.buzzPeriod = time.Duration(uint64(m["buzz_period"].(int)))
	}
	if _, ok := m["web_api"]; ok {
		webApiUnmarshaler := func(i interface{}) error {
			b, _ := yaml.Marshal(m["web_api"])
			return yaml.Unmarshal(b, i)
		}
		d.WebAPI = &web_api.WebAPI{}
		if err := webApiUnmarshaler(d.WebAPI); err != nil {
			return err
		}
	}
	return nil
}

func (d *DoorBuzzer) Initialize() error {
	// Start web API if needed
	// Set up a POST method to buzz
	if d.WebAPI != nil {
		d.WebAPI.Initialize()
		d.WebAPI.HandleFunc("/buzz", func(w http.ResponseWriter, req *http.Request) {
			d.Buzz()
		})
		go d.WebAPI.Start()
	}
	return nil
}

func (d *DoorBuzzer) Buzz() error {
	if err := d.Relay.On(0); err != nil {
		return fmt.Errorf("Failed to turn on relay: %v", err)
	}
	time.Sleep(d.buzzPeriod)
	if err := d.Relay.Off(0); err != nil {
		return fmt.Errorf("Failed to turn off relay: %v", err)
	}
	return nil
}
