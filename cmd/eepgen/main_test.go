package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const tinyEEPXML = `<?xml version="1.0"?>
<eep>
  <rorg>
    <number>0xA5</number>
    <title>RORG title</title>
    <func>
      <number>0x02</number>
      <title>FUNC title</title>
      <type>
        <number>0x05</number>
        <title>Temperature &amp;amp; Humidity</title>
        <case>
          <datafield>
            <data>Temperature   Sensor</data>
            <shortcut>TMP</shortcut>
            <bitoffs>8</bitoffs>
            <bitsize>8</bitsize>
            <range><min>0</min><max>255</max></range>
            <scale><min>-40</min><max>60</max></scale>
            <unit>°C</unit>
            <enum>
              <item><value>0</value><description>Off - disabled</description></item>
              <item><value>1</value><description>On, enabled</description></item>
            </enum>
          </datafield>
          <datafield>
            <data>Temperature Sensor</data>
            <shortcut>TMP</shortcut>
            <bitoffs>8</bitoffs>
            <bitsize>8</bitsize>
          </datafield>
          <datafield>
            <data>ignored bad size</data>
            <bitoffs>0</bitoffs>
            <bitsize>0</bitsize>
          </datafield>
        </case>
      </type>
    </func>
  </rorg>
</eep>`

func TestRunGeneratesProfilesFromFlags(t *testing.T) {
	xmlPath := writeXML(t)
	outDir := t.TempDir()

	if err := run(flag.NewFlagSet("eepgen", flag.ContinueOnError), []string{"-xml", xmlPath, "-out", outDir}); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	got := readGenerated(t, outDir)
	for _, want := range []string{
		`Registry["A5-02-05"]`,
		`Title: "Temperature & Humidity"`,
		`Name: "Temperature Sensor"`,
		`Shortcut: "TMP"`,
		`BitOff: 8`,
		`BitSize: 8`,
		`Unit: "°C"`,
		`ScaleMin: -40`,
		`ScaleMax: 60`,
		`RawMin: 0`,
		`RawMax: 255`,
		`{Raw: 0, Name: "Off", Description: "Off - disabled"}`,
		`{Raw: 1, Name: "On", Description: "On, enabled"}`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("generated file missing %q\n%s", want, got)
		}
	}
	if strings.Contains(got, "ignored bad size") || strings.Count(got, `Name: "Temperature Sensor"`) != 1 {
		t.Fatalf("generated file did not filter invalid/duplicate fields\n%s", got)
	}
}

func TestRunUsesDefaultPaths(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.WriteFile("eep268.xml", []byte(tinyEEPXML), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := run(flag.NewFlagSet("eepgen", flag.ContinueOnError), []string{}); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join("pkg", "eep", "profiles", "profiles_gen.go")); err != nil {
		t.Fatalf("default output was not written: %v", err)
	}
}

func TestRunReturnsFlagErrors(t *testing.T) {
	if err := run(flag.NewFlagSet("eepgen", flag.ContinueOnError), []string{"-nope"}); err == nil {
		t.Fatal("run() error = nil, want flag parse error")
	}
}

func TestRunReturnsGenerateErrors(t *testing.T) {
	outDir := t.TempDir()
	if err := run(flag.NewFlagSet("eepgen", flag.ContinueOnError), []string{"-xml", filepath.Join(t.TempDir(), "missing.xml"), "-out", outDir}); err == nil {
		t.Fatal("run() error = nil, want missing XML error")
	}
}

func writeXML(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "eep.xml")
	if err := os.WriteFile(path, []byte(tinyEEPXML), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func readGenerated(t *testing.T, outDir string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(outDir, "profiles_gen.go"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
