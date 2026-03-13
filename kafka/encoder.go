package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/martketplace-vkr/pkg/logger/log"
	"google.golang.org/protobuf/proto"
)

type protoEncoder struct {
	msg     proto.Message
	b       []byte
	err     error
	encoded bool
}

func ProtoEncoder(msg proto.Message) sarama.Encoder {
	return &protoEncoder{
		msg: msg,
	}
}

func (p *protoEncoder) Encode() ([]byte, error) {
	if !p.encoded {
		p.b, p.err = proto.Marshal(p.msg)
		p.encoded = true
	}

	if p.err != nil {
		return nil, p.err
	}

	return p.b, nil
}

func (p *protoEncoder) Length() int {
	b, err := p.Encode()
	if err != nil {
		log.Error(fmt.Errorf("proto encode: %w", err))
		return 0
	}

	return len(b)
}
