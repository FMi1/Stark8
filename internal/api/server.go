package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"stark8/internal/k8sclient"
	"stark8/internal/proxy"
	"stark8/internal/token"
	"stark8/internal/utils"

	"strings"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router       *gin.Engine
	proxyRouter  *gin.Engine
	proxyHub     *proxy.ProxyHub
	TokenMaker   token.Maker
	config       utils.Config
	k8sClientset *k8sclient.Client
}

// NewServer creates a new instance of the Server struct.
// It initializes the proxyHub field with a new instance of the ProxyHub struct.
// It then creates two new routers (appRouter and proxyRouter) and sets up a middleware
// function to handle requests for any subdomain. If the request host contains the
// provided hostname, the request is passed to the proxyRouter. If the request host
// matches the provided hostname, the request is passed to the next handler. If the
// request host does not contain the provided hostname, a 404 error is returned.
func NewServer(config utils.Config) (*Server, error) {
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

	kubernetsClient, err := k8sclient.NewClient()
	if err != nil {
		return nil, err
	}

	server.k8sClientset = kubernetsClient
	// Create a new router for handling requests that are not for any subdomain.
	appRouter := gin.Default()

	// Create a new router for handling requests that are for any subdomain.
	proxyRouter := gin.Default()
	// appRouter.Use(gin.Recovery())

	// appRouter.Use(authMiddleware(tokenMaker))

	// Set up a middleware function to handle requests for any subdomain.

	// Set up the proxy router to handle all requests (i.e. requests for any subdomain)
	// by calling the proxyHandler method of the Server struct.

	appRouter.Use(func(c *gin.Context) {
		// Get the host from the request.
		host := c.Request.Host
		fmt.Println("host request: ", host)
		fmt.Println("host config: ", config.Hostname)
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
	})

	authRouter := appRouter.Group("/")

	authRouter.Use(authMiddleware(tokenMaker))
	proxyRouter.Use(authMiddleware(tokenMaker))

	appRouter.POST("/users/login", server.loginUserRequest)
	appRouter.GET("/users/login", server.getLoginUserRequest)
	appRouter.Static("/static", "./static")
	// Set up the app router to handle requests for /services by calling the getServicesRequest
	// and createServiceRequest methods of the Server struct.
	authRouter.GET("/", server.getHomeRequest)
	authRouter.POST("/new", server.createStark8Request)
	authRouter.GET("/namespaces/:namespace/services", server.getModalServicesRequest)
	authRouter.GET("/namespaces/:namespace/services/:service", server.getModalSettingsRequest)

	authRouter.GET("/new", server.getModalStruct)
	authRouter.GET("/namespaces", server.getModalNamespacesRequest)

	authRouter.POST("/namespaces/:namespace/services/:service", server.createStark8Request)

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

func (s *Server) Start() {
	s.router.Run(":8080")
}
