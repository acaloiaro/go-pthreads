package pthread

/*
#include <pthread.h>
#include <signal.h>
#include <unistd.h>
#include <stdio.h>

extern void createThreadCallback();
static void sig_func(int sig);

static void createThread(pthread_t* pid) {
	pthread_create(pid, NULL, (void*)createThreadCallback, NULL);
}

static void sig_func(int sig)
{
	pthread_exit(NULL);
}

static void register_sig_handler() {
	signal(SIGABRT, sig_func);
}
*/
import "C"
import "unsafe"

type Thread uintptr
type ThreadCallback func()

var create_callback chan ThreadCallback

func init() {
	C.register_sig_handler()
	create_callback = make(chan ThreadCallback, 1)
}

//export createThreadCallback
func createThreadCallback() {
	C.register_sig_handler()
	C.pthread_setcanceltype(C.PTHREAD_CANCEL_ASYNCHRONOUS, nil)
	(<-create_callback)()
}

// initializes a thread using pthread_create
func Create(cb ThreadCallback) Thread {
	var pid C.pthread_t
	pidptr := &pid
	create_callback <- cb

	C.createThread(pidptr)

	return Thread(uintptr(unsafe.Pointer(&pid)))
}

// signals the thread in question to terminate
func (t Thread) Kill() {
	C.pthread_detach(t.c());
	C.pthread_kill(t.c(), C.SIGABRT)
}

// helper function to convert the Thread object into a C.pthread_t object
func (t Thread) c() C.pthread_t {
	return *(*C.pthread_t)(unsafe.Pointer(t))
}
