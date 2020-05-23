// Code generated by go-swagger; DO NOT EDIT.

package image

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/metal-stack/metal-go/api/models"
)

// UpdateImageReader is a Reader for the UpdateImage structure.
type UpdateImageReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateImageReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewUpdateImageOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 409:
		result := NewUpdateImageConflict()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		result := NewUpdateImageDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUpdateImageOK creates a UpdateImageOK with default headers values
func NewUpdateImageOK() *UpdateImageOK {
	return &UpdateImageOK{}
}

/*UpdateImageOK handles this case with default header values.

OK
*/
type UpdateImageOK struct {
	Payload *models.V1ImageResponse
}

func (o *UpdateImageOK) Error() string {
	return fmt.Sprintf("[POST /v1/image][%d] updateImageOK  %+v", 200, o.Payload)
}

func (o *UpdateImageOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1ImageResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateImageConflict creates a UpdateImageConflict with default headers values
func NewUpdateImageConflict() *UpdateImageConflict {
	return &UpdateImageConflict{}
}

/*UpdateImageConflict handles this case with default header values.

Conflict
*/
type UpdateImageConflict struct {
	Payload *models.HttperrorsHTTPErrorResponse
}

func (o *UpdateImageConflict) Error() string {
	return fmt.Sprintf("[POST /v1/image][%d] updateImageConflict  %+v", 409, o.Payload)
}

func (o *UpdateImageConflict) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HttperrorsHTTPErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateImageDefault creates a UpdateImageDefault with default headers values
func NewUpdateImageDefault(code int) *UpdateImageDefault {
	return &UpdateImageDefault{
		_statusCode: code,
	}
}

/*UpdateImageDefault handles this case with default header values.

Error
*/
type UpdateImageDefault struct {
	_statusCode int

	Payload *models.HttperrorsHTTPErrorResponse
}

// Code gets the status code for the update image default response
func (o *UpdateImageDefault) Code() int {
	return o._statusCode
}

func (o *UpdateImageDefault) Error() string {
	return fmt.Sprintf("[POST /v1/image][%d] updateImage default  %+v", o._statusCode, o.Payload)
}

func (o *UpdateImageDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HttperrorsHTTPErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
