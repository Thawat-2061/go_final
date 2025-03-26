package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	CustomerID uint   `json:"customer_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone_number,omitempty"`
	Address    string `json:"address,omitempty"`
	CreatedAt  string `json:"created_at"`
}

type ChangePasswordRequest struct {
	CustomerID  string `json:"customer_id" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// Then use ShouldBindJSON() for all fields
