package restfulServer

// RouteSetup represents a function that sets up routes on the server
type RouteSetup func(s *Server)

// SetupAllRoutes configures all application routes in an organized manner
func SetupAllRoutes(s *Server, routeSetups ...RouteSetup) {
	s.log.Info("Setting up all application routes...")

	// Execute all route setups
	for _, setupFunc := range routeSetups {
		setupFunc(s)
	}

	s.log.Info("All application routes configured successfully")
}
