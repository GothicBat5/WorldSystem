import hashlib
import hmac
import secrets
import base64
import ssl
import binascii


print("=" * 55)
print("        Python Security Toolkit")
print("=" * 55)

message = input("\nEnter a secret message: ")

print("\nGenerating secure random key...")
secret_key = secrets.token_bytes(32)


# Base64 Encoding
base64_encoded = base64.b64encode(message.encode())


# Hexadecimal Encoding
hex_encoded = binascii.hexlify(message.encode())


# SHA-256 Hash
sha256_hash = hashlib.sha256(message.encode()).hexdigest()


# HMAC Signature
signature = hmac.new(secret_key,
    message.encode(),
    hashlib.sha256
).hexdigest()


# Random Secure Token
random_token = secrets.token_hex(16)


# SSL Information
ssl_version = ssl.OPENSSL_VERSION

print("\n" + "=" * 55)
print("              RESULTS")
print("=" * 55)

print(f"\nOriginal Message: ")
print(message)

print("\nBase64: ")
print(base64_encoded.decode())

print("\nHexadecimal: ")
print(hex_encoded.decode())

print("\nSHA-256: ")
print(sha256_hash)

print("\nHMAC Signature: ")
print(signature)

print("\nSecure Random Token: ")
print(random_token)

print("\nSSL Library: ")
print(ssl_version)

print("\nRandom Secret Key (Hex): ")
print(secret_key.hex())

print("\n" + "=" * 55)
print("Security Toolkit Complete.")
print("=" * 55)
