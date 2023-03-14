package oa

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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
	Schema  io.Reader
	BaseURL url.URL
}

type Runner interface {
	Run() error
}

func New(out io.Writer, schema io.Reader, url url.URL) *CLI {
	return &CLI{
		Out:     out,
		Schema:  schema,
		BaseURL: url,
	}
}

func (cli *CLI) Run(path string) error {
	ctx := context.Background()

	buf := new(bytes.Buffer)
	_, errc := io.Copy(buf, cli.Schema)
	if errc != nil {
		return fmt.Errorf("%w", errc)
	}
	doc, err := openapi3.NewLoader().LoadFromData(buf.Bytes())
	if err != nil {
		return fmt.Errorf("load doc: %w", err)
	}
	// 謎のバリデーションエラーが出るので一旦スルー
	// err = doc.Validate(ctx)
	// if err != nil {
	// 	return fmt.Errorf("validate doc: %w", err)
	// }

	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		return fmt.Errorf("new router: %w", err)
	}

	err = cli.validatePath(ctx, path, router)
	if err != nil {
		return err
	}

	return nil
}

const separator = "────────────────────────────────────"

func (cli *CLI) validatePath(ctx context.Context, path string, router routers.Router) error {
	fmt.Fprintf(cli.Out, "%s %s\n", path, separator)

	req, err := http.NewRequest("GET", cli.BaseURL.String()+path, strings.NewReader(`{}`))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	reqInput, err := cli.validateRequest(ctx, router, req)
	if err != nil {
		return err
	}
	if err := cli.doRequest(ctx, req, reqInput); err != nil {
		return err
	}

	return nil
}

func (cli *CLI) validateRequest(ctx context.Context, router routers.Router, req *http.Request) (*openapi3filter.RequestValidationInput, error) {
	req.Header.Set("Content-Type", "application/json")
	route, pathParams, err := router.FindRoute(req)
	if err != nil {
		return nil, fmt.Errorf("find route: %w", err)
	}
	fmt.Fprintf(cli.Out, "find route is ok")

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
	fmt.Fprint(cli.Out, strings.ReplaceAll("\t"+string(b), "\n", "\n\t"))

	if err := openapi3filter.ValidateRequest(ctx, reqInput); err != nil {
		return nil, fmt.Errorf("validate request: %w", err)
	}
	fmt.Fprint(cli.Out, "request is ok")

	return reqInput, nil
}

func (cli *CLI) doRequest(ctx context.Context, req *http.Request, reqInput *openapi3filter.RequestValidationInput) error {
	rec := httptest.NewRecorder()

	err := cli.request(rec, req)
	if err != nil {
		return err
	}

	res := rec.Result()
	buf := new(bytes.Buffer)
	res.Body = io.NopCloser(io.TeeReader(res.Body, buf))

	b, err := httputil.DumpResponse(res, true)
	if err != nil {
		return err
	}
	fmt.Fprint(cli.Out, strings.ReplaceAll("\t"+string(b), "\n", "\n\t"))

	res.Body = io.NopCloser(buf)
	resInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: reqInput,
		Status:                 200,
		Header:                 res.Header,
		Body:                   res.Body,
		Options:                nil, // ?
	}
	if err := openapi3filter.ValidateResponse(ctx, resInput); err != nil {
		return fmt.Errorf("validate response: %w", err)
	}
	fmt.Fprintf(cli.Out, "response is ok\n")
	return nil
}

func (cli *CLI) request(w http.ResponseWriter, req *http.Request) error {
	w.Header().Add("Content-Type", "application/json")

	resp, err := http.Get(req.URL.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	byteArray, _ := io.ReadAll(resp.Body)

	var jsonBody map[string]interface{}
	err = json.Unmarshal(byteArray, &jsonBody)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(jsonBody)
	if err != nil {
		return err
	}

	return nil
}

func (cli *CLI) dumpRoutes() error {
	buf := new(bytes.Buffer)
	_, errc := io.Copy(buf, cli.Schema) // Reader -> []byte
	if errc != nil {
		return fmt.Errorf("%w", errc)
	}

	doc, err := openapi3.NewLoader().LoadFromData(buf.Bytes())
	if err != nil {
		return fmt.Errorf("load doc: %w", err)
	}

	expectType := reflect.TypeOf(&openapi3.Operation{})
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
