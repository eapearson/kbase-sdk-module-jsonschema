package paths

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/kbase/kbase-sdk-module-jsonschema/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type RequestError struct {
	StatusCode int
	Message    string
}

func (e *RequestError) Error() string {
	return e.Message
}

type JSONSchema interface{}

func getAbsoluteSchema(path string, schema string, model int, revision int, addition int) ([]byte, error) {
	resolvedSchemaName := strings.Join([]string{schema, strings.Join([]string{strconv.Itoa(model), strconv.Itoa(revision), strconv.Itoa(addition)}, "-"), "json"}, ".")
	path = filepath.Join(path, resolvedSchemaName)

	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, &RequestError{404, "schema not found"}
		} else {
			return nil, &RequestError{400, "error opening schema file"}
		}
	} else {
		return content, nil
	}
}

func getMostRecentAdditionToSchema(path string, schema string, model int, revision int) ([]byte, error) {
	schemaFilePattern := strings.Join([]string{schema, strings.Join([]string{strconv.Itoa(model), strconv.Itoa(revision), "*"}, "-"), "json"}, ".")
	matches, err := filepath.Glob(filepath.Join(path, schemaFilePattern))
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, &RequestError{404, "schema not found"}
	}
	maxAddition := 0
	for _, fileName := range matches {
		parts := strings.SplitN(filepath.Base(fileName), ".", 3)
		versionMatch := strings.SplitN(parts[1], "-", 3)
		addition, err := strconv.Atoi(versionMatch[2])
		if err != nil {
			return nil, err
		}
		if addition > maxAddition {
			maxAddition = addition
		}
	}
	return getAbsoluteSchema(path, schema, model, revision, maxAddition)
}

func getMostRecentRevisionToSchema(path string, schema string, model int) ([]byte, error) {
	schemaFilePattern := strings.Join([]string{schema, strings.Join([]string{strconv.Itoa(model), "*", "*"}, "-"), "json"}, ".")
	matches, err := filepath.Glob(filepath.Join(path, schemaFilePattern))
	if len(matches) == 0 {
		return nil, &RequestError{404, "schema not found"}
	}
	if err != nil {
		return nil, err
	}
	maxRevision := 0
	maxAddition := 0
	for _, fileName := range matches {
		parts := strings.SplitN(filepath.Base(fileName), ".", 3)
		versionMatch := strings.SplitN(parts[1], "-", 3)
		revision, err := strconv.Atoi(versionMatch[1])
		if err != nil {
			return nil, err
		}
		addition, err := strconv.Atoi(versionMatch[2])
		if err != nil {
			return nil, err
		}
		if revision > maxRevision {
			maxRevision = revision
			maxAddition = addition
		} else {
			if addition > maxAddition {
				maxAddition = addition
			}
		}
	}
	return getAbsoluteSchema(path, schema, model, maxRevision, maxAddition)
}

func getMostRecentSchema(path string, schema string) ([]byte, error) {
	schemaFilePattern := strings.Join([]string{schema, strings.Join([]string{"*", "*", "*"}, "-"), "json"}, ".")
	matches, err := filepath.Glob(filepath.Join(path, schemaFilePattern))
	if len(matches) == 0 {
		return nil, &RequestError{404, "schema not found"}
	}
	if err != nil {
		return nil, err
	}
	maxModel := 0
	maxRevision := 0
	maxAddition := 0
	for _, fileName := range matches {
		parts := strings.SplitN(filepath.Base(fileName), ".", 3)
		versionMatch := strings.SplitN(parts[1], "-", 3)
		model, err := strconv.Atoi(versionMatch[0])
		if err != nil {
			return nil, err
		}
		revision, err := strconv.Atoi(versionMatch[1])
		if err != nil {
			return nil, err
		}
		addition, err := strconv.Atoi(versionMatch[2])
		if err != nil {
			return nil, err
		}
		if model > maxModel {
			maxModel = model
			maxRevision = revision
			maxAddition = addition
		} else if revision > maxRevision {
			maxRevision = revision
			maxAddition = addition
		} else {
			if addition > maxAddition {
				maxAddition = addition
			}
		}
	}
	return getAbsoluteSchema(path, schema, maxModel, maxRevision, maxAddition)
}

func getSchema(path string, schema string, version string) ([]byte, error) {
	rootPath := os.Getenv("SCHEMA_ROOT")
	if rootPath == "" {
		return nil, errors.New("SCHEMA_ROOT environment variable not set")
	}

	if version == "" {
		return getMostRecentSchema(filepath.Join(rootPath, path), schema)
	}

	versionParts := strings.Split(version, "-")

	switch len(versionParts) {
	case 3:
		model, err := strconv.Atoi(versionParts[0])
		if err != nil {
			return nil, errors.New("invalid model (first component) in version")
		}
		revision, err := strconv.Atoi(versionParts[1])
		if err != nil {
			return nil, errors.New("invalid revision (second component) in version")
		}
		addition, err := strconv.Atoi(versionParts[2])
		if err != nil {
			return nil, errors.New("invalid addition (third component) in version")
		}
		return getAbsoluteSchema(filepath.Join(rootPath, path), schema, model, revision, addition)
	case 2:
		model, err := strconv.Atoi(versionParts[0])
		if err != nil {
			return nil, errors.New("invalid model (first component) in version")
		}
		revision, err := strconv.Atoi(versionParts[1])
		if err != nil {
			return nil, errors.New("invalid revision (second component) in version")
		}
		return getMostRecentAdditionToSchema(filepath.Join(rootPath, path), schema, model, revision)
	case 1:
		model, err := strconv.Atoi(versionParts[0])
		if err != nil {
			return nil, errors.New("invalid model (first component) in version")
		}
		return getMostRecentRevisionToSchema(filepath.Join(rootPath, path), schema, model)
	default:
		return nil, errors.New("invalid schema version")

	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_, err := w.Write([]byte(message))
	if err != nil {
		utils.InternalError()
	}
}

func HandleGetSchema(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	content, err := getSchema(params["path"], params["schema"], params["version"])
	if err != nil {
		// if cannot write, not much to really do.
		// TODO: log here
		//if errors.Is(err, &RequestError{}) {
		//	writeError(w, err.StatusCode, err.Error())
		//}

		if t, ok := err.(*RequestError); ok {
			writeError(w, t.StatusCode, t.Error())
		} else {
			writeError(w, http.StatusBadRequest, err.Error())
		}
		//switch err {
		//case RequestError:
		//	writeClientError(w, err.Serr.Error())
		//}
	} else {
		_, err := w.Write(content)
		if err != nil {
			// nothing to do if wan't write!
			// perhaps log
		}
	}
}
