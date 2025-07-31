package main

// BrowserResult представляет результат работы с браузером
type BrowserResult struct {
	Success   bool   `json:"success"`
	URL       string `json:"url"`
	Username  string `json:"username"`
	Message   string `json:"message"`
	Error     string `json:"error,omitempty"`
	Timestamp int64  `json:"timestamp"`
}
