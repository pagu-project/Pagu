import bech32m
import hashlib
from ripemd.ripemd160 import ripemd160

def decode_to_base256_with_type(text):
    (hrp, data, spec) = bech32m.bech32_decode(text)
    assert spec == bech32m.Encoding.BECH32M

    regrouped = bech32m.convertbits(data[1:], 5, 8, False)
    return (hrp, data[0], regrouped)

def encode_from_base256_with_type(hrp, typ, data):
    converted = bech32m.convertbits(list(data), 8, 5, True)
    converted = [typ] + converted
    return bech32m.bech32_encode(hrp, converted, bech32m.Encoding.BECH32M)