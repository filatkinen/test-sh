package main

import (
	"bytes"
	"errors"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClientAdd(t *testing.T) {
	t.Run("Server close connection", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			// wait till closed server connection
			time.Sleep(time.Second * 1)

			// first send -got nil
			in.WriteString("hello\n")
			err1 := client.Send()
			require.NoError(t, err)

			// second send - got error
			in.WriteString("hello\n")
			err2 := client.Send()
			require.True(t, errors.Is(err1, ErrErrorsSend) || errors.Is(err2, ErrErrorsSend))
		}()
		go func() {
			defer wg.Done()
			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			require.NoError(t, conn.Close())
		}()
		wg.Wait()
	})
}
