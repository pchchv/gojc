package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

func TestServerPing(t *testing.T) {
	res, err := http.Get("http://" + getEnvValue("HOST") + ":" + getEnvValue("PORT") + "/ping")
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("status not OK")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error(err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	b := string(body)
	if !strings.Contains(b, "JSON") {
		t.Fatal()
	}
}

func TestLoadPing(t *testing.T) {
	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := 5 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://" + getEnvValue("HOST") + ":" + getEnvValue("PORT") + "/ping",
	})
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()
	log.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}
