package model

import "time"

type CreateWalletRequest struct {
	CreatorId string
}

type GetWalletRequest struct {
	CreatorId string
}

type GetWalletIdRequest struct {
	CreatorId string
}

type WalletResponse struct {
	Id        string     `json:"id"`
	CreatorId string     `json:"creator_id"`
	Balance   int        `json:"balance"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
