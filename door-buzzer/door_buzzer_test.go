package door_buzzer

import (
	"testing"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/gurupras/home-automation/web-api"
	"github.com/gurupras/thermabox"
	"github.com/stretchr/testify/require"
)

func TestYAMLUnmarshal(t *testing.T) {
	require := require.New(t)

	str := `
relay:
  active_high: false
  pins: [14]
buzz_period: 1000000000
web_api:
  port: 8080
`
	db := DoorBuzzer{}
	db.Relay = &thermabox.FakeRelay{}
	err := yaml.Unmarshal([]byte(str), &db)
	require.Nil(err)

	fk := thermabox.NewFakeRelay(false, []int{14})
	require.Equal(fk, db.Relay)
	require.Equal(time.Duration(1*time.Second), db.buzzPeriod)

	wb := &web_api.WebAPI{8080}
	require.Equal(wb, db.WebAPI)
}
