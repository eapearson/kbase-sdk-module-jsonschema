package paths

import "net/http"
import "github.com/kbase/kbase-sdk-module-jsonschema/utils"

func HandleGetAbout(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("This is the jsonschema server"))

	if err != nil {
		// if cannot write, not much to really do.
		// TODO: log here
		println("Oops, can't write")
		utils.InternalError()
	}
}
