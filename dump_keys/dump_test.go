package main

import (
	"os"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
)

func TestCreateDump(t *testing.T) {
	dmp := DumpKeys{}
	lightningdir, _ = os.Getwd()
	bitcoinNet = &chaincfg.RegressionNetParams

	keys, err := dmp.Call()
	if err != nil {
		t.Errorf("no go: %s", err.Error())
	}

	if keys.(DumpKeysResult).Xpub != "tpubDAuGhPdBpUdigMh4vzHj8vP1d6nHj3najPGdA1umyBvfCDCotARa22t4kjMjDvZNcVQEJ68u9JN44zekTsqFxaExQNcdtn2kMbsnhNXLscc" {
		t.Error("xpub not created")
	}
}

func TestCreateWithPriv(t *testing.T) {
	dmp := DumpKeys{true}
	lightningdir, _ = os.Getwd()
	bitcoinNet = &chaincfg.RegressionNetParams

	keys, err := dmp.Call()
	if err != nil {
		t.Errorf("no go: %s\n", err.Error())
	}
	if keys.(DumpKeysResult).Xpub != "tpubDAuGhPdBpUdigMh4vzHj8vP1d6nHj3najPGdA1umyBvfCDCotARa22t4kjMjDvZNcVQEJ68u9JN44zekTsqFxaExQNcdtn2kMbsnhNXLscc" {
		t.Error("xpub not created")
	}
	if keys.(DumpKeysResult).Xpriv != "tprv8eDEYyawg6x3ntfH3Ld8jWiu45GMZibgA5fqsVsUYv8GMix3FmbyqYGCacbUGtKk5RLy3fXo4U6UuKVNSfnSV4R42HaMAYSzPnDs8DbgXaR" {
		t.Error("xpriv not created")
	}
}
