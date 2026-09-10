package users_service

import (
	"context"
	"fmt"
	"uuid"
)

func (s UserService) DelUser(
	ctx context.Context,
	id uuid.UUID,
) error {
	if err := s.userRepository.DelUser(ctx, id); err != nil {
		return fmt.Errorf("del user from repository: %w", err)
	}

	return nil
}
