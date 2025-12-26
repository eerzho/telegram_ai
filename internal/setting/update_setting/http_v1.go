package updatesetting

import (
	"net/http"
	"strconv"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	httpjson "github.com/eerzho/goiler/pkg/http_json"
	httpserver "github.com/eerzho/goiler/pkg/http_server"
	"github.com/eerzho/telegram_ai/internal/domain"
)

func HTTPv1(usecase *Usecase) httpserver.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "update_setting.HTTPv1"

		ctx := r.Context()
		defer r.Body.Close()

		input, err := httpjson.Decode[Input](r)
		if err != nil {
			return errorhelp.WithOP(op, err)
		}
		input.UserID, err = strconv.ParseInt(r.PathValue("user_id"), 10, 64)
		if err != nil {
			return errorhelp.WithOP(op, domain.ErrSettingNotFound)
		}
		input.ChatID, err = strconv.ParseInt(r.PathValue("chat_id"), 10, 64)
		if err != nil {
			return errorhelp.WithOP(op, domain.ErrSettingNotFound)
		}

		output, err := usecase.Execute(ctx, input)
		if err != nil {
			return errorhelp.WithOP(op, err)
		}

		httpjson.Encode(w, http.StatusOK, output)
		return nil
	}
}
