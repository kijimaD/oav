package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

//go:embed openapi.yml
var spec []byte

type Res struct {
	Pets []Pet `json: pets`
}

type Pet struct {
	ID int `json:id`
}

func main() {
	http.HandleFunc("/pets", pets)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func pets(w http.ResponseWriter, r *http.Request) {
	resp := Res{
		Pets: []Pet{
			{
				ID: 1,
			},
			{
				ID: 2,
			},
		},
	}

	json.NewEncoder(w).Encode(resp)
}

func handle() {
	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

func run() error {
	ctx := context.Background()

	doc, err := openapi3.NewLoader().LoadFromData(spec)
	if err != nil {
		return fmt.Errorf("load doc: %w", err)
	}
	if doc.Validate(ctx); err != nil {
		return fmt.Errorf("validate doc: %w", err)
	}

	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		return fmt.Errorf("new router: %w", err)
	}

	err = testPath("/pets", router, ctx)
	if err != nil {
		return err
	}

	return nil
}

func testPath(path string, router routers.Router, ctx context.Context) error {
	baseURL := "http://localhost:8080"

	log.Printf("%s ----------------------------------------", path)
	req, err := http.NewRequest("GET", baseURL+path, strings.NewReader(`{}`))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	if err := doRequest(ctx, router, req); err != nil {
		log.Println(strings.ReplaceAll(err.Error(), "\n", "\n\t"))
	}

	return nil
}

func doRequest(ctx context.Context, router routers.Router, req *http.Request) error {
	req.Header.Set("Content-Type", "application/json")
	route, pathParams, err := router.FindRoute(req)
	if err != nil {
		return fmt.Errorf("find route: %w", err)
	}
	log.Println("find route is ok")

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
		log.Printf("[ERROR] dump request: %+v", err)
	}
	fmt.Println(strings.ReplaceAll("\t"+string(b), "\n", "\n\t"))

	if err := openapi3filter.ValidateRequest(ctx, reqInput); err != nil {
		log.Printf("validate request is failed: %T", err)
		return fmt.Errorf("validate request: %w", err)
	}
	log.Println("request is ok")

	rec := httptest.NewRecorder()
	func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		resp, _ := http.Get(fmt.Sprintf("%s", req.URL.String()))
		defer resp.Body.Close()
		byteArray, _ := ioutil.ReadAll(resp.Body)

		var jsonBody map[string]interface{}
		err = json.Unmarshal(byteArray, &jsonBody)
		if err != nil {
			panic(err)
		}
		json.NewEncoder(w).Encode(jsonBody)
	}(rec, req)

	res := rec.Result()
	buf := new(bytes.Buffer)
	res.Body = io.NopCloser(io.TeeReader(res.Body, buf))

	b, err = httputil.DumpResponse(res, true)
	if err != nil {
		log.Printf("[ERROR] dump request: %+v", err)
		return err
	}
	fmt.Println(strings.ReplaceAll("\t"+string(b), "\n", "\n\t"))

	res.Body = io.NopCloser(buf)
	resInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: reqInput,
		Status:                 200,
		Header:                 res.Header,
		Body:                   res.Body,
		Options:                nil, // ?
	}
	if err := openapi3filter.ValidateResponse(ctx, resInput); err != nil {
		log.Printf("valicate response is failed: %T", err)
		return fmt.Errorf("validate response: %w", err)
	}
	log.Println("response is ok")
	return nil
}

func dumpRoutes(doc *openapi3.T) {
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
			fmt.Printf("%-10s\t%-10s\t%s\n", k, rf.Name, op.OperationID)
		}
	}
}
