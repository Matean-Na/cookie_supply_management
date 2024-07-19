package seeds

import (
	"cookie_supply_management/internal/models"
	"cookie_supply_management/pkg/base/base_seed"
	"cookie_supply_management/pkg/security"
	"encoding/json"
)

type UserSeed struct {
	*base_seed.GenericSeed
}

func NewUserSeed() *UserSeed {
	s := base_seed.NewGenericSeed(models.User{})
	return &UserSeed{&s}
}

func (s *UserSeed) Seed() (base_seed.Summary, error) {
	defer s.Summarize()
	data, err := s.LoadFixture()
	if err != nil {
		return s.Error(err)
	}

	var fixtureData map[string][]models.User
	if err = json.Unmarshal(data, &fixtureData); err != nil {
		return s.Error(err)
	}

	for _, user := range fixtureData["data"] {
		var exists bool = s.Exists(map[string]interface{}{
			"username": user.Username,
		})

		if exists {
			s.Summary.Exist++
			continue
		}

		hashPassword, err := security.Hash(user.Password)
		if err != nil {
			s.LogFail(err.Error())
			s.Summary.Errors++
			continue
		}

		user.Password = hashPassword

		r := s.Query.Create(&user)
		if err := r.Error; err != nil {
			s.LogFail(err.Error())
			s.Summary.Errors++
			continue
		}

		s.Summary.Created += r.RowsAffected
	}

	return s.Summary, nil
}
