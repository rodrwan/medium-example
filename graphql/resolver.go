package graphql

import (
	"context"

	mediumexample "github.com/rodrwan/medium-example"
)

// Resolver represent a GraphQL resolver that need to implemente graphql interface.
type Resolver struct {
	Service mediumexample.UserService
}

// Mutation Expose a MutationResolver
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query Expose a QueryResolver
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, input *mediumexample.UserData) (*mediumexample.User, error) {
	u := &mediumexample.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Phone:     input.Phone,
		Birthdate: input.Birthdate,
		Address: &mediumexample.Address{
			AddressLine: input.Address.AddressLine,
			City:        input.Address.City,
			Locality:    input.Address.Locality,
			Region:      input.Address.Region,
			Country:     input.Address.Country,
			PostalCode:  input.Address.PostalCode,
		},
	}

	if err := r.Service.Create(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Users(ctx context.Context) ([]*mediumexample.User, error) {
	users, err := r.Service.Select(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (mediumexample.User, error) {
	user, err := r.Service.GetByID(ctx, id)
	if err != nil {
		return mediumexample.User{}, err
	}

	return *user, nil
}
