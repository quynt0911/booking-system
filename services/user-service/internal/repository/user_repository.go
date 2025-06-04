package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"booking-system/services/user-service/internal/model"
)

type UserRepository interface {
	Create(user *model.User) error
	GetByID(id string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id string) error
	List(limit, offset int) ([]*model.User, int, error)
	UpdatePassword(id, hashedPassword string) error
	UpdateVerificationStatus(id string, isVerified bool) error
	GetExpertsBySpecialization(specialization string) ([]*model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (id, email, password, first_name, last_name, role, phone, 
			avatar_url, bio, gender, date_of_birth, is_verified, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := r.db.Exec(query,
		user.ID, user.Email, user.Password, user.FirstName, user.LastName,
		user.Role, user.Phone, user.AvatarURL, user.Bio, user.Gender,
		user.DateOfBirth, user.IsVerified, user.IsActive, user.CreatedAt, user.UpdatedAt,
	)

	return err
}

func (r *userRepository) GetByID(id string) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, email, password, first_name, last_name, role, phone, 
			avatar_url, bio, gender, date_of_birth, is_verified, is_active, 
			created_at, updated_at
		FROM users 
		WHERE id = $1 AND is_active = true
	`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.Role, &user.Phone, &user.AvatarURL, &user.Bio, &user.Gender,
		&user.DateOfBirth, &user.IsVerified, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, email, password, first_name, last_name, role, phone, 
			avatar_url, bio, gender, date_of_birth, is_verified, is_active, 
			created_at, updated_at
		FROM users 
		WHERE email = $1 AND is_active = true
	`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.Role, &user.Phone, &user.AvatarURL, &user.Bio, &user.Gender,
		&user.DateOfBirth, &user.IsVerified, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Update(user *model.User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users 
		SET first_name = $2, last_name = $3, phone = $4, avatar_url = $5, 
			bio = $6, gender = $7, date_of_birth = $8, updated_at = $9
		WHERE id = $1 AND is_active = true
	`

	result, err := r.db.Exec(query,
		user.ID, user.FirstName, user.LastName, user.Phone, user.AvatarURL,
		user.Bio, user.Gender, user.DateOfBirth, user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) Delete(id string) error {
	query := `UPDATE users SET is_active = false, updated_at = $2 WHERE id = $1`
	
	result, err := r.db.Exec(query, id, time.Now())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) List(limit, offset int) ([]*model.User, int, error) {
	var users []*model.User
	var total int

	// Get total count
	countQuery := `SELECT COUNT(*) FROM users WHERE is_active = true`
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get users with pagination
	query := `
		SELECT id, email, first_name, last_name, role, phone, 
			avatar_url, bio, gender, date_of_birth, is_verified, is_active, 
			created_at, updated_at
		FROM users 
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.FirstName, &user.LastName,
			&user.Role, &user.Phone, &user.AvatarURL, &user.Bio, &user.Gender,
			&user.DateOfBirth, &user.IsVerified, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *userRepository) UpdatePassword(id, hashedPassword string) error {
	query := `UPDATE users SET password = $2, updated_at = $3 WHERE id = $1 AND is_active = true`
	
	result, err := r.db.Exec(query, id, hashedPassword, time.Now())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) UpdateVerificationStatus(id string, isVerified bool) error {
	query := `UPDATE users SET is_verified = $2, updated_at = $3 WHERE id = $1 AND is_active = true`
	
	result, err := r.db.Exec(query, id, isVerified, time.Now())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) GetExpertsBySpecialization(specialization string) ([]*model.User, error) {
	var users []*model.User

	query := `
		SELECT u.id, u.email, u.first_name, u.last_name, u.role, u.phone, 
			u.avatar_url, u.bio, u.gender, u.date_of_birth, u.is_verified, u.is_active, 
			u.created_at, u.updated_at
		FROM users u
		INNER JOIN experts e ON u.id = e.user_id
		WHERE u.role = 'expert' AND u.is_active = true AND e.is_available = true
		AND ($1 = '' OR e.specialization ILIKE '%' || $1 || '%')
		ORDER BY e.rating DESC, e.total_reviews DESC
	`

	rows, err := r.db.Query(query, specialization)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.FirstName, &user.LastName,
			&user.Role, &user.Phone, &user.AvatarURL, &user.Bio, &user.Gender,
			&user.DateOfBirth, &user.IsVerified, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}