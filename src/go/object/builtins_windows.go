package object

import "syscall"

func doSyscall(trap, a1, a2, a3 uintptr) (uintptr, uintptr, syscall.Errno) {
	return syscall.Syscall(trap, a1, a2, a3, 0)
}
