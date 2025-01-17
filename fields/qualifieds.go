package fields

import (
	"bytes"
	"encoding"
	"fmt"
	"reflect"

	"git.sr.ht/~whereswaldon/forest-go/serialize"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

const minSizeofQualified = sizeofDescriptor

// concrete qualified data types
type QualifiedHash struct {
	Descriptor HashDescriptor `arbor:"order=0,recurse=serialize"`
	Blob       Blob           `arbor:"order=1"`
}

const minSizeofQualifiedHash = sizeofHashDescriptor

const qualifiedTextSeparator = "__"

func marshalTextQualified(first, second encoding.TextMarshaler) ([]byte, error) {
	buf := new(bytes.Buffer)
	b, err := first.MarshalText()
	if err != nil {
		return nil, err
	}
	_, _ = buf.Write(b)
	_, _ = buf.Write([]byte(qualifiedTextSeparator))
	b, err = second.MarshalText()
	if err != nil {
		return nil, err
	}
	_, _ = buf.Write(b)
	return buf.Bytes(), nil
}

// NewQualifiedHash returns a valid QualifiedHash from the given data
func NewQualifiedHash(t HashType, content []byte) (*QualifiedHash, error) {
	hd, err := NewHashDescriptor(t, len(content))
	if err != nil {
		return nil, err
	}
	return &QualifiedHash{*hd, Blob(content)}, nil
}

func NullHash() *QualifiedHash {
	return &QualifiedHash{
		Descriptor: HashDescriptor{
			Type:   HashTypeNullHash,
			Length: 0,
		},
		Blob: []byte{},
	}
}

func (q *QualifiedHash) UnmarshalBinary(b []byte) error {
	unused, err := serialize.ArborDeserialize(reflect.ValueOf(&q.Descriptor), b)
	if err != nil {
		return err
	}
	return q.Blob.UnmarshalBinary(unused[:q.Descriptor.Length])
}

func (q *QualifiedHash) BytesConsumed() int {
	return sizeofHashDescriptor + q.Blob.BytesConsumed()
}

func (q *QualifiedHash) Equals(other *QualifiedHash) bool {
	return q.Descriptor.Equals(&other.Descriptor) && q.Blob.Equals(&other.Blob)
}

func (q *QualifiedHash) MarshalText() ([]byte, error) {
	return marshalTextQualified(&q.Descriptor, q.Blob)
}

func (q *QualifiedHash) UnmarshalText(b []byte) error {
	return unmarshalTextDelimited(b, qualifiedTextSeparator, &q.Descriptor, &q.Blob)
}

func (q *QualifiedHash) MarshalString() (string, error) {
	s, e := q.MarshalText()
	return string(s), e
}

func (q *QualifiedHash) Validate() error {
	if err := q.Descriptor.Validate(); err != nil {
		return err
	}
	if int(q.Descriptor.Length) != len(q.Blob) {
		return fmt.Errorf("Descriptor length %d does not match value length %d", q.Descriptor.Length, len(q.Blob))
	}
	return nil
}

type QualifiedContent struct {
	Descriptor ContentDescriptor `arbor:"order=0,recurse=serialize"`
	Blob       Blob              `arbor:"order=1"`
}

const minSizeofQualifiedContent = sizeofContentDescriptor

// NewQualifiedContent returns a valid QualifiedContent from the given data
func NewQualifiedContent(t ContentType, content []byte) (*QualifiedContent, error) {
	hd, err := NewContentDescriptor(t, len(content))
	if err != nil {
		return nil, err
	}
	return &QualifiedContent{*hd, Blob(content)}, nil
}

func (q *QualifiedContent) Equals(other *QualifiedContent) bool {
	return q.Descriptor.Equals(&other.Descriptor) && q.Blob.Equals(&other.Blob)
}

func (q *QualifiedContent) UnmarshalBinary(b []byte) error {
	unused, err := serialize.ArborDeserialize(reflect.ValueOf(&q.Descriptor), b)
	if err != nil {
		return err
	}
	return q.Blob.UnmarshalBinary(unused[:q.Descriptor.Length])
}

func (q *QualifiedContent) BytesConsumed() int {
	return sizeofContentDescriptor + q.Blob.BytesConsumed()
}

func (q *QualifiedContent) MarshalText() ([]byte, error) {
	switch q.Descriptor.Type {
	case ContentTypeUTF8String:
		fallthrough
	case ContentTypeJSON:
		descText, err := (&q.Descriptor).MarshalText()
		if err != nil {
			return nil, err
		}
		buf := bytes.NewBuffer(descText)
		_, err = buf.Write(q.Blob)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	default:
		return marshalTextQualified(&q.Descriptor, q.Blob)
	}
}

func (q *QualifiedContent) Validate() error {
	if err := q.Descriptor.Validate(); err != nil {
		return err
	}
	if int(q.Descriptor.Length) != len(q.Blob) {
		return fmt.Errorf("Descriptor length %d does not match value length %d", q.Descriptor.Length, len(q.Blob))
	}
	return nil
}

type QualifiedKey struct {
	Descriptor KeyDescriptor `arbor:"order=0,recurse=serialize"`
	Blob       Blob          `arbor:"order=1"`
}

const minSizeofQualifiedKey = sizeofKeyDescriptor

// NewQualifiedKey returns a valid QualifiedKey from the given data
func NewQualifiedKey(t KeyType, content []byte) (*QualifiedKey, error) {
	hd, err := NewKeyDescriptor(t, len(content))
	if err != nil {
		return nil, err
	}
	return &QualifiedKey{*hd, Blob(content)}, nil
}

func (q *QualifiedKey) Equals(other *QualifiedKey) bool {
	return q.Descriptor.Equals(&other.Descriptor) && q.Blob.Equals(&other.Blob)
}

func (q *QualifiedKey) UnmarshalBinary(b []byte) error {
	unused, err := serialize.ArborDeserialize(reflect.ValueOf(&q.Descriptor), b)
	if err != nil {
		return err
	}
	return q.Blob.UnmarshalBinary(unused[:q.Descriptor.Length])
}

func (q *QualifiedKey) BytesConsumed() int {
	return sizeofKeyDescriptor + q.Blob.BytesConsumed()
}

func (q *QualifiedKey) MarshalText() ([]byte, error) {
	return marshalTextQualified(&q.Descriptor, q.Blob)
}

func (q *QualifiedKey) Validate() error {
	if err := q.Descriptor.Validate(); err != nil {
		return err
	}
	if int(q.Descriptor.Length) != len(q.Blob) {
		return fmt.Errorf("Descriptor length %d does not match value length %d", q.Descriptor.Length, len(q.Blob))
	}
	return nil
}

func (q *QualifiedKey) AsEntity() (*openpgp.Entity, error) {
	buf := bytes.NewBuffer(q.Blob)
	entity, err := openpgp.ReadEntity(packet.NewReader(buf))
	if err != nil {
		return nil, fmt.Errorf("Error reading public key: %v", err)
	}
	return entity, nil
}

type QualifiedSignature struct {
	Descriptor SignatureDescriptor `arbor:"order=0,recurse=serialize"`
	Blob       Blob                `arbor:"order=1"`
}

const minSizeofQualifiedSignature = sizeofSignatureDescriptor

// NewQualifiedSignature returns a valid QualifiedSignature from the given data
func NewQualifiedSignature(t SignatureType, content []byte) (*QualifiedSignature, error) {
	hd, err := NewSignatureDescriptor(t, len(content))
	if err != nil {
		return nil, err
	}
	return &QualifiedSignature{*hd, Blob(content)}, nil
}

func (q *QualifiedSignature) Equals(other *QualifiedSignature) bool {
	return q.Descriptor.Equals(&other.Descriptor) && q.Blob.Equals(&other.Blob)
}

func (q *QualifiedSignature) UnmarshalBinary(b []byte) error {
	unused, err := serialize.ArborDeserialize(reflect.ValueOf(&q.Descriptor), b)
	if err != nil {
		return err
	}
	return q.Blob.UnmarshalBinary(unused[:q.Descriptor.Length])
}

func (q *QualifiedSignature) BytesConsumed() int {
	return sizeofSignatureDescriptor + q.Blob.BytesConsumed()
}

func (q *QualifiedSignature) MarshalText() ([]byte, error) {
	return marshalTextQualified(&q.Descriptor, q.Blob)
}

func (q *QualifiedSignature) Validate() error {
	if err := q.Descriptor.Validate(); err != nil {
		return err
	}
	if int(q.Descriptor.Length) != len(q.Blob) {
		return fmt.Errorf("Descriptor length %d does not match value length %d", q.Descriptor.Length, len(q.Blob))
	}
	return nil
}
