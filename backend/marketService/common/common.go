package common


type MarketStatus string

const (
    StatusDraft    MarketStatus = "draft"
    StatusOpen     MarketStatus = "open"
    StatusClosed   MarketStatus = "closed"
    StatusResolved MarketStatus = "resolved"
)