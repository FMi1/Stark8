package api

import (
	"errors"
	"log"
	"net/http"
	"stark8/internal/k8sclient"
	"stark8/internal/proxy"
	"stark8/internal/store"
	"stark8/internal/token"
	"stark8/internal/utils"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	router       *gin.Engine
	proxyRouter  *gin.Engine
	proxyHub     *proxy.ProxyHub
	TokenMaker   token.Maker
	config       utils.Config
	k8sClientset *k8sclient.Client
	db           store.Database
}

func NewServer(config utils.Config) (*Server, error) {

	// gin.SetMode(gin.ReleaseMode)
	tokenMaker, err := token.NewPasetoMaker(config.TokenSimmetricKey)
	if err != nil {
		return nil, errors.New("cannot create token maker")
	}
	// Create a new Server struct and initialize the proxyHub field with a new instance
	// of the ProxyHub struct.
	server := &Server{
		proxyHub:   proxy.NewProxyHub(config.Hostname),
		TokenMaker: tokenMaker,
		config:     config,
	}

	db, err := store.NewBbolt()
	if err != nil {
		return nil, err
	}
	server.db = db

	kubernetsClient, err := k8sclient.NewClient()
	if err != nil {
		return nil, err
	}

	server.k8sClientset = kubernetsClient

	// Create a new router for handling requests that are not for any subdomain.
	appRouter := gin.Default()
	// Create a new router for handling requests that are for any subdomain.
	proxyRouter := gin.Default()

	// Set up the validator.
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("customPassword", utils.CustomPasswordValidator)
	}

	// Set up the proxy router to handle all requests (i.e. requests for any subdomain)
	// by calling the proxyHandler method of the Server struct.
	appRouter.Use(func(c *gin.Context) {
		// Get the host from the request.
		host := c.Request.Host
		// If the host contains the provided hostname, pass the request to the proxyRouter.
		if strings.Contains(host, "."+config.Hostname) {
			proxyRouter.HandleContext(c)
			c.Abort() // Stop the request from being handled by any other middleware or handlers.
			return
		}

		// If the host matches the provided hostname, pass the request to the next handler.
		if host == config.Hostname {
			c.Next()
			return
		}

		// If the host does not contain the provided hostname, return a 404 error.
		c.JSON(http.StatusNotFound, gin.H{"Error": "Error hostname not found"})
		c.Abort()
	})

	authRouter := appRouter.Group("/")

	authRouter.Use(authMiddleware(tokenMaker, config.Hostname))
	proxyRouter.Use(authMiddleware(tokenMaker, config.Hostname))

	appRouter.POST("/users/login", server.loginUserRequest)
	appRouter.GET("/users/login", server.getLoginUserRequest)

	appRouter.GET("/users/sign_up", server.getSignUpUserRequest)
	appRouter.POST("/users/sign_up", server.createUserRequest)

	appRouter.GET("/users/logout", server.logoutUserRequest)

	appRouter.Static("/static", "./static")
	// Set up the app router to handle requests for /services by calling the getServicesRequest
	// and createServiceRequest methods of the Server struct.
	authRouter.GET("/", server.getHomeRequest)

	authRouter.GET("/namespaces", server.getModalNamespacesRequest)
	authRouter.GET("/namespaces/:namespace/services", server.getModalServicesRequest)
	authRouter.GET("/namespaces/:namespace/services/:service", server.getModalSettingsRequest)
	authRouter.POST("/namespaces/:namespace/services/:service", server.createStark8Request)

	authRouter.GET("/new", server.getModalStruct)
	authRouter.POST("/new", server.createStark8Request)

	authRouter.POST("/logos", server.getLogosRequest)

	authRouter.GET("/stark8s", server.getStark8sRequest)

	proxyRouter.Any("/*any", server.proxyHandler)

	// Set the proxy router field of the Server struct to the newly created proxy router.
	server.proxyRouter = proxyRouter

	// Set the router field of the Server struct to the newly created app router.
	server.router = appRouter

	// Return the newly created Server struct.
	return server, nil
}

func (s *Server) proxyHandler(c *gin.Context) {
	subdomain := strings.SplitN(c.Request.Host, ".", 2)[0]
	proxy, err := s.proxyHub.GetProxy(subdomain)
	if err != nil {
		log.Printf("Proxy not found for service: %s", subdomain)
		c.JSON(http.StatusNotFound, gin.H{"Error": "Proxy Not Found"})
		return
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

func (s *Server) Start(cert string, key string) {

	s.router.RunTLS(":8443", cert, key)
}
