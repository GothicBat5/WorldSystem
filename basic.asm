section .data
    msg db "Result: ", 0

section .bss
    buffer resb 10

section .text
    global _start

_start:
    ; Example: just print something

    mov eax, 4      ; sys_write
    mov ebx, 1      ; stdout
    mov ecx, msg
    mov edx, 8
    int 0x80

    ; Exit
    mov eax, 1
    xor ebx, ebx
    int 0x80

section .text
    global _start

_start:
    ; Load numbers into registers
    mov eax, 1000
    mov ebx, 25
    add eax, ebx
    mov al, [num1]   ; AL = 10
    mov bl, [num2]   ; BL = 4

    fld dword [num1]   ; load float
    fld dword [num2]
    fadd               ; add
    fstp dword [result]

    ; Addition
    add al, bl       ; AL = AL + BL (10 + 4 = 14)

    ; Subtraction
    sub al, bl       ; AL = AL - BL (14 - 4 = 10)

    ; Multiply
    mov eax, 10
    mov ebx, 4
    mul ebx          ; EAX = 10 * 4 = 40

    ; Division
    mov eax, 40
    mov ebx, 4
    div ebx          ; EAX = 10, EDX = remainder

    cmp eax, ebx
    jg greater
    jl lesser

    ; Exit program
    mov eax, 1       ; sys_exit
    xor ebx, ebx     ; exit code 0
    int 0x80
