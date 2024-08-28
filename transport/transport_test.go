package transport_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/GnarloqGames/genesis-avalon-kit/proto"
	"github.com/GnarloqGames/genesis-avalon-kit/transport"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	pproto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TransportMock struct{}

func (tm *TransportMock) Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	if subj == "test" {
		resp := &proto.BuildResponse{
			Response: "OK",
		}

		rawResp, err := pproto.Marshal(resp)
		if err != nil {
			return nil, err
		}

		return &nats.Msg{
			Subject: subj,
			Data:    rawResp,
		}, nil
	}

	return nil, fmt.Errorf("invalid request")
}

func (tm *TransportMock) Publish(subj string, data []byte) error {
	return nil
}

func TestRequest(t *testing.T) {
	conn := &TransportMock{}

	req := &proto.BuildRequest{
		Header: &proto.RequestHeader{
			Timestamp: timestamppb.Now(),
		},
		Name: "test",
	}

	var res proto.BuildResponse

	rawRes, err := transport.Request(conn, "test", req, &res, 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, "test", rawRes.Subject)
	require.Equal(t, "OK", res.Response)
}
