package pagedata

import (
	"fmt"
	"github.com/cecobask/vmware-coding-challenge/pkg/entity"
	"github.com/cecobask/vmware-coding-challenge/pkg/logger"
	"github.com/go-chi/render"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

var successResponse = entity.GetPageDataResponse{
	Data: []entity.PageData{
		{
			Url:            "https://example.com",
			Views:          10,
			RelevanceScore: 1,
		},
	},
	Count: 1,
}

func Test_doRequest(t *testing.T) {
	type args struct {
		statusCode int
		response   entity.GetPageDataResponse
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success doing http request",
			args: args{
				statusCode: http.StatusOK,
				response:   successResponse,
			},
		},
		{
			name: "error issuing http get request",
			args: args{
				statusCode: http.StatusMovedPermanently,
				response: entity.GetPageDataResponse{
					Error: fmt.Errorf("error issuing http get request"),
				},
			},
		},
		{
			name: "received non-healthy http status code",
			args: args{
				statusCode: http.StatusInternalServerError,
				response: entity.GetPageDataResponse{
					Error: fmt.Errorf("received non-healthy http status code %d", http.StatusInternalServerError),
				},
			},
		},
		{
			name: "error while unmarshalling http response",
			args: args{
				statusCode: http.StatusOK,
				response: entity.GetPageDataResponse{
					Error: fmt.Errorf("error while unmarshaling http response: json: cannot unmarshal string into Go value of type entity.GetPageDataResponse"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := mockHttpServer(tt.args.statusCode, tt.args.response)
			defer server.Close()
			if gotPdr := doRequest(server.URL); gotPdr.String() != tt.args.response.String() {
				t.Errorf("doRequest() = %v, want %v", gotPdr, tt.args.response)
			}
		})
	}
}

func Test_controller_getPageData(t *testing.T) {
	type args struct {
		pdrChan chan entity.GetPageDataResponse
		wg      func(pdrChan chan entity.GetPageDataResponse) *sync.WaitGroup
	}
	type config struct {
		statusCode int
		response   entity.GetPageDataResponse
	}
	waitGroupFn := func(pdrChan chan entity.GetPageDataResponse) *sync.WaitGroup {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			wg.Wait()
			close(pdrChan)
		}()
		return &wg
	}
	errorResponse := entity.GetPageDataResponse{
		Error: fmt.Errorf("received non-healthy http status code %d", http.StatusInternalServerError),
	}
	tests := []struct {
		name   string
		args   args
		config config
		want   entity.GetPageDataResponse
	}{
		{
			name: "success getting page data",
			args: args{
				pdrChan: make(chan entity.GetPageDataResponse),
				wg:      waitGroupFn,
			},
			config: config{
				statusCode: http.StatusOK,
				response:   successResponse,
			},
			want: successResponse,
		},
		{
			name: "error getting page data",
			args: args{
				pdrChan: make(chan entity.GetPageDataResponse),
				wg:      waitGroupFn,
			},
			config: config{
				statusCode: http.StatusInternalServerError,
				response:   errorResponse,
			},
			want: errorResponse,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := mockHttpServer(tt.config.statusCode, tt.config.response)
			defer server.Close()
			c := &controller{logger: logger.NewLogger()}
			go c.getPageData(server.URL, tt.args.pdrChan, tt.args.wg(tt.args.pdrChan))
			if gotPdr := <-tt.args.pdrChan; gotPdr.String() != tt.want.String() {
				t.Errorf("getPageData() = %v, want %v", gotPdr, tt.want)
			}
		})
	}
}

func Test_controller_GetPageData(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "success getting page data",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &controller{logger: logger.NewLogger()}
			got, err := c.GetPageData(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPageData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("GetPageData() got = %v, want non-nil response", got)
			}
		})
	}
}

func mockHttpServer(statusCode int, response entity.GetPageDataResponse) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if response.Error != nil && strings.Contains(response.Error.Error(), "error while unmarshaling http response") {
			render.JSON(w, r, "unmarshallable")
			return
		}
		render.Status(r, statusCode)
		render.JSON(w, r, response)
	}))
}
