package getsetting

import (
	"net/http"
	"strconv"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	httpjson "github.com/eerzho/goiler/pkg/http_json"
	httpserver "github.com/eerzho/goiler/pkg/http_server"
	"github.com/eerzho/telegram_ai/internal/domain"
)

// HTTPv1 godoc
//
// @tags setting
// @summary get setting
//
// @param user_id path int true "UserID"
// @param chat_id path int true "ChatID"
//
// @produce json
// @success 200 {object} Output
// @failure 400 {object} httpjson.Error
// @failure 500 {object} httpjson.Error
//
// @router /v1/settings/{user_id}/{chat_id} [get].
func HTTPv1(usecase *Usecase) httpserver.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "get_setting.HTTPv1"

		ctx := r.Context()

		userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
		if err != nil {
			return errorhelp.WithOP(op, domain.ErrSettingNotFound)
		}
		chatID, err := strconv.ParseInt(r.PathValue("chat_id"), 10, 64)
		if err != nil {
			return errorhelp.WithOP(op, domain.ErrSettingNotFound)
		}
		input := Input{UserID: userID, ChatID: chatID}

		output, err := usecase.Execute(ctx, input)
		if err != nil {
			return errorhelp.WithOP(op, err)
		}

		httpjson.Encode(w, http.StatusOK, output)
		return nil
	}
}
