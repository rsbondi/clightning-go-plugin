package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/niftynei/glightning/glightning"
	"github.com/niftynei/glightning/jrpc2"
	"log"
	"strings"
)

type ListSendpaysExt struct {
	PaymentHashes []string `json:"payment_hashes"` // TODO: ? dropping bolt11 here, maybe add another call?
}

func (h *ListSendpaysExt) New() interface{} {
	return &ListSendpaysExt{}
}

func (h *ListSendpaysExt) Name() string {
	return "listsendpaysext"
}

// TODO: msat string values
// TODO: bolt11 ? why not in glightning ?
// TODO: maybe this should just call multiple rpc calls instead of db query?
func (h *ListSendpaysExt) Call() (jrpc2.Result, error) {
	dbpath := lightningdir + "/lightningd.sqlite3"
	db, err := sql.Open("sqlite3", dbpath)
	defer db.Close()
	if err != nil {
		log.Printf("cannot open database: %s", err.Error())
	}
	q := `SELECT id, status, lower(hex(destination)), msatoshi, lower(hex(payment_hash)), timestamp, 
		lower(hex(payment_preimage)), msatoshi_sent FROM payments 
		WHERE hex(payment_hash) COLLATE NOCASE in (?` + strings.Repeat(",?", len(h.PaymentHashes)-1) + ")"

	log.Printf("querying for payment_hashes %s", h.PaymentHashes)
	ihash := make([]interface{}, len(h.PaymentHashes))
	for i, v := range h.PaymentHashes {
		ihash[i] = v
	}

	rows, err := db.Query(q, ihash...)
	if err != nil {
		log.Printf("cannot execute query: %s", err.Error())
	}
	defer rows.Close()

	var result struct {
		Payments []glightning.SendPayFields `json:"payments"`
	}

	for rows.Next() {
		f := glightning.SendPayFields{}
		err = rows.Scan(&f.Id, &f.Status, &f.Destination, &f.MilliSatoshi, &f.PaymentHash, &f.CreatedAt, &f.PaymentPreimage, &f.MilliSatoshiSent)
		if err != nil {
			log.Printf("cannot read database row: %s", err.Error())
		}

		result.Payments = append(result.Payments, f)
	}

	return result, nil
}
