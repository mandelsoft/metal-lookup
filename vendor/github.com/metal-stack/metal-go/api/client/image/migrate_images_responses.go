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

// MigrateImagesReader is a Reader for the MigrateImages structure.
type MigrateImagesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *MigrateImagesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewMigrateImagesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewMigrateImagesDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewMigrateImagesOK creates a MigrateImagesOK with default headers values
func NewMigrateImagesOK() *MigrateImagesOK {
	return &MigrateImagesOK{}
}

/*MigrateImagesOK handles this case with default header values.

OK
*/
type MigrateImagesOK struct {
	Payload []*models.V1ImageResponse
}

func (o *MigrateImagesOK) Error() string {
	return fmt.Sprintf("[GET /v1/image/migrate][%d] migrateImagesOK  %+v", 200, o.Payload)
}

func (o *MigrateImagesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewMigrateImagesDefault creates a MigrateImagesDefault with default headers values
func NewMigrateImagesDefault(code int) *MigrateImagesDefault {
	return &MigrateImagesDefault{
		_statusCode: code,
	}
}

/*MigrateImagesDefault handles this case with default header values.

Error
*/
type MigrateImagesDefault struct {
	_statusCode int

	Payload *models.HttperrorsHTTPErrorResponse
}

// Code gets the status code for the migrate images default response
func (o *MigrateImagesDefault) Code() int {
	return o._statusCode
}

func (o *MigrateImagesDefault) Error() string {
	return fmt.Sprintf("[GET /v1/image/migrate][%d] migrateImages default  %+v", o._statusCode, o.Payload)
}

func (o *MigrateImagesDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HttperrorsHTTPErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
