section .data
    msg db "Result: ", 0
    buffer db "   ", 10, 0 
    num1 dd 10
    num2 dd 4

section .bss

section .text
    global _start

_start:
    mov eax, [num1]
    mov ebx, [num2]  
    add eax, ebx    
    mov ecx, buffer + 2 
    mov ebx, 10      

.convert_loop:
    mov edx, 0          
    div ebx         
    add dl, '0'  
    mov [ecx], dl  
    dec ecx            
    test eax, eax     
    jnz .convert_loop 

    inc ecx      
    mov edx, buffer + 3 
    sub edx, ecx

    mov eax, 4      
    mov ebx, 1       
    mov ecx, msg     
    mov edx, 8       
    int 0x80

    mov eax, 4
    mov ebx, 1
    mov ecx, ecx      
    mov edx, 3         
    int 0x80

    mov eax, 1   
    xor ebx, ebx        
    int 0x80
