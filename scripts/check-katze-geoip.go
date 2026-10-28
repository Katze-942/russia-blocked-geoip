// Run from the checked-out geoip module so its protobuf schema is used.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	v2ray "github.com/Loyalsoldier/geoip/plugin/v2ray"
	"google.golang.org/protobuf/proto"
)

func load(path string) map[string]*v2ray.GeoIP {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	var list v2ray.GeoIPList
	if err := proto.Unmarshal(data, &list); err != nil {
		log.Fatalf("%s: %v", path, err)
	}
	entries := make(map[string]*v2ray.GeoIP)
	for _, entry := range list.GetEntry() {
		code := strings.ToUpper(entry.GetCountryCode())
		if code == "" || entries[code] != nil {
			log.Fatalf("%s: empty or repeated category %q", path, code)
		}
		if len(entry.GetCidr()) == 0 {
			log.Fatalf("%s: category %s has no CIDRs", path, code)
		}
		for _, cidr := range entry.GetCidr() {
			bits := len(cidr.GetIp()) * 8
			if (bits != 32 && bits != 128) || cidr.GetPrefix() > uint32(bits) {
				log.Fatalf("%s: invalid CIDR in category %s", path, code)
			}
		}
		entries[code] = entry
	}
	return entries
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("usage: check-katze-geoip config.json output-directory")
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	var config struct {
		Output []struct {
			Type   string `json:"type"`
			Action string `json:"action"`
			Args   struct {
				OutputName string   `json:"outputName"`
				WantedList []string `json:"wantedList"`
			} `json:"args"`
		} `json:"output"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}
	wanted := make(map[string]bool)
	outputs := 0
	for _, output := range config.Output {
		if output.Type != "v2rayGeoIPDat" || output.Action != "output" || output.Args.OutputName != "geoip-katze.dat" {
			continue
		}
		outputs++
		for _, name := range output.Args.WantedList {
			code := strings.ToUpper(strings.TrimSpace(name))
			if code == "" || wanted[code] {
				log.Fatalf("empty or repeated requested category %q", code)
			}
			wanted[code] = true
		}
	}
	if len(config.Output) != 1 || outputs != 1 || len(wanted) == 0 {
		log.Fatal("expected only geoip-katze.dat with a nonempty wantedList")
	}
	subset := load(filepath.Join(os.Args[2], "geoip-katze.dat"))
	for code := range wanted {
		if subset[code] == nil {
			log.Fatalf("missing required category %s", code)
		}
	}
	if len(subset) != len(wanted) {
		log.Fatal("geoip-katze.dat contains unexpected categories")
	}
	names := make([]string, 0, len(wanted))
	for code := range wanted {
		names = append(names, strings.ToLower(code))
	}
	sort.Strings(names)
	if err := os.WriteFile(filepath.Join(os.Args[2], "geoip-katze.txt"), []byte(strings.Join(names, "\n")+"\n"), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("geoip-katze.dat: verified %d categories and their CIDRs\n", len(wanted))
}
