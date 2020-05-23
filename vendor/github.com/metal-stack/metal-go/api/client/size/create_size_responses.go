// Code generated by go-swagger; DO NOT EDIT.

package size

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/metal-stack/metal-go/api/models"
)

// CreateSizeReader is a Reader for the CreateSize structure.
type CreateSizeReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateSizeReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 201:
		result := NewCreateSizeCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 409:
		result := NewCreateSizeConflict()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		result := NewCreateSizeDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewCreateSizeCreated creates a CreateSizeCreated with default headers values
func NewCreateSizeCreated() *CreateSizeCreated {
	return &CreateSizeCreated{}
}

/*CreateSizeCreated handles this case with default header values.

Created
*/
type CreateSizeCreated struct {
	Payload *models.V1SizeResponse
}

func (o *CreateSizeCreated) Error() string {
	return fmt.Sprintf("[PUT /v1/size][%d] createSizeCreated  %+v", 201, o.Payload)
}

func (o *CreateSizeCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1SizeResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateSizeConflict creates a CreateSizeConflict with default headers values
func NewCreateSizeConflict() *CreateSizeConflict {
	return &CreateSizeConflict{}
}

/*CreateSizeConflict handles this case with default header values.

Conflict
*/
type CreateSizeConflict struct {
	Payload *models.HttperrorsHTTPErrorResponse
}

func (o *CreateSizeConflict) Error() string {
	return fmt.Sprintf("[PUT /v1/size][%d] createSizeConflict  %+v", 409, o.Payload)
}

func (o *CreateSizeConflict) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HttperrorsHTTPErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateSizeDefault creates a CreateSizeDefault with default headers values
func NewCreateSizeDefault(code int) *CreateSizeDefault {
	return &CreateSizeDefault{
		_statusCode: code,
	}
}

/*CreateSizeDefault handles this case with default header values.

Error
*/
type CreateSizeDefault struct {
	_statusCode int

	Payload *models.HttperrorsHTTPErrorResponse
}

// Code gets the status code for the create size default response
func (o *CreateSizeDefault) Code() int {
	return o._statusCode
}

func (o *CreateSizeDefault) Error() string {
	return fmt.Sprintf("[PUT /v1/size][%d] createSize default  %+v", o._statusCode, o.Payload)
}

func (o *CreateSizeDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HttperrorsHTTPErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
