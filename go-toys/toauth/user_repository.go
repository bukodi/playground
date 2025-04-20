package toauth

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"time"
)

// UserToken represents an OAuth token for a specific provider
type UserToken struct {
	ID           uint      `gorm:"primaryKey" json:"-"`
	UserID       string    `gorm:"index" json:"-"`
	ProviderName string    `gorm:"index" json:"provider_name"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Expiry       time.Time `json:"expiry"`
}

// User represents a user with multiple OAuth tokens
type User struct {
	ID     string      `gorm:"primaryKey" json:"id"`
	Email  string      `gorm:"uniqueIndex" json:"email"`
	Tokens []UserToken `gorm:"foreignKey:UserID" json:"tokens"`
}

// UserRepository manages users and their OAuth tokens
type UserRepository struct {
	db       *gorm.DB
	filePath string // kept for backward compatibility
}

// NewUserRepository creates a new user repository
func NewUserRepository(filePath string) (*UserRepository, error) {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Create a SQLite database in the same directory
	dbPath := filepath.Join(dir, "users.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&User{}, &UserToken{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	repo := &UserRepository{
		db:       db,
		filePath: filePath, // Keep filePath for backward compatibility
	}

	// If the JSON file exists and the database is empty, migrate data from JSON to SQLite
	if _, err := os.Stat(filePath); err == nil {
		var count int64
		db.Model(&User{}).Count(&count)
		if count == 0 {
			if err := repo.migrateFromJSON(filePath); err != nil {
				return nil, fmt.Errorf("failed to migrate data from JSON: %w", err)
			}
		}
	}

	return repo, nil
}

// migrateFromJSON migrates data from the JSON file to the SQLite database
func (r *UserRepository) migrateFromJSON(filePath string) error {
	// Read the JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse the JSON data
	var users []*User
	if err := json.Unmarshal(data, &users); err != nil {
		return fmt.Errorf("failed to unmarshal users: %w", err)
	}

	// Migrate each user to the database
	for _, user := range users {
		// Create the user in the database
		if err := r.db.Create(user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	return nil
}

// GetUserByID returns a user by ID
func (r *UserRepository) GetUserByID(id string) (*User, error) {
	var user User
	if err := r.db.Preload("Tokens").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %s", id)
	}
	return &user, nil
}

// GetUserByEmail returns a user by email
func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Preload("Tokens").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %s", email)
	}
	return &user, nil
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(id, email string) (*User, error) {
	// Check if user already exists
	var count int64
	r.db.Model(&User{}).Where("id = ?", id).Count(&count)
	if count > 0 {
		return nil, fmt.Errorf("user already exists: %s", id)
	}

	// Create a new user
	user := &User{
		ID:     id,
		Email:  email,
		Tokens: []UserToken{},
	}

	// Save the user to the database
	if err := r.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetOrCreateUserByEmail gets a user by email or creates a new one if it doesn't exist
func (r *UserRepository) GetOrCreateUserByEmail(email string) (*User, error) {
	// Try to get the user by email
	user, err := r.GetUserByEmail(email)
	if err == nil {
		return user, nil
	}

	// Create a new user with a generated ID
	id := fmt.Sprintf("user_%d", time.Now().UnixNano())
	return r.CreateUser(id, email)
}

// StoreToken stores an OAuth token for a user
func (r *UserRepository) StoreToken(userID, providerName string, token *oauth2.Token) error {
	// Get the user
	_, err := r.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %s", userID)
	}

	// Check if the user already has a token for this provider
	var existingToken UserToken
	result := r.db.Where("user_id = ? AND provider_name = ?", userID, providerName).First(&existingToken)
	if result.Error == nil {
		// Update the existing token
		existingToken.AccessToken = token.AccessToken
		existingToken.RefreshToken = token.RefreshToken
		existingToken.Expiry = token.Expiry
		if err := r.db.Save(&existingToken).Error; err != nil {
			return fmt.Errorf("failed to update token: %w", err)
		}
		return nil
	}

	// Create a new token
	userToken := UserToken{
		UserID:       userID,
		ProviderName: providerName,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}

	// Add the token to the database
	if err := r.db.Create(&userToken).Error; err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}

	return nil
}

// GetToken returns an OAuth token for a user and provider
func (r *UserRepository) GetToken(userID, providerName string) (*oauth2.Token, error) {
	// Get the user
	_, err := r.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", userID)
	}

	// Find the token for the provider
	var userToken UserToken
	if err := r.db.Where("user_id = ? AND provider_name = ?", userID, providerName).First(&userToken).Error; err != nil {
		return nil, fmt.Errorf("token not found for provider: %s", providerName)
	}

	// Convert to oauth2.Token
	return &oauth2.Token{
		AccessToken:  userToken.AccessToken,
		RefreshToken: userToken.RefreshToken,
		Expiry:       userToken.Expiry,
	}, nil
}

// DefaultUserRepository is the default user repository
var DefaultUserRepository *UserRepository

// Initialize the default user repository
func init() {
	// Create the default user repository
	repo, err := NewUserRepository("./testdata/users.json")
	if err != nil {
		// Just log the error and continue
		fmt.Printf("Failed to create user repository: %v\n", err)
		return
	}

	DefaultUserRepository = repo
}
