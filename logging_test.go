package pg_test

import (
	"testing"

	"github.com/pandudpn/go-payment-gateway"
)

func TestLogger_Enabled(t *testing.T) {
	pg.NewLogger()

	pg.Log.Errorf("err")
	pg.Log.Error("err")
	pg.Log.Println("println")
	pg.Log.Printf("printf")
	pg.Log.Print("print")
	pg.Log.Warnf("warnf")
	pg.Log.Warn("warn")
}

func TestLogger_Disabled(t *testing.T) {
	pg.DisableLogging()

	pg.Log.Errorf("err")
	pg.Log.Error("err")
	pg.Log.Println("println")
	pg.Log.Printf("printf")
	pg.Log.Print("print")
	pg.Log.Warnf("warnf")
	pg.Log.Warn("warn")
}
