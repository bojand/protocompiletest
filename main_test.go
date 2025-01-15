package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bufbuild/protocompile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	compiler := protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(&protocompile.SourceResolver{
			ImportPaths: []string{"./testdata/proto"},
		}),
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

	actualJSONData, err := marshalJSON(actualMsg, registry)
	assert.NoError(t, err)
	assert.NotEmpty(t, actualJSONData)
	fmt.Println("!!! ACTUAL JSON DATA:")
	fmt.Println(string(actualJSONData))

	// actual

	orderCreatedAt := time.Date(2023, time.July, 15, 10, 0, 0, 0, time.UTC)
	orderUpdatedAt := time.Date(2023, time.July, 15, 11, 0, 0, 0, time.UTC)

	expectedMsg := shopv2.Order{
		Version:       1,
		Id:            "444",
		CreatedAt:     timestamppb.New(orderCreatedAt),
		LastUpdatedAt: timestamppb.New(orderUpdatedAt),
	}

	expectData, err := proto.Marshal(&expectedMsg)
	require.NoError(t, err)

	assert.Equal(t, expectData, actualData)

	expectedJSONData, err := marshalJSON(&expectedMsg, registry)
	assert.NoError(t, err)
	assert.NotEmpty(t, expectedJSONData)
	fmt.Println("!!! EXPECTED JSON DATA:")
	fmt.Println(string(expectedJSONData))

	assert.Equal(t, expectedJSONData, actualJSONData)
}

func marshalJSON(m proto.Message, registry *protoregistry.Types) ([]byte, error) {
	marshaller := protojson.MarshalOptions{
		EmitDefaultValues: true,
		Resolver:          registry,
	}

	return marshaller.Marshal(m)
}
