package oa

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

type CLI struct {
	Out     io.Writer
	Schema  bytes.Buffer
	BaseURL url.URL
}

type Runner interface {
	Run() error
}

func New(out io.Writer, schema bytes.Buffer, url url.URL) *CLI {
	return &CLI{
		Out:     out,
		Schema:  schema,
		BaseURL: url,
	}
}

func (cli *CLI) Run(path string, method string, body string, token string) error {
	ctx := context.Background()

	doc, err := openapi3.NewLoader().LoadFromData(cli.Schema.Bytes())
	if err != nil {
		return fmt.Errorf("load doc: %w", err)
	}

	err = doc.Validate(ctx)
	if err != nil {
		return fmt.Errorf("validate doc: %w", err)
	}

	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		return fmt.Errorf("new router: %w", err)
	}

	err = cli.validatePath(ctx, path, router, method, body, token)
	if err != nil {
		return err
	}

	return nil
}

const separator = "────────────────────────────────────"

func (cli *CLI) validatePath(ctx context.Context, path string, router routers.Router, method string, body string, token string) error {
	fmt.Fprintf(cli.Out, "%s\n%s\n\n", separator, path)

	req, err := http.NewRequest(method, cli.BaseURL.String()+path, strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	reqInput, err := cli.validateRequest(ctx, router, req, token)
	if err != nil {
		return err
	}
	if err := cli.doRequest(ctx, req, reqInput); err != nil {
		return err
	}

	return nil
}

const MsgRouteGood = "✓ find route\n"
const MsgReqGood = "\n✓ request valid\n"
const MsgResGood = "✓ response valid\n\n"

func (cli *CLI) validateRequest(ctx context.Context, router routers.Router, req *http.Request, token string) (*openapi3filter.RequestValidationInput, error) {
	req.Header.Set("Content-Type", "application/json")
	bearer := "Bearer " + token
	req.Header.Add("Authorization", bearer)
	route, pathParams, err := router.FindRoute(req)
	if err != nil {
		return nil, fmt.Errorf("find route: %w", err)
	}
	fmt.Fprint(cli.Out, MsgRouteGood)

	reqInput := &openapi3filter.RequestValidationInput{
		Request:     req,
		PathParams:  pathParams,
		QueryParams: req.URL.Query(),
		Route:       route,
		// Options: nil, // ?
		// ParamDecoder: nil, // ?
	}

	b, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] dump request: %+v", err)
	}
	fmt.Fprintf(cli.Out, "%s", strings.ReplaceAll("\t"+string(b), "\n", "\n\t"))

	if err := openapi3filter.ValidateRequest(ctx, reqInput); err != nil {
		return nil, fmt.Errorf("validate request: %w", err)
	}
	fmt.Fprint(cli.Out, MsgReqGood)

	return reqInput, nil
}

func (cli *CLI) doRequest(ctx context.Context, req *http.Request, reqInput *openapi3filter.RequestValidationInput) error {
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	buf := bytes.Buffer{}
	if err := json.Indent(&buf, body, "\t", "  "); err != nil {
		return err
	}
	fmt.Fprintf(cli.Out, "\t%s", buf.String())

	resInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: reqInput,
		Status:                 response.StatusCode,
		Header:                 response.Header,
		Options:                nil, // ?
	}
	resInput.SetBodyBytes(body)

	if err := openapi3filter.ValidateResponse(ctx, resInput); err != nil {
		return fmt.Errorf("validate response: %w", err)
	}
	fmt.Fprint(cli.Out, MsgResGood)
	return nil
}

func (cli *CLI) DumpRoutes() error {
	doc, err := openapi3.NewLoader().LoadFromData(cli.Schema.Bytes())
	if err != nil {
		return fmt.Errorf("load doc: %w", err)
	}

	expectType := reflect.TypeOf(&openapi3.Operation{})
	fmt.Fprintf(cli.Out, "%-10s\t%-10s\t%s\n", "Endpoint", "Method", "ID")
	fmt.Fprintf(cli.Out, "%-10s\t%-10s\t%s\n", strings.Repeat("─", 10), strings.Repeat("─", 10), strings.Repeat("─", 10))
	for k, path := range doc.Paths {
		rv := reflect.ValueOf(path).Elem()
		rt := reflect.TypeOf(path).Elem()
		for i := 0; i < rt.NumField(); i++ {
			rf := rt.Field(i)
			if !rf.Type.AssignableTo(expectType) {
				continue
			}
			rfv := rv.Field(i)
			if rfv.IsNil() {
				continue
			}
			op := rfv.Interface().(*openapi3.Operation)
			fmt.Fprintf(cli.Out, "%-10s\t%-10s\t%s\n", k, rf.Name, op.OperationID)
		}
	}

	return nil
}
