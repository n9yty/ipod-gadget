package iap

import (
	"bytes"
	"encoding/binary"
	"log"
	"time"
)

type LingoCmdHandler func(IapPacket, chan<- IapPacket)

func StaticHandler(resp ...IapPacket) LingoCmdHandler {
	return func(input IapPacket, output chan<- IapPacket) {
		for _, r := range resp {
			output <- r
		}
	}
}

func Ack(in IapPacket) IapPacket {
	payload := bytes.Buffer{}
	if in.LingoCmdId.Id1 == 0x04 {
		payload.WriteByte(0x0)
		binary.Write(&payload, binary.BigEndian, in.LingoCmdId.Id2)
		return IapPacket{LingoCmdId{0x04, 0x0001}, payload.Bytes()}

	} else {
		payload.Write(in.Payload[:2])
		payload.Write([]byte{in.LingoCmdId.Id1, uint8(in.LingoCmdId.Id2)})
		return IapPacket{LingoCmdId{0x00, 0x02}, payload.Bytes()}
	}

}

// LINGO 0x00
func GetIpodOptionsForLingo(msg IapPacket, resp chan<- IapPacket) {
	resp <- msg
}

var lingoMap = map[LingoCmdId]LingoCmdHandler{

	//--------------
	//General Lingo
	//--------------

	//EnterRemoteUIMode
	LingoCmdId{0x00, 0x05}: func(in IapPacket, out chan<- IapPacket) {
		out <- Ack(in)
	},

	//ExitRemoteUIMode
	LingoCmdId{0x00, 0x06}: func(in IapPacket, out chan<- IapPacket) {
		out <- Ack(in)
	},

	//RequestiPodName
	LingoCmdId{0x00, 0x07}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		payload.WriteString("Ipod Gadget")
		payload.WriteByte(0x00)

		//ReturniPodName
		out <- IapPacket{LingoCmdId{0x00, 0x08}, payload.Bytes()}
	},

	//RequestiPodSerialNum
	LingoCmdId{0x00, 0x0b}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		payload.WriteString("Serial")
		payload.WriteByte(0x00)

		//ReturnIpodSerialNum
		out <- IapPacket{LingoCmdId{0x00, 0x0c}, payload.Bytes()}
	},

	//RequestLingoProtocolVersion
	LingoCmdId{0x00, 0x0f}: StaticHandler(

		//ReturnLingoProtocolVersion
		// 0x0a = Lingo for this information (Digital Audio Lingo)
		// 0x01 = Major Version of protocol
		// 0x02 = Minor Version of protocol
		IapPacket{LingoCmdId{0x00, 0x10}, []byte{0x0a, 0x01, 0x02}},
	),

	//RequestTransportMaxPayloadSize
	LingoCmdId{0x00, 0x11}: func(in IapPacket, out chan<- IapPacket) {
		
		//ReturnTransportMaxPayloadSize - 0xFFF9 = 65529
		out <- IapPacket{LingoCmdId{0x00, 0x12}, []byte{in.Payload[0], in.Payload[1], 0x01, 0xfa}}
	},

	//IdentifyDeviceLingoes
	LingoCmdId{0x00, 0x13}: func(in IapPacket, out chan<- IapPacket) {
		//Ack
		out <- Ack(in)

		//GetDevAuthenticationInfo
		out <- IapPacket{LingoCmdId{0x00, 0x14}, []byte{}}
	},

	//RetDevAuthenticationInfo
	LingoCmdId{0x00, 0x15}: func(in IapPacket, out chan<- IapPacket) {
		out <- Ack(in)
		if in.Payload[4] == in.Payload[5] {
			//AckDevAuthenticationInfo
			out <- IapPacket{LingoCmdId{0x00, 0x16}, []byte{in.Payload[0], in.Payload[1], 0x00}}
			//GetDevAuthenticationSignature
			out <- IapPacket{LingoCmdId{0x00, 0x17}, []byte{in.Payload[0], in.Payload[1], 0x7F, 0x59, 0x27, 0x9B, 0x9D, 0x3B, 0x80, 0xE4, 0x63, 0x5C, 0xBB, 0x0E, 0xEF, 0xA5, 0x28, 0x96, 0x5E, 0x30, 0x33, 0x05, 0x01}}
		}
	},

	//RetDevAuthenticationSignature
	LingoCmdId{0x00, 0x18}: func(in IapPacket, out chan<- IapPacket) {
		//AckDevAuthenticationStatus
		out <- IapPacket{LingoCmdId{0x00, 0x19}, []byte{in.Payload[0], in.Payload[1], 0x00}}

		//GetAccSampleRateCaps
		out <- IapPacket{LingoCmdId{0x0a, 0x02}, []byte{in.Payload[0], in.Payload[1] + 1}}
	},

	//StartIDPS
	LingoCmdId{0x00, 0x38}: func(in IapPacket, out chan<- IapPacket) {
		//Ack
		out <- Ack(in)
	},

	//SetFIDTokenValues
	LingoCmdId{0x00, 0x39}: func(in IapPacket, out chan<- IapPacket) {
		//RetFIDTokenValuesACKs
		out <- IapPacket{LingoCmdId{0x00, 0x3a}, []byte{in.Payload[0], in.Payload[1], 0x08, 0x03, 0x00, 0x00, 0x00, 0x03, 0x00, 0x01, 0x00, 0x04, 0x00, 0x02, 0x00, 0x01, 0x04, 0x00, 0x02, 0x00, 0x04, 0x04, 0x00, 0x02, 0x00, 0x05, 0x04, 0x00, 0x02, 0x00, 0x06, 0x04, 0x00, 0x02, 0x00, 0x07, 0x04, 0x00, 0x02, 0x00, 0x0C}}
	},

	//EndIDPS
	LingoCmdId{0x00, 0x3b}: func(in IapPacket, out chan<- IapPacket) {
		//IDSPStatus
		out <- IapPacket{LingoCmdId{0x00, 0x3c}, []byte{in.Payload[0], in.Payload[1], 0x00}}
		
		//GetDevAuthenticationInfo
		out <- IapPacket{LingoCmdId{0x00, 0x14}, []byte{0x03, 0x0B}}
	},

	//SetEventNotification
	LingoCmdId{0x00, 0x49}: func(in IapPacket, out chan<- IapPacket) {
		//Ack
		out <- Ack(in)
	},

	//GetiPodOptionsForLingo
	LingoCmdId{0x00, 0x4b}: func(in IapPacket, out chan<- IapPacket) {
		//RetiPodOptionsForLingo
		if in.Payload[2] == 0x00 {
			out <- IapPacket{LingoCmdId{0x00, 0x4c}, []byte{in.Payload[0], in.Payload[1], in.Payload[2], 0x00, 0x00, 0x00, 0x00, 0x58, 0xA8, 0x53, 0x7f}}
		} else if in.Payload[2] == 0x02 {
			out <- IapPacket{LingoCmdId{0x00, 0x4c}, []byte{in.Payload[0], in.Payload[1], in.Payload[2], 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xB3}}
		} else if in.Payload[2] == 0x03 {
			out <- IapPacket{LingoCmdId{0x00, 0x4c}, []byte{in.Payload[0], in.Payload[1], in.Payload[2], 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03}}
		} else if in.Payload[2] == 0x04 {
			out <- IapPacket{LingoCmdId{0x00, 0x4c}, []byte{in.Payload[0], in.Payload[1], in.Payload[2], 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x26}}
		} else if in.Payload[2] == 0x0A {
			out <- IapPacket{LingoCmdId{0x00, 0x4c}, []byte{in.Payload[0], in.Payload[1], in.Payload[2], 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}}
		} else {
			log.Printf("Unknown lingo requesting options %d \n", in.Payload[2])
			out <- IapPacket{LingoCmdId{0x00, 0x4c}, []byte{in.Payload[0], in.Payload[1], in.Payload[2], 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}
		}
		
	},
	
	//GetSupportedEventNotification
	LingoCmdId{0x00, 0x4f}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		binary.Write(&payload, binary.BigEndian, uint64(0x30AC)) //Flow Control bit set
		
		//RetSupportedEventNotification
		out <- IapPacket{LingoCmdId{0x00, 0x51}, payload.Bytes()}
	},

	//-------------------
	//Digital Audio Lingo
	//-------------------

	//RetAccSampleRateCaps
	LingoCmdId{0x0a, 0x03}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		binary.Write(&payload, binary.BigEndian, uint32(44100)) //44.1 kHz
		binary.Write(&payload, binary.BigEndian, uint32(0))
		binary.Write(&payload, binary.BigEndian, uint32(0))
		out <- IapPacket{LingoCmdId{0x0a, 0x04}, payload.Bytes()}
	},

	//--------------
	//Extended Lingo
	//--------------

	//Reserved
	LingoCmdId{0x04, 0x00}: func(in IapPacket, out chan<- IapPacket) {
	},
	
	//ResetDBSelection
	LingoCmdId{0x04, 0x0016}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		payload.WriteByte(0x0)
		binary.Write(&payload, binary.BigEndian, uint16(in.LingoCmdId.Id2))

		//Extended Lingo Ack
		out <- IapPacket{LingoCmdId{0x04, 0x0001}, payload.Bytes()}
	},

	//SelectDBRecord
	LingoCmdId{0x04, 0x0017}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		payload.WriteByte(0x0)
		binary.Write(&payload, binary.BigEndian, uint16(in.LingoCmdId.Id2))

		//Extended Lingo Ack
		out <- IapPacket{LingoCmdId{0x04, 0x0001}, payload.Bytes()}
	},

	//GetNumberCategorizedDBRecords
	LingoCmdId{0x04, 0x0018}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		binary.Write(&payload, binary.BigEndian, uint32(1))

		//ReturnNumberCategorizedDBRecords
		out <- IapPacket{LingoCmdId{0x04, 0x0019}, payload.Bytes()}
	},

	//RetrieveCategorizedDatabaseRecord
	LingoCmdId{0x04, 0x001A}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		binary.Write(&payload, binary.BigEndian, uint32(0))
		payload.WriteString("Testing")
		payload.WriteByte(0x0)

		//ReturnCategorizedDatabaseRecord
		out <- IapPacket{LingoCmdId{0x04, 0x001B}, payload.Bytes()}
	},

	//GetPlayStatus
	LingoCmdId{0x04, 0x001C}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		binary.Write(&payload, binary.BigEndian, uint32(1000*60*5)) //Track Length
		binary.Write(&payload, binary.BigEndian, uint32(1000*60*2)) //Track Position
		
		//Play states
		// 0x00 = Stopped
		// 0x01 = Playing
		// 0x02 = Paused
		// 0x03 - 0xFE = Reserved
		// 0xFF = Error
		payload.WriteByte(0x01)

		//ReturnPlayStatus
		out <- IapPacket{LingoCmdId{0x04, 0x001D}, payload.Bytes()}
	},

	//GetCurrentPlayingTrackIndex
	LingoCmdId{0x04, 0x001E}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		binary.Write(&payload, binary.BigEndian, uint32(0)) //Track index zero
		
		//ReturnCurrentPlayingTrackIndex
		out <- IapPacket{LingoCmdId{0x04, 0x001F}, payload.Bytes()}
	},

	//GetIndexedPlayingTrackTitle
	LingoCmdId{0x04, 0x0020}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})		
		payload.WriteString("Testing Now")
		payload.WriteByte(0x0)

		//ReturnIndexedPlayingTrackTitle
		out <- IapPacket{LingoCmdId{0x04, 0x0021}, payload.Bytes()}
	},

	//GetIndexedPlayingTrackArtistName
	LingoCmdId{0x04, 0x0022}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		payload.WriteString("Testing artist")
		payload.WriteByte(0x00)

		//ReturnIndexedPlayingTrackArtistName
		out <- IapPacket{LingoCmdId{0x04, 0x0023}, payload.Bytes()}
	},

	//GetIndexedPlayingTrackAlbumName
	LingoCmdId{0x04, 0x0024}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		payload.WriteString("Testing album")
		payload.WriteByte(0x00)

		//ReturnIndexedPlayingTrackAlbumName
		out <- IapPacket{LingoCmdId{0x04, 0x0025}, payload.Bytes()}
	},

	//SetPlayStatusChangeNotification
	LingoCmdId{0x04, 0x0026}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		payload.WriteByte(0x00)
		binary.Write(&payload, binary.BigEndian, uint16(in.LingoCmdId.Id2))

		//Extended Lingo Ack
		out <- IapPacket{LingoCmdId{0x04, 0x0001}, payload.Bytes()}
	},

	//Play Selection
	LingoCmdId{0x04, 0x0028}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		payload.WriteByte(0x0)
		binary.Write(&payload, binary.BigEndian, uint16(in.LingoCmdId.Id2))

		//Extended Lingo Ack
		out <- IapPacket{LingoCmdId{0x04, 0x0001}, payload.Bytes()}
	},

	//Play control
	LingoCmdId{0x04, 0x0029}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		payload.WriteByte(0x0)
		binary.Write(&payload, binary.BigEndian, uint16(in.LingoCmdId.Id2))

		//Extended Lingo Ack
		out <- IapPacket{LingoCmdId{0x04, 0x0001}, payload.Bytes()}

		//Update Track Position with PlayStatusChangeNotification
		go func() {
			time.Sleep(1 * time.Second)
			out <- IapPacket{LingoCmdId{0x04, 0x0027}, []byte{0x06, 0x0A}}

			for i := 0; i < 20; i++ {
				payload := bytes.Buffer{}				
				payload.WriteByte(0x04) // 0x04 = Track time offset (milliseconds)
				binary.Write(&payload, binary.BigEndian, uint32(1000*(60*2+i)))
				out <- IapPacket{LingoCmdId{0x04, 0x0027}, payload.Bytes()}
				time.Sleep(1 * time.Second)
			}
		}()

	},

	//Get Shuffle
	LingoCmdId{0x04, 0x002C}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		
		// Shuffle modes
		// 0x00 = Shuffle Off
		// 0x01 = Shuffle Tracks
		// 0x02 = Shuffle Albums
		payload.WriteByte(0x0)

		//ReturnShuffle
		out <- IapPacket{LingoCmdId{0x04, 0x002D}, payload.Bytes()}
	},

	//Set Shuffle
	LingoCmdId{0x04, 0x002E}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		payload.WriteByte(0x0)
		binary.Write(&payload, binary.BigEndian, uint16(in.LingoCmdId.Id2))

		//Extended Lingo Ack
		out <- IapPacket{LingoCmdId{0x04, 0x0001}, payload.Bytes()}
	},

	//Get Repeat
	LingoCmdId{0x04, 0x002F}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})

		// Repeat modes
		// 0x00 = Repeat Off
		// 0x01 = Repeat One Track
		// 0x02 = Repeat All Tracks
		payload.WriteByte(0x0)

		//Extended Lingo Ack
		out <- IapPacket{LingoCmdId{0x04, 0x0030}, payload.Bytes()}
	},

	//Set Repeat
	LingoCmdId{0x04, 0x0031}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		payload.WriteByte(0x0)
		binary.Write(&payload, binary.BigEndian, uint16(in.LingoCmdId.Id2))

		//Extended Lingo Ack
		out <- IapPacket{LingoCmdId{0x04, 0x0001}, payload.Bytes()}
	},

	//SetDisplayImage
	LingoCmdId{0x04, 0x0032}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		payload.WriteByte(0x00)
		binary.Write(&payload, binary.BigEndian, uint16(in.LingoCmdId.Id2))

		//Extended Lingo Ack
		out <- IapPacket{LingoCmdId{0x04, 0x0001}, payload.Bytes()}
	},
	
	//GetMonoDisplayImageLimits
	LingoCmdId{0x04, 0x0033}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		binary.Write(&payload, binary.BigEndian, uint16(4))	//Width
		binary.Write(&payload, binary.BigEndian, uint16(4)) //Heigh
		payload.WriteByte(0x01) //Monochrome, 2 bits per pixel

		//ReturnMonoDisplayImageLimits
		out <- IapPacket{LingoCmdId{0x04, 0x0034}, payload.Bytes()}
	},

	//GetNumPlayingTracks
	LingoCmdId{0x04, 0x0035}: func(in IapPacket, out chan<- IapPacket) {
		payload := bytes.Buffer{}
		payload.Write([]byte{in.Payload[0], in.Payload[1]})
		binary.Write(&payload, binary.BigEndian, uint32(1)) //One Track

		//ReturnNumPlayingTracks
		out <- IapPacket{LingoCmdId{0x04, 0x0036}, payload.Bytes()}
	},
}

func Route(input <-chan IapPacket, output chan<- IapPacket) {
	for inputMsg := range input {
		if handler, ok := lingoMap[inputMsg.LingoCmdId]; ok {
			handler(inputMsg, output)
		} else {
			log.Printf("No handler for %02X %02X", inputMsg.LingoCmdId.Id1, inputMsg.LingoCmdId.Id2)
		}
	}
	close(output)
}
