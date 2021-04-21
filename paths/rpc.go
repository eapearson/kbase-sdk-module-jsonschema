package paths

import (
	"encoding/json"
	"errors"
	"github.com/kbase/kbase-sdk-module-jsonschema/utils"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type RPCRequest struct {
	Version string          `json:"version"`
	ID      string          `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type RPCResult interface {
}

type RPCResultResponse struct {
	Version string    `json:"version"`
	ID      string    `json:"id,omitempty"`
	Result  RPCResult `json:"result"`
}

type RPCErrorDetail interface{}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   RPCErrorDetail
}

type RPCErrorResponse struct {
	Version string   `json:"version"`
	ID      string   `json:"id,omitempty"`
	Error   RPCError `json:"error"`
}

func makeErrorResponse(id string, error RPCError) RPCErrorResponse {
	return RPCErrorResponse{
		Version: "1.1",
		ID:      id,
		Error:   error,
	}
}

func makeResultResponse(id string, result RPCResult) RPCResultResponse {
	return RPCResultResponse{
		Version: "1.1",
		ID:      id,
		Result:  result,
	}
}

type Params interface {
}

type GetSchemaParams struct {
	Path    string `json:"path"`
	Version string `json:"version"`
	Name    string `json:"name"`
}

type GetSchemaResult struct {
	Schema string          `json:"schema"`
	Params GetSchemaParams `json:"params"`
	Path   string          `json:"path"`
}

func getSchemaMethod(rpc RPCRequest) (*GetSchemaResult, error) {
	println(rpc.Method)
	params := GetSchemaParams{}
	parseErr := json.Unmarshal(rpc.Params, &params)
	if parseErr != nil {
		//writeErrorResponse()
		// TODO: do something!
		return nil, errors.New("Error parsing rpc request: " + parseErr.Error())
	}

	// handle the params.

	// basically, we build a path too the schema
	// TODO: handle partial or missing version.

	fileName := strings.Join([]string{params.Name, params.Version, "json"}, ".")

	rootPath := os.Getenv("SCHEMA_ROOT")
	if rootPath == "" {
		// return error
		return nil, errors.New("SCHEMA_ROOT environment variable not set")
	}

	// ensure file exists...
	// ensure file is under the rootPath directory
	fullPath := filepath.Join(rootPath, params.Path, fileName)

	schema, err := ioutil.ReadFile(fullPath)

	if err != nil {
		//return nil, errors.New("Can't open schema file")
		return nil, &utils.CantOpenFileError{fullPath, err.Error()}
	}

	//schema := JSONSchema{}
	//parseErr = json.Unmarshal(data, )

	return &GetSchemaResult{
		Schema: string(schema[:]),
		Params: params,
		Path:   fullPath,
	}, nil
}

func writeErrorResponse(w http.ResponseWriter, rpc RPCRequest, code int, message string, errorInfo interface{}) {
	w.WriteHeader(http.StatusBadRequest)
	// w.Write([]byte(http.StatusText(http.StatusBadRequest)))
	response := makeErrorResponse(rpc.ID, RPCError{
		Code:    code,
		Message: message,
		Error:   errorInfo,
	})
	body, _ := json.Marshal(response)
	_, err := w.Write(body)
	if err != nil {
		utils.InternalError()
	}
}

func writeResultResponse(w http.ResponseWriter, rpc RPCRequest, result RPCResult) {
	w.WriteHeader(http.StatusOK)
	response := makeResultResponse(rpc.ID, result)
	body, _ := json.Marshal(response)
	_, err := w.Write(body)
	if err != nil {
		utils.InternalError()
	}
}

type MethodNotFoundData struct {
	Method string `json:"method"`
}

func writeMethodNotFound(w http.ResponseWriter, rpc RPCRequest, method string) {
	data := &MethodNotFoundData{method}

	writeErrorResponse(w, rpc, -32601, "Method not found", data)
}

type MethodErrorData struct {
	Message string `json:"message"`
}

func writeMethodError(w http.ResponseWriter, rpc RPCRequest, message string) {
	data := &MethodErrorData{message}

	writeErrorResponse(w, rpc, -32001, "Method error", data)
}

func HandleRPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotImplemented)
		_, err := w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
		if err != nil {
			utils.InternalError()
		}
		return

	}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		if err != nil {
			utils.InternalError()
		}
	}
	rpc := RPCRequest{}
	parseErr := json.Unmarshal(reqBody, &rpc)
	if parseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		// w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		response := makeErrorResponse(rpc.ID, RPCError{
			Code:    123,
			Message: "Parse error",
		})
		body, _ := json.Marshal(response)
		_, err := w.Write(body)
		if err != nil {
			utils.InternalError()
		}
		return
	}
	//var result *RPCResult
	switch rpc.Method {
	case "get-schema":
		result, err := getSchemaMethod(rpc)
		if err != nil {
			// TODO: pass data
			writeMethodError(w, rpc, err.Error())
		} else {
			writeResultResponse(w, rpc, result)
		}
	default:
		writeMethodNotFound(w, rpc, rpc.Method)
	}
}
