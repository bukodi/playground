import os
import hashlib
import array

start_dir = "./../.."
checksum = bytearray(32)

def process_dir(dir):
    for item in os.listdir(dir):
        if item.startswith("."): continue

        full_path = os.path.join(dir, item)

        if os.path.isdir(full_path):
            process_dir(full_path)
            continue

        relative_to_start = os.path.relpath(full_path, start_dir)
        hash = hashlib.sha256()
        hash.update(relative_to_start.encode("utf-8"))
        hash.update(bytearray(open(full_path, "rb").read()))
        item_hash = hash.digest()
        print( item_hash.hex() + " " + relative_to_start)
        xor_item( item_hash )

def xor_item(item_hash):
    global checksum
    checksum = bytes(a ^ b for (a, b) in zip(checksum, item_hash))

process_dir(start_dir)
print( "Checksum : " + checksum.hex() )
