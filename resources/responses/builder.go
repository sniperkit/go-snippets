package responses

//
// Using these tips as guidelines to these interfaces: http://www.restapitutorial.com/lessons/restquicktips.html
//

//Builder is a fluent API to create a response
type Builder interface {
	Success() Response                         //200
	SuccessWithBody(body interface{}) Response //200
	Created() Response                         //201   //Created(newResourceLocation string) Response
	NoContent() Response                       //204
	BadRequest(msg string) Response            //400
	Unauthorized(msg string) Response          //401
	Forbidden() Response                       //403
	NotFound() Response                        //404
	MethodNotAllowed() Response                //405   //MethodNotAllowed(allowedMethods []string) Response
	Conflict() Response                        //409
	InternalServerError(err error) Response    //500

	CustomWithBody(status Status, body interface{}) Response
	CustomError(status Status, err error) Response

	TODO(message string) Response //500  //Custom error for unimplemented responses
}

type builder struct {
	status Status
	body   interface{}
	format Format
}

func (b *builder) toResponse() Response {
	return Response{
		Status: b.status,
		Body:   b.body,
		Format: b.format,
	}
}

func (b *builder) Success() Response {
	return b.SuccessWithBody(nil)
}

func (b *builder) SuccessWithBody(body interface{}) Response {
	if body != nil {
		b.status = StatusOK
	} else {
		b.status = StatusNoContent
	}
	b.body = body
	return b.toResponse()
}

func (b *builder) Created() Response {
	b.status = StatusCreated
	return b.toResponse()
}

func (b *builder) NoContent() Response {
	b.status = StatusNoContent
	return b.toResponse()
}

func (b *builder) BadRequest(msg string) Response {
	b.status = StatusBadRequest
	b.body = &Error{msg}
	return b.toResponse()
}

func (b *builder) Unauthorized(msg string) Response {
	b.status = StatusUnauthorized
	b.body = &Error{msg}
	return b.toResponse()
}

func (b *builder) Forbidden() Response {
	b.status = StatusForbidden
	return b.toResponse()
}

func (b *builder) NotFound() Response {
	b.status = StatusNotFound
	return b.toResponse()
}

func (b *builder) MethodNotAllowed() Response {
	b.status = StatusMethodNotAllowed
	return b.toResponse()
}

func (b *builder) Conflict() Response {
	b.status = StatusConflict
	return b.toResponse()
}

func (b *builder) InternalServerError(err error) Response {
	b.status = StatusInternalServerError
	b.body = &Error{err.Error()}
	return b.toResponse()
}

func (b *builder) CustomWithBody(status Status, body interface{}) Response {
	b.status = status
	b.body = body
	return b.toResponse()
}

func (b *builder) CustomError(status Status, err error) Response {
	b.status = status
	b.body = &Error{err.Error()}
	return b.toResponse()
}

func (b *builder) TODO(message string) Response {
	b.status = StatusInternalServerError
	b.body = &Error{message}
	return b.toResponse()
}
