package entities

type CreateSaleTypeRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type IdResponse struct {
	Id int32 `json:"id"`
}

type CreateBranchDirectorBody struct {
	CreateUserBody
	BranchId int32 `json:"branch_id"`
}
