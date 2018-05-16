package bitbucket

import (
	"encoding/json"
	"errors"
	"fmt"

	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Client struct
type Client struct {
	Auth         *auth
	Repositories *Repositories
	Pagelen      uint64
}

type auth struct {
	accessToken string
}

// NewAccessToken set token
func NewAccessToken(token string) *Client {
	a := &auth{accessToken: token}
	return injectClient(a)
}

// DEFAULT_PAGE_LENGHT default page size
const defaultPageLength = 10

func injectClient(a *auth) *Client {
	c := &Client{Auth: a, Pagelen: defaultPageLength}
	c.Repositories = &Repositories{
		c:            c,
		PullRequests: &PullRequests{c: c},
		Branch:       &Branch{c: c},
		Path:         &Path{c: c},
	}
	return c
}

func (c *Client) execute(method string, urlStr string, text string, content string) (interface{}, error) {
	// Use pagination if changed from default value
	const decRadix = 10
	if strings.Contains(urlStr, "/repositories/") {
		if c.Pagelen != defaultPageLength {
			urlObj, err := url.Parse(urlStr)
			if err != nil {
				return nil, err
			}
			q := urlObj.Query()
			q.Set("pagelen", strconv.FormatUint(c.Pagelen, decRadix))
			urlObj.RawQuery = q.Encode()
			urlStr = urlObj.String()
		}
	}

	body := strings.NewReader(text)
	req, err := http.NewRequest(method, urlStr, body)

	req.Header.Set("Content-Type", content)
	req.Header.Set("Authorization", " Bearer "+c.Auth.accessToken)

	if err != nil {
		return nil, err
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if (resp.StatusCode != http.StatusOK) && (resp.StatusCode != http.StatusCreated) && (resp.StatusCode >= 500) {
		return nil, fmt.Errorf(resp.Status)
	}

	if resp.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}

	resBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = json.Unmarshal(resBodyBytes, &result)
	if err != nil {
		return nil, err
	}

	resultMap, isMap := result.(map[string]interface{})
	if isMap {
		nextIn := resultMap["next"]
		valuesIn := resultMap["values"]
		if nextIn != nil && valuesIn != nil {
			nextURL := nextIn.(string)
			if nextURL != "" {
				valuesSlice := valuesIn.([]interface{})
				if valuesSlice != nil {
					nextResult, err := c.execute(method, nextURL, text, "application/json")
					if err != nil {
						return nil, err
					}
					nextResultMap, isNextMap := nextResult.(map[string]interface{})
					if !isNextMap {
						return nil, fmt.Errorf("next page result is not map, it's %T", nextResult)
					}
					nextValuesIn := nextResultMap["values"]
					if nextValuesIn == nil {
						return nil, fmt.Errorf("next page result has no values")
					}
					nextValuesSlice, isSlice := nextValuesIn.([]interface{})
					if !isSlice {
						return nil, fmt.Errorf("next page result 'values' is not slice")
					}
					valuesSlice = append(valuesSlice, nextValuesSlice...)
					resultMap["values"] = valuesSlice
					delete(resultMap, "page")
					delete(resultMap, "pagelen")
					delete(resultMap, "size")
					result = resultMap
				}
			}
		}
	}

	// Check for errors
	value, b := result.(map[string]interface{})["errors"]

	if b {
		// build content with pretty print
		content, _ := json.MarshalIndent(value, "", "\t")
		return nil, errors.New(string(content))
	}

	// Only for debug
	if req.Header.Get("X-Remote-Debug") == "true" {
		// build result with pretty print
		res, _ := json.MarshalIndent(result, "", "\t")
		fmt.Println("RESULT", string(res))
	}

	return result, nil
}

func (c *Client) requestURL(template string, args ...interface{}) string {
	if len(args) == 1 && args[0] == "" {
		return getAPIBaseURL() + template
	}
	return getAPIBaseURL() + fmt.Sprintf(template, args...)
}
