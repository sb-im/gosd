package app

// @title Gosd RESTful API
// @version 3.0.0
// @description superdock web service restful api

// @host localhost:8000
// @BasePath /gosd/api/v3
// @query.collection.format multi

// @securitydefinitions.bearerAuth BearerAuth
// @scheme bearer
// @bearerFormat JWT

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl /oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @tag.name status
// @tag.description Server Time, Running Status
// @tag.name auth
// @tag.description auth, login, logout, token
// @tag.name team
// @tag.description a Team, all resource belongs to team
// @tag.name user
// @tag.description User
// @tag.name node
// @tag.description Node
// @tag.name task
// @tag.description Task
