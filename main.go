package main

import (
	"crypto/tls"
	"io"
	"log/slog"
	"os"

	"net/http"
	"net/url"

	"strings"
)

func main() {
	logopts := &slog.HandlerOptions{Level: slog.LevelInfo}
	loghandler := slog.NewTextHandler(os.Stdout, logopts)
	logger := slog.New(loghandler)
	slog.SetDefault(logger)

	remoteRaw := os.Getenv("REMOTE_URL")
	remote, err := url.Parse(remoteRaw)
	if err != nil {
		slog.Error(err.Error(), "remote", remoteRaw)
		os.Exit(1)
	}

	proxyHandler := &ProxyHandler{remote: remote}

	slog.Info("Starting proxy")
	err = http.ListenAndServe(":8080", proxyHandler)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

type ProxyHandler struct {
	remote *url.URL
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(
		r.Context(),
		"Request",
		slog.String("http.method", r.Method),
		slog.String("http.url", r.URL.String()),
	)
	r.URL.Scheme = h.remote.Scheme
	r.URL.Host = h.remote.Host
	r.URL.Path = h.remote.Path + r.URL.Path

	rule_ApiProjectsFix(r)

	proxy := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	resp, err := proxy.RoundTrip(r)
	if err != nil {
		slog.ErrorContext(
			r.Context(),
			err.Error(),
			slog.String("http.method", r.Method),
			slog.String("http.url", r.URL.String()),
		)
		http.Error(w, "Server Error", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func rule_ApiProjectsFix(r *http.Request) {
	prefix := "/api/v4/projects/"

	if strings.HasPrefix(r.URL.Path, prefix) {
		projectName := r.URL.Path[len(prefix):]

		oldPath := r.URL.Path
		newPath := "/api/v4/projects/" + url.PathEscape(projectName)

		slog.InfoContext(
			r.Context(),
			"Rewriting API Projects URL",
			"rule.api-projects-fix.old-path", oldPath,
			"rule.api-projects-fix.new-path", newPath,
		)

		r.URL.Path = newPath
	}
}
