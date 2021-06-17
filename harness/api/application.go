package api

import (
	"fmt"

	"github.com/micahlmartin/terraform-provider-harness/harness/api/graphql"
)

// Helper type for accessing all application related crud methods
type ApplicationClient struct {
	APIClient *Client
}

// Get the client for interacting with Harness Applications
func (c *Client) Applications() *ApplicationClient {
	return &ApplicationClient{
		APIClient: c,
	}
}

// CRUD
func (ac *ApplicationClient) GetApplicationById(id string) (*graphql.Application, error) {
	query := &GraphQLQuery{
		Query: fmt.Sprintf(`query($applicationId: String!) {
			application(applicationId: $applicationId) {
				%s
			}
		}`, standardApplicationFields),
		Variables: map[string]interface{}{
			"applicationId": id,
		},
	}

	res := struct {
		Application graphql.Application
	}{}
	err := ac.APIClient.ExecuteGraphQLQuery(query, &res)

	if err != nil {
		return nil, err
	}

	return &res.Application, nil
}

func (ac *ApplicationClient) GetApplicationByName(name string) (*graphql.Application, error) {
	query := &GraphQLQuery{
		Query: fmt.Sprintf(`query($name: String!) {
			applicationByName(name: $name) {
				%s	
			}
		}`, standardApplicationFields),
		Variables: map[string]interface{}{
			"name": name,
		},
	}

	res := &struct {
		ApplicationByName graphql.Application
	}{}
	err := ac.APIClient.ExecuteGraphQLQuery(query, &res)

	if err != nil {
		return nil, err
	}

	return &res.ApplicationByName, nil
}

func (ac *ApplicationClient) CreateApplication(input *graphql.Application) (*graphql.Application, error) {

	query := &GraphQLQuery{
		Query: `mutation createapp($app: CreateApplicationInput!) {
			createApplication(input: $app) {
				clientMutationId
				application {
					id
					name
					description
				}
			}
		}`,
		Variables: map[string]interface{}{
			"app": &input,
		},
	}

	res := &struct {
		CreateApplication graphql.CreateApplicationPayload
	}{}
	err := ac.APIClient.ExecuteGraphQLQuery(query, &res)

	if err != nil {
		return nil, err
	}

	return res.CreateApplication.Application, nil
}

func (ac *ApplicationClient) DeleteApplication(id string) error {

	query := &GraphQLQuery{
		Query: `mutation deleteApp($app: DeleteApplicationInput!) {
			deleteApplication(input: $app) {
				clientMutationId
			}
		}`,
		Variables: map[string]interface{}{
			"app": &graphql.DeleteApplicationInput{
				ApplicationId: id,
			},
		},
	}

	err := ac.APIClient.ExecuteGraphQLQuery(query, &struct{}{})

	return err
}

func (ac *ApplicationClient) UpdateApplication(input *graphql.UpdateApplicationInput) (*graphql.Application, error) {

	query := &GraphQLQuery{
		Query: fmt.Sprintf(`mutation updateapp($app: UpdateApplicationInput!) {
			updateApplication(input: $app) {
				clientMutationId
				application {
					%s
				}
			}
		}`, standardApplicationFields),
		Variables: map[string]interface{}{
			"app": &input,
		},
	}

	res := struct {
		UpdateApplication graphql.UpdateApplicationPayload
	}{}

	err := ac.APIClient.ExecuteGraphQLQuery(query, &res)

	if err != nil {
		return nil, err
	}

	return res.UpdateApplication.Application, nil
}

const (
	standardApplicationFields = `
	id
	name
	description
	createdBy {
		id
		name
		email
	}
	gitSyncConfig {
		branch
		gitConnector {
			id
			name
			branch
		}
		repositoryName
		syncEnabled
	}
	tags {
		name
		value
	}	
	`
)
