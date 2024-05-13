package main

import (
	"log/slog"
	"os"

	"net/http"
	"net/http/httputil"
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

	proxyHandler := &ProxyHandler{
		proxy:  httputil.NewSingleHostReverseProxy(remote),
		remote: remote,
	}

	slog.Info("Starting proxy")
	err = http.ListenAndServe(":8080", proxyHandler)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

type ProxyHandler struct {
	proxy  *httputil.ReverseProxy
	remote *url.URL
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(
		r.Context(),
		"Request",
		slog.String("http.method", r.Method),
		slog.String("http.url", r.URL.String()),
	)
	r.Host = h.remote.Host

	rule_ApiProjectsFix(r)

	h.proxy.ServeHTTP(w, r)
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
