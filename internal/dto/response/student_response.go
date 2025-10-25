package response

import "time"

type StudentResponse struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    NIM       string    `json:"nim"`
    Major     string    `json:"major"`
    CreatedAt time.Time `json:"created_at"`
}