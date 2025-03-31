package toauth

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// UserToken represents an OAuth token for a specific provider
type UserToken struct {
	ProviderName string    `json:"provider_name"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Expiry       time.Time `json:"expiry"`
}

// User represents a user with multiple OAuth tokens
type User struct {
	ID     string      `json:"id"`
	Email  string      `json:"email"`
	Tokens []UserToken `json:"tokens"`
}

// UserRepository manages users and their OAuth tokens
type UserRepository struct {
	users    map[string]*User // map of user ID to user
	filePath string           // path to the JSON file
	mu       sync.RWMutex     // mutex for thread safety
}

// NewUserRepository creates a new user repository
func NewUserRepository(filePath string) (*UserRepository, error) {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	repo := &UserRepository{
		users:    make(map[string]*User),
		filePath: filePath,
	}

	// Load existing users from the file if it exists
	if _, err := os.Stat(filePath); err == nil {
		if err := repo.load(); err != nil {
			return nil, fmt.Errorf("failed to load users: %w", err)
		}
	}

	return repo, nil
}

// load loads users from the JSON file
func (r *UserRepository) load() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var users []*User
	if err := json.Unmarshal(data, &users); err != nil {
		return fmt.Errorf("failed to unmarshal users: %w", err)
	}

	// Populate the users map
	r.users = make(map[string]*User)
	for _, user := range users {
		r.users[user.ID] = user
	}

	return nil
}

// save saves users to the JSON file
func (r *UserRepository) save() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Convert the map to a slice
	var users []*User
	for _, user := range r.users {
		users = append(users, user)
	}

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal users: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// GetUserByID returns a user by ID
func (r *UserRepository) GetUserByID(id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", id)
	}

	return user, nil
}

// GetUserByEmail returns a user by email
func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, fmt.Errorf("user not found: %s", email)
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(id, email string) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if user already exists
	if _, ok := r.users[id]; ok {
		return nil, fmt.Errorf("user already exists: %s", id)
	}

	// Create a new user
	user := &User{
		ID:     id,
		Email:  email,
		Tokens: []UserToken{},
	}

	// Add the user to the map
	r.users[id] = user

	// Save the users to the file
	if err := r.save(); err != nil {
		return nil, fmt.Errorf("failed to save users: %w", err)
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
	r.mu.Lock()
	defer r.mu.Unlock()

	// Get the user
	user, ok := r.users[userID]
	if !ok {
		return fmt.Errorf("user not found: %s", userID)
	}

	// Create a new token
	userToken := UserToken{
		ProviderName: providerName,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}

	// Check if the user already has a token for this provider
	for i, t := range user.Tokens {
		if t.ProviderName == providerName {
			// Update the existing token
			user.Tokens[i] = userToken
			// Save the users to the file
			if err := r.save(); err != nil {
				return fmt.Errorf("failed to save users: %w", err)
			}
			return nil
		}
	}

	// Add the token to the user
	user.Tokens = append(user.Tokens, userToken)

	// Save the users to the file
	if err := r.save(); err != nil {
		return fmt.Errorf("failed to save users: %w", err)
	}

	return nil
}

// GetToken returns an OAuth token for a user and provider
func (r *UserRepository) GetToken(userID, providerName string) (*oauth2.Token, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Get the user
	user, ok := r.users[userID]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", userID)
	}

	// Find the token for the provider
	for _, t := range user.Tokens {
		if t.ProviderName == providerName {
			// Convert to oauth2.Token
			return &oauth2.Token{
				AccessToken:  t.AccessToken,
				RefreshToken: t.RefreshToken,
				Expiry:       t.Expiry,
			}, nil
		}
	}

	return nil, fmt.Errorf("token not found for provider: %s", providerName)
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
