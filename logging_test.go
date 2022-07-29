package pg

import (
	"testing"
)

func TestLogger_Enabled(t *testing.T) {
	NewLogger()

	Log.Errorf("err")
	Log.Error("err")
	Log.Println("println")
	Log.Printf("printf")
	Log.Print("print")
	Log.Warnf("warnf")
	Log.Warn("warn")
}

func TestLogger_Disabled(t *testing.T) {
	DisableLogging()

	Log.Errorf("err")
	Log.Error("err")
	Log.Println("println")
	Log.Printf("printf")
	Log.Print("print")
	Log.Warnf("warnf")
	Log.Warn("warn")
}
