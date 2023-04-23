package appvuln

import (
	"bytes"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/ac0d3r/macoder/pkg/macho"
	"howett.net/plist"
)

type AppScanx struct {
}

func New() *AppScanx {
	return &AppScanx{}
}

type Info struct {
	Path           string     `json:"path"`
	ExecutablePath string     `json:"executable_path"`
	Injectable     bool       `json:"injectable"` // 应用可以注入动态库
	Dylibs         []VulnItem `json:"dylibs"`
}

type VulnItem struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

func (d *AppScanx) Scan() (res []*Info, err error) {
	fs, err := os.ReadDir("/Applications")
	if err != nil {
		return nil, err
	}

	for _, f := range fs {
		if f.IsDir() && strings.HasSuffix(f.Name(), ".app") {
			vuln, err := d.ScanSingleApp(path.Join("/Applications", f.Name()))
			if err != nil {
				return nil, err
			}
			if !vuln.Injectable && len(vuln.Dylibs) == 0 {
				continue
			}
			res = append(res, vuln)
		}
	}
	return
}

func (d *AppScanx) ScanSingleApp(path_ string) (info *Info, err error) {
	info = &Info{
		Path:   path_,
		Dylibs: make([]VulnItem, 0),
	}
	if err = d.SetAppInfo(path_, info); err != nil {
		return nil, err
	}

	machoInfo, err := macho.Parse(info.ExecutablePath)
	if err != nil {
		return nil, err
	}

	for _, m := range machoInfo {
		if (IsRuntime(m.CodeSign.Flags) || IsLibraryValidation(m.CodeSign.Flags)) &&
			(!hasKey(m.CodeSign.Entitlements, disableLibraryValidation) || !hasKey(m.CodeSign.Entitlements, allowDyldEnvironment)) {
			continue
		}

		info.Injectable = true
		// parse rpaths
		parseRPaths(m.Dylibs.RPaths, path.Dir(info.ExecutablePath))
		// weak dylib
		weaks := make([]string, 0)
		for _, weak := range m.Dylibs.Weaks {
			switch true {
			case path.IsAbs(weak.Name):
				weaks = append(weaks, weak.Name)
			case strings.HasPrefix(weak.Name, "@"):
				weaks = append(weaks, joinRPath(weak.Name, info.ExecutablePath, m.Dylibs.RPaths)...)
			}
		}
		for _, weak := range weaks {
			if InSIPPath(weak) || pathExist(weak) {
				continue
			}
			info.Dylibs = append(info.Dylibs, VulnItem{
				Type: "weak",
				Path: weak,
			})
		}
		// @rpath dylib
		for _, dylib := range m.Dylibs.Loads {
			if strings.HasPrefix(dylib.Name, rpathPre) {
				paths := joinRPath(dylib.Name, info.ExecutablePath, m.Dylibs.RPaths)
				if len(paths) == 0 {
					continue
				}
				// find the first exist path
				index := -1
				for i := range paths {
					if pathExist(paths[i]) {
						index = i
						break
					}
				}
				if index == -1 {
					index = len(paths)
				}
				for i := range paths[:index] {
					if InSIPPath(paths[i]) {
						continue
					}
					info.Dylibs = append(info.Dylibs, VulnItem{
						Type: "rpath",
						Path: paths[i],
					})
				}
			} else { // dylib proxying
				dylib.Name = parseProxyPath(dylib.Name,
					path.Dir(info.ExecutablePath))
				if !InSIPPath(dylib.Name) {
					info.Dylibs = append(info.Dylibs, VulnItem{
						Type: "proxy",
						Path: dylib.Name,
					})
				}
			}
		}
	}
	return
}

var (
	errNotApplication             = errors.New("not application")
	errNotFoundCFBundleExecutable = errors.New("not found 'CFBundleExecutable' in info.plist")
)

func (d *AppScanx) SetAppInfo(path_ string, info *Info) error {
	fi, err := os.Stat(path_)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return errNotApplication
	}

	// get application info.plist
	data, err := os.ReadFile(path.Join(path_, "Contents", "Info.plist"))
	if err != nil {
		return err
	}

	appinfo := make(map[string]any)
	decoder := plist.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(appinfo); err != nil {
		return err
	}
	v, ok := appinfo["CFBundleExecutable"]
	if !ok {
		return errNotFoundCFBundleExecutable
	}
	exe, ok := v.(string)
	if !ok {
		return errNotFoundCFBundleExecutable
	}

	info.ExecutablePath = path.Join(path_, "Contents", "MacOS", exe)
	return nil
}

const (
	REQUIRE_LV = 0x2000
	RUNTIME    = 0x10000
)

func IsRuntime(f uint32) bool {
	return f&RUNTIME != 0
}

func IsLibraryValidation(f uint32) bool {
	return f&REQUIRE_LV != 0
}

const (
	disableLibraryValidation = "com.apple.security.cs.disable-library-validation"
	allowDyldEnvironment     = "com.apple.security.cs.allow-dyld-environment-variables"
)

func hasKey(entitlements map[string]any, key string) bool {
	if e, ok := entitlements[key]; ok {
		if ee, ok := e.(bool); ok {
			return ee
		}
	}
	return false
}

var (
	// TODO 粗略判断
	// ls -lO /
	// restricted
	sippaths = []string{
		"/System/",
		"/bin/",
		"/etc/",
		"/sbin/",
		"/tmp/",
		"/usr/",
		"/var/",
	}
)

func InSIPPath(path string) bool {
	for _, sip := range sippaths {
		if strings.HasPrefix(path, sip) {
			return true
		}
	}

	return false
}

func pathExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
