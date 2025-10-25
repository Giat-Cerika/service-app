package response

import "time"

type AdminResponse struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
}