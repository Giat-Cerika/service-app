package response

type BaseResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

type PaginationResponse struct {
    Success    bool        `json:"success"`
    Message    string      `json:"message"`
    Data       interface{} `json:"data,omitempty"`
    Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
    CurrentPage  int `json:"current_page"`
    PerPage      int `json:"per_page"`
    TotalData    int `json:"total_data"`
    TotalPages   int `json:"total_pages"`
    PreviousPage int `json:"previous_page"`
    NextPages    int `json:"next_pages"`
}