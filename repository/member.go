package repository

import (
	"context"
	"errors"

	"github.com/faruryo/toban-api/models"
	"gorm.io/gorm"
)

func (r repository) GetMemberByID(ctx context.Context, id uint) (*models.Member, error) {
	return getMemberByID(r.db, id)
}

func getMemberByID(db *gorm.DB, id uint) (*models.Member, error) {
	var member models.Member
	err := db.First(&member, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNoSuchEntity
	}
	if err != nil {
		return nil, err
	}

	return &member, nil
}

func (r repository) GetAllMembers(ctx context.Context) ([]*models.Member, error) {
	var members []*models.Member
	if err := r.db.Find(&members).Error; err != nil {
		return nil, err
	}

	return members, nil
}

func (r repository) CreateMember(ctx context.Context, member *models.Member) (*models.Member, error) {
	if member.ID != 0 {
		return nil, ErrBadRequestIDMustBeZero
	}
	if !member.CreatedAt.IsZero() {
		return nil, ErrBadRequestUpdateCreatedAt
	}
	if !member.UpdatedAt.IsZero() {
		return nil, ErrBadRequestUpdateUpdatedAt
	}

	if err := r.db.Create(member).Error; err != nil {
		return nil, err
	}

	return member, nil
}

func (r repository) UpdateMember(ctx context.Context, input *models.UpdateMemberInput) (*models.Member, error) {
	if input.ID == 0 {
		return nil, ErrBadRequestIDMustNotBeZero
	}

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	output, err := getMemberByID(tx, input.ID)
	if err != nil {
		return nil, err
	}

	if input.SlackID != nil {
		output.SlackID = *input.SlackID
	}
	if input.Name != nil {
		output.Name = *input.Name
	}

	if err := tx.Save(&output).Error; err != nil {
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return output, nil
}

func (r repository) DeleteMemberByID(ctx context.Context, id uint) (bool, error) {
	if id == 0 {
		return false, ErrBadRequestIDMustNotBeZero
	}

	var member models.Member
	if err := r.db.Delete(member, id).Error; err != nil {
		return false, err
	}

	return true, nil
}
