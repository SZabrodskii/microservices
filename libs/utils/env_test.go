package utils

import (
	"os"
	"reflect"
	"testing"
)

func TestGetEnv(t *testing.T) {
	t.Run("should return int8", func(t *testing.T) {
		if val := GetEnv("TI", int8(123)); val != 123 {
			t.Errorf("default should be int8 %s", reflect.TypeOf(val).Kind().String())
		}
	})

	t.Run("should return int16", func(t *testing.T) {
		if GetEnv("TI", int16(123)) != 123 {
			t.Error("default should be int16 123")
		}
	})

	t.Run("should return int64", func(t *testing.T) {
		if GetEnv("TI", int16(123)) != 123 {
			t.Error("default should be int64 123")
		}
	})

	t.Run("should return string", func(t *testing.T) {
		os.Setenv("TS", "test string")

		if GetEnv("TS", "default string") != "test string" {
			t.Error("Should return string 'test string")
		}
	})

	t.Run("Should return default string", func(t *testing.T) {
		if val := GetEnv("TS", "default string"); val != "default string" {
			t.Errorf("Should return string 'default string' got %v", val)
		}
	})

	t.Run("should return default fallback", func(t *testing.T) {
		if GetEnv("TI", 123) != 123 {
			t.Error("default should be int 123")
		}
	})

	t.Run("should get as int", func(t *testing.T) {
		os.Setenv("TI", "123456789")

		if GetEnv("TI", 0) != 123456789 {
			t.Error("should be int 123456789")
		}
	})

	t.Run("should return bool", func(t *testing.T) {
		if GetEnv("TI", true) != true {
			t.Error("default should be bool true")
		}
	})

	t.Run("should return bool", func(t *testing.T) {
		os.Setenv("TI", "false")
		if GetEnv("TI", false) != false {
			t.Error("Should be bool false")
		}
	})

}
