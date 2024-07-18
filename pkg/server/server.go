package server

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"strings"

	"cassette/config"
	"cassette/pkg/repository"
	"cassette/pkg/storage"

	"github.com/rs/cors"
)

const (
        cookieName = "cassette-session"
        webhookURL = "https://mentisnativo.bubbleapps.io/version-test/api/1.1/wf/session-data/initialize"
)

var (
        //go:embed cassette.js
        jsCassette string

        //go:embed record.umd.min.cjs
        jsRecord string
)

type Server struct {
        *config.Config

        handler    http.Handler
        filesystem fs.FS
        Repository repository.Repository // Add repository field
        Storage    storage.Storage       // Add storage field
}

func New(config *config.Config, repo repository.Repository, storage storage.Storage) *Server {
        mux := http.NewServeMux()

        cors := cors.New(cors.Options{
                AllowOriginFunc: func(origin string) bool {
                        return true
                },
                AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
                AllowedHeaders:   []string{"*"},
                AllowCredentials: true,
        })

        s := &Server{
                Config:     config,
                handler:    cors.Handler(mux),
                filesystem: os.DirFS("./public"),
                Repository: repo,
                Storage:    storage,
        }

        mux.HandleFunc("/events", s.handleEvents)
        mux.HandleFunc("/cassette.min.cjs", s.handleScript)
		mux.HandleFunc("/agoraStream/{session}", s.handleAgoraStreamUrl)


        mux.HandleFunc("/sessions", s.handleAuth(s.handleSessions))
        mux.HandleFunc("/sessions/", s.handleAuth(s.handleSession)) // Added trailing slash
        mux.HandleFunc("/sessions/{session}/events", s.handleAuth(s.handleSessionEvents))
        mux.HandleFunc("/sessions/{session}/delete", s.handleAuth(s.handleSessionDelete)) // Unique path for delete

        mux.HandleFunc("/", s.handleAuth(s.handleUI))

        return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        s.handler.ServeHTTP(w, r)
}

func (s *Server) handleAuth(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
                if s.Username == "" || s.Password == "" {
                        next(w, r)
                        return
                }

                if username, password, ok := r.BasicAuth(); ok {
                        if username == s.Username && password == s.Password {
                                next(w, r)
                                return
                        }
                }

                w.Header().Set("WWW-Authenticate", `Basic realm="Cassette - Admin UI", charset="UTF-8"`)
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
        }
}

func (s *Server) handleUI(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/help" || r.URL.Path == "/help/" {
                http.ServeFileFS(w, r, s.filesystem, "index.html")
                return
        }

        handler := http.FileServerFS(s.filesystem)
        handler.ServeHTTP(w, r)
}

func (s *Server) handleScript(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/javascript")

        var result bytes.Buffer

        result.WriteString(jsRecord)
        result.WriteString("\n")
        result.WriteString(jsCassette)

        w.Write(result.Bytes())
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
        var body struct {
                Events      []storage.Event `json:"events"`
                UserEmail   string          `json:"userEmail"`
                QaId        string          `json:"qaId"`
                QaSessionId string          `json:"qaSessionId"`
				AgoraStreamUrl string       `json:"agoraStreamUrl"`
        }

        if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
        }

        var err error
        var session *repository.Session

        // Check if a session exists for the given qaSessionId
        sessions, err := s.Repository.FindSessionsByQaSessionId(body.QaSessionId)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        if len(sessions) > 0 {
                session = &sessions[0]
        } else {
                info := &repository.SessionInfo{
                        Origin:      getOrigin(r),
                        Address:     getAddress(r),
                        UserAgent:   r.UserAgent(),
                        UserEmail:   body.UserEmail,
                        QaId:        body.QaId,
                        QaSessionId: body.QaSessionId,
						AgoraStreamUrl: body.AgoraStreamUrl,
                }

                session, err = s.Repository.CreateSession(info)
                if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }
        }

        if err := s.Storage.AppendEvents(session.ID, body.Events...); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        setSessionID(w, r, session.ID)

        if err := s.sendSessionIDToWebhook(session.ID); err != nil {
                fmt.Printf("Error sending session ID to webhook: %v\n", err)
        }

        w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleAgoraStreamUrl(w http.ResponseWriter, r *http.Request) {
    // Extract session ID from URL path
    id := r.URL.Path[len("/agoraStream/"):]

    // Extract session ID correctly
    if idx := strings.Index(id, "/"); idx != -1 {
        id = id[:idx]
    }

    // Decode JSON request body
    var body struct {
        AgoraStreamUrl string `json:"agoraStreamUrl"`
    }
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Retrieve session from repository
    session, err := s.Repository.Session(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    // Update AgoraStreamUrl field
    session.AgoraStreamUrl = body.AgoraStreamUrl

    // Update session in repository
    if err := s.Repository.UpdateSessionAgoraStreamURL(id, body.AgoraStreamUrl); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (s *Server) sendSessionIDToWebhook(sessionID string) error {
        jsonData := map[string]string{"session_id": sessionID}
        jsonValue, err := json.Marshal(jsonData)
        if err != nil {
                return fmt.Errorf("error marshalling JSON: %v", err)
        }

        resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonValue))
        if err != nil {
                return fmt.Errorf("error sending request to webhook: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
                return fmt.Errorf("webhook returned non-OK status: %v", resp.Status)
        }

        return nil
}

func (s *Server) handleSessions(w http.ResponseWriter, r *http.Request) {
        sessions, err := s.Repository.Sessions()
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        json.NewEncoder(w).Encode(sessions)
}

func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
        id := r.URL.Path[len("/sessions/"):]
        session, err := s.Repository.Session(id)
        if err != nil {
                http.Error(w, err.Error(), http.StatusNotFound)
                return
        }

        json.NewEncoder(w).Encode(session)
}

func (s *Server) handleSessionDelete(w http.ResponseWriter, r *http.Request) {
        id := r.URL.Path[len("/sessions/"):]

        // Extract session ID correctly
        if idx := strings.Index(id, "/"); idx != -1 {
                id = id[:idx]
        }

        if err := s.Storage.DeleteSession(id); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        if err := s.Repository.DeleteSession(id); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleSessionEvents(w http.ResponseWriter, r *http.Request) {
        id := r.URL.Path[len("/sessions/"):]

        // Extract session ID correctly
        if idx := strings.Index(id, "/"); idx != -1 {
                id = id[:idx]
        }

        events, err := s.Storage.Events(id)
        if err != nil {
                http.Error(w, err.Error(), http.StatusNotFound)
                return
        }

        json.NewEncoder(w).Encode(events)
}

func getOrigin(r *http.Request) string {
        if val := r.Header.Get("Origin"); val != "" {
                return val
        }
        if val := r.Header.Get("Referer"); val != "" {
                return val
        }
        return ""
}

func getAddress(r *http.Request) string {
        host, _, _ := net.SplitHostPort(r.RemoteAddr)
        return host
}

func getSessionID(r *http.Request) string {
        cookie, _ := r.Cookie(cookieName)
        if cookie != nil && cookie.Value != "" {
                return cookie.Value
        }
        return ""
}



func setSessionID(w http.ResponseWriter, r *http.Request, id string) {
        http.SetCookie(w, &http.Cookie{
                Name:  cookieName,
                Value: id,
                Path:  "/",
                Secure:   r.URL.Scheme == "https",
                HttpOnly: true,
        })
}