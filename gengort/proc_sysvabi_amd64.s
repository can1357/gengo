//go:build !windows && amd64
#include "textflag.h"

#define CDECL NOSPLIT|NOFRAME

// func invoke0(proc uintptr) uintptr
TEXT ·invoke0(SB),CDECL,$0
   SUBQ $0x8, SP
   MOVQ proc+0*8+0x8(FP), AX
   CALL AX
   ADDQ $0x8, SP
   MOVQ AX, ret+1*8(FP)
   RET

// func invoke1(proc uintptr, a uintptr) uintptr
TEXT ·invoke1(SB),CDECL,$0
   SUBQ $0x8, SP
   MOVQ proc+0*8+0x8(FP), AX
   MOVQ _a+1*8+0x8(FP), DI
   CALL AX
   ADDQ $0x8, SP
   MOVQ AX, ret+2*8(FP)
   RET

// func invoke2(proc uintptr, a, b uintptr) uintptr
TEXT ·invoke2(SB),CDECL,$0 
   SUBQ $0x8, SP
   MOVQ proc+0*8+0x8(FP), AX
   MOVQ _a+1*8+0x8(FP), DI
   MOVQ _b+2*8+0x8(FP), SI
   CALL AX
   ADDQ $0x8, SP
   MOVQ AX, ret+3*8(FP)
   RET


// func invoke3(proc uintptr, a, b, c uintptr) uintptr
TEXT ·invoke3(SB),CDECL,$0
   SUBQ $0x8, SP
   MOVQ proc+0*8+0x8(FP), AX
   MOVQ _a+1*8+0x8(FP), DI
   MOVQ _b+2*8+0x8(FP), SI
   MOVQ _c+3*8+0x8(FP), DX
   CALL AX
   ADDQ $0x8, SP
   MOVQ AX, ret+4*8(FP)
   RET

// func invoke4(proc uintptr, a, b, c, d uintptr) uintptr
TEXT ·invoke4(SB),CDECL,$0 
   SUBQ $0x8, SP
   MOVQ proc+0*8+0x8(FP), AX
   MOVQ _a+1*8+0x8(FP), DI
   MOVQ _b+2*8+0x8(FP), SI
   MOVQ _c+3*8+0x8(FP), DX
   MOVQ _d+4*8+0x8(FP), CX
   CALL AX
   ADDQ $0x8, SP
   MOVQ AX, ret+5*8(FP)
   RET

// func invoke5(proc uintptr, a, b, c, d, e uintptr) uintptr
TEXT ·invoke5(SB),CDECL,$0
   SUBQ $0x8, SP
   MOVQ proc+0*8+0x8(FP), AX
   MOVQ _a+1*8+0x8(FP), DI
   MOVQ _b+2*8+0x8(FP), SI
   MOVQ _c+3*8+0x8(FP), DX
   MOVQ _d+4*8+0x8(FP), CX
   MOVQ _e+5*8+0x8(FP), R8
   CALL AX
   ADDQ $0x8, SP
   MOVQ AX, ret+6*8(FP)
   RET
   


// func invoke6(proc uintptr, a, b, c, d, e, f uintptr) uintptr
TEXT ·invoke6(SB),CDECL,$0
   SUBQ $0x8, SP
   MOVQ proc+0*8+0x8(FP), AX
   MOVQ _a+1*8+0x8(FP), DI
   MOVQ _b+2*8+0x8(FP), SI
   MOVQ _c+3*8+0x8(FP), DX
   MOVQ _d+4*8+0x8(FP), CX
   MOVQ _e+5*8+0x8(FP), R8
   MOVQ _f+6*8+0x8(FP), R9
   CALL AX
   ADDQ $0x8, SP
   MOVQ AX, ret+7*8(FP)
   RET

// func invoke7(proc uintptr, a, b, c, d, e, f, g uintptr) uintptr
TEXT ·invoke7(SB),CDECL,$0 
   SUBQ $0x8, SP
   MOVQ proc+0*8+0x8(FP), AX
   MOVQ _a+1*8+0x8(FP), DI
   MOVQ _b+2*8+0x8(FP), SI
   MOVQ _c+3*8+0x8(FP), DX
   MOVQ _d+4*8+0x8(FP), CX
   MOVQ _e+5*8+0x8(FP), R8
   MOVQ _f+6*8+0x8(FP), R9
   MOVQ _g+7*8+0x8(FP), R10
   MOVQ R10, 0(SP)
   CALL AX
   ADDQ $0x8, SP
   MOVQ AX, ret+8*8(FP)
   RET

