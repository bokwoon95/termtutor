package main

/*
// Start with the basic example from https://docs.julialang.org/en/release-0.6/manual/embedding/
//
// Obviously the paths below may need to be modified to match your julia install location and version number.
//
#cgo CFLAGS: -std=gnu99 -I'/Applications/Julia-1.0.app/Contents/Resources/julia/include/julia' -DJULIA_ENABLE_THREADING=1 -fPIC
#cgo LDFLAGS: -L'/Applications/Julia-1.0.app/Contents/Resources/julia/lib' -Wl,-rpath,'/Applications/Julia-1.0.app/Contents/Resources/julia/lib' -Wl,-rpath,'/Applications/Julia-1.0.app/Contents/Resources/julia/lib/julia' -ljulia
#include <julia.h>
*/
import "C"
import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	/* required: setup the Julia context */
	C.jl_init()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)

	/* run Julia commands */
	C.jl_eval_string(C.CString(`println(sqrt(2.0))`))
	C.jl_eval_string(C.CString(text))
	if C.jl_exception_occurred() {
		fmt.Printf("%s \n", C.jl_typeof_str(C.jl_exception_occurred()));
	}

	/* strongly recommended: notify Julia that the
	   program is about to terminate. this allows
	   Julia time to cleanup pending write requests
	   and run all finalizers
	*/
	C.jl_atexit_hook(0)
}
