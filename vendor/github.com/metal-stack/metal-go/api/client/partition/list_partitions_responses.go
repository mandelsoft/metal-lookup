// Code generated by go-swagger; DO NOT EDIT.

package partition

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/metal-stack/metal-go/api/models"
)

// ListPartitionsReader is a Reader for the ListPartitions structure.
type ListPartitionsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListPartitionsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewListPartitionsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewListPartitionsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListPartitionsOK creates a ListPartitionsOK with default headers values
func NewListPartitionsOK() *ListPartitionsOK {
	return &ListPartitionsOK{}
}

/*ListPartitionsOK handles this case with default header values.

OK
*/
type ListPartitionsOK struct {
	Payload []*models.V1PartitionResponse
}

func (o *ListPartitionsOK) Error() string {
	return fmt.Sprintf("[GET /v1/partition][%d] listPartitionsOK  %+v", 200, o.Payload)
}

func (o *ListPartitionsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListPartitionsDefault creates a ListPartitionsDefault with default headers values
func NewListPartitionsDefault(code int) *ListPartitionsDefault {
	return &ListPartitionsDefault{
		_statusCode: code,
	}
}

/*ListPartitionsDefault handles this case with default header values.

Error
*/
type ListPartitionsDefault struct {
	_statusCode int

	Payload *models.HttperrorsHTTPErrorResponse
}

// Code gets the status code for the list partitions default response
func (o *ListPartitionsDefault) Code() int {
	return o._statusCode
}

func (o *ListPartitionsDefault) Error() string {
	return fmt.Sprintf("[GET /v1/partition][%d] listPartitions default  %+v", o._statusCode, o.Payload)
}

func (o *ListPartitionsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HttperrorsHTTPErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
