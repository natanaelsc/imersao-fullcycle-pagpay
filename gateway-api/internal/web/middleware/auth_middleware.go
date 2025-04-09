package middleware

import (
	"net/http"

	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/domain"
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/service"
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/web"
)

type AuthMiddleware struct {
	accountService *service.AccountService
}

func NewAuthMiddleware(accountService *service.AccountService) *AuthMiddleware {
	return &AuthMiddleware{
		accountService: accountService,
	}
}

func (m *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(web.HEADER_X_API_KEY)
		if apiKey == "" {
			http.Error(w, domain.ErrRequiredAPIKey.Error(), http.StatusUnauthorized)
			return
		}

		_, err := m.accountService.FindByAPIKey(apiKey)
		if err != nil {
			if err == domain.ErrAccountNotFound {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
