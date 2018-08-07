// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// FindUsersHandlerFunc turns a function with the right signature into a find users handler
type FindUsersHandlerFunc func(FindUsersParams) middleware.Responder

// Handle executing the request and returning a response
func (fn FindUsersHandlerFunc) Handle(params FindUsersParams) middleware.Responder {
	return fn(params)
}

// FindUsersHandler interface for that can handle valid find users params
type FindUsersHandler interface {
	Handle(FindUsersParams) middleware.Responder
}

// NewFindUsers creates a new http.Handler for the find users operation
func NewFindUsers(ctx *middleware.Context, handler FindUsersHandler) *FindUsers {
	return &FindUsers{Context: ctx, Handler: handler}
}

/*FindUsers swagger:route GET /user users findUsers

FindUsers find users API

*/
type FindUsers struct {
	Context *middleware.Context
	Handler FindUsersHandler
}

func (o *FindUsers) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewFindUsersParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
