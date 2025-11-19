package check

import "context"

type Usecase struct {
	version string
}

func NewUsecase(
	version string,
) *Usecase {
	return &Usecase{
		version: version,
	}
}

func (h *Usecase) Execute(ctx context.Context, input Input) (Output, error) {
	return Output{
		Status:  "ok",
		Version: h.version,
	}, nil
}
