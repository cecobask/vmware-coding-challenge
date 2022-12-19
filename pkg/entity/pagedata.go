package entity

import (
	"fmt"
	"github.com/go-chi/render"
	"net/http"
	"sort"
	"strconv"
)

type PageData struct {
	Url            string  `json:"url"`
	Views          int     `json:"views"`
	RelevanceScore float64 `json:"relevanceScore"`
}

type GetPageDataRequest struct {
	SortKey string
	Limit   int
}

type GetPageDataResponse struct {
	Data  []PageData          `json:"data"`
	Count int                 `json:"count"`
	Req   *GetPageDataRequest `json:"-"`
	Error error               `json:"-"`
}

const (
	limitMax              = 200
	limitMin              = 1
	queryLimit            = "limit"
	querySortKey          = "sortKey"
	sortKeyRelevanceScore = "relevanceScore"
	sortKeyViews          = "views"
)

var allowedSortKeys = map[string]bool{
	sortKeyRelevanceScore: true,
	sortKeyViews:          true,
}

func NewGetPageDataRequest(req *http.Request) (*GetPageDataRequest, error) {
	if req.Method != http.MethodGet {
		details := fmt.Sprintf("the only allowed http method for this endpoint is %s", http.MethodGet)
		return nil, ErrorBadRequest(details)
	}
	sortKeyQuery, ok := req.URL.Query()[querySortKey]
	if !ok || len(sortKeyQuery) != 1 {
		details := fmt.Sprintf("the query parameter %s is required and needs exactly one value", querySortKey)
		return nil, ErrorBadRequest(details)
	}
	if !allowedSortKeys[sortKeyQuery[0]] {
		details := fmt.Sprintf("the query parameter %s can only be one of the following - %s / %s", querySortKey, sortKeyRelevanceScore, sortKeyViews)
		return nil, ErrorBadRequest(details)
	}
	limitQuery, ok := req.URL.Query()[queryLimit]
	if !ok || len(limitQuery) != 1 {
		details := fmt.Sprintf("the query parameter %s is required and needs exactly one value", queryLimit)
		return nil, ErrorBadRequest(details)
	}
	limit, err := strconv.Atoi(limitQuery[0])
	if err != nil || limit >= limitMax || limit <= limitMin {
		details := fmt.Sprintf("the query parameter %s can only be of integer type, between 2 and 199", queryLimit)
		return nil, ErrorBadRequest(details)
	}
	return &GetPageDataRequest{
		SortKey: sortKeyQuery[0],
		Limit:   limit,
	}, nil
}

func (resp *GetPageDataResponse) Render(w http.ResponseWriter, r *http.Request) error {
	sort.Slice(resp.Data, func(a, b int) bool {
		var result bool
		switch resp.Req.SortKey {
		case sortKeyRelevanceScore:
			result = resp.Data[a].RelevanceScore < resp.Data[b].RelevanceScore
		case sortKeyViews:
			result = resp.Data[a].Views < resp.Data[b].Views
		}
		return result
	})
	if len(resp.Data) > resp.Req.Limit {
		resp.Data = resp.Data[:resp.Req.Limit]
		resp.Count = resp.Req.Limit
	}
	render.JSON(w, r, &resp)
	return nil
}