// func invoke8(proc uintptr, a, b, c, d, e, f, g, h uintptr) uintptr
TEXT ·invoke8(SB),CDECL,$0
   SUBQ $0x18, SP
   MOVQ proc+0*8+0x18(FP), AX
   MOVQ _a+1*8+0x18(FP), DI
   MOVQ _b+2*8+0x18(FP), SI
   MOVQ _c+3*8+0x18(FP), DX
   MOVQ _d+4*8+0x18(FP), CX
   MOVQ _e+5*8+0x18(FP), R8
   MOVQ _f+6*8+0x18(FP), R9
   MOVQ _g+7*8+0x18(FP), R10
   MOVQ _h+8*8+0x18(FP), R11
   MOVQ R10, 0(SP)
   MOVQ R11, 8(SP)
   CALL AX
   ADDQ $0x18, SP
   MOVQ AX, ret+8*8(FP)
   RET

// func invoke9(proc uintptr, a, b, c, d, e, f, g, h, i uintptr) uintptr
TEXT ·invoke9(SB),CDECL,$0
   SUBQ $0x18, SP
   MOVQ proc+0*8+0x18(FP), AX
   MOVQ _a+1*8+0x18(FP), DI
   MOVQ _b+2*8+0x18(FP), SI
   MOVQ _c+3*8+0x18(FP), DX
   MOVQ _d+4*8+0x18(FP), CX
   MOVQ _e+5*8+0x18(FP), R8
   MOVQ _f+6*8+0x18(FP), R9
  
   MOVUPS _g+7*8+0x18(FP), X0
   MOVQ   _i+9*8+0x18(FP), R10
   MOVUPS X0, 0(SP)
   MOVQ   R10, 16(SP)
   
   CALL AX
   ADDQ $0x18, SP
   MOVQ AX, ret+9*8(FP)
   RET

// func invoke10(proc uintptr, a, b, c, d, e, f, g, h, i, j uintptr) uintptr
TEXT ·invoke10(SB),CDECL,$0
   SUBQ $0x28, SP
   MOVQ proc+0*8+0x28(FP), AX
   MOVQ _a+1*8+0x28(FP), DI
   MOVQ _b+2*8+0x28(FP), SI
   MOVQ _c+3*8+0x28(FP), DX
   MOVQ _d+4*8+0x28(FP), CX
   MOVQ _e+5*8+0x28(FP), R8
   MOVQ _f+6*8+0x28(FP), R9
  
   MOVUPS _g+7*8+0x28(FP), X0
   MOVUPS _i+9*8+0x28(FP), X1
   MOVUPS X0, 0(SP)
   MOVUPS X1, 16(SP)
   
   CALL AX
   ADDQ $0x28, SP
   MOVQ AX, ret+10*8(FP)
   RET

// func invoke11(proc uintptr, a, b, c, d, e, f, g, h, i, j, k uintptr) uintptr
TEXT ·invoke11(SB),CDECL,$0
   SUBQ $0x28, SP
   MOVQ proc+0*8+0x28(FP), AX
   MOVQ _a+1*8+0x28(FP), DI
   MOVQ _b+2*8+0x28(FP), SI
   MOVQ _c+3*8+0x28(FP), DX
   MOVQ _d+4*8+0x28(FP), CX
   MOVQ _e+5*8+0x28(FP), R8
   MOVQ _f+6*8+0x28(FP), R9
  
   MOVUPS _g+7*8+0x28(FP), X0
   MOVUPS _i+9*8+0x28(FP), X1
   MOVQ   _k+11*8+0x28(FP), R10
   MOVUPS X0, 0(SP)
   MOVUPS X1, 16(SP)
   MOVQ   R10, 32(SP)
   
   CALL AX
   ADDQ $0x28, SP
   MOVQ AX, ret+11*8(FP)
   RET

// func invoke12(proc uintptr, a, b, c, d, e, f, g, h, i, j, k, l uintptr) uintptr
TEXT ·invoke12(SB),CDECL,$0
   SUBQ $0x38, SP
   MOVQ proc+0*8+0x38(FP), AX
   MOVQ _a+1*8+0x38(FP), DI
   MOVQ _b+2*8+0x38(FP), SI
   MOVQ _c+3*8+0x38(FP), DX
   MOVQ _d+4*8+0x38(FP), CX
   MOVQ _e+5*8+0x38(FP), R8
   MOVQ _f+6*8+0x38(FP), R9
  
   MOVUPS _g+7*8+0x38(FP), X0
   MOVUPS _i+9*8+0x38(FP), X1
   MOVUPS _k+11*8+0x38(FP), X2
   MOVUPS X0, 0(SP)
   MOVUPS X1, 16(SP)
   MOVUPS X2, 32(SP)
   
   CALL AX
   ADDQ $0x38, SP
   MOVQ AX, ret+12*8(FP)
   RET
