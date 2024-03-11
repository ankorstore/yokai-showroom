package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/ankorstore/yokai-showroom/grpc-demo/proto"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/prometheus/client_golang/prometheus"
)

// TransformerCounter is a metrics collector that counts each transformer usage.
var TransformerCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "transformer_total",
		Help: "Total of TransformTextService transformer usage",
	},
	[]string{
		"transformer",
	},
)

// TransformTextServiceService is the gRPC server service for TransformTextService.
type TransformTextServiceService struct {
	proto.UnimplementedTransformTextServiceServer
	config *config.Config
}

// NewTransformTextServiceService returns a new [TransformTextServiceService] instance.
func NewTransformTextServiceService(cfg *config.Config) *TransformTextServiceService {
	return &TransformTextServiceService{
		config: cfg,
	}
}

// TransformText transforms the text provided in a [proto.TransformTextRequest] by applying the provided transformer.
func (s *TransformTextServiceService) TransformText(ctx context.Context, in *proto.TransformTextRequest) (*proto.TransformTextResponse, error) {
	ctx, span := trace.CtxTracerProvider(ctx).Tracer("TransformTextService").Start(ctx, "TransformText")
	defer span.End()

	transformedText := s.transform(in)

	log.CtxLogger(ctx).Info().Msgf("TransformText: %s -> %s", in.Text, transformedText)

	return &proto.TransformTextResponse{
		Text: transformedText,
	}, nil
}

// TransformAndSplitText splits the text provided in a streamed [proto.TransformTextRequest] in words, and transform each word by applying the provided transformer.
func (s *TransformTextServiceService) TransformAndSplitText(stream proto.TransformTextService_TransformAndSplitTextServer) error {
	ctx := stream.Context()

	ctx, span := trace.CtxTracerProvider(ctx).Tracer("TransformTextService").Start(ctx, "TransformAndSplitText")
	defer span.End()

	logger := log.CtxLogger(ctx)

	for {
		req, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			logger.Info().Msg("TransformTextAndSplit: end of rpc")

			return nil
		}

		if err != nil {
			logger.Error().Err(err).Msgf("TransformTextAndSplit: error while receiving: %v", err)
		}

		logger.Info().Msgf("TransformTextAndSplit: -> %s", req.Text)

		split := strings.Split(s.transform(req), " ")

		for _, word := range split {
			err = stream.Send(&proto.TransformTextResponse{
				Text: word,
			})

			if err != nil {
				logger.Error().Err(err).Msgf("TransformTextAndSplit: error while sending: %v", err)

				return err
			}

			span.AddEvent(fmt.Sprintf("send word: %s", word))

			logger.Info().Msgf("TransformTextAndSplit: <- %s", word)
		}
	}
}

//nolint:exhaustive
func (s *TransformTextServiceService) transform(in *proto.TransformTextRequest) string {
	switch in.Transformer {
	case proto.Transformer_TRANSFORMER_UPPERCASE:
		TransformerCounter.WithLabelValues("uppercase").Inc()

		return strings.ToUpper(in.Text)
	case proto.Transformer_TRANSFORMER_LOWERCASE:
		TransformerCounter.WithLabelValues("lowercase").Inc()

		return strings.ToLower(in.Text)
	default:
		TransformerCounter.WithLabelValues("unspecified").Inc()

		return in.Text
	}
}
