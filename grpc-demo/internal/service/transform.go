package service

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/ankorstore/yokai-showroom/grpc-demo/proto"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
)

type TransformTextServiceService struct {
	proto.UnimplementedTransformTextServiceServer
	config *config.Config
}

func NewTransformTextServiceService(cfg *config.Config) *TransformTextServiceService {
	return &TransformTextServiceService{
		config: cfg,
	}
}

func (s *TransformTextServiceService) TransformText(ctx context.Context, in *proto.TransformTextRequest) (*proto.TransformTextResponse, error) {
	transformedText := s.transform(in)

	log.CtxLogger(ctx).Debug().Msgf("TransformText: %s -> %s", in.Text, transformedText)

	return &proto.TransformTextResponse{
		Text: transformedText,
	}, nil
}

func (s *TransformTextServiceService) TransformAndSplitText(stream proto.TransformTextService_TransformAndSplitTextServer) error {
	logger := log.CtxLogger(stream.Context())

	for {
		req, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			logger.Debug().Msg("TransformTextAndSplit: end of rpc")

			return nil
		}

		if err != nil {
			logger.Error().Err(err).Msgf("TransformTextAndSplit: error while receiving: %v", err)
		}

		logger.Debug().Msgf("TransformTextAndSplit: -> %s", req.Text)

		split := strings.Split(s.transform(req), " ")

		for _, word := range split {
			err = stream.Send(&proto.TransformTextResponse{
				Text: word,
			})

			if err != nil {
				logger.Error().Err(err).Msgf("TransformTextAndSplit: error while sending: %v", err)

				return err
			}

			logger.Debug().Msgf("TransformTextAndSplit: <- %s", word)
		}
	}
}

func (s *TransformTextServiceService) transform(in *proto.TransformTextRequest) string {
	switch in.Transformer {
	case proto.Transformer_TRANSFORMER_UPPERCASE:
		return strings.ToUpper(in.Text)
	case proto.Transformer_TRANSFORMER_LOWERCASE:
		return strings.ToLower(in.Text)
	default:
		return in.Text
	}
}
