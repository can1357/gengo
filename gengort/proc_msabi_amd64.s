//go:build windows && amd64
#include "textflag.h"

#define CDECL NOSPLIT|NOFRAME

// func invoke0(proc uintptr) uintptr
TEXT ·invoke0(SB),CDECL,$0
   SUBQ $0x28, SP

   MOVQ proc+0*8+0x28(FP), AX
   CALL AX

   MOVQ AX, ret+1*8+0x28(FP)

   ADDQ $0x28, SP
   RET

// func invoke1(proc uintptr, a uintptr) uintptr
TEXT ·invoke1(SB),CDECL,$0
   SUBQ $0x28, SP

   MOVQ proc+0*8+0x28(FP), AX
   MOVQ _a+1*8+0x28(FP), CX
   CALL AX

   MOVQ AX, ret+2*8+0x28(FP)

   ADDQ $0x28, SP
   RET

// func invoke2(proc uintptr, a, b uintptr) uintptr
TEXT ·invoke2(SB),CDECL,$0 
   SUBQ $0x28, SP

   MOVQ proc+0*8+0x28(FP), AX
   MOVQ _a+1*8+0x28(FP), CX
   MOVQ _b+2*8+0x28(FP), DX
   CALL AX

   MOVQ AX, ret+3*8+0x28(FP)

   ADDQ $0x28, SP
   RET


// func invoke3(proc uintptr, a, b, c uintptr) uintptr
TEXT ·invoke3(SB),CDECL,$0
   SUBQ $0x28, SP

   MOVQ proc+0*8+0x28(FP), AX
   MOVQ _a+1*8+0x28(FP), CX
   MOVQ _b+2*8+0x28(FP), DX
   MOVQ _c+3*8+0x28(FP), R8

   CALL AX

   MOVQ AX, ret+4*8+0x28(FP)

   ADDQ $0x28, SP
   RET

// func invoke4(proc uintptr, a, b, c, d uintptr) uintptr
TEXT ·invoke4(SB),CDECL,$0
   SUBQ $0x28, SP
   
   MOVQ proc+0*8+0x28(FP), AX
   MOVQ _a+1*8+0x28(FP), CX
   MOVQ _b+2*8+0x28(FP), DX
   MOVQ _c+3*8+0x28(FP), R8
   MOVQ _d+4*8+0x28(FP), R9

   CALL AX

   MOVQ AX, ret+5*8+0x28(FP)

   ADDQ $0x28, SP
   RET

// func invoke5(proc uintptr, a, b, c, d, e uintptr) uintptr
TEXT ·invoke5(SB),CDECL,$0
   SUBQ $0x38, SP

   MOVQ proc+0*8+0x38(FP), AX
   MOVQ _a+1*8+0x38(FP), CX
   MOVQ _b+2*8+0x38(FP), DX
   MOVQ _c+3*8+0x38(FP), R8
   MOVQ _d+4*8+0x38(FP), R9
   
   MOVQ e+5*8+0x38(FP), R10
   MOVQ R10, 0x20(SP)

   CALL AX

   MOVQ AX, ret+6*8+0x38(FP)

   ADDQ $0x38, SP
   RET
   


// func invoke6(proc uintptr, a, b, c, d, e, f uintptr) uintptr
TEXT ·invoke6(SB),CDECL,$0
   SUBQ $0x38, SP

   MOVQ proc+0*8+0x38(FP), AX
   MOVQ _a+1*8+0x38(FP), CX
   MOVQ _b+2*8+0x38(FP), DX
   MOVQ _c+3*8+0x38(FP), R8
   MOVQ _d+4*8+0x38(FP), R9

   MOVUPS _e+5*8+0x38(FP), X0
   MOVUPS X0, 0x20(SP)

   CALL AX

   MOVQ AX, ret+7*8+0x38(FP)

   ADDQ $0x38, SP
   RET

// func invoke7(proc uintptr, a, b, c, d, e, f, g uintptr) uintptr
TEXT ·invoke7(SB),CDECL,$0
   SUBQ $0x48, SP

   MOVQ proc+0*8+0x48(FP), AX
   MOVQ _a+1*8+0x48(FP), CX
   MOVQ _b+2*8+0x48(FP), DX
   MOVQ _c+3*8+0x48(FP), R8
   MOVQ _d+4*8+0x48(FP), R9

   MOVUPS _e+5*8+0x48(FP), X0
   MOVQ   _f+6*8+0x48(FP), X1
   MOVUPS X0,  0x20(SP)
   MOVUPS X1,  0x28(SP)

   CALL AX

   MOVQ AX, ret+8*8+0x48(FP)

   ADDQ $0x48, SP
   RET

