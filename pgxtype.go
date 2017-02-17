package pgxtype

import (
  "strconv"
  "encoding/binary"
  "fmt"
)

type FieldDescription interface {
  Name() string
  Table() uint32
  AttributeNumber() int16
  DataType() uint32
  DataTypeSize() int16
  DataTypeName() string
  Modifier() int32
  FormatCode() int16
}

// Remember need to delegate for server controlled format like inet

// Or separate interfaces for raw bytes and preprocessed by pgx?

// Or interface{} like database/sql - and just pre-process into more things

type ScannerV3 interface {
  ScanPgxV3(fieldDescription FieldDescription, src interface{}) error
}

// Encoders could also return interface{} to delegate to internal pgx

type TextEncoderV3 interface {
  EncodeTextPgxV3(oid uint32) (interface{}, error)
}

type BinaryEncoderV3 interface {
  EncodeBinaryPgxV3(oid uint32) (interface{}, error)
}

const (
  Int4OID = 23
)

type NullInt32 struct {
  Value  int32
  Status byte
}

func (s *NullInt32) ScanPgxV3(fieldDescription FieldDescription, buf []byte) error {
  if fieldDescription.DataType() != Int4OID {
    return 0, fmt.Errorf("cannot decode %s (oid: %d)", fieldDescription.DataTypeName(), fieldDescription.DataType()))
  }

  if buf == nil {
    s.Int32, s.Status = 0, false
    return nil
  }

  if len(buf) != 4 {
    return 0, fmt.Errorf("invalid length for %s (oid: %d): %d bytes", fieldDescription.DataTypeName(), fieldDescription.DataType(), len(buf)))
  }

  s.Value = int32(binary.BigEndian.Uint32(buf))
  s.Status = 1

  return nil
}

func (s NullInt32) EncodeTextPgxV3(oid uint32) ([]byte, error) {
  return []byte(strconv.FormatInt(int64(s.Value), 10)), nil
}
