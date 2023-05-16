package web

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/joao-macieira/fclx/chatservice/internal/usecases/chatcompletion"
)

type WebChatGPTHandler struct {
	CompletionUseCase chatcompletion.ChatCompletionUseCase
	Config            chatcompletion.ChatCompletionConfigInputDTO
	AuthToken         string
}

func NewWebChatGPTHandler(usecase chatcompletion.ChatCompletionUseCase, config chatcompletion.ChatCompletionConfigInputDTO, authToken string) *WebChatGPTHandler {
	return &WebChatGPTHandler{
		CompletionUseCase: usecase,
		Config:            config,
		AuthToken:         authToken,
	}
}

func (handler *WebChatGPTHandler) Handle(web http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		web.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if request.Header.Get("Authorization") != handler.AuthToken {
		web.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(request.Body)

	if err != nil {
		http.Error(web, err.Error(), http.StatusInternalServerError)
		return
	}

	if !json.Valid(body) {
		http.Error(web, "Invalid json", http.StatusBadRequest)
		return
	}

	var dto chatcompletion.ChatCompletionInputDTO

	err = json.Unmarshal(body, &dto)

	if err != nil {
		http.Error(web, err.Error(), http.StatusBadRequest)
		return
	}

	dto.Config = handler.Config

	result, err := handler.CompletionUseCase.Execute(request.Context(), dto)

	if err != nil {
		http.Error(web, err.Error(), http.StatusInternalServerError)
		return
	}

	web.WriteHeader(http.StatusOK)
	web.Header().Set("Content-Type", "application/json")
	json.NewEncoder(web).Encode(result)
}
