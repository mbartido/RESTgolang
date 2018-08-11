// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"simpleAPI/swagger3/models"
	"simpleAPI/swagger3/restapi/operations"
	"simpleAPI/swagger3/restapi/operations/users"
)

// Users struct which has array of users
type Users struct {
	Users []*models.User `json:"users"`
}

//go:generate swagger generate server --target .. --name UserList --spec ../swagger.yml

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

	// initialize Users array
	var uList Users

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	api.UsersAddOneHandler = users.AddOneHandlerFunc(func(params users.AddOneParams) middleware.Responder {
		// File opening
		jsonFile, err := os.Open("users.json")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Opened users.json successfully")
		defer jsonFile.Close()
		// read file as byte array
		byteValue, _ := ioutil.ReadAll(jsonFile)

		// Read json into uList and append users to main structure of json
		json.Unmarshal(byteValue, &uList)
		result := make([]*models.User, 0)
		mainStruct := Users{}
		for i := 0; i < len(uList.Users); i++ {
			result = append(result, uList.Users[i])
			mainStruct.Users = append(mainStruct.Users, uList.Users[i])
		}

		// Add user to mainStruct
		// Main struct is for holding users with new added user
		resultLen := int64(len(result)) + 1 // because 0 indexed
		u := &models.User{ID: resultLen, Name: params.Body.Name}
		mainStruct.Users = append(mainStruct.Users, u)

		// Error checking
		error := errors.New(500, "user must be present")
		endList, err := json.MarshalIndent(mainStruct, "", "  ")
		if err != nil {
			return users.NewAddOneDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(error.Error())})
		}
		ioutil.WriteFile("users.json", endList, 0644)
		return users.NewAddOneCreated().WithPayload(params.Body)
	})

	api.UsersDestroyOneHandler = users.DestroyOneHandlerFunc(func(params users.DestroyOneParams) middleware.Responder {
		// File opening
		jsonFile, err := os.Open("users.json")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Opened users.json successfully")
		defer jsonFile.Close()
		// read file as byte array
		byteValue, _ := ioutil.ReadAll(jsonFile)

		// Read json into uList and append users to main structure of json
		json.Unmarshal(byteValue, &uList)
		result := make([]*models.User, 0)
		mainStruct := Users{}
		for i := 0; i < len(uList.Users); i++ {
			result = append(result, uList.Users[i])
			mainStruct.Users = append(mainStruct.Users, uList.Users[i])
		}

		// Error checking
		error := errors.NotFound("not found: user %d", params.ID)
		if int(params.ID) > len(mainStruct.Users) || int(params.ID) <= 0 {
			return users.NewUpdateOneDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(error.Error())})
		}

		// Delete user in mainStruct using slices
		i1 := mainStruct.Users[:params.ID-1]
		i2 := mainStruct.Users[params.ID:]
		for i := 0; i < len(i2); i++ { // adjust id of right side
			i2[i].ID = i2[i].ID - 1
		}
		s := append(i1, i2...)
		mainStruct.Users = s
		endList, _ := json.MarshalIndent(mainStruct, "", "  ")
		// Write to file
		ioutil.WriteFile("users.json", endList, 0644)
		return users.NewDestroyOneNoContent()
	})

	api.UsersFindUsersHandler = users.FindUsersHandlerFunc(func(params users.FindUsersParams) middleware.Responder {
		// File opening
		jsonFile, err := os.Open("users.json")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Opened users.json successfully")
		defer jsonFile.Close()
		// read file as byte array
		byteValue, _ := ioutil.ReadAll(jsonFile)

		// Reads json into result and then we use map of result to return
		json.Unmarshal(byteValue, &uList)
		result := make([]*models.User, 0)
		for i := 0; i < len(uList.Users); i++ {
			result = append(result, uList.Users[i])
		}
		return users.NewFindUsersOK().WithPayload(result)
	})

	api.UsersGetOneHandler = users.GetOneHandlerFunc(func(params users.GetOneParams) middleware.Responder {
		// File opening
		jsonFile, err := os.Open("users.json")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Opened users.json successfully")
		defer jsonFile.Close()
		// read file as byte array
		byteValue, _ := ioutil.ReadAll(jsonFile)

		// Reading json into list where we'll find the user
		json.Unmarshal(byteValue, &uList)
		list := make([]*models.User, 0)
		for i := 0; i < len(uList.Users); i++ {
			list = append(list, uList.Users[i])
		}

		// Error checking
		error := errors.NotFound("not found: user %d", params.ID)
		if int(params.ID) > len(list) || int(params.ID) <= 0 {
			return users.NewGetOneDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(error.Error())})
		}
		return users.NewGetOneOK().WithPayload(list[params.ID-1]) // ID -1 because list is 0 indexed
	})

	api.UsersUpdateOneHandler = users.UpdateOneHandlerFunc(func(params users.UpdateOneParams) middleware.Responder {
		// File opening
		jsonFile, err := os.Open("users.json")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Opened users.json successfully")
		defer jsonFile.Close()
		// read file as byte array
		byteValue, _ := ioutil.ReadAll(jsonFile)

		// Read json into uList and append users to main structure of json
		json.Unmarshal(byteValue, &uList)
		result := make([]*models.User, 0)
		mainStruct := Users{}
		for i := 0; i < len(uList.Users); i++ {
			result = append(result, uList.Users[i])
			mainStruct.Users = append(mainStruct.Users, uList.Users[i])
		}

		// Error checking
		error := errors.NotFound("not found: user %d", params.ID)
		if int(params.ID) > len(mainStruct.Users) || int(params.ID) <= 0 {
			return users.NewUpdateOneDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(error.Error())})
		}

		// Update user in mainStruct
		mainStruct.Users[params.ID-1].Name = params.Body.Name
		endList, _ := json.MarshalIndent(mainStruct, "", "  ")
		// Write to file
		ioutil.WriteFile("users.json", endList, 0644)
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
