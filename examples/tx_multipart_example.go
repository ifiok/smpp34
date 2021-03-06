package main

import (
	"crypto/rand"
	"fmt"
	smpp "github.com/mergenchik/smpp34"
	"math"
)

func main() {
	// connect and bind
	tx, err := smpp.NewTransmitter(
		"localhost",
		9000,
		5,
		smpp.Params{
			"system_type": "CMT",
			"system_id":   "hugo",
			"password":    "ggoohu",
		},
	)
	if err != nil {
		fmt.Println("Connection Err:", err)
		return
	}
	msg := "Very Long Message, 1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890END"
	sm_len := len(msg)
	fmt.Println("Message Bytes count:", sm_len)

	if sm_len > 140 {
		sm_runes := []byte(msg)
		rune_len := len(sm_runes)
		send_params := smpp.Params{smpp.DATA_CODING: smpp.ENCODING_DEFAULT, smpp.ESM_CLASS: smpp.ESM_CLASS_GSMFEAT_UDHI}
		total_parts := byte(int(math.Ceil(float64(rune_len) / 134)))
		partNum := 1
		uid := make([]byte, 1)
		_, err := rand.Read(uid)
		if err != nil {
			// fmt.Println("QuerySM error:", err)
			fmt.Println("Rand.Read error:", err)
			return
		}
		for i := 0; i < rune_len; i += 134 {
			start := i
			end := i + 134
			if end > rune_len {
				end = rune_len
			}
			part := []byte{0x05, 0x00, 0x03, uid[0], total_parts, byte(partNum)}
			part = append(part, []byte(sm_runes[start:end])...)
			fmt.Println("Part:", part)
			// Send SubmitSm
			seq, err := tx.SubmitSmEncoded("test", "test2", part, &send_params)
			// Pdu gen errors
			if err != nil {
				fmt.Println("SubmitSm err:", err)
			}
			// Should save this to match with message_id
			fmt.Println("seq:", seq)
			partNum++

		}

	} else {
		send_params := smpp.Params{}
		// Send SubmitSm
		seq, err := tx.SubmitSm("test", "test2", msg, &send_params)

		// Pdu gen errors
		if err != nil {
			fmt.Println("SubmitSm err:", err)
		}
		// Should save this to match with message_id
		fmt.Println("seq:", seq)

	}

	for {
		pdu, err := tx.Read() // This is blocking
		if err != nil {
			fmt.Println("Read Err:", err)
			break
		}

		// EnquireLinks are auto handles
		switch pdu.GetHeader().Id {
		case smpp.SUBMIT_SM_RESP:
			// message_id should match this with seq message
			fmt.Println("MSG ID:", pdu.GetField("message_id").Value())
			fmt.Printf("PDU Header: %v", pdu.GetHeader())
			fmt.Println()
		default:
			// ignore all other PDUs or do what you link with them
			fmt.Println("PDU ID:", pdu.GetHeader().Id)
		}
	}

	fmt.Println("ending...")
}
