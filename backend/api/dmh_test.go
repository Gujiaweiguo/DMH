package main

import (
	"os"
	"testing"

	"dmh/api/internal/config"
)

func TestApplyEnvOverrides(t *testing.T) {
	os.Setenv("APP_HOST", "0.0.0.0")
	os.Setenv("APP_PORT", "9999")
	os.Setenv("JWT_SECRET", "jwt-test")
	os.Setenv("DB_HOST", "127.0.0.1")
	os.Setenv("DB_PORT", "3307")
	os.Setenv("DB_USER", "u1")
	os.Setenv("DB_PASSWORD", "p1")
	os.Setenv("DB_NAME", "d1")
	t.Cleanup(func() {
		os.Unsetenv("APP_HOST")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	})

	var c config.Config
	applyEnvOverrides(&c)

	if c.Host != "0.0.0.0" || c.Port != 9999 || c.Auth.AccessSecret != "jwt-test" {
		t.Fatalf("basic env overrides not applied")
	}
	if c.Mysql.DataSource == "" {
		t.Fatalf("db dsn should be generated")
	}
}

func TestApplyEnvOverrides_PartialDBVars(t *testing.T) {
	os.Setenv("DB_HOST", "192.168.1.1")
	os.Setenv("DB_PORT", "3308")
	t.Cleanup(func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	})

	var c config.Config
	applyEnvOverrides(&c)

	if c.Mysql.DataSource == "" {
		t.Fatalf("db dsn should be generated with partial env vars")
	}
}

func TestApplyEnvOverrides_EmptyEnv(t *testing.T) {
	t.Cleanup(func() {
		os.Unsetenv("APP_HOST")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	})

	os.Unsetenv("APP_HOST")
	os.Unsetenv("APP_PORT")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")

	var c config.Config
	c.Host = "original-host"
	c.Port = 8888
	originalDSN := "user:pass@tcp(localhost:3306)/db"
	c.Mysql.DataSource = originalDSN

	applyEnvOverrides(&c)

	if c.Host != "original-host" {
		t.Fatalf("host should not change when env vars are empty")
	}
	if c.Mysql.DataSource != originalDSN {
		t.Fatalf("datasource should not change when DB env vars are empty")
	}
}

func TestApplyEnvOverrides_InvalidPort(t *testing.T) {
	os.Setenv("APP_PORT", "invalid")
	os.Setenv("DB_HOST", "127.0.0.1")
	t.Cleanup(func() {
		os.Unsetenv("APP_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	})

	var c config.Config
	c.Port = 8888
	applyEnvOverrides(&c)

	if c.Port != 8888 {
		t.Fatalf("port should not change when APP_PORT is invalid")
	}
}
