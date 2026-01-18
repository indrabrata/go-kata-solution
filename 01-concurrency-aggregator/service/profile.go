package service

type ProfileService struct {
	// Define fields here
}

func NewProfileService() *ProfileService {
	return &ProfileService{}
}

func (s *ProfileService) GetProfile() (string, error) {
	// Implement logic here
	return "Alice ", nil
}
