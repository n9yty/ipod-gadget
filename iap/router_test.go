package iap

import "testing"
import "bytes"
import _ "encoding/binary"
import "log"

func TestRouter(t *testing.T) {
	inputMsg := make(chan IapPacket)
	outputMsg := make(chan IapPacket)

	doneChan := make(chan interface{})

	go func() {
		for resp := range outputMsg {
			report := BuildReport(resp)
			//log.Printf("Output %#v \n", report)

			buf := bytes.Buffer{}

			report.Ser(&buf)
			
			log.Printf("Len: %d %#v \n", buf.Len(), buf.Bytes())
			//t.Logf("%#v", buf.Bytes())

			//var report2 Report
			//report2.Deser(&buf)
			//t.Logf("Back: %#v", report2)

		}
		doneChan <- struct{}{}
	}()

	go Route(inputMsg, outputMsg)

	//go func() {
	//inputMsg <- IapPacket{LingoCmdId: LingoCmdId{0x00, 0x00}}
	/*
	inputMsg <- IapPacket{LingoCmdId: LingoCmdId{0x00, 0x38}, Payload: []uint8{0x00, 0x00, 0xC4}}
	inputMsg <- IapPacket{LingoCmdId: LingoCmdId{0x00, 0x11}, Payload: []uint8{0x00, 0x01, 0xEA}}
	inputMsg <- IapPacket{LingoCmdId: LingoCmdId{0x00, 0x4B}, Payload: []uint8{0x00, 0x02, 0x00, 0xAE}}
	inputMsg <- IapPacket{LingoCmdId: LingoCmdId{0x00, 0x4B}, Payload: []uint8{0x00, 0x03, 0x02, 0xAB}}
	inputMsg <- IapPacket{LingoCmdId: LingoCmdId{0x00, 0x4B}, Payload: []uint8{0x00, 0x04, 0x03, 0xA9}}
	inputMsg <- IapPacket{LingoCmdId: LingoCmdId{0x00, 0x4B}, Payload: []uint8{0x00, 0x05, 0x04, 0xA7}}
	inputMsg <- IapPacket{LingoCmdId: LingoCmdId{0x00, 0x4B}, Payload: []uint8{0x00, 0x06, 0x0A, 0xA0}}
	*/
	//inputMsg <- IapPacket{LingoCmdId: LingoCmdId{0x00, 0x39}, Payload: []uint8{0x00, 0x07}}
	//inputMsg <- IapPacket{LingoCmdId: LingoCmdId{0x00, 0x3B}, Payload: []uint8{0x00, 0x08, 0x00, 0xB7}}

	//inputMsg <- IapPacket{LingoCmdId: LingoCmdId{0x00, 0x15}, Payload: []uint8{0x03, 0x0B, 0x02, 0x00, 0x00, 0x07, 0x30}}

	//inputMsg <- IapPacket{LingoCmdId: LingoCmdId{0x00, 0x15}, Payload: []uint8{0x03, 0x0B, 0x02, 0x00, 0x07, 0x07, 0x8A}}

	inputMsg <- IapPacket{LingoCmdId: LingoCmdId{0x00, 0x4F}, Payload: []uint8{0x00, 0x09, 0xA4}}


	//inputMsg <- IapPacket{LingoCmdId: LingoCmdId{0x01, 0x4b}}
	//inputMsg <- IapPacket{LingoCmdId: LingoCmdId{0x00, 0x13}, Payload: []uint8{0x00, 0x01, 0x02}}
	close(inputMsg)

	<-doneChan
	//}()

}
