//go:build windows && amd64
#include "textflag.h"

#define CDECL NOSPLIT|NOFRAME

#define INVOKE_BEG(N) \
   LEAQ a0+0(FP), AX        /* ax=frame */ \
   PUSHQ BP              /* create new frame */ \
   MOVQ SP, BP           /* */ \
   SUBQ $(0x20+8*(N)), SP /* */ \
   ANDQ $~15, SP         /* */
#define INVOKE_FIN(N) \
   MOVQ (AX), AX      /* ax=proc */ \
   CALL AX               /* call proc */ \
   MOVQ BP, SP           /* restore old frame */ \
   POPQ BP               /* */ \
   MOVQ AX, ret+((N)+1)*8(FP) /* ret=AX */ \
   RET                   /* return */

// func invoke0(proc uintptr) uintptr
TEXT ·invoke0(SB),CDECL,$0
   INVOKE_BEG(0)
   INVOKE_FIN(0)
// func invoke1(proc uintptr, a uintptr) uintptr
TEXT ·invoke1(SB),CDECL,$0
   INVOKE_BEG(1)
   MOVQ 1*8(AX), CX
   INVOKE_FIN(1)
// func invoke2(proc uintptr, a, b uintptr) uintptr
TEXT ·invoke2(SB),CDECL,$0
   INVOKE_BEG(2)
   MOVQ 1*8(AX), CX
   MOVQ 2*8(AX), DX
   INVOKE_FIN(2)
// func invoke3(proc uintptr, a, b, c uintptr) uintptr
TEXT ·invoke3(SB),CDECL,$0 
   INVOKE_BEG(3)
   MOVQ 1*8(AX), CX
   MOVQ 2*8(AX), DX
   MOVQ 3*8(AX), R8
   INVOKE_FIN(3)
// func invoke4(proc uintptr, a, b, c, d uintptr) uintptr
TEXT ·invoke4(SB),CDECL,$0 
   INVOKE_BEG(4)
   MOVQ 1*8(AX), CX
   MOVQ 2*8(AX), DX
   MOVQ 3*8(AX), R8
   MOVQ 4*8(AX), R9
   INVOKE_FIN(4)
// func invoke5(proc uintptr, a, b, c, d, e uintptr) uintptr
TEXT ·invoke5(SB),CDECL,$0 
   INVOKE_BEG(5)
   MOVQ   1*8(AX), CX
   MOVQ   2*8(AX), DX
   MOVQ   3*8(AX), R8
   MOVQ   4*8(AX), R9
   MOVUPS 4*8(AX), X0
   MOVUPS X0,      3*8(SP)
   INVOKE_FIN(5)
// func invoke6(proc uintptr, a, b, c, d, e, f uintptr) uintptr
TEXT ·invoke6(SB),CDECL,$0 
   INVOKE_BEG(6)
   MOVQ   1*8(AX), CX
   MOVQ   2*8(AX), DX
   MOVQ   3*8(AX), R8
   MOVQ   4*8(AX), R9
   MOVUPS 5*8(AX), X0
   MOVUPS X0,      4*8(SP)
   INVOKE_FIN(6)
// func invoke7(proc uintptr, a, b, c, d, e, f, g uintptr) uintptr
TEXT ·invoke7(SB),CDECL,$0 
   INVOKE_BEG(7)
   MOVQ   1*8(AX), CX
   MOVQ   2*8(AX), DX
   MOVQ   3*8(AX), R8
   MOVQ   4*8(AX), R9
   MOVUPS 4*8(AX), X0
   MOVUPS 6*8(AX), X1
   MOVUPS X0,      3*8(SP)
   MOVUPS X1,      5*8(SP)
   INVOKE_FIN(7)
// func invoke8(proc uintptr, a, b, c, d, e, f, g, h uintptr) uintptr
TEXT ·invoke8(SB),CDECL,$0 
   INVOKE_BEG(8)
   MOVQ   1*8(AX), CX
   MOVQ   2*8(AX), DX
   MOVQ   3*8(AX), R8
   MOVQ   4*8(AX), R9
   MOVUPS 5*8(AX), X0
   MOVUPS 7*8(AX), X1
   MOVUPS X0,      4*8(SP)
   MOVUPS X1,      6*8(SP)
   INVOKE_FIN(8)
// func invoke9(proc uintptr, a, b, c, d, e, f, g, h, i uintptr) uintptr
TEXT ·invoke9(SB),CDECL,$0
   INVOKE_BEG(9)
   MOVQ   1*8(AX), CX
   MOVQ   2*8(AX), DX
   MOVQ   3*8(AX), R8
   MOVQ   4*8(AX), R9
   MOVUPS 4*8(AX), X0
   MOVUPS 6*8(AX), X1
   MOVUPS 8*8(AX), X2
   MOVUPS X0,      3*8(SP)
   MOVUPS X1,      5*8(SP)
   MOVUPS X2,      7*8(SP)
   INVOKE_FIN(9)
// func invoke10(proc uintptr, a, b, c, d, e, f, g, h, i, j uintptr) uintptr
TEXT ·invoke10(SB),CDECL,$0
   INVOKE_BEG(10)
   MOVQ   1*8(AX), CX
   MOVQ   2*8(AX), DX
   MOVQ   3*8(AX), R8
   MOVQ   4*8(AX), R9
   MOVUPS 5*8(AX), X0
   MOVUPS 7*8(AX), X1
   MOVUPS 9*8(AX), X2
   MOVUPS X0,      4*8(SP)
   MOVUPS X1,      6*8(SP)
   MOVUPS X2,      8*8(SP)
   INVOKE_FIN(10)
// func invoke11(proc uintptr, a, b, c, d, e, f, g, h, i, j, k uintptr) uintptr
TEXT ·invoke11(SB),CDECL,$0
   INVOKE_BEG(11)
   MOVQ   1*8(AX),  CX
   MOVQ   2*8(AX),  DX
   MOVQ   3*8(AX),  R8
   MOVQ   4*8(AX),  R9
   MOVUPS 4*8(AX),  X0
   MOVUPS 6*8(AX),  X1
   MOVUPS 8*8(AX),  X2
   MOVUPS 10*8(AX), X3
   MOVUPS X0,       3*8(SP)
   MOVUPS X1,       5*8(SP)
   MOVUPS X2,       7*8(SP)
   MOVUPS X3,       9*8(SP)
   INVOKE_FIN(11)
// func invoke12(proc uintptr, a, b, c, d, e, f, g, h, i, j, k, l uintptr) uintptr
TEXT ·invoke12(SB),CDECL,$0
   INVOKE_BEG(12)
   MOVQ   1*8(AX),  CX
   MOVQ   2*8(AX),  DX
   MOVQ   3*8(AX),  R8
   MOVQ   4*8(AX),  R9
   MOVUPS 5*8(AX),  X0
   MOVUPS 7*8(AX),  X1
   MOVUPS 9*8(AX),  X2
   MOVUPS 11*8(AX), X3
   MOVUPS X0,       4*8(SP)
   MOVUPS X1,       6*8(SP)
   MOVUPS X2,       8*8(SP)
   MOVUPS X3,       10*8(SP)
   INVOKE_FIN(12)
   