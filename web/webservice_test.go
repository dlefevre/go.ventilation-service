package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dlefevre/go.ventilation-service/controller"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// Set the environment variable for the configuration path
	os.Setenv("VENTILATIONSERVICE_CONFIG_PATH", "..")
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func setup() {
	dc := controller.GetVentilationControllerService()
	dc.Start()
	ws := GetWebService()
	ws.Start()

	time.Sleep(500 * time.Microsecond)
}

func teardown() {
	dc := controller.GetVentilationControllerService()
	dc.Stop()
	ws := GetWebService()
	ws.Stop()
}

func reqHelper(t *testing.T, path string, body string) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:8000%s", path), strings.NewReader(body))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Add("x-api-key", "test")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response: %v", err)
	}
	var myResponse SimpleResponse
	if err := json.Unmarshal(respBody, &myResponse); err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}
	if myResponse.Result != "ok" {
		t.Fatalf("Expected result to be ok, got %s", myResponse.Result)
	}
}

func TestStartStop(t *testing.T) {
	setup()
	teardown()
}

func TestSpeed(t *testing.T) {
	setup()
	defer teardown()
	reqHelper(t, "/speed", `{"speed": "1"}`)
	reqHelper(t, "/speed", `{"speed": "2"}`)
	reqHelper(t, "/speed", `{"speed": "3"}`)
	reqHelper(t, "/speed", `{"speed": "low"}`)
	reqHelper(t, "/speed", `{"speed": "medium"}`)
	reqHelper(t, "/speed", `{"speed": "high"}`)
	time.Sleep(20 * time.Second)
}

func TestTimer(t *testing.T) {
	setup()
	defer teardown()
	reqHelper(t, "/timer", `{"duration": 15}`)
	reqHelper(t, "/timer", `{"duration": 30}`)
	reqHelper(t, "/timer", `{"duration": 60}`)
	time.Sleep(10 * time.Second)
}

func TestOther(t *testing.T) {
	setup()
	defer teardown()
	reqHelper(t, "/away", `{"away": "true"}`)
	reqHelper(t, "/auto", `{"auto": "true"}`)
	time.Sleep(8 * time.Second)
}
