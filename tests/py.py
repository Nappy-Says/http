
arrOfBytes = [
    60,63,120,109,108,32,118,101,114,115,105,111,110,61,34,49,46,
    48,34,32,101,110,99,111,100,105,110,103,61,34,117,116,102,45,
    56,34,63,62,32,60,115,118,103,32,118,101,114,115,105,111,110,
    61,34,49,46,49,34,32,120,109,108,110,115,61,34,104,116,116,
    112,58,47,47,119,119,119,46,119,51,46,111,114,103,47,50,48,
    48,48,47,115,118,103,34,62,118,49,60,47,115,118,103,62
]

newFileByteArray = bytearray(arrOfBytes)

with open("svg.svg", "wb") as binary_file:
    binary_file.write(newFileByteArray)