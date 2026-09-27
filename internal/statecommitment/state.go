package statecommitment

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"
	"sort"
)

const (
	Version                uint8 = 1
	checkpointMagic              = "SYJSTATE1"
	maxCheckpointAccounts        = 1_000_000
	maxCheckpointFieldSize       = 1 << 20
)

type Snapshot struct {
	Balances      map[string]uint64
	Nonces        map[string]uint64
	GenesisSupply uint64
	MiningIssued  uint64
	Supply        uint64
}

func NewSnapshot() *Snapshot {
	return &Snapshot{
		Balances: make(map[string]uint64),
		Nonces:   make(map[string]uint64),
	}
}

func (s *Snapshot) Clone() *Snapshot {
	out := NewSnapshot()

	out.GenesisSupply = s.GenesisSupply
	out.MiningIssued = s.MiningIssued
	out.Supply = s.Supply

	for k, v := range s.Balances {
		out.Balances[k] = v
	}

	for k, v := range s.Nonces {
		out.Nonces[k] = v
	}

	return out
}

func (s *Snapshot) Root() string {
	h := sha256.New()

	h.Write([]byte("SYJ-STATE-ROOT-V1\x00"))
	h.Write([]byte{Version})

	writeU64(h, s.GenesisSupply)
	writeU64(h, s.MiningIssued)
	writeU64(h, s.Supply)

	seen := make(map[string]struct{}, len(s.Balances)+len(s.Nonces))

	for k := range s.Balances {
		seen[k] = struct{}{}
	}

	for k := range s.Nonces {
		seen[k] = struct{}{}
	}

	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	writeU64(h, uint64(len(keys)))

	for _, k := range keys {
		h.Write([]byte("SYJ-STATE-ACCOUNT-V1\x00"))
		writeBytes(h, []byte(k))
		writeU64(h, s.Balances[k])
		writeU64(h, s.Nonces[k])
	}

	return fmtHex(h.Sum(nil))
}

func Encode(s *Snapshot) []byte {
	var b bytes.Buffer

	b.WriteString(checkpointMagic)
	b.WriteByte(Version)

	writeU64(&b, s.GenesisSupply)
	writeU64(&b, s.MiningIssued)
	writeU64(&b, s.Supply)

	seen := make(map[string]struct{}, len(s.Balances)+len(s.Nonces))

	for k := range s.Balances {
		seen[k] = struct{}{}
	}

	for k := range s.Nonces {
		seen[k] = struct{}{}
	}

	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	writeU64(&b, uint64(len(keys)))

	for _, k := range keys {
		writeBytes(&b, []byte(k))
		writeU64(&b, s.Balances[k])
		writeU64(&b, s.Nonces[k])
	}

	b.WriteString(s.Root())

	return b.Bytes()
}

func Decode(data []byte) (*Snapshot, error) {
	r := bytes.NewReader(data)

	magic := make([]byte, len(checkpointMagic))
	if _, err := io.ReadFull(r, magic); err != nil {
		return nil, errors.New("invalid state checkpoint magic")
	}

	if string(magic) != checkpointMagic {
		return nil, errors.New("invalid state checkpoint magic")
	}

	v, err := r.ReadByte()
	if err != nil {
		return nil, errors.New("unsupported state checkpoint version")
	}

	if v != Version {
		return nil, errors.New("unsupported state checkpoint version")
	}

	s := NewSnapshot()

	if s.GenesisSupply, err = readU64(r); err != nil {
		return nil, err
	}

	if s.MiningIssued, err = readU64(r); err != nil {
		return nil, err
	}

	if s.Supply, err = readU64(r); err != nil {
		return nil, err
	}

	n, err := readU64(r)
	if err != nil {
		return nil, err
	}

	if n > maxCheckpointAccounts {
		return nil, errors.New("state checkpoint account count exceeds limit")
	}

	for i := uint64(0); i < n; i++ {
		k, err := readBytes(r)
		if err != nil {
			return nil, err
		}

		key := string(k)

		if _, exists := s.Balances[key]; exists {
			return nil, errors.New("duplicate account in state checkpoint")
		}

		if _, exists := s.Nonces[key]; exists {
			return nil, errors.New("duplicate account in state checkpoint")
		}

		bal, err := readU64(r)
		if err != nil {
			return nil, err
		}

		nonce, err := readU64(r)
		if err != nil {
			return nil, err
		}

		s.Balances[key] = bal
		s.Nonces[key] = nonce
	}

	if r.Len() != 64 {
		return nil, errors.New("state checkpoint missing root checksum")
	}

	root := make([]byte, 64)
	if _, err := io.ReadFull(r, root); err != nil {
		return nil, err
	}

	if string(root) != s.Root() {
		return nil, errors.New("state checkpoint root mismatch")
	}

	return s, nil
}

func writeU64(w interface{ Write([]byte) (int, error) }, v uint64) {
	var b [8]byte

	binary.BigEndian.PutUint64(b[:], v)

	_, _ = w.Write(b[:])
}

func writeBytes(w interface{ Write([]byte) (int, error) }, v []byte) {
	writeU64(w, uint64(len(v)))
	_, _ = w.Write(v)
}

func readU64(r *bytes.Reader) (uint64, error) {
	var b [8]byte

	if _, err := io.ReadFull(r, b[:]); err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(b[:]), nil
}

func readBytes(r *bytes.Reader) ([]byte, error) {
	n, err := readU64(r)
	if err != nil {
		return nil, err
	}

	if n > maxCheckpointFieldSize {
		return nil, errors.New("checkpoint field too large")
	}

	b := make([]byte, n)

	if _, err := io.ReadFull(r, b); err != nil {
		return nil, err
	}

	return b, nil
}

func fmtHex(v []byte) string {
	const hexdigits = "0123456789abcdef"

	out := make([]byte, len(v)*2)

	for i, b := range v {
		out[i*2] = hexdigits[b>>4]
		out[i*2+1] = hexdigits[b&15]
	}

	return string(out)
}
