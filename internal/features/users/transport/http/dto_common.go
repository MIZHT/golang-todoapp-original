package users_transport_http

import "github.com/MIZHT/golang-todoapp/internal/core/domain"

type UserDTOResponse struct {
	ID          int     `json:"id" example:"0"`
	Vesrion     int     `json:"version" example:"0"`
	FullName    string  `json:"full_name"    example:"Jorge Bush"`
	PhoneNumber *string `json:"phone_number" example:"+79877128288"`
}

func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:          user.ID,
		Vesrion:     user.Vesrion,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))

	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}

	return usersDTO

}