// func invoke8(proc uintptr, a, b, c, d, e, f, g, h uintptr) uintptr
TEXT ·invoke8(SB),CDECL,$0
   SUBQ $0x48, SP

   MOVQ proc+0*8+0x48(FP), AX
   MOVQ _a+1*8+0x48(FP), CX
   MOVQ _b+2*8+0x48(FP), DX
   MOVQ _c+3*8+0x48(FP), R8
   MOVQ _d+4*8+0x48(FP), R9

   MOVUPS _e+5*8+0x48(FP), X0
   MOVUPS _g+7*8+0x48(FP), X1
   MOVUPS X0, 0x20(SP)
   MOVUPS X1, 0x38(SP)

   CALL AX

   MOVQ AX, ret+9*8+0x48(FP)

   ADDQ $0x48, SP
   RET

// func invoke9(proc uintptr, a, b, c, d, e, f, g, h, i uintptr) uintptr
TEXT ·invoke9(SB),CDECL,$0
   SUBQ $0x58, SP

   MOVQ proc+0*8+0x58(FP), AX
   MOVQ _a+1*8+0x58(FP), CX
   MOVQ _b+2*8+0x58(FP), DX
   MOVQ _c+3*8+0x58(FP), R8
   MOVQ _d+4*8+0x58(FP), R9

   MOVUPS _e+5*8+0x58(FP), X0
   MOVUPS _g+7*8+0x58(FP), X1
   MOVUPS _h+8*8+0x58(FP), X2
   MOVUPS X0,  0x20(SP)
   MOVUPS X1,  0x38(SP)
   MOVUPS X2,  0x38(SP)

   CALL AX

   MOVQ AX, ret+10*8+0x58(FP)

   ADDQ $0x58, SP
   RET

// func invoke10(proc uintptr, a, b, c, d, e, f, g, h, i, j uintptr) uintptr
TEXT ·invoke10(SB),CDECL,$0
   SUBQ $0x58, SP

   MOVQ proc+0*8+0x58(FP), AX
   MOVQ _a+1*8+0x58(FP), CX
   MOVQ _b+2*8+0x58(FP), DX
   MOVQ _c+3*8+0x58(FP), R8
   MOVQ _d+4*8+0x58(FP), R9

   MOVUPS _e+5*8+0x58(FP), X0
   MOVUPS _g+7*8+0x58(FP), X1
   MOVUPS _i+9*8+0x58(FP), X2
   MOVUPS X0,  0x20(SP)
   MOVUPS X1,  0x38(SP)
   MOVUPS X2,  0x48(SP)

   CALL AX

   MOVQ AX, ret+11*8+0x58(FP)

   ADDQ $0x58, SP
   RET
// func invoke11(proc uintptr, a, b, c, d, e, f, g, h, i, j, k uintptr) uintptr
TEXT ·invoke11(SB),CDECL,$0
   SUBQ $0x68, SP

   MOVQ proc+0*8+0x68(FP), AX
   MOVQ _a+1*8+0x68(FP), CX
   MOVQ _b+2*8+0x68(FP), DX
   MOVQ _c+3*8+0x68(FP), R8
   MOVQ _d+4*8+0x68(FP), R9

   MOVUPS _e+5*8+0x68(FP), X0
   MOVUPS _g+7*8+0x68(FP), X1
   MOVUPS _i+9*8+0x68(FP), X2
   MOVUPS _j+10*8+0x68(FP), X3
   MOVUPS X0,  0x20(SP)
   MOVUPS X1,  0x38(SP)
   MOVUPS X2,  0x48(SP)
   MOVUPS X3,  0x48(SP)

   CALL AX

   MOVQ AX, ret+12*8+0x68(FP)

   ADDQ $0x68, SP
   RET
// func invoke12(proc uintptr, a, b, c, d, e, f, g, h, i, j, k, l uintptr) uintptr
TEXT ·invoke12(SB),CDECL,$0
   SUBQ $0x68, SP

   MOVQ proc+0*8+0x68(FP), AX
   MOVQ _a+1*8+0x68(FP), CX
   MOVQ _b+2*8+0x68(FP), DX
   MOVQ _c+3*8+0x68(FP), R8
   MOVQ _d+4*8+0x68(FP), R9

   MOVUPS _e+5*8+0x68(FP), X0
   MOVUPS _g+7*8+0x68(FP), X1
   MOVUPS _i+9*8+0x68(FP), X2
   MOVUPS _k+11*8+0x68(FP), X3
   MOVUPS X0,  0x20(SP)
   MOVUPS X1,  0x38(SP)
   MOVUPS X2,  0x48(SP)
   MOVUPS X3,  0x58(SP)

   CALL AX

   MOVQ AX, ret+13*8+0x68(FP)

   ADDQ $0x68, SP
   RET
