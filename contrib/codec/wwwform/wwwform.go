package wwwform

import (
	"github.com/go-playground/form/v4"
	"github.com/lxzan/xray/codec"
	"io"
	"net/url"
)

var FormCodec = new(Codec)

type Codec struct{}

func (c *Codec) NewEncoder(w io.Writer) codec.Encoder {
	return &Encoder{writer: w}
}

func (c *Codec) NewDecoder(r io.Reader) codec.Decoder {
	return &Decoder{reader: r}
}

func (c *Codec) Encode(v any) ([]byte, error) {
	values, err := form.NewEncoder().Encode(v)
	if err != nil {
		return nil, err
	}
	return []byte(values.Encode()), nil
}

func (c *Codec) EncodeToString(v any) (string, error) {
	values, err := form.NewEncoder().Encode(v)
	if err != nil {
		return "", err
	}
	return values.Encode(), nil
}

func (c *Codec) DecodeFromString(data string, v any) error {
	values, err := url.ParseQuery(data)
	if err != nil {
		return err
	}
	return form.NewDecoder().Decode(v, values)
}

func (c *Codec) Decode(data []byte, v any) error {
	values, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}
	return form.NewDecoder().Decode(v, values)
}

type Encoder struct {
	writer io.Writer
}

func (c *Encoder) Encode(v any) error {
	values, err := form.NewEncoder().Encode(v)
	if err != nil {
		return err
	}
	_, err = c.writer.Write([]byte(values.Encode()))
	return err
}

type Decoder struct {
	reader io.Reader
}

func (c *Decoder) Decode(v any) error {
	data, err := io.ReadAll(c.reader)
	if err != nil {
		return err
	}
	values, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}
	return form.NewDecoder().Decode(v, values)
}
