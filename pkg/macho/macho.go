package macho

import (
	"bytes"

	"github.com/blacktop/go-macho"
	"github.com/blacktop/go-macho/types"
	"howett.net/plist"
)

type MachOInfo struct {
	Magic    string   `json:"magic"`
	CPU      string   `json:"cpu"`
	Type     string   `json:"type"`
	Dylibs   Dylibs   `json:"dylibs"`
	CodeSign CodeSign `json:"codesign"`
}

func Parse(path string) ([]*MachOInfo, error) {
	var res = make([]*MachOInfo, 0)
	// first check for fat file
	f, err := macho.OpenFat(path)
	if err != nil && err != macho.ErrNotFat {
		return nil, err
	}

	if err == macho.ErrNotFat {
		m, err := macho.Open(path)
		if err != nil {
			return nil, err
		}
		defer m.Close()

		r, err := parseMacho(m)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	} else {
		defer f.Close()

		for i := range f.Arches {
			r, err := parseMacho(f.Arches[i].File)
			if err != nil {
				return nil, err
			}
			res = append(res, r)
		}
	}
	return res, nil
}

type Dylibs struct {
	Dylinker string   `json:"dylinker"`
	RPaths   []string `json:"rpaths"`
	Loads    []Dylib  `json:"loads"`
	Weaks    []Dylib  `json:"weaks"`
}

type Dylib struct {
	Name           string `json:"name"`
	Time           uint32 `json:"time"`
	CurrentVersion string `json:"current_version"`
	CompatVersion  string `json:"compat_version"`
}

type CodeSign struct {
	ID             string         `json:"id"`
	TeamID         string         `json:"team_id"`
	Flags          uint32         `json:"flags"`
	FlagsString    string         `json:"flags_string"`
	RuntimeVersion string         `json:"runtime_version"`
	Entitlements   map[string]any `json:"entitlements"`
}

func parseMacho(m *macho.File) (*MachOInfo, error) {
	res := &MachOInfo{
		Magic: m.Magic.String(),
		Type:  m.Type.String(),
		CPU:   m.CPU.String(),
		Dylibs: Dylibs{
			Loads:  make([]Dylib, 0),
			Weaks:  make([]Dylib, 0),
			RPaths: make([]string, 0),
		},
		CodeSign: CodeSign{
			Entitlements: make(map[string]any),
		},
	}

	for _, v := range m.Loads {
		switch v.Command() {
		case types.LC_LOAD_DYLINKER:
			res.Dylibs.Dylinker = v.String()
		case types.LC_LOAD_DYLIB,
			types.LC_LOAD_WEAK_DYLIB:
			var d Dylib
			switch dylib := v.(type) {
			case *macho.LoadDylib:
				d = Dylib{
					Name:           dylib.Name,
					Time:           dylib.Timestamp,
					CurrentVersion: dylib.CurrentVersion.String(),
					CompatVersion:  dylib.CompatVersion.String(),
				}
				res.Dylibs.Loads = append(res.Dylibs.Loads, d)
			case *macho.WeakDylib:
				d = Dylib{
					Name:           dylib.Name,
					Time:           dylib.Timestamp,
					CurrentVersion: dylib.CurrentVersion.String(),
					CompatVersion:  dylib.CompatVersion.String(),
				}
				res.Dylibs.Weaks = append(res.Dylibs.Weaks, d)
			}
		case types.LC_RPATH:
			res.Dylibs.RPaths = append(res.Dylibs.RPaths, v.String())
		case types.LC_REEXPORT_DYLIB:
		case types.LC_CODE_SIGNATURE:
			cs, ok := v.(*macho.CodeSignature)
			if !ok {
				continue
			}
			if len(cs.Entitlements) > 0 {
				if err := parsePlist(cs.Entitlements, res.CodeSign.Entitlements); err != nil {
					return nil, err
				}
			}
			if len(cs.CodeDirectories) > 0 {
				res.CodeSign.ID = cs.CodeDirectories[0].ID
				res.CodeSign.TeamID = cs.CodeDirectories[0].TeamID
				res.CodeSign.RuntimeVersion = cs.CodeDirectories[0].RuntimeVersion
				res.CodeSign.Flags = uint32(cs.CodeDirectories[0].Header.Flags)
				res.CodeSign.FlagsString = cs.CodeDirectories[0].Header.Flags.String()
			}
		}
	}
	return res, nil
}

func parsePlist(data string, e map[string]any) error {
	decoder := plist.NewDecoder(bytes.NewReader([]byte(data)))
	if err := decoder.Decode(e); err != nil {
		return err
	}
	return nil
}
