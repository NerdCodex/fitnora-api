from cryptography.fernet import Fernet
import base64
import hashlib

password = bytes(input("Enter the password: ").encode())

key = base64.urlsafe_b64encode(hashlib.sha256(password).digest())

cipher = Fernet(key)

encrypted = b"gAAAAABp1UloL49oU2rVxk7W4my2G7W9nnnOIbidfupZcLIkSwTjj5GvBO8xAd-zZux6O0KuGbe7Ez0BJMwTwRd6ZsOUYzYT52MS6XIXoOav62in7HxfhXwsA3t55n727KBPFjhCc0-i"

decrypted = cipher.decrypt(encrypted)

print(decrypted.decode())