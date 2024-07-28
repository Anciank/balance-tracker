package models

type Balance struct {
    ID        int     `json:"id"`
    UserID    int     `json:"user_id"`
    Amount    float64 `json:"amount"`
    CreatedAt string   `json:"created_at"`
    UpdatedAt string   `json:"updated_at"`
}

