package videonetics

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/pion/rtp"
	"github.com/rs/zerolog/log"
	"github.com/vtpl1/vrtc3/pkg/core"
	"github.com/vtpl1/vrtc3/pkg/h264"
	"github.com/vtpl1/vrtc3/pkg/h264/annexb"
	"github.com/vtpl1/vrtc3/pkg/h265"
	pb "github.com/vtpl1/vrtc3/pkg/videonetics/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (c *Conn) Reconnect() error {
	log.Info().Msgf("[videonetics] Reconnect start")
	defer func() {
		log.Info().Msgf("[videonetics] Reconnect end")
	}()

	// close current session
	_ = c.close()

	// start new session
	if err := c.dial(); err != nil {
		return err
	}
	if err := c.describe(); err != nil {
		return err
	}
	return nil
}

func (c *Conn) close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *Conn) dial() (err error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(c.host, opts...)
	if err != nil {
		log.Err(err).Msgf("[%v] failed to dial for %v", c.host, c.channel)
		return
	}
	log.Info().Msgf("[%v] success to dial for %v", c.host, c.channel)
	c.stateMu.Lock()
	c.state = StateConn
	c.stateMu.Unlock()
	c.conn = conn
	return
}

func binSize(val int) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(val))
	return buf
}

// Get the fmtpline
func (c *Conn) describe() (err error) {
	log.Info().Msgf("[videonetics] Describe start")
	defer func() {
		log.Info().Msgf("[videonetics] Describe end")
	}()
	serviceClient := pb.NewStreamServiceClient(c.conn)
	stream, err := serviceClient.ReadFrame(*c.ctx, &pb.ReadFrameRequest{Channel: &pb.Channel{
		SiteId:     c.channel.SiteID,
		ChannelId:  c.channel.ChannelID,
		AppId:      c.channel.AppID,
		LiveOrRec:  c.channel.LiveOrRec,
		StreamType: c.channel.StreamType,
		StartTs:    c.channel.StartTS,
		SessionId:  c.channel.SessionID,
	}})
	if err != nil {
		log.Info().Msg("Failed to FrameRead 1: " + err.Error() + ", ")
		serviceClient = nil
		stream = nil
		return
	}
	c.stream = stream
	totalFrameReceived := 0
	var sps []byte = nil
	var pps []byte = nil
	var vps []byte = nil
dd:
	for {
		response, err := c.stream.Recv()
		if err != nil || response == nil {
			log.Info().Msg("Failed to FrameRead 3: " + err.Error() + ", ")
			c.stream = nil
			return err
		}
		if totalFrameReceived > 32 {
			log.Info().Msgf("[videonetics] sps pps vps not received yet %v", totalFrameReceived)
			return errors.New("sps pps vps not received yet")
		}
		totalFrameReceived++
		mediaType := response.GetFrame().GetMediaType()
		frameType := response.GetFrame().GetFrameType()
		buffer := response.GetFrame().GetBuffer()
		fmt.Printf("[videonetics] mediaType %v frameType %v\n", mediaType, frameType)
		switch mediaType {
		case 2:
			switch h264.NALUType(buffer) {
			case h264.NALUTypeSPS:
				sps = buffer[4:]
			case h264.NALUTypePPS:
				pps = buffer[4:]
			}
			if sps != nil && pps != nil {
				// sps_pps := append(sps, pps...)
				// fmt.Printf("sps: %v pps: %v sps_pps: %v len %v", sps, pps, sps_pps, len(sps_pps))
				avccbuffer := append(binSize(len(sps)), sps...)
				avccbuffer = append(avccbuffer, binSize(len(pps))...)
				avccbuffer = append(avccbuffer, pps...)

				codec := h264.AVCCToCodec(avccbuffer)
				codec.PayloadType = 0

				c.Medias = append(c.Medias, &core.Media{
					Kind:      core.KindVideo,
					Direction: core.DirectionRecvonly,
					Codecs: []*core.Codec{
						codec,
					},
				})
				log.Info().Msgf("[videonetics] Codec H264 %v %v", codec, codec.FmtpLine)
				break dd
			}
		case 8:
			switch h265.NALUType(buffer) {
			case h265.NALUTypeSPS:
				sps = buffer[4:]
			case h265.NALUTypePPS:
				pps = buffer[4:]
			case h265.NALUTypeVPS:
				vps = buffer[4:]
			}
			if sps != nil && pps != nil && vps != nil {
				avccbuffer := append(binSize(len(sps)), sps...)
				avccbuffer = append(avccbuffer, binSize(len(pps))...)
				avccbuffer = append(avccbuffer, pps...)
				avccbuffer = append(avccbuffer, binSize(len(vps))...)
				avccbuffer = append(avccbuffer, vps...)

				codec := h265.AVCCToCodec(avccbuffer)
				codec.PayloadType = 0

				c.Medias = append(c.Medias, &core.Media{
					Kind:      core.KindVideo,
					Direction: core.DirectionRecvonly,
					Codecs: []*core.Codec{
						codec,
					},
				})
				log.Info().Msgf("[videonetics] Codec H265 %v %v", codec, codec.FmtpLine)
				break dd
			}
		}

		// c.Medias = append(c.Medias, )
		// if totalFrameReceived > 32 {
		// 	break
		// }
	}

	return
}

