package config

import "testing"
import "github.com/herb-go/util/cli/app"

func TestConfig(t *testing.T) {
	a := app.NewApplication(Config)
	err := a.ShowIntro()
	if err != nil {
		t.Fatal(err)
	}
}
