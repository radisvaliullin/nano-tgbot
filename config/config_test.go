package config

import "testing"

func Test_NewAppConfig(t *testing.T) {

	c, err := NewAppConfig()
	if err != nil {
		t.Fatal("config test: new app config err - ", err)
	}
	t.Logf("config test: new app config: log - %+v; bot - %+v", c.Log, c.Bot)
}
