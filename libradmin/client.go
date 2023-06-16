package libradmin

import (
	"encoding/binary"
	"log"
	"net"
	"sync"
)

type RadminClient struct {
	net.Conn
	lastReadChannel uint32 // last channel the response was written to
	lastReadStatus  int32  // status for the last reponse that was written to the STATUS channel
	mu              sync.Mutex
	lastmu          sync.Mutex
}

/*
NewRadminClient returns a Radmin Client object ready to recieve radmin commands
*/
func NewRadminClient(socketAddr string) (r *RadminClient, err error) {
	r = &RadminClient{}

	r.Conn, err = net.Dial("unix", socketAddr)
	if err != nil {
		return r, err
	}

	// communicate magic header
	err = r.magicInit()
	if err != nil {
		return r, err
	}

	return
}

/*
NewRadminClientWithConn returns non-conventional (non unix-socket) RadminClient
for Testing/Mock Purpose
eg:
  - unix-socket listened over TCP using socat
    socat TCP-LISTEN:18121,fork,reuseaddr UNIX-CONNECT:/usr/local/var/run/radiusd/radiusd.sock
*/
func NewRadminClientWithConn(connType, connAddr string) (r *RadminClient, err error) {
	r = &RadminClient{}

	r.Conn, err = net.Dial(connType, connAddr)
	if err != nil {
		log.Fatal(err)
	}
	// communicate  magic header
	err = r.magicInit()
	if err != nil {
		return r, err
	}

	return
}

func (r *RadminClient) LastReadChannel() (channel uint32) {
	return r.lastReadChannel
}

func (r *RadminClient) LastReadStatus() (channel int32) {
	return r.lastReadStatus
}

func (r *RadminClient) channel_read(data *[]byte) (n int, err error) {
	r.lastmu.Lock()
	defer r.lastmu.Unlock()

	if n, err = fr_channel_read(r.Conn, &r.lastReadChannel, data); err != nil {
		return
	}
	if r.lastReadChannel == FR_CHANNEL_CMD_STATUS {
		b := make([]byte, 4)
		r.lastReadStatus = int32(binary.BigEndian.Uint32(b))
	}
	return
}

/*
Radius Control Socket expects a Magic Number to be sent on the INIT channel
at the very beginning of the connection before accepting any further commands
*/
func (r *RadminClient) magicInit() (err error) {
	mheader := []uint32{FR_CHANNEL_INIT_ACK, 8, FR_CHANNEL_MAGIC, 0}
	if err = binary.Write(r.Conn, binary.BigEndian, mheader); err != nil {
		return
	}

	// read back magic response
	var b []byte
	_, err = r.channel_read(&b)
	// discard response and nullify memory pointer
	b = nil
	return
}

func (r *RadminClient) Write(p []byte) (n int, err error) {
	hdr := []uint32{uint32(FR_CHANNEL_STDIN), uint32(len(p))}
	if err = binary.Write(r.Conn, binary.BigEndian, hdr); err != nil {
		return -1, err
	}
	return r.Conn.Write(p)
}

func (r *RadminClient) Read(p []byte) (n int, err error) {
	var data []byte
	n, err = r.channel_read(&data)
	if n <= 0 {
		return n, err
	}
	copy(p, data)
	return
}

func (r *RadminClient) Execute(command []byte) (result [][]byte, status int, err error) {
	// concurrency safe
	r.mu.Lock()
	defer r.mu.Unlock()

	var n int
	result = [][]byte{}
	p := make([]byte, 65536)

	r.Write([]byte(command))
	for {
		if n, err = r.Read(p); err != nil {
			return
		}
		tmp := make([]byte, n-1)
		copy(tmp, p[:n-1])

		if int(r.LastReadChannel()) == FR_CHANNEL_CMD_STATUS {
			status = int(r.lastReadStatus)
			break
		}
		result = append(result, tmp)
	}

	return
}
