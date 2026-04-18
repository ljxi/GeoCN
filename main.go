package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/oschwald/maxminddb-golang"
)

// ---------- MMDB record structs ----------

type ASNRecord struct {
	Number uint   `maxminddb:"autonomous_system_number" json:"number"`
	Name   string `maxminddb:"autonomous_system_organization" json:"name"`
	Info   string `maxminddb:"-" json:"info"`
}

type CityNames struct {
	De   string `maxminddb:"de" json:"de,omitempty"`
	En   string `maxminddb:"en" json:"en,omitempty"`
	Es   string `maxminddb:"es" json:"es,omitempty"`
	Fr   string `maxminddb:"fr" json:"fr,omitempty"`
	Ja   string `maxminddb:"ja" json:"ja,omitempty"`
	PtBR string `maxminddb:"pt-BR" json:"pt-BR,omitempty"`
	Ru   string `maxminddb:"ru" json:"ru,omitempty"`
	ZhCN string `maxminddb:"zh-CN" json:"zh-CN,omitempty"`
}

type CityCountry struct {
	ISOCode string    `maxminddb:"iso_code" json:"iso_code"`
	Names   CityNames `maxminddb:"names" json:"names"`
}

type CitySubdivision struct {
	ISOCode string    `maxminddb:"iso_code" json:"iso_code"`
	Names   CityNames `maxminddb:"names" json:"names"`
}

type CityCity struct {
	Names CityNames `maxminddb:"names" json:"names"`
}

type CityRecord struct {
	City              CityCity          `maxminddb:"city" json:"city"`
	Country           CityCountry       `maxminddb:"country" json:"country"`
	RegisteredCountry CityCountry       `maxminddb:"registered_country" json:"registered_country"`
	Subdivisions      []CitySubdivision `maxminddb:"subdivisions" json:"subdivisions"`
}

// ---------- Response structs ----------

type CountryInfo struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type IPResponse struct {
	IP                string                 `json:"ip"`
	Addr              string                 `json:"addr"`
	AS                *ASNRecord             `json:"as,omitempty"`
	Country           *CountryInfo           `json:"country,omitempty"`
	RegisteredCountry *CountryInfo           `json:"registered_country,omitempty"`
	Subdivision       string                 `json:"subdivision,omitempty"`
	City              string                 `json:"city,omitempty"`
	Area              string                 `json:"area,omitempty"`
	GeoCN             map[string]interface{} `json:"geo_cn,omitempty"`
}

// ---------- Global state ----------

var (
	asnDB  *maxminddb.Reader
	cityDB *maxminddb.Reader
	cnDB   *maxminddb.Reader

	divisionFull  map[int]string  // code -> full name
	divisionShort map[int]string  // code -> short name
	asnInfo       map[uint]string // AS number -> info
)

// toStringMap converts an interface{} from maxminddb into map[string]interface{}.
// maxminddb unmarshals MMDB maps as map[string]interface{}.
func toStringMap(v interface{}) map[string]interface{} {
	if v == nil {
		return nil
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return m
}

// extractUint32 safely extracts a uint32 value from an interface{} that may be
// uint, uint16, uint32, uint64, or int types as returned by maxminddb.
func extractUint32(v interface{}) uint32 {
	switch val := v.(type) {
	case uint32:
		return val
	case uint64:
		return uint32(val)
	case uint:
		return uint32(val)
	case uint16:
		return uint32(val)
	case int:
		if val > 0 {
			return uint32(val)
		}
	}
	return 0
}

func loadASNInfo(path string) (map[uint]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := make(map[uint]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}
		n, err := strconv.ParseUint(strings.TrimSpace(parts[0]), 10, 64)
		if err != nil {
			continue
		}
		m[uint(n)] = strings.TrimSpace(parts[1])
	}
	return m, scanner.Err()
}

func loadDivisionCode(path string, sep string) (map[int]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := make(map[int]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, sep, 2)
		if len(parts) != 2 {
			continue
		}
		code, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			continue
		}
		m[code] = strings.TrimSpace(parts[1])
	}
	return m, scanner.Err()
}

func resolveDivision(code uint32) map[string]interface{} {
	if code == 0 {
		return nil
	}
	c := int(code)
	provinceCode := (c / 10000) * 10000
	cityCode := (c / 100) * 100

	var full, short []string

	// Province
	if v, ok := divisionFull[provinceCode]; ok {
		full = append(full, v)
	}
	if v, ok := divisionShort[provinceCode]; ok {
		short = append(short, v)
	}

	// City (skip if same as province, e.g. 直辖市)
	if cityCode != provinceCode {
		if v, ok := divisionFull[cityCode]; ok {
			full = append(full, v)
		}
		if v, ok := divisionShort[cityCode]; ok {
			short = append(short, v)
		}
	}

	// District (skip if same as city)
	if c != cityCode && c != provinceCode {
		if v, ok := divisionFull[c]; ok {
			full = append(full, v)
		}
		if v, ok := divisionShort[c]; ok {
			short = append(short, v)
		}
	}

	if len(full) == 0 && len(short) == 0 {
		return nil
	}
	// 直辖市：省级和市级名称相同时，去掉重复的第一项
	if len(full) >= 2 && full[0] == full[1] {
		full = full[1:]
	}
	if len(short) >= 2 && short[0] == short[1] {
		short = short[1:]
	}
	return map[string]interface{}{"full": full, "short": short}
}

