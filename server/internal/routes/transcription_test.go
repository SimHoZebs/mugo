package routes_test

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/simhozebs/mugo/internal/routes"
	"github.com/stretchr/testify/assert"
)

func TestTranscribeAudio(t *testing.T) {
	whisperServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/asr", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		err := r.ParseMultipartForm(2 << 20)
		assert.NoError(t, err)

		file, _, err := r.FormFile("audio_file")
		assert.NoError(t, err)
		assert.NotNil(t, file)
		_ = file.Close()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"text":"hello from whisper"}`))
	}))
	defer whisperServer.Close()

	t.Setenv("TRANSCRIPTION_SERVER_URL", whisperServer.URL)

	_, api := humatest.New(t)
	routes.RegisterTranscriptionEndpoints(api, "/transcription")

	resp := api.Post("/transcription", struct {
		AudioBase64 string `json:"audio_base64"`
		FileName    string `json:"file_name"`
		Language    string `json:"language"`
	}{
		AudioBase64: base64.StdEncoding.EncodeToString([]byte("test-audio")),
		FileName:    "sample.wav",
		Language:    "en",
	})

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "hello from whisper")
}

func TestTranscribeAudio_InvalidBase64(t *testing.T) {
	_, api := humatest.New(t)
	routes.RegisterTranscriptionEndpoints(api, "/transcription")

	resp := api.Post("/transcription", struct {
		AudioBase64 string `json:"audio_base64"`
	}{
		AudioBase64: "%%%",
	})

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestTranscribeAudio_ServiceError(t *testing.T) {
	whisperServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "service failed", http.StatusInternalServerError)
	}))
	defer whisperServer.Close()

	t.Setenv("TRANSCRIPTION_SERVER_URL", whisperServer.URL)

	_, api := humatest.New(t)
	routes.RegisterTranscriptionEndpoints(api, "/transcription")

	resp := api.Post("/transcription", struct {
		AudioBase64 string `json:"audio_base64"`
	}{
		AudioBase64: base64.StdEncoding.EncodeToString([]byte("test-audio")),
	})

	assert.Equal(t, http.StatusBadGateway, resp.Code)
}
