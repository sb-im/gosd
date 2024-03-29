package app

// @title Gosd RESTful API
// @version 3.0.0
// @description superdock web service restful api

// @host localhost:8000
// @BasePath /gosd/api/v3
// @query.collection.format multi

// @securitydefinitions.basic BasicAuth

// @securitydefinitions.apiKey JWTSecret
// @in header
// @name Authorization

// @securitydefinitions.apiKey APIKeyHeader
// @in header
// @name X-API-Key

// @tag.name status
// @tag.description Server Time, Running Status
// @tag.name auth
// @tag.description auth, login, logout, token
// @tag.name team
// @tag.description a Team, all resource belongs to team
// @tag.name user
// @tag.description User management
// @tag.name node
// @tag.description Node management
// @tag.name task
// @tag.description Task management
// @tag.name running
// @tag.description Task Running management
