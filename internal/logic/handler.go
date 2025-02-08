package logic

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/hashicorp/vault/shamir"
	"github.com/hoisie/mustache"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mclyashko/everlock/internal/crypt"
	"github.com/mclyashko/everlock/internal/db"
	"github.com/mclyashko/everlock/internal/model"
)

const (
	maxNicknameSizeAllowed    = 16
	maxMessageSizeAllowed     = 1000
	minKeyholdersAllowed      = 2
	maxKeyholdersAllowed      = 256
	aesKeySize                = 32
	maxEncryptedMessageLength = 1024
	messageIDKey              = "messageID"
	nicknameKey               = "nickname"
	messageKey                = "message"
	keyholdersKey             = "keyholders"
	minKeyholdersKey          = "minKeyholders"
)

// предоставляет доступ к шаблону главной страницы
func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logErrorAndRespond(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	template, err := mustache.ParseFile(filepath.Join("..", "..", "internal", "template", "index.html"))
	if err != nil {
		logErrorAndRespond(w, fmt.Sprintf("error loading main page template, error: %v", err), http.StatusInternalServerError)
		return
	}

	renderedTemplate := template.Render()
	w.Header().Set("Content-Type", "text/html")
	if _, err = w.Write([]byte(renderedTemplate)); err != nil {
		logErrorAndRespond(w, fmt.Sprintf("error writing main page response, error: %v", err), http.StatusInternalServerError)
	}
}

// обрабабатывает форму добавления сообщения
func SubmitMessageHandler(p *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logErrorAndRespond(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	nickname, messageText, keyholders, minKeyholders, err := parseFormData(r)
	if err != nil {
		logErrorAndRespond(w, err.Error(), http.StatusBadRequest)
		return
	}

	aesKey, err := crypt.GenerateRandomAESKey(aesKeySize)
	if err != nil {
		logErrorAndRespond(w, fmt.Sprintf("failed to generate encryption key, error: %v", err), http.StatusInternalServerError)
		return
	}

	encryptedMessage, err := crypt.EncryptAES([]byte(messageText), aesKey)
	if err != nil {
		logErrorAndRespond(w, fmt.Sprintf("failed to encrypt message, error: %v", err), http.StatusInternalServerError)
		return
	}
	if len(encryptedMessage) > maxEncryptedMessageLength {
		logErrorAndRespond(w, fmt.Sprintf("encryptedMessage length is to large: %d, max : %d", len(encryptedMessage), maxEncryptedMessageLength), http.StatusBadRequest)
		return
	}

	keyHash := sha256.Sum256(aesKey)

	keyShares, err := shamir.Split(aesKey, keyholders, minKeyholders)
	if err != nil {
		logErrorAndRespond(w, "failed to split secret key", http.StatusBadRequest)
		return
	}

	message := model.NewMessage(nickname, encryptedMessage, keyHash, keyholders, minKeyholders, keyShares)

	if err = db.SaveNewMessage(p, message); err != nil {
		logErrorAndRespond(w, fmt.Sprintf("transaction commit failed, error: %v", err), http.StatusInternalServerError)
		return
	}

	template, err := mustache.ParseFile(filepath.Join("..", "..", "internal", "template", "message_success.html"))
	if err != nil {
		logErrorAndRespond(w, fmt.Sprintf("error loading success template, error: %v", err), http.StatusInternalServerError)
		return
	}

	renderedTemplate := template.Render(map[string]interface{}{
		messageIDKey: message.ID,
	})

	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write([]byte(renderedTemplate))
	if err != nil {
		logErrorAndRespond(w, fmt.Sprintf("error writing response, error: %v", err), http.StatusInternalServerError)
	}
}

func logErrorAndRespond(w http.ResponseWriter, errorMessage string, statusCode int) {
	log.Printf("%s", errorMessage)
	http.Error(w, "Internal Server Error", statusCode)
}

func parseFormData(r *http.Request) (string, string, int, int, error) {
	if err := r.ParseForm(); err != nil {
		return "", "", 0, 0, fmt.Errorf("error parsing form data")
	}

	nickname := r.FormValue(nicknameKey)
	messageText := r.FormValue(messageKey)
	keyholders := r.FormValue(keyholdersKey)
	minKeyholders := r.FormValue(minKeyholdersKey)

	if len(nickname) > maxNicknameSizeAllowed {
		return "", "", 0, 0, fmt.Errorf("nickname is too large: %s", nickname)
	}

	if len(messageText) > maxMessageSizeAllowed {
		return "", "", 0, 0, fmt.Errorf("message is too large, message size: %d", len(messageText))
	}

	keyholdersInt, err := strconv.Atoi(keyholders)
	if err != nil || keyholdersInt < minKeyholdersAllowed || keyholdersInt > maxKeyholdersAllowed {
		return "", "", 0, 0, fmt.Errorf("invalid number of keyholders: %s, error: %v", keyholders, err)
	}

	minKeyholdersInt, err := strconv.Atoi(minKeyholders)
	if err != nil || minKeyholdersInt < minKeyholdersAllowed || minKeyholdersInt > maxKeyholdersAllowed || minKeyholdersInt > keyholdersInt {
		return "", "", 0, 0, fmt.Errorf("invalid minimum keyholders: %s, error: %v", minKeyholders, err)
	}

	return nickname, messageText, keyholdersInt, minKeyholdersInt, nil
}
