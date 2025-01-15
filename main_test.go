package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/reporter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	shopv2 "github.com/bojand/protocompiletest/testdata/proto/gen/shop/v2"
)

func TestPlainProto(t *testing.T) {

	ctx := context.Background()

	filePaths := []string{
		"shop/v2/order.proto",
	}

	logger := zaptest.NewLogger(t)

	errorReporter := func(err reporter.ErrorWithPos) error {
		position := err.GetPosition()
		logger.Warn("failed to parse proto file to descriptor",
			zap.String("file", position.Filename),
			zap.Int("line", position.Line),
			zap.Error(err))
		return nil
	}

	compiler := protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(&protocompile.SourceResolver{
			ImportPaths: []string{"./testdata/proto"},
		}),
		Reporter: reporter.NewReporter(errorReporter, nil),
	}

	compiledFiles, err := compiler.Compile(ctx, filePaths...)
	require.NoError(t, err)
	assert.NotEmpty(t, compiledFiles)

	registry := new(protoregistry.Types)
	for _, descriptor := range compiledFiles {
		mds := descriptor.Messages()

		for i := 0; i < mds.Len(); i++ {
			mt := dynamicpb.NewMessageType(mds.Get(i))
			registry.RegisterMessage(mt)
		}
	}

	// expected

	inputData := `{"version":1,"id":"444","createdAt":"2023-07-15T10:00:00Z","lastUpdatedAt":"2023-07-15T11:00:00Z"}`

	messageType, err := registry.FindMessageByURL("shop.v2.Order")
	require.NoError(t, err)
	require.NotEmpty(t, messageType)

	md := messageType.Descriptor()
	actualMsg := dynamicpb.NewMessage(md)
	unmarshaller := protojson.UnmarshalOptions{
		DiscardUnknown: false,
		Resolver:       registry,
	}
	err = unmarshaller.Unmarshal([]byte(inputData), actualMsg)
	assert.NoError(t, err)

	actualData, err := proto.Marshal(actualMsg)
	assert.NoError(t, err)
	assert.NotEmpty(t, actualData)

	jsonData, err := marshalJSON(actualMsg, registry)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)
	fmt.Println("!!! ACTUAL JSON DATA:")
	fmt.Println(string(jsonData))

	// actual

	orderCreatedAt := time.Date(2023, time.July, 15, 10, 0, 0, 0, time.UTC)
	orderUpdatedAt := time.Date(2023, time.July, 15, 11, 0, 0, 0, time.UTC)
	// orderDeliveredAt := time.Date(2023, time.July, 15, 12, 0, 0, 0, time.UTC)
	// orderCompletedAt := time.Date(2023, time.July, 15, 13, 0, 0, 0, time.UTC)

	expectedMsg := shopv2.Order{
		Version:       1,
		Id:            "444",
		CreatedAt:     timestamppb.New(orderCreatedAt),
		LastUpdatedAt: timestamppb.New(orderUpdatedAt),
		// DeliveredAt:   timestamppb.New(orderDeliveredAt),
		// CompletedAt:   timestamppb.New(orderCompletedAt),
	}

	expectData, err := proto.Marshal(&expectedMsg)
	require.NoError(t, err)

	assert.Equal(t, expectData, actualData)

	jsonData, err = marshalJSON(&expectedMsg, registry)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)
	fmt.Println("!!! EXPECTED JSON DATA:")
	fmt.Println(string(jsonData))
}

func marshalJSON(m proto.Message, registry *protoregistry.Types) ([]byte, error) {
	marshaller := protojson.MarshalOptions{
		EmitDefaultValues: true,
		Resolver:          registry,
	}

	return marshaller.Marshal(m)
}
