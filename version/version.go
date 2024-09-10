package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/mod/semver"
)

// Constants for repository details
const (
	repoOwner = "bootdotdev"
	repoName  = "bootdev"
)

// VersionInfo holds information about the current and latest version
type VersionInfo struct {
	CurrentVersion   string
	LatestVersion    string
	IsOutdated       bool
	IsUpdateRequired bool
	FailedToFetch    error
}

// FetchUpdateInfo fetches the latest version info and checks if an update is needed
func FetchUpdateInfo(currentVersion string) VersionInfo {
	latest, err := getLatestVersion()
	if err != nil {
		return VersionInfo{
			FailedToFetch: err,
		}
	}
	isUpdateRequired := isUpdateRequired(currentVersion, latest)
	isOutdated := isOutdated(currentVersion, latest)
	return VersionInfo{
		IsUpdateRequired: isUpdateRequired,
		IsOutdated:       isOutdated,
		CurrentVersion:   currentVersion,
		LatestVersion:    latest,
	}
}

// PromptUpdateIfAvailable prints a message if an update is available
func (v *VersionInfo) PromptUpdateIfAvailable() {
	if v.IsOutdated {
		fmt.Fprintln(os.Stderr, "A new version of the bootdev CLI is available!")
		fmt.Fprintln(os.Stderr, "Please run the following command to update:")
		fmt.Fprintf(os.Stderr, "  bootdev upgrade\n\n")
	}
}

// isOutdated returns true if the current version is older than the latest version
func isOutdated(current string, latest string) bool {
	return semver.Compare(current, latest) < 0
}

// isUpdateRequired returns true if the latest version has a higher major or minor number than the current version
func isUpdateRequired(current string, latest string) bool {
	latestMajorMinor := semver.MajorMinor(latest)
	currentMajorMinor := semver.MajorMinor(current)
	return semver.Compare(currentMajorMinor, latestMajorMinor) < 0
}

// getLatestVersion retrieves the latest version from the Go proxy
func getLatestVersion() (string, error) {
	const goproxyDefault = "https://proxy.golang.org"
	var goproxy = goproxyDefault

	cmd := exec.Command("go", "env", "GOPROXY")
	output, err := cmd.Output()
	if err == nil {
		goproxy = strings.TrimSpace(string(output))
	}

	proxies := strings.Split(goproxy, ",")
	if !contains(proxies, goproxyDefault) {
		proxies = append(proxies, goproxyDefault)
	}

	for _, proxy := range proxies {
		proxy = strings.TrimSpace(proxy)
		proxy = strings.TrimRight(proxy, "/")
		if proxy == "direct" || proxy == "off" {
			continue
		}

		url := fmt.Sprintf("%s/github.com/%s/%s/@latest", proxy, repoOwner, repoName)
		if !isValidURL(url) {
			continue
		}

		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		var version struct{ Version string }
		if err = json.Unmarshal(body, &version); err != nil {
			continue
		}

		return version.Version, nil
	}

	return "", fmt.Errorf("failed to fetch latest version")
}

// isValidURL checks if the URL is a well-formed URL
func isValidURL(url string) bool {
	_, err := http.NewRequest("GET", url, nil)
	return err == nil
}

// contains checks if a slice contains a specific element
func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}
	return false
}
