/*
 * Revision History:
 *     Initial: 2018/07/05        Tong Yuehong
 */

package parser

import (
	"github.com/TechCatsLab/rumour"
	"github.com/TechCatsLab/rumour/pkg/message"
	"github.com/TechCatsLab/rumour/pkg/log"
	"errors"
)

type PacketParser struct {
}

// NewPacketParser new a packetParser.
func NewPacketParser() rumour.Parser {
	return &PacketParser{}
}

// Parse parse the packet.
func (p *PacketParser) Parse(data []byte) (rumour.Message, error) {
	var m message.Message

	err := m.Unmarshal(data)
	if err != nil {
		log.Error("[Parser packet] Parse err", log.Err(err))
		return nil, err
	}

	if m.From =="" || m.To == "" || m.Type == 0 {
		log.Error("[Parser Packet] Parser param err")
		return nil, errors.New("param err")
	}

	return &m, nil
}
