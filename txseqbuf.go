package main

import "time"

// This value is sent to the transceiver and - according to my observations - it will use
// this as it's RX buf length. Note that if it is set to larger than 500-600ms then audio TX
// won't work (small radio memory?) - HA2NON
const txSeqBufLength = 300 * time.Millisecond

type txSeqBufEntry struct {
	seq     seqNum
	data    []byte
	addedAt time.Time
}

type txSeqBufStruct struct {
	entries []txSeqBufEntry
}

func (s *txSeqBufStruct) add(seq seqNum, p []byte) {
	s.entries = append(s.entries, txSeqBufEntry{
		seq:     seq,
		data:    p,
		addedAt: time.Now(),
	})
	s.purgeOldEntries()
}

func (s *txSeqBufStruct) purgeOldEntries() {
	// previous comment:
	// We keep much more entries than the specified length of the TX seqbuf, so we can serve
	// any requests coming from the server.
	//
	//  NOTE: need to wade deeper through the entire code base and see if the following is a more useful comment, or we need a better "why" here:
	//
	// Prune the oldest item in the txSeqBuf if it's older than ten time the txSeqBuf length... so "about 1 second"
	//  this enables serving any requests from the server, even if older than the max... up to a point

	for len(s.entries) > 0 && time.Since(s.entries[0].addedAt) > txSeqBufLength*10 {
		s.entries = s.entries[1:]
	}
}

func (s *txSeqBufStruct) get(seq seqNum) (d []byte) {
	if len(s.entries) == 0 {
		return nil
	}

	// Searching from backwards, as we expect most queries for latest entries.
	for i := len(s.entries) - 1; i >= 0; i-- {
		if s.entries[i].seq == seq {
			d = make([]byte, len(s.entries[i].data))
			copy(d, s.entries[i].data)
			break
		}
	}
	s.purgeOldEntries()
	return d
}
