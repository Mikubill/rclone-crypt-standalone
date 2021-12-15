
# rclone-crypt-standalone

An encryption and decryption tool based on rclone crypt and nacl/secretbox. It can be used to encrypt/decrypt files encrypted by rclone, or you can generate your own encrypted file and decrypt it on rclone.

Currently supported:

* Encrypted files
* Decrypt files
* Obscure password
* Unobscure password

## Example
```
Usage:
  rcs [command]

Available Commands:
  decrypt     decrypt files
  encrypt     encrypt files
  help        Help about any command
  obscure     Obscure password for use in somewhere.
  unobscure   Unobscure password for use in somewhere.
```

encrypt file
```
rcs encrypt -m off -p <password> -s <password2> potato
```

decrypt file
```
rcs decrypt -m off -p <password> -s <password2> potato.bin
```

obscure string
```
rcs obscure abcd
```

unobscure string
```
rcs unobscure E9RDEKki5MKohzOGuga7wyGsPy0
```