// GetMedias implements core.Producer.
// Subtle: this method shadows the method (Connection).GetMedias of Conn.Connection.
func (c *Conn) GetMedias() []*core.Media {
	log.Info().Msgf("[videonetics] Medias: %v", c.Medias)
	return c.Medias
}

// GetTrack implements core.Producer.
// Subtle: this method shadows the method (Connection).GetTrack of Conn.Connection.
func (c *Conn) GetTrack(media *core.Media, codec *core.Codec) (*core.Receiver, error) {
	core.Assert(media.Direction == core.DirectionRecvonly)
	log.Info().Msgf("[videonetics] GetTrack start %v %v", media, codec)
	defer func() {
		log.Info().Msgf("[videonetics] GetTrack end %v %v", media, codec)
	}()
	for _, track := range c.Receivers {
		if track.Codec == codec {
			return track, nil
		}
	}

	c.stateMu.Lock()
	defer c.stateMu.Unlock()

	if c.state == StatePlay {
		if err := c.Reconnect(); err != nil {
			return nil, err
		}
	}

	track := core.NewReceiver(media, codec)
	var i int = 0
	track.ID = byte(i)
	c.Receivers = append(c.Receivers, track)

	for _, receiver := range c.Receivers {
		c.handler = func(packet *rtp.Packet) {
			receiver.WriteRTP(packet)
		}
		if receiver.Codec.Name == "H264" {
			c.handler = h264.RTPPay(1412, c.handler)
		} else if receiver.Codec.Name == "H265" {
			c.handler = h265.RTPPay(1412, c.handler)
		}
	}
	return track, nil
}

// Start implements core.Producer.
func (c *Conn) Start() (err error) {
	log.Info().Msgf("[videonetics] Start start")
	defer func() {
		log.Info().Msgf("[videonetics] Start end")
	}()

	ok := false
	c.stateMu.Lock()
	switch c.state {
	case StateNone:
		err = nil
	case StateConn:
		c.state = StatePlay
		ok = true
	}
	c.stateMu.Unlock()

	if !ok {
		return
	}
	log.Info().Msgf("[videonetics] Handle Start %v", err)
	err = c.handle()
	log.Info().Msgf("[videonetics] Handle Return %v", err)
	return
}

func (c *Conn) handle() (err error) {
	return c.readFramePVA()
}

func (c *Conn) readFramePVA() (err error) {
	for {
		response, err := c.stream.Recv()
		if err != nil || response == nil {
			log.Info().Msg("Failed to FrameRead 3: " + err.Error() + ", ")
			c.stream = nil
			return err
		}
		// if response.GetFramePva().GetFrame().FrameType > 2 {
		// 	continue
		// }
		payload := response.GetFrame().Buffer
		size := len(payload)
		c.Recv += int(size)
		timeStamp := core.TimeStamp90000(response.GetFrame().Timestamp)
		// fmt.Printf("...... readFramePVA:  %d\n", size)
		// timeStamp := uint32(response.GetFramePva().GetFrame().Timestamp)
		packet := &rtp.Packet{
			Header:  rtp.Header{Timestamp: timeStamp, SSRC: 9582},
			Payload: annexb.EncodeToAVCC(payload),
		}
		// for _, receiver := range c.Receivers {
		// 	if receiver.Codec.Name == "H264" {
		// 		h264.RTPPay(1412, debug.Logger(func(packet *rtp.Packet) bool {
		// 			receiver.WriteRTP(packet)
		// 			return true
		// 		}))(packet)
		// 	} else if receiver.Codec.Name == "H265" {
		// 		h265.RTPPay(1412, debug.Logger(func(packet *rtp.Packet) bool {
		// 			receiver.WriteRTP(packet)
		// 			return true
		// 		}))(packet)
		// 	}
		// }
		c.handler(packet)
	}

}

// Stop implements core.Producer.
// Subtle: this method shadows the method (Connection).Stop of Conn.Connection.
func (c *Conn) Stop() (err error) {
	log.Info().Msgf("[videonetics] Stop start")
	defer func() {
		log.Info().Msgf("[videonetics] Stop end")
	}()
	for _, receiver := range c.Receivers {
		receiver.Close()
	}
	for _, sender := range c.Senders {
		sender.Close()
	}

	c.stateMu.Lock()
	if c.state != StateNone {
		c.state = StateNone
		err = c.close()
	}
	c.stateMu.Unlock()

	return
}
