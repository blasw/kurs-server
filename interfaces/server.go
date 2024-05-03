package interfaces

type Server interface {
	Run(string)
	SetupRoutes()
}
