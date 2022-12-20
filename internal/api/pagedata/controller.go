package pagedata

import (
	"encoding/json"
	"fmt"
	"github.com/cecobask/vmware-coding-challenge/pkg/entity"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	retryDelay    = time.Second * 3
	retryAttempts = 5
)

type Controller interface {
	GetPageData(req *entity.GetPageDataRequest) (*entity.GetPageDataResponse, error)
}

func NewController(logger *zap.Logger) Controller {
	return &controller{
		logger: logger,
	}
}

func (c *controller) GetPageData(req *entity.GetPageDataRequest) (*entity.GetPageDataResponse, error) {
	pdrChan := make(chan entity.GetPageDataResponse)
	var wg sync.WaitGroup
	for _, url := range urls() {
		wg.Add(1)
		go c.getPageData(url, pdrChan, &wg)
	}
	go func() {
		wg.Wait()
		close(pdrChan)
	}()
	resp := &entity.GetPageDataResponse{Req: req}
	foundErr := false
	for pdr := range pdrChan {
		if pdr.Error != nil {
			foundErr = true
			continue
		}
		resp.Data = append(resp.Data, pdr.Data...)
	}
	resp.Count = len(resp.Data)
	if len(resp.Data) == 0 && foundErr {
		return nil, entity.ErrorInternalServer("all http requests failed to yield results")
	}
	c.logger.Info(fmt.Sprintf("retrieved %d data entries before sort and limit were applied", resp.Count))
	return resp, nil
}

func urls() []string {
	return []string{
		"https://raw.githubusercontent.com/assignment132/assignment/main/duckduckgo.json",
		"https://raw.githubusercontent.com/assignment132/assignment/main/google.json",
		"https://raw.githubusercontent.com/assignment132/assignment/main/wikipedia.json",
	}
}

func (c *controller) getPageData(url string, pdrChan chan entity.GetPageDataResponse, wg *sync.WaitGroup) {
	defer wg.Done()
	var pdr entity.GetPageDataResponse
	for i := 0; i < retryAttempts; i++ {
		if i > 0 {
			c.logger.Info("retrying http request", zap.String("method", http.MethodGet), zap.String("url", url), zap.Int("attempt", i))
			time.Sleep(retryDelay)
		}
		pdr = doRequest(url)
		if pdr.Error != nil {
			c.logger.Error("error processing http request", zap.String("method", http.MethodGet), zap.String("url", url))
			continue
		}
		break
	}
	pdrChan <- pdr
}

func doRequest(url string) (pdr entity.GetPageDataResponse) {
	resp, err := http.Get(url)
	if err != nil {
		pdr.Error = fmt.Errorf("error issuing http get request")
		return pdr
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		pdr.Error = fmt.Errorf("received non-healthy http status code %d", resp.StatusCode)
		return pdr
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		pdr.Error = fmt.Errorf("error while reading http response: %w", err)
		return pdr
	}
	if err = json.Unmarshal(body, &pdr); err != nil {
		pdr.Error = fmt.Errorf("error while unmarshaling http response: %w", err)
		return pdr
	}
	return pdr
}

type controller struct {
	logger *zap.Logger
}

type controllerMock struct {
	GetPageDataFn func(req *entity.GetPageDataRequest) (*entity.GetPageDataResponse, error)
}

func (c *controllerMock) GetPageData(req *entity.GetPageDataRequest) (*entity.GetPageDataResponse, error) {
	if c.GetPageDataFn != nil {
		return c.GetPageDataFn(req)
	}
	return &entity.GetPageDataResponse{}, nil
}
