package output

type StreamAnswer struct {
	TextChan <-chan string
	ErrChan  <-chan error
}