// fixSpecialRegion prepends "中国" for HK/MO/TW
func fixCountryName(name string) string {
	switch name {
	case "香港", "澳门", "台湾":
		return "中国" + name
	}
	return name
}

func queryIP(ipStr string) (*IPResponse, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP: %s", ipStr)
	}

	resp := &IPResponse{IP: ipStr}

	// ASN lookup
	var asn ASNRecord
	network, _, err := asnDB.LookupNetwork(ip, &asn)
	if err == nil && asn.Number != 0 {
		asn.Info = asnInfo[asn.Number]
		resp.AS = &asn
		if network != nil {
			resp.Addr = network.String()
		}
	}

	// City lookup
	var city CityRecord
	network2, _, err := cityDB.LookupNetwork(ip, &city)
	if err == nil {
		if network2 != nil {
			resp.Addr = network2.String()
		}

		countryName := fixCountryName(city.Country.Names.ZhCN)
		resp.Country = &CountryInfo{
			Code: city.Country.ISOCode,
			Name: countryName,
		}

		regName := fixCountryName(city.RegisteredCountry.Names.ZhCN)
		resp.RegisteredCountry = &CountryInfo{
			Code: city.RegisteredCountry.ISOCode,
			Name: regName,
		}

		if city.City.Names.ZhCN != "" {
			resp.City = city.City.Names.ZhCN
		} else if city.City.Names.En != "" {
			resp.City = city.City.Names.En
		}

		if len(city.Subdivisions) > 0 {
			if city.Subdivisions[0].Names.ZhCN != "" {
				resp.Subdivision = city.Subdivisions[0].Names.ZhCN
			} else if city.Subdivisions[0].Names.En != "" {
				resp.Subdivision = city.Subdivisions[0].Names.En
			}
		}
	}

	// GeoCN lookup for Chinese IPs
	if cnDB != nil && city.Country.ISOCode == "CN" {
		var cnRaw interface{}
		network3, _, err := cnDB.LookupNetwork(ip, &cnRaw)
		if err == nil {
			resp.Addr = network3.String()
			cnMap := toStringMap(cnRaw)
			if cnMap != nil {
				divCode := extractUint32(cnMap["division_code"])
				if divCode != 0 {
					div := resolveDivision(divCode)
					if div != nil {
						cnMap["division"] = div
						if short, ok := div["short"].([]string); ok {
							if len(short) >= 1 {
								resp.Subdivision = short[0]
							}
							if len(short) >= 2 {
								resp.City = short[1]
							}
							if len(short) >= 3 {
								resp.Area = short[2]
							}
						}
					}
				}
				resp.GeoCN = cnMap
			}
		}
	}

	return resp, nil
}

func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.SplitN(xff, ",", 2)
		ip := strings.TrimSpace(parts[0])
		if net.ParseIP(ip) != nil {
			return ip
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var ipStr string

	// Priority: query param > path > client IP
	if q := r.URL.Query().Get("ip"); q != "" {
		ipStr = q
	} else if path := strings.TrimPrefix(r.URL.Path, "/"); path != "" {
		ipStr = path
	} else {
		ipStr = getClientIP(r)
	}

	if net.ParseIP(ipStr) == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid IP address"})
		return
	}

	resp, err := queryIP(ipStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	enc.Encode(resp)
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	var err error

	asnPath := envOrDefault("ASN_DB", "db/GeoLite2-ASN.mmdb")
	cityPath := envOrDefault("CITY_DB", "db/GeoLite2-City.mmdb")
	cnPath := envOrDefault("CN_DB", "db/GeoCN.mmdb")
	divFullPath := envOrDefault("DIV_FULL", "data/full.txt")
	divShortPath := envOrDefault("DIV_SHORT", "data/short.txt")
	listen := envOrDefault("LISTEN", ":80")

	asnDB, err = maxminddb.Open(asnPath)
	if err != nil {
		log.Fatalf("Failed to open ASN database: %v", err)
	}
	defer asnDB.Close()

	cityDB, err = maxminddb.Open(cityPath)
	if err != nil {
		log.Fatalf("Failed to open City database: %v", err)
	}
	defer cityDB.Close()

	cnDB, err = maxminddb.Open(cnPath)
	if err != nil {
		log.Printf("Warning: GeoCN database not available: %v", err)
		cnDB = nil
	} else {
		defer cnDB.Close()
	}

	asnInfoPath := envOrDefault("ASN_INFO", "data/asn.txt")
	asnInfo, err = loadASNInfo(asnInfoPath)
	if err != nil {
		log.Printf("Warning: ASN info not available: %v", err)
		asnInfo = make(map[uint]string)
	}

	divisionFull, err = loadDivisionCode(divFullPath, "\t")
	if err != nil {
		log.Fatalf("Failed to load full division codes: %v", err)
	}

	divisionShort, err = loadDivisionCode(divShortPath, "  ")
	if err != nil {
		log.Fatalf("Failed to load short division codes: %v", err)
	}

	log.Printf("Listening on %s", listen)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(listen, nil))
}
