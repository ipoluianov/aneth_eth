package common

type ViewProvider interface {
	GetResult(code string) *Result
}
