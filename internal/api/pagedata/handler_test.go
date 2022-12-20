package pagedata

import (
	"github.com/cecobask/vmware-coding-challenge/pkg/entity"
	"github.com/cecobask/vmware-coding-challenge/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetPageData(t *testing.T) {
	type fields struct {
		controller Controller
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
	}{
		{
			name: "success getting page data",
			fields: fields{
				controller: &controllerMock{
					GetPageDataFn: func(req *entity.GetPageDataRequest) (*entity.GetPageDataResponse, error) {
						return &entity.GetPageDataResponse{
							Data: []entity.PageData{
								{
									Url:            "https://example.com",
									Views:          100,
									RelevanceScore: 0.5,
								},
							},
							Count: 1,
							Req: &entity.GetPageDataRequest{
								SortKey: entity.SortKeyViews,
								Limit:   10,
							},
						}, nil
					},
				},
			},
			args: args{
				r: httptest.NewRequest(http.MethodGet, "https://example.com?sortKey=views&limit=10", http.NoBody),
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "error validating http request",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "https://example.com", http.NoBody),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "error getting page data",
			fields: fields{
				controller: &controllerMock{
					GetPageDataFn: func(req *entity.GetPageDataRequest) (*entity.GetPageDataResponse, error) {
						return nil, entity.ErrorInternalServer("generic message")
					},
				},
			},
			args: args{
				r: httptest.NewRequest(http.MethodGet, "https://example.com?sortKey=views&limit=10", http.NoBody),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				controller: tt.fields.controller,
				logger:     logger.NewLogger(),
			}
			responseRecorder := httptest.NewRecorder()
			h.GetPageData(responseRecorder, tt.args.r)
			if responseRecorder.Code != tt.wantStatus {
				t.Errorf("GetPageData() got = %d, want %d", responseRecorder.Code, tt.wantStatus)
			}
		})
	}
}
