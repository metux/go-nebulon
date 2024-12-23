// Nebulon blockstore API v1
//
//	Schemes: http
//	BasePath: /v1
//	Version: 1.0.0
//	Host: localhost:8080
//
//	Consumes:
//	- application/json
//	- application/octet-stream
//
//	Produces:
//	- application/json
//	- application/octet-stream
//
// swagger:meta
package rest_v1

// swagger:model

// these structs only serving the api docs - not used in the actual code

// Raw block data (application/octet-stream)
//
// swagger:response blockDataResponse
type swagger_blockDataResponse struct {
}

// Object not found
//
// swagger:response notFoundResponse
type swagger_notFoundResponse struct {
}

// internal error
//
// swagger:response internalErrorResponse
type swagger_internalErrorResponse struct {
}

// content not acceptable
//
// swagger:response notAcceptableResponse
type swagger_notAcceptableResponse struct {
}

// object created
//
// swagger:response objectCreatedResponse
type swagger_objectCreatedResponse struct {
}

// keeping block
//
// swagger:response keepingBlock
type swagger_keepingBlock struct {
}

// reference list
//
// swagger:response blockrefListResponse
type swagger_blockrefListResponse struct {
}

// bad request
//
// happens eg. when on invalid/missing URL parameters
//
// swagger:response badRequest
type swagger_badRequest struct {
}

// hello reply
//
// swagger:response helloResponse
type swagger_helloResponse struct {
}

// BlockRef reftype
//
// swagger:parameters ServeBlockGet ServeBlockPut ServeBlockKeep
type swagger_paramBlockRefstruct struct {
	RefType string
	RefOID  string
}
