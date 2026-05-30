section .data
    msg db "Result: ", 0
    ; We need a buffer to hold the ASCII characters of the result
    ; Max 3 digits for small numbers + newline
    buffer db "   ", 10, 0 
    
    ; Define your numbers here
    num1 dd 10
    num2 dd 4

section .bss
    ; No extra buffer needed if we use the one in .data for printing

section .text
    global _start

_start:
    ; --- 1. Perform Integer Math ---
    mov eax, [num1]   ; Load 10 into EAX
    mov ebx, [num2]   ; Load 4 into EBX

    add eax, ebx      ; EAX = 14
    ; You can uncomment these to test other ops:
    ; sub eax, ebx    ; EAX = 10
    ; mul ebx         ; EAX = 40 (result in EAX)
    ; div ebx         ; EAX = 2 (quotient), EDX = remainder

    ; --- 2. Convert Number (EAX) to String ---
    ; Simple conversion logic for positive numbers
    mov ecx, buffer + 2 ; Point ECX to the last digit position
    mov ebx, 10         ; Divisor

.convert_loop:
    mov edx, 0          ; Clear EDX before division
    div ebx             ; EAX = EAX / 10, EDX = Remainder (the digit)
    
    add dl, '0'         ; Convert remainder to ASCII
    mov [ecx], dl       ; Store digit in buffer
    dec ecx             ; Move pointer left

    test eax, eax       ; Check if quotient is 0
    jnz .convert_loop   ; If not 0, continue loop

    inc ecx             ; ECX now points to the first digit

    ; --- 3. Print the Result ---
    ; Calculate length: (buffer + 2) - ecx + 1 (for newline)
    ; For "14\n", length is 3. Let's just hardcode length for simplicity here 
    ; or calculate: mov edx, 3 (since we know it's 2 digits + newline)
    
    ; Let's calculate dynamically just to be safe:
    mov edx, buffer + 3 ; End address (buffer + 2 is last char, +1 for newline)
    sub edx, ecx        ; Length = End - Start

    ; Syscall: write
    mov eax, 4          ; sys_write
    mov ebx, 1          ; stdout
    mov ecx, msg        ; Print "Result: "
    mov edx, 8          ; Length of "Result: "
    int 0x80

    ; Syscall: write result number
    mov eax, 4
    mov ebx, 1
    mov ecx, ecx        ; Address of first digit
    mov edx, 3          ; Length: 2 digits + 1 newline
    int 0x80

    ; --- 4. Exit Program ---
    mov eax, 1          ; sys_exit
    xor ebx, ebx        ; Exit code 0
    int 0x80
