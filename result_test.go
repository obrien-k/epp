package epp

import (
	"bytes"
	"testing"

	"github.com/nbio/st"
)

func TestDecodeResult(t *testing.T) {
	var r Result
	var buf bytes.Buffer
	d := NewDecoder(&buf)

	buf.Reset()
	buf.WriteString(`<result code="1000"><msg>Command completed successfully</msg></result>`)
	d.Reset()
	err := decodeResult(&d, &r)
	st.Expect(t, err, nil)
	st.Expect(t, r.Code, 1000)
	st.Expect(t, r.Message, "Command completed successfully")
	st.Expect(t, r.IsError(), false)
	st.Expect(t, r.IsFatal(), false)

	buf.Reset()
	buf.WriteString(`<result code="2001"><msg>Command syntax error</msg></result>`)
	d.Reset()
	err = decodeResult(&d, &r)
	st.Expect(t, err, &r)
	st.Expect(t, r.Code, 2001)
	st.Expect(t, r.Message, "Command syntax error")
	st.Expect(t, r.IsError(), true)
	st.Expect(t, r.IsFatal(), false)

	buf.Reset()
	buf.WriteString(`<result code="2501"><msg>Authentication error; server closing connection</msg></result>`)
	d.Reset()
	err = decodeResult(&d, &r)
	st.Expect(t, err, &r)
	st.Expect(t, r.Code, 2501)
	st.Expect(t, r.Message, "Authentication error; server closing connection")
	st.Expect(t, r.IsError(), true)
	st.Expect(t, r.IsFatal(), true)
}