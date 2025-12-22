package healthcheck

import (
	"net/http"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	httpjson "github.com/eerzho/goiler/pkg/http_json"
	httpserver "github.com/eerzho/goiler/pkg/http_server"
)

// HTTPv1 godoc
//
// @summary health check
// @tags monitoring
//
// @produce json
// @success 200 {object} Output
// @failure 400 {object} httpjson.Error
// @failure 500 {object} httpjson.Error
//
// @router /_hc [get].
func HTTPv1(usecase *Usecase) httpserver.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "health_check.HTTPv1"

		ctx := r.Context()
		defer r.Body.Close()

		output, err := usecase.Execute(ctx, Input{})
		if err != nil {
			return errorhelp.WithOP(op, err)
		}

		httpjson.Encode(w, http.StatusOK, output)
		return nil
	}
}
