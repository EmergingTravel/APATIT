package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"apatit/internal/config"
)

// TestApplication_Profile run APATIT for profiling.
func TestApplication_Profile(t *testing.T) {

	const profileDuration = 120 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), profileDuration)
	defer cancel()

	t.Logf("Starting application for profiling, duration: %s", profileDuration)

	cfg, err := config.New()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	app, err := newApp(cfg)
	if err != nil {
		t.Fatalf("Application initialization failed: %v", err)
	}

	if err := app.Run(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Application terminated with an unexpected error: %v", err)
		}
	}

	t.Logf("Application ran for %s. Profiling test finished.", profileDuration)
}
