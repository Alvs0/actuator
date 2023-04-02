package impl

import (
	"bytes"
	"context"
	"fmt"
	processor "github.com/Alvs0/actuator/service/processor/api"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net/http"
)

type ProcessorService struct {
	ProcessorQuery ProcessorQuery

	processor.UnimplementedProcessorServer
}

func Authorize(ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	var values []string
	var token string

	requestMetadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("[Generator] Missing authorization token")
	}

	values = requestMetadata.Get("Authorization")
	if values != nil && len(values) > 0 {
		token = values[0]
	}

	httpReq, err := http.NewRequest("POST", "http://localhost:8082/validate", bytes.NewBuffer(nil))
	if err != nil {
		log.Warn("[Generator] failed creating http request to validate token")
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

	httpClient := &http.Client{}
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("[Generator] error validating token. cause: %v", err.Error())
	}

	defer httpResp.Body.Close()
	if httpResp.StatusCode == http.StatusBadRequest || httpResp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("[Generator] unauthorized")
	}

	return handler(ctx, req)
}
