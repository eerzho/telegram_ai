package createsetting

import (
	"net/http"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	httpjson "github.com/eerzho/goiler/pkg/http_json"
	httpserver "github.com/eerzho/goiler/pkg/http_server"
)

func HTTPv1(usecase *Usecase) httpserver.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "create_setting.HTTPv1"

		ctx := r.Context()
		defer r.Body.Close()

		input, err := httpjson.Decode[Input](r)
		if err != nil {
			return errorhelp.WithOP(op, err)
		}

		output, err := usecase.Execute(ctx, input)
		if err != nil {
			return errorhelp.WithOP(op, err)
		}

		httpjson.Encode(w, http.StatusCreated, output)
		return nil
	}
}
