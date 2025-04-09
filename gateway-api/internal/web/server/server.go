package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/service"
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/web"
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/web/handlers"
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/web/middleware"
)

type Server struct {
	router         *chi.Mux
	server         *http.Server
	accountService *service.AccountService
	invoiceService *service.InvoiceService
	port           string
}

func NewServer(accountService *service.AccountService, invoiceService *service.InvoiceService, port string) *Server {
	return &Server{
		router:         chi.NewRouter(),
		accountService: accountService,
		invoiceService: invoiceService,
		port:           port,
	}
}

func (s *Server) ConfigureRoutes() {
	accountHandler := handlers.NewAccountHandler(s.accountService)
	invoiceHandler := handlers.NewInvoiceHandler(s.invoiceService)
	authMiddleware := middleware.NewAuthMiddleware(s.accountService)

	s.router.Post(web.ROUTER_ACCOUNTS, accountHandler.Create)
	s.router.Get(web.ROUTER_ACCOUNTS, accountHandler.Get)

	s.router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Auth)
		s.router.Post(web.ROUTER_INVOICES, invoiceHandler.Create)
		s.router.Get(web.ROUTER_INVOICES_ID, invoiceHandler.GetByID)
		s.router.Get(web.ROUTER_INVOICES, invoiceHandler.ListByAccount)
	})
}

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}

	return s.server.ListenAndServe()
}
