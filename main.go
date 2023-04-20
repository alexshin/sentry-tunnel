package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type DsnBody struct {
	Dsn string `json:"dsn"`
}

type ErrorMsg struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func getError(msg string) []byte {
	m, _ := json.Marshal(&ErrorMsg{Error: true, Message: msg})
	log.Println(msg)
	return m
}

func getRoot(host string, projectIds []string, sentrySchema string) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(405)
			w.Write(getError(fmt.Sprint("Method", r.Method, "is not supported")))
			return
		}

		body, err := io.ReadAll(r.Body)
		var dsn DsnBody
		if err != nil || len(body) == 0 {
			w.WriteHeader(400)
			w.Write(getError(fmt.Sprint("Cannot read request body")))
			return
		}

		fLine := strings.Split(string(body), "\n")[0]

		err = json.Unmarshal([]byte(fLine), &dsn)
		if err != nil {
			w.WriteHeader(400)
			w.Write(getError(fmt.Sprint("Cannot deserialize body")))
			return
		}

		u, err := url.Parse(dsn.Dsn)
		if err != nil {
			w.WriteHeader(400)
			w.Write(getError(fmt.Sprint("DSN contains wrong URL format")))
			return
		}

		if u.Hostname() != host {
			w.WriteHeader(400)
			w.Write(getError(fmt.Sprint("Hostname is not allowed:", u.Hostname())))
			return
			// panic("Invalid Sentry Host")
		}

		projectId := strings.TrimSuffix(u.Path, "/")
		projectId = strings.TrimPrefix(projectId, "/")
		if !contains(projectIds, projectId) {
			w.WriteHeader(400)
			w.Write(getError(fmt.Sprint("Project is not allowed:", projectId)))
			return
		}

		reqURL := fmt.Sprintf("%s://%s/api/%s/envelope/", sentrySchema, host, projectId)
		res, err := http.Post(reqURL, "application/x-sentry-envelope", strings.NewReader(string(body)))
		if err != nil {
			w.WriteHeader(500)
			w.Write(getError(fmt.Sprintf("Error occurred: %s", err)))
			return
		}

		w.WriteHeader(res.StatusCode)
		resBody, _ := io.ReadAll(res.Body)
		w.Write(resBody)
	}
}

func getHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func main() {
	sentryHost := os.Getenv("SENTRY_HOST")
	sentrySchema := os.Getenv("SENTRY_SCHEMA")
	if sentrySchema != "http" && sentrySchema != "https" {
		sentrySchema = "https"
	}
	projectIds := strings.Split(os.Getenv("SENTRY_PROJECT_IDS"), ",")
	routePath := os.Getenv("APP_ROUTE_PATH")
	if routePath == "" {
		routePath = "/bugs"
	}

	appHost := os.Getenv("APP_HOST")
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "3333"
	}
	addr := fmt.Sprintf("%s:%s", appHost, appPort)

	if sentryHost == "" || len(projectIds) == 0 {
		log.Fatal("Env variables SENTRY_HOST and SENTRY_PROJECT_IDS are required")
	}

	http.HandleFunc(routePath, getRoot(sentryHost, projectIds, sentrySchema))
	http.HandleFunc("/health-check", getHealthcheck)

	log.Println("Listening on ", addr)
	log.Println("- SENTRY_HOST:", sentryHost)
	log.Println("- SENTRY_SCHEMA:", sentrySchema)
	log.Println("- SENTRY_PROJECT_IDS:", projectIds)
	log.Println("- APP_ROUTE_PATH:", routePath)
	log.Println("-----------------")

	err := http.ListenAndServe(addr, nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Fatalln("server closed")
	} else if err != nil {
		log.Fatalln("error starting server: ", err)
	}
}
