package routes

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/simhozebs/mugo/internal/api"
	"github.com/simhozebs/mugo/internal/config"
)

func RegisterTranscriptionEndpoints(humaAPI huma.API, prefix string) {
	transcriptionGroup := huma.NewGroup(humaAPI, prefix)

	huma.Register(transcriptionGroup, huma.Operation{
		OperationID: "transcribe-audio",
		Method:      http.MethodPost,
		Path:        "",
		Summary:     "Transcribe audio",
		Description: "Transcribes base64-encoded audio using external whisper service",
		Tags:        []string{"Transcription"},
	}, func(ctx context.Context, input *api.TranscriptionRequest) (*api.TranscriptionResponse, error) {
		if input.Body.AudioBase64 == "" {
			return nil, huma.Error400BadRequest("audio_base64 is required")
		}

		audioBytes, err := base64.StdEncoding.DecodeString(input.Body.AudioBase64)
		if err != nil {
			return nil, huma.Error400BadRequest("audio_base64 must be valid base64")
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		fileName := input.Body.FileName
		if fileName == "" {
			fileName = "audio.wav"
		}

		fileWriter, err := writer.CreateFormFile("audio_file", filepath.Base(fileName))
		if err != nil {
			return nil, fmt.Errorf("failed to build transcription request: %w", err)
		}

		if _, err := fileWriter.Write(audioBytes); err != nil {
			return nil, fmt.Errorf("failed to attach audio to transcription request: %w", err)
		}

		if input.Body.Language != "" {
			if err := writer.WriteField("language", input.Body.Language); err != nil {
				return nil, fmt.Errorf("failed to set language hint: %w", err)
			}
		}

		if err := writer.WriteField("task", "transcribe"); err != nil {
			return nil, fmt.Errorf("failed to set transcription task: %w", err)
		}

		if err := writer.Close(); err != nil {
			return nil, fmt.Errorf("failed to finalize transcription request: %w", err)
		}

		transcriptionURL := strings.TrimRight(config.GetTranscriptionServerURL(), "/") + "/asr"
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, transcriptionURL, &body)
		if err != nil {
			return nil, fmt.Errorf("failed to create transcription request: %w", err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, huma.Error502BadGateway("Transcription service is unavailable", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
			return nil, huma.Error502BadGateway("Transcription service returned an error", fmt.Errorf("status %d: %s", resp.StatusCode, strings.TrimSpace(string(errBody))))
		}

		var transcriptionResult struct {
			Text string `json:"text"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&transcriptionResult); err != nil {
			return nil, huma.Error502BadGateway("Failed to decode transcription response", err)
		}

		output := &api.TranscriptionResponse{}
		output.Body.Text = transcriptionResult.Text
		return output, nil
	})
}
