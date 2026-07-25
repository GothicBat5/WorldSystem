import secrets
import string 
 #what is this all about!?? haha
characters = string.ascii_letters + string.digits + "!@#$%"
password = "".join(secrets.choice(characters) for _ in range(16))

print("Generated Password: ")
print(password)
