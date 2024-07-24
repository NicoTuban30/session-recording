package main

import (
	"cassette/config"
	"cassette/pkg/server"
	"cassette/pkg/storage/filesystem"
	"io"
	"log"
	"net/http"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    session := r.FormValue("session")
    videoFile, _, err := r.FormFile("video")
    if err != nil {
        http.Error(w, "Failed to upload video", http.StatusInternalServerError)
        return
    }
    defer videoFile.Close()

    videoData, err := io.ReadAll(videoFile)
    if err != nil {
        http.Error(w, "Failed to read video", http.StatusInternalServerError)
        return
    }

    fs, err := filesystem.New("sessions")
    if err != nil {
        http.Error(w, "Failed to initialize filesystem storage", http.StatusInternalServerError)
        return
    }

    err = fs.SaveVideo(session, videoData)
    if err != nil {
        http.Error(w, "Failed to save video", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}
func main() {
	http.HandleFunc("/upload", uploadHandler)
	cfg, err := config.FromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// Create the server with the config, repository, and storage
	s := server.New(cfg, cfg.Repository, cfg.Storage)

	// Start the HTTP server
	if err := http.ListenAndServe(":3000", s); err != nil {
		log.Fatal(err)
	}
}
