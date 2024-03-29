// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.2 DO NOT EDIT.
package openapi

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Entry defines model for Entry.
type Entry struct {
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	CreatorId  *string    `json:"creator_id,omitempty"`
	DeleteTime *time.Time `json:"delete_time,omitempty"`
	Id         *string    `json:"id,omitempty"`
	Text       *string    `json:"text,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

// Error defines model for Error.
type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// ListEntriesParams defines parameters for ListEntries.
type ListEntriesParams struct {
	PageSize  int32   `json:"page_size"`
	PageToken *string `json:"page_token,omitempty"`
}

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// ListEntries request
	ListEntries(ctx context.Context, params *ListEntriesParams, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) ListEntries(ctx context.Context, params *ListEntriesParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListEntriesRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewListEntriesRequest generates requests for ListEntries
func NewListEntriesRequest(server string, params *ListEntriesParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v2/entries")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	queryValues := queryURL.Query()

	if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page_size", runtime.ParamLocationQuery, params.PageSize); err != nil {
		return nil, err
	} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
		return nil, err
	} else {
		for k, v := range parsed {
			for _, v2 := range v {
				queryValues.Add(k, v2)
			}
		}
	}

	if params.PageToken != nil {

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page_token", runtime.ParamLocationQuery, *params.PageToken); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

	}

	queryURL.RawQuery = queryValues.Encode()

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// ListEntries request
	ListEntriesWithResponse(ctx context.Context, params *ListEntriesParams, reqEditors ...RequestEditorFn) (*ListEntriesResponse, error)
}

type ListEntriesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Entries       *[]Entry `json:"entries,omitempty"`
		HasNextPage   *bool    `json:"has_next_page,omitempty"`
		NextPageToken *string  `json:"next_page_token,omitempty"`
		TotalSize     *int32   `json:"total_size,omitempty"`
	}
	JSONDefault *Error
}

// Status returns HTTPResponse.Status
func (r ListEntriesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListEntriesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// ListEntriesWithResponse request returning *ListEntriesResponse
func (c *ClientWithResponses) ListEntriesWithResponse(ctx context.Context, params *ListEntriesParams, reqEditors ...RequestEditorFn) (*ListEntriesResponse, error) {
	rsp, err := c.ListEntries(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListEntriesResponse(rsp)
}

// ParseListEntriesResponse parses an HTTP response from a ListEntriesWithResponse call
func ParseListEntriesResponse(rsp *http.Response) (*ListEntriesResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &ListEntriesResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Entries       *[]Entry `json:"entries,omitempty"`
			HasNextPage   *bool    `json:"has_next_page,omitempty"`
			NextPageToken *string  `json:"next_page_token,omitempty"`
			TotalSize     *int32   `json:"total_size,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /api/v2/entries)
	ListEntries(ctx echo.Context, params ListEntriesParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ListEntries converts echo context to params.
func (w *ServerInterfaceWrapper) ListEntries(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListEntriesParams
	// ------------- Required query parameter "page_size" -------------

	err = runtime.BindQueryParameter("form", true, true, "page_size", ctx.QueryParams(), &params.PageSize)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter page_size: %s", err))
	}

	// ------------- Optional query parameter "page_token" -------------

	err = runtime.BindQueryParameter("form", true, false, "page_token", ctx.QueryParams(), &params.PageToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter page_token: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ListEntries(ctx, params)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/api/v2/entries", wrapper.ListEntries)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/5RWS28jyQ3+K0Qnx07Lmb0Eug2SWcBBduFk8jiMDYOqpiTOVFeVSZYwjqH/HrC69XDG",
	"u5i9SN1dVayPHz8+XrqQp5ITJdNu/dJp2NOE7fFDMnn2hyK5kBhT+xyE0Gh8RPO3bZbJn7oRjf5gPFHX",
	"d/ZcqFt3asJp1x37+UyWRx79zDfLI0Uyemynv9vmL9gy+mpvLtQy/kbcx/OXvPlMwdzKB5Esb5CSx9fQ",
	"OdkP7y4mORntSNzCRKq4ozcwHvtO6Kmy0NitP802L/sfvkFz7DulUIXt+aOHbYayIRSS99X2l7cfT7D+",
	"+p9/dv0cZLc0r15g7s1Kd3TDnLa5QWSLvvJzNtrk/AXe3912fXcgUc6pW3d/HG6GG/crF0pYuFt3P7RP",
	"fVfQ9g3SCguvDu9WlEwWwnbUwjCSBuFis61/kFVJChgjFDKFreQJbE+gz2rkj2jtvSoJ7FEBQyBVsHyf",
	"fsYJlEYIOY08UbI6AakN8BNSoIQKRlPJAoo7NmMFxcKUekgUQPY5haqgNF1tYAOcyAZ4T4kwARrsBA88",
	"ImDdVeoBAzCGGrkdHeDPVXDDVgXyyBliFpp6yJJQCGhHBhRpQZco9BCqaFXgESIFqzrAXyorTAxWpbD2",
	"UGo8cELxu0iyO92DcQo81mRwQOGq8Lmq5QFuE+wxwN5BoCpBiWiEMHKwOjkdt7MW3RccubAGTjvAZO7N",
	"xffIuxrx7HnZo5AJnkj0/TDlSGpMwFMhGdmZ+jcfcJodwshPFScYGZ0ZQYUn9+1AkQ1STmBZLItTwltK",
	"4/n2Ae4ESSmZw6TE0wVAlYRwyLFaQYMDJUrogGdy/WfCKm7jNl0sb0kW1rcYOLK+uqTd4D/9Jb4BNI8Y",
	"yQM79s5jIEFzx/x/gI9VC6WRneWILp4xxyy9K1ApmKu5edmk4l73cKA9hxoRvCLIWCeIvCHJA/yUZcNA",
	"lXXK43UYfLkJO2LgxDjcp/v0kcYWiaqwJRdfzJss7QDli2KkmtRpAM+NCZvBhXzW2APVV9kyhxxidR26",
	"Oge426NSjHNiFJLleKO5hZcMtlgDb+pMOJ7u8X3X5w8Ul9DxgUSwf3215wnw2J8TMfFmP8C/DArFSMlI",
	"nypByVrJM+mURAM4FXjKAk+6E5cnSye3GpN9A3KWRaopgAmruS9wYEMa4MeqgYCsVYOx8jkLvFJooEjC",
	"Dc6s39OBydVSsYkn1EkxwYQ7d5niEq0B/l7no1OOHrc5elRn7Vyg9OfiA1iDJ8m8c5Hn7PYijqXInLPR",
	"xeIBBk79BcqSuImVT4DVMQS2OrJDVUWodtLZEsj5plektfsGuLsOTGNuwViEjOt0Vblm0dT+St9eeof7",
	"1LXOIegt4Hbs1t3fWO3D0ii8h4hzQKLd+tNLx94mnirJc9d3CVsTK7ijR+X/erO89E+TSkunw+/qzMf+",
	"V8xb/kKpu7b3/737wS/X4pnfGty7m5t5MEhGqfU6LCVyaI6uPqs3vJcre6/niatOyUZTe/i90LZbd79b",
	"XSa21TKureZZ7TKwoAi29z3qY6Kv9lhezxybnCNh8i3n5cXLN6eqbBhnkr+Py29mlWOb87ZYo/0mXn7V",
	"6zaMLaavB4ma6GvxCuyletlzPS41LV0PSp8ePIBKcjgprUpcBiJdr1Yj/2nYRExfjMJ+SGTd8eH4vwAA",
	"AP///vpQ9DwLAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}

