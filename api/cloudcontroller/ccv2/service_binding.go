package ccv2

import (
	"bytes"
	"encoding/json"

	"code.cloudfoundry.org/cli/api/cloudcontroller"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccerror"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv2/internal"
)

// ServiceBinding represents a Cloud Controller Service Binding.
type ServiceBinding struct {
	// GUID is the unique Service Binding identifier.
	GUID string
	// Name is the name of the service binding
	Name string
	// AppGUID is the associated application GUID.
	AppGUID string
	// ServiceInstanceGUID is the associated service GUID.
	ServiceInstanceGUID string
}

// UnmarshalJSON helps unmarshal a Cloud Controller Service Binding response.
func (serviceBinding *ServiceBinding) UnmarshalJSON(data []byte) error {
	var ccServiceBinding struct {
		Metadata internal.Metadata
		Entity   struct {
			AppGUID             string `json:"app_guid"`
			ServiceInstanceGUID string `json:"service_instance_guid"`
			Name                string `json:"name"`
		} `json:"entity"`
	}
	err := cloudcontroller.DecodeJSON(data, &ccServiceBinding)
	if err != nil {
		return err
	}

	serviceBinding.AppGUID = ccServiceBinding.Entity.AppGUID
	serviceBinding.GUID = ccServiceBinding.Metadata.GUID
	serviceBinding.ServiceInstanceGUID = ccServiceBinding.Entity.ServiceInstanceGUID
	serviceBinding.Name = ccServiceBinding.Entity.Name
	return nil
}

// serviceBindingRequestBody represents the body of the service binding create
// request.
type serviceBindingRequestBody struct {
	ServiceInstanceGUID string                 `json:"service_instance_guid"`
	AppGUID             string                 `json:"app_guid"`
	Name                string                 `json:"name,omitempty"`
	Parameters          map[string]interface{} `json:"parameters"`
}

func (client *Client) CreateServiceBinding(appGUID string, serviceInstanceGUID string, bindingName string, parameters map[string]interface{}) (ServiceBinding, Warnings, error) {
	requestBody := serviceBindingRequestBody{
		ServiceInstanceGUID: serviceInstanceGUID,
		AppGUID:             appGUID,
		Name:                bindingName,
		Parameters:          parameters,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return ServiceBinding{}, nil, err
	}

	request, err := client.newHTTPRequest(requestOptions{
		RequestName: internal.PostServiceBindingRequest,
		Body:        bytes.NewReader(bodyBytes),
	})
	if err != nil {
		return ServiceBinding{}, nil, err
	}

	var serviceBinding ServiceBinding
	response := cloudcontroller.Response{
		Result: &serviceBinding,
	}

	err = client.connection.Make(request, &response)
	if err != nil {
		return ServiceBinding{}, response.Warnings, err
	}

	return serviceBinding, response.Warnings, nil
}

// GetServiceBinding returns back a service binding with the proviced guid
func (client *Client) GetServiceBinding(guid string) (ServiceBinding, Warnings, error) {
	request, err := client.newHTTPRequest(requestOptions{
		RequestName: internal.GetServiceBindingRequest,
		URIParams:   Params{"service_binding_guid": guid},
	})
	if err != nil {
		return ServiceBinding{}, nil, err
	}

	var serviceBinding ServiceBinding
	response := cloudcontroller.Response{
		Result: &serviceBinding,
	}

	err = client.connection.Make(request, &response)
	return serviceBinding, response.Warnings, err
}

// GetServiceBindings returns back a list of Service Bindings based off of the
// provided filters.
func (client *Client) GetServiceBindings(filters ...Filter) ([]ServiceBinding, Warnings, error) {
	request, err := client.newHTTPRequest(requestOptions{
		RequestName: internal.GetServiceBindingsRequest,
		Query:       ConvertFilterParameters(filters),
	})
	if err != nil {
		return nil, nil, err
	}

	var fullBindingsList []ServiceBinding
	warnings, err := client.paginate(request, ServiceBinding{}, func(item interface{}) error {
		if binding, ok := item.(ServiceBinding); ok {
			fullBindingsList = append(fullBindingsList, binding)
		} else {
			return ccerror.UnknownObjectInListError{
				Expected:   ServiceBinding{},
				Unexpected: item,
			}
		}
		return nil
	})

	return fullBindingsList, warnings, err
}

// DeleteServiceBinding will destroy the requested Service Binding.
func (client *Client) DeleteServiceBinding(serviceBindingGUID string) (Warnings, error) {
	request, err := client.newHTTPRequest(requestOptions{
		RequestName: internal.DeleteServiceBindingRequest,
		URIParams:   map[string]string{"service_binding_guid": serviceBindingGUID},
	})
	if err != nil {
		return nil, err
	}

	var response cloudcontroller.Response
	err = client.connection.Make(request, &response)
	return response.Warnings, err
}

// GetServiceInstanceServiceBindings returns back a list of Service Bindings for the provided service instance GUID.
func (client *Client) GetServiceInstanceServiceBindings(serviceInstanceGUID string) ([]ServiceBinding, Warnings, error) {
	request, err := client.newHTTPRequest(requestOptions{
		RequestName: internal.GetServiceInstanceServiceBindingsRequest,
		URIParams:   map[string]string{"service_instance_guid": serviceInstanceGUID},
	})
	if err != nil {
		return nil, nil, err
	}

	var fullBindingsList []ServiceBinding
	warnings, err := client.paginate(request, ServiceBinding{}, func(item interface{}) error {
		if binding, ok := item.(ServiceBinding); ok {
			fullBindingsList = append(fullBindingsList, binding)
		} else {
			return ccerror.UnknownObjectInListError{
				Expected:   ServiceBinding{},
				Unexpected: item,
			}
		}
		return nil
	})

	return fullBindingsList, warnings, err
}

// GetUserProvidedServiceInstanceServiceBindings returns back a list of Service Bindings for the provided user provided service instance GUID.
func (client *Client) GetUserProvidedServiceInstanceServiceBindings(userProvidedServiceInstanceGUID string) ([]ServiceBinding, Warnings, error) {
	request, err := client.newHTTPRequest(requestOptions{
		RequestName: internal.GetUserProvidedServiceInstanceServiceBindingsRequest,
		URIParams:   map[string]string{"user_provided_service_instance_guid": userProvidedServiceInstanceGUID},
	})
	if err != nil {
		return nil, nil, err
	}

	var fullBindingsList []ServiceBinding
	warnings, err := client.paginate(request, ServiceBinding{}, func(item interface{}) error {
		if binding, ok := item.(ServiceBinding); ok {
			fullBindingsList = append(fullBindingsList, binding)
		} else {
			return ccerror.UnknownObjectInListError{
				Expected:   ServiceBinding{},
				Unexpected: item,
			}
		}
		return nil
	})

	return fullBindingsList, warnings, err
}
