package http

// nolint: revive
import (
	"context"
	"fmt"
	"github.com/Rasikrr/my_project/configs"
	_ "github.com/Rasikrr/my_project/docs"
	mcqTest "github.com/Rasikrr/my_project/internal/ports/http/handlers/mcq_test"
	"github.com/Rasikrr/my_project/internal/ports/http/middlewares"
	mcqS "github.com/Rasikrr/my_project/internal/services/mcq"
	"github.com/pkg/errors"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	name         = "http server"
	idleTimeout  = time.Minute * 3
	readTimeout  = time.Second * 60
	writeTimeout = time.Second * 60
	filesDir     = "files"
)

// @title Learning Platform API
// @version 1.0
// @description This is docs for Learning Platform API
// @termsOfService http://swagger.io/terms/

// @contact.name Pick me team
// @contact.url https://github.com/Rasikrr/my_project
// @contact.email leaning_platform@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /

// Server @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
type Server struct {
	name string
	port string
	host string
	srv  *http.Server
}

func NewServer(
	cfg *configs.Config,
	mcqService mcqS.Service,
) *Server {
	router := http.NewServeMux()

	// Middlewares
	authMiddleware := middlewares.NewAuthMiddleware()

	// Controllers
	answerController := mcqTest.NewController(mcqService, authMiddleware)
	answerController.Init(router)

	// CORS
	r := middlewares.CORSMiddleware(router)

	srv := &http.Server{
		Addr:         address(cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	s := &Server{
		name: name,
		port: cfg.Server.Port,
		host: cfg.Server.Host,
		srv:  srv,
	}
	s.initSwagger(router)
	if err := s.initStatic(router); err != nil {
		log.Fatal(err)
	}
	return s
}

func address(host, port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}

func (s *Server) initSwagger(router *http.ServeMux) {
	addr := hostPort(s.host, s.port)
	url := fmt.Sprintf("http://%s/swagger/doc.json", addr)
	router.Handle("/swagger/",
		httpSwagger.Handler(httpSwagger.URL(url)),
	)
}

// nolint: unparam
func (s *Server) initStatic(router *http.ServeMux) error {
	router.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./files"))))
	router.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./files/avatars"))))
	router.Handle("/courses/", http.StripPrefix("/courses/", http.FileServer(http.Dir("./files/courses"))))
	router.Handle("/faq/", http.StripPrefix("/faq/", http.FileServer(http.Dir("./files/faq"))))
	return nil
}

func (s *Server) Start() error {
	if err := s.createFilesDir(); err != nil {
		return errors.Wrap(err, "create files dir")
	}
	log.Println("starting http server")
	log.Printf("swagger url: http://%s/swagger/index.html\n", hostPort(s.host, s.port))
	log.Printf("Server running on http://%s", hostPort(s.host, s.port))
	return s.srv.ListenAndServe()
}

func (s *Server) createFilesDir() error {
	err := os.MkdirAll(filesDir, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func hostPort(host, port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}
