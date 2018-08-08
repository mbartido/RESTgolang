// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"
	"sync"
	"sync/atomic"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"simpleAPI/swagger3/models"
	"simpleAPI/swagger3/restapi/operations"
	"simpleAPI/swagger3/restapi/operations/users"
)

//go:generate swagger generate server --target .. --name UserList --spec ../swagger.yml

var userList = make(map[int64]*models.User)
var lastID int64
var usersLock = &sync.Mutex{}

func newUserID() int64 {
	return atomic.AddInt64(&lastID, 1)
}

func addUser(user *models.User) error {
	if user == nil {
		return errors.New(500, "user must be present")
	}

	usersLock.Lock()
	defer usersLock.Unlock()

	newID := newUserID()
	user.ID = newID
	userList[newID] = user

	return nil
}

func updateUser(id int64, user *models.User) error {
	if user == nil {
		return errors.New(500, "user must be present")
	}

	usersLock.Lock()
	defer usersLock.Unlock()

	_, exists := userList[id]
	if !exists {
		return errors.NotFound("not found: user %d", id)
	}

	user.ID = id
	userList[id] = user
	return nil
}

func deleteUser(id int64) error {
	usersLock.Lock()
	defer usersLock.Unlock()

	_, exists := userList[id]
	if !exists {
		return errors.NotFound("not found: user %d", id)
	}

	delete(userList, id)
	return nil
}

func allUsers(since int64, limit int32) (result []*models.User) {
	result = make([]*models.User, 0)
	for id, user := range userList {
		if len(result) >= int(limit) {
			return
		}
		if since == 0 || id > since {
			result = append(result, user)
		}
	}
	return
}

func oneUser(one int64) *models.User {
	count := int64(0)
	for id := range userList {
		if id == one {
			count = one
		}
	}
	return userList[count]
}

func configureFlags(api *operations.UserListAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.UserListAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.UsersAddOneHandler = users.AddOneHandlerFunc(func(params users.AddOneParams) middleware.Responder {
		if err := addUser(params.Body); err != nil {
			return users.NewAddOneDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return users.NewAddOneCreated().WithPayload(params.Body)
	})
	api.UsersDestroyOneHandler = users.DestroyOneHandlerFunc(func(params users.DestroyOneParams) middleware.Responder {
		if err := deleteUser(params.ID); err != nil {
			return users.NewDestroyOneDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return users.NewDestroyOneNoContent()
	})
	api.UsersFindUsersHandler = users.FindUsersHandlerFunc(func(params users.FindUsersParams) middleware.Responder {
		mergedParams := users.NewFindUsersParams()
		mergedParams.Since = swag.Int64(0)
		if params.Since != nil {
			mergedParams.Since = params.Since
		}
		if params.Limit != nil {
			mergedParams.Limit = params.Limit
		}
		return users.NewFindUsersOK().WithPayload(allUsers(*mergedParams.Since, *mergedParams.Limit))
	})
	api.UsersGetOneHandler = users.GetOneHandlerFunc(func(params users.GetOneParams) middleware.Responder {
		// if err := oneUser(params.ID, nil); err != nil {
		// 	return users.NewGetOneDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		// }
		// return users.NewGetOneOK().WithPayload(oneUser(params.ID))
		return users.NewGetOneOK().WithPayload(oneUser(params.ID))
	})
	api.UsersUpdateOneHandler = users.UpdateOneHandlerFunc(func(params users.UpdateOneParams) middleware.Responder {
		if err := updateUser(params.ID, params.Body); err != nil {
			return users.NewUpdateOneDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return users.NewUpdateOneOK().WithPayload(params.Body)
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
