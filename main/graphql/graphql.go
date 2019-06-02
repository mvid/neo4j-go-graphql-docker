package graphql

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/jensneuse/graphql-go-tools/pkg/middleware"
	"github.com/jensneuse/graphql-go-tools/pkg/proxy"
	httpproxy "github.com/jensneuse/graphql-go-tools/pkg/proxy/http"
)

var singleton *httpproxy.Proxy

func getBackendURL() (*url.URL, error) {
	hostURLStr, ok := os.LookupEnv("GRAPHQL_BACKEND_URL")
	if !ok {
		return nil, fmt.Errorf("no graphql backend url set")
	}
	return url.Parse(hostURLStr)
}

func setAuthHeader(hs *http.Header) error {
	authStr, ok := os.LookupEnv("NEO4J_AUTH")
	if !ok {
		return fmt.Errorf("no neo credentials set")
	}

	b64auth := base64.StdEncoding.EncodeToString([]byte(authStr))
	hs.Set("Authorization", fmt.Sprintf("Basic %s", b64auth))

	return nil
}

func getProxy() (*httpproxy.Proxy, error) {
	if singleton != nil {
		return singleton, nil
	}

	schemaPath, err := filepath.Rel("app/graphql", "public.graphql")
	if err != nil {
		return nil, err
	}
	schema, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}
	hostURL, err := getBackendURL()
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	if err := setAuthHeader(&headers); err != nil {
		return nil, err
	}
	requestConfig := proxy.RequestConfig{
		Schema: &schema,
		BackendURL: *hostURL,
		BackendHeaders: headers,
	}
	configProvider := proxy.NewStaticRequestConfigProvider(requestConfig)
	prx := httpproxy.NewDefaultProxy(configProvider, &middleware.ContextMiddleware{})

	return prx, nil
}

func InitializeSchema() error {

	// set auth headers
	headers := http.Header{}
	if err := setAuthHeader(&headers); err != nil {
		return err
	}
	// get url
	hostURL, err := getBackendURL()
	if err != nil {
		return err
	}
	idlURLPath, err := url.Parse("/graphql/idl")
	if err != nil {
		return err
	}
	idlURL := hostURL.ResolveReference(idlURLPath)

	// get schema
	schemaPath, err := filepath.Rel("app/graphql","private.graphql")
	if err != nil {
		return err
	}
	schema, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return err
	}

	// send schema
	_, err = http.NewRequest("POST", idlURL.String(), bytes.NewReader(schema))
	if err != nil {
		return err
	}

	return nil
}

func Handler(rw http.ResponseWriter, req *http.Request) {
	// get proxy instance
	prx, err := getProxy()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		if _, err := rw.Write([]byte(fmt.Sprintf("error %s", err))); err != nil {
			panic(err)
		}
	}

	// do the user auth check, in the example case, look at a specific header
	email := req.Header.Get("Authentication")

	// create the context to pass to the proxy
	ctxWithEmail := context.WithValue(req.Context(), "email", []byte(email))
	prx.ServeHTTP(rw, req.WithContext(ctxWithEmail))
}