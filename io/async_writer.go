package io

import (
	"io"
	"runtime"
	"sync"

	xcontext "go.nanasi880.dev/x/context"
)

type AsyncWriter interface {
	io.WriteCloser
	Await() error
}

type asyncWriter struct {
	w         io.WriteCloser
	mutex     *sync.Mutex
	writeCond *sync.Cond
	awaitCond *sync.Cond
	buffers   [2][]byte
	err       error
	close     xcontext.CancelableContext
	done      xcontext.CancelableContext
}

func NewAsyncWriter(w io.WriteCloser) AsyncWriter {
	mutex := new(sync.Mutex)
	asyncWrite := &asyncWriter{
		w:         w,
		mutex:     mutex,
		writeCond: sync.NewCond(mutex),
		awaitCond: sync.NewCond(mutex),
		buffers: [2][]byte{
			make([]byte, 0, 4096),
			make([]byte, 0, 4096),
		},
		close: xcontext.NewCancelable(),
		done:  xcontext.NewCancelable(),
	}
	runtime.SetFinalizer(asyncWrite, func(obj *asyncWriter) {
		_ = obj.Close()
	})
	go asyncWrite.start()
	return asyncWrite
}

func (w *asyncWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.close.IsDone() {
		return 0, io.EOF
	}

	w.buffers[0] = append(w.buffers[0], p...)
	w.writeCond.Signal()

	return len(p), nil
}

func (w *asyncWriter) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	err := w.awaitInternal()
	w.close.Cancel()
	<-w.done.Done()
	return err
}

func (w *asyncWriter) Await() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.awaitInternal()
}

func (w *asyncWriter) awaitInternal() error {
	w.writeCond.Signal()
	w.awaitCond.Wait()

	err := w.err
	w.err = nil

	return err
}

func (w *asyncWriter) start() {
	defer w.done.Cancel()
	for {
		done := w.startMain()
		if done {
			return
		}
	}
}

func (w *asyncWriter) startMain() bool {
	// wait write signal
	w.mutex.Lock()
	w.writeCond.Wait()

	// swap buffer
	w.buffers[0], w.buffers[1] = w.buffers[1], w.buffers[0]
	w.buffers[0] = w.buffers[0][:0]

	w.mutex.Unlock()

	var (
		buf = w.buffers[1]
		err error
	)

	if len(buf) > 0 {
		_, err = w.w.Write(buf)
	}

	if err != nil {
		w.mutex.Lock()
		if w.err == nil {
			w.err = err
		}
		w.mutex.Unlock()
	}

	w.mutex.Lock()
	w.awaitCond.Signal()
	w.mutex.Unlock()

	return w.close.IsDone()
}
