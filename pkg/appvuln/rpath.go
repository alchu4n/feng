package appvuln

import "strings"

// rpath
const (
	// TODO other rpaths
	rpathPre          = "@rpath"
	executablePathPre = "@executable_path"
	loaderPathPre     = "@loader_path"
)

func parseProxyPath(path string, ExePath string) string {
	if strings.HasPrefix(path, executablePathPre) {
		return strings.ReplaceAll(path, executablePathPre, ExePath)
	} else if strings.HasPrefix(path, loaderPathPre) {
		return strings.ReplaceAll(path, loaderPathPre, ExePath)
	}
	return path
}

func parseRPaths(rpath []string, ExePath string) {
	for i := range rpath {
		switch true {
		case strings.HasPrefix(rpath[i], executablePathPre):
			rpath[i] = strings.ReplaceAll(rpath[i], executablePathPre, ExePath)
		case strings.HasPrefix(rpath[i], loaderPathPre):
			rpath[i] = strings.ReplaceAll(rpath[i], loaderPathPre, ExePath)
		}
	}
}

func joinRPath(p, exec string, rpaths []string) []string {
	paths := make([]string, 0)
	switch true {
	case strings.HasPrefix(p, rpathPre):
		for i := range rpaths {
			paths = append(paths, strings.ReplaceAll(p, rpathPre, rpaths[i]))
		}
	case strings.HasPrefix(p, executablePathPre):
		paths = append(paths, strings.ReplaceAll(p, executablePathPre, exec))
	}
	return paths
}
