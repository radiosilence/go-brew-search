package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/user/go-brew-search/internal/cache"
)

const (
	formulaeAPIURL = "https://formulae.brew.sh/api/formula.json"
	casksAPIURL    = "https://formulae.brew.sh/api/cask.json"
)

type Package struct {
	Token       string `json:"token,omitempty"`       // for casks
	Name        string `json:"name,omitempty"`         // for formulae
	FullName    string `json:"full_name,omitempty"`    // for formulae
	Description string `json:"desc,omitempty"`         // for both
	Homepage    string `json:"homepage,omitempty"`     // for both
	Version     string `json:"version,omitempty"`      // for both
	Type        string // "formula" or "cask"
}

type Client struct {
	cache      *cache.Manager
	httpClient *http.Client
}

func New(cacheManager *cache.Manager) *Client {
	return &Client{
		cache: cacheManager,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) FetchAllPackages() ([]Package, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var allPackages []Package
	var fetchErr error

	wg.Add(2)

	// Fetch formulae
	go func() {
		defer wg.Done()
		formulae, err := c.fetchFormulae()
		if err != nil {
			mu.Lock()
			fetchErr = fmt.Errorf("failed to fetch formulae: %w", err)
			mu.Unlock()
			return
		}
		mu.Lock()
		allPackages = append(allPackages, formulae...)
		mu.Unlock()
	}()

	// Fetch casks
	go func() {
		defer wg.Done()
		casks, err := c.fetchCasks()
		if err != nil {
			mu.Lock()
			fetchErr = fmt.Errorf("failed to fetch casks: %w", err)
			mu.Unlock()
			return
		}
		mu.Lock()
		allPackages = append(allPackages, casks...)
		mu.Unlock()
	}()

	wg.Wait()

	if fetchErr != nil {
		return nil, fetchErr
	}

	return allPackages, nil
}

func (c *Client) fetchFormulae() ([]Package, error) {
	// Check cache first
	var formulae []map[string]interface{}
	if err := c.cache.Get("formulae", &formulae); err == nil {
		return c.parseFormulae(formulae), nil
	}

	// Fetch from API
	resp, err := c.httpClient.Get(formulaeAPIURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &formulae); err != nil {
		return nil, err
	}

	// Cache the result
	c.cache.Set("formulae", formulae)

	return c.parseFormulae(formulae), nil
}

func (c *Client) fetchCasks() ([]Package, error) {
	// Check cache first
	var casks []map[string]interface{}
	if err := c.cache.Get("casks", &casks); err == nil {
		return c.parseCasks(casks), nil
	}

	// Fetch from API
	resp, err := c.httpClient.Get(casksAPIURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &casks); err != nil {
		return nil, err
	}

	// Cache the result
	c.cache.Set("casks", casks)

	return c.parseCasks(casks), nil
}

func (c *Client) parseFormulae(formulae []map[string]interface{}) []Package {
	packages := make([]Package, 0, len(formulae))
	for _, f := range formulae {
		pkg := Package{
			Type: "formula",
		}

		if name, ok := f["name"].(string); ok {
			pkg.Name = name
			pkg.Token = name // Use name as token for formulae
		}

		if fullName, ok := f["full_name"].(string); ok {
			pkg.FullName = fullName
		}

		if desc, ok := f["desc"].(string); ok {
			pkg.Description = desc
		}

		if homepage, ok := f["homepage"].(string); ok {
			pkg.Homepage = homepage
		}

		if versions, ok := f["versions"].(map[string]interface{}); ok {
			if stable, ok := versions["stable"].(string); ok {
				pkg.Version = stable
			}
		}

		if pkg.Name != "" {
			packages = append(packages, pkg)
		}
	}
	return packages
}

func (c *Client) parseCasks(casks []map[string]interface{}) []Package {
	packages := make([]Package, 0, len(casks))
	for _, cs := range casks {
		pkg := Package{
			Type: "cask",
		}

		if token, ok := cs["token"].(string); ok {
			pkg.Token = token
			pkg.Name = token // Use token as name for casks
		}

		if name, ok := cs["name"].([]interface{}); ok && len(name) > 0 {
			if n, ok := name[0].(string); ok {
				pkg.FullName = n
			}
		}

		if desc, ok := cs["desc"].(string); ok {
			pkg.Description = desc
		}

		if homepage, ok := cs["homepage"].(string); ok {
			pkg.Homepage = homepage
		}

		if version, ok := cs["version"].(string); ok {
			pkg.Version = version
		}

		if pkg.Token != "" {
			packages = append(packages, pkg)
		}
	}
	return packages
}