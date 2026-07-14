package ddf

import (
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/edlundin/enocean-esp3/pkg/eep"
	"github.com/edlundin/enocean-esp3/pkg/enums"
)

var productIDPattern = regexp.MustCompile(`^0x[0-9A-Fa-f]{12}$`)

type File struct {
	XMLName xml.Name `xml:"Enocean_Devices"`
	Version string   `xml:"schemaVersion,attr"`
	Devices []Device `xml:"Device"`
}

type Device struct {
	ProductID string   `xml:"Product_ID,attr"`
	TX        Endpoint `xml:"TX"`
	RX        Endpoint `xml:"RX"`
}

type Endpoint struct {
	EURID  []EEPRef `xml:"EURID>EEP"`
	BaseID []EEPRef `xml:"BaseID>EEP"`
	EEP    []EEPRef `xml:"EEP"`
}

type EEPRef struct {
	LinkEntry string `xml:"LinkEntry,attr"`
	Rorg      string `xml:"Rorg,attr"`
	Func      string `xml:"Func,attr"`
	Type      string `xml:"Type,attr"`
}

func Parse(r io.Reader) (File, error) {
	var f File
	if err := xml.NewDecoder(r).Decode(&f); err != nil {
		return File{}, err
	}
	if strings.TrimSpace(f.Version) == "" {
		return File{}, fmt.Errorf("schemaVersion is required")
	}
	if len(f.Devices) == 0 {
		return File{}, fmt.Errorf("at least one Device is required")
	}
	for _, d := range f.Devices {
		if !productIDPattern.MatchString(d.ProductID) {
			return File{}, fmt.Errorf("invalid Product_ID %q", d.ProductID)
		}
	}
	return f, nil
}

func (r EEPRef) EEP() (eep.EEP, error) {
	rorg, err := parseHexByte(r.Rorg)
	if err != nil {
		return eep.EEP{}, fmt.Errorf("invalid Rorg: %w", err)
	}
	fn, err := parseOptionalHexByte(r.Func)
	if err != nil {
		return eep.EEP{}, fmt.Errorf("invalid Func: %w", err)
	}
	typ, err := parseOptionalHexByte(r.Type)
	if err != nil {
		return eep.EEP{}, fmt.Errorf("invalid Type: %w", err)
	}
	return eep.FromTriplet(enums.Rorg(rorg), fn, typ)
}

func parseOptionalHexByte(s string) (byte, error) {
	if strings.TrimSpace(s) == "" {
		return 0, nil
	}
	return parseHexByte(s)
}

func parseHexByte(s string) (byte, error) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "0x")
	v, err := strconv.ParseUint(s, 16, 8)
	return byte(v), err
}
