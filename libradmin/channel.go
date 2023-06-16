package libradmin

import (
	"encoding/binary"
	"io"
)

// const FR_CHANNEL_HDR_SIZE = 8

type rchannel_t struct {
	channel uint32
	length  uint32
}

func lo_read(conn io.Reader, inbuf *[]byte) (n int, err error) {
	n, err = conn.Read(*inbuf)
	return n, err
}

func fr_channel_read(conn io.Reader, channel *uint32, inbuf *[]byte) (n int, err error) {
	p := make([]uint32, 2)
	if err = binary.Read(conn, binary.BigEndian, &p); err != nil {
		return
	}
	length := p[1]
	*channel = p[0]

	data := make([]byte, length)
	*inbuf = data
	n, err = lo_read(conn, &data)
	return n, err
}

/* currentrly unused
func fr_channel_write(conn io.Writer, channel *uint32, data *[]byte) (n int, err error) {
	p := make([]uint32, 2)
	p[0] = *channel
	p[1] = uint32(len(*data))

	if err = binary.Write(conn, binary.BigEndian, p); err != nil {
		return
	}

	n, err = conn.Write(*data)
	return
}
*/
