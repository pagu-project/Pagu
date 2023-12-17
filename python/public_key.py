import hashlib
from ripemd.ripemd160 import ripemd160
import address
import utils

PublicKeySize = 96
PublicKeyHRP = "tpublic"
PublicKeyTypeBLS = 1


class PublicKey:
    def __init__(self, data=None):
        if data is None:
            self.data = bytearray(PublicKeySize)
        else:
            if len(data) != PublicKeySize:
                raise ValueError("Public key data must be 96 bytes long")
            self.data = bytearray(data)

    @classmethod
    def from_string(cls, text):
        hrp, typ, data = utils.decode_to_base256_with_type(text)

        if hrp != PublicKeyHRP:
            raise ValueError(f"Invalid hrp: {hrp}")

        if typ != PublicKeyTypeBLS:
            raise ValueError(f"Invalid public key type: {typ}")

        return cls(data)

    def bytes(self):
        return self.data

    def string(self):
        return utils.encode_from_base256_with_type(PublicKeyHRP, PublicKeyTypeBLS, self.data)

    def account_address(self):
        return self._make_address(address.AddressType.BLSAccount)

    def validator_address(self):
        return self._make_address(address.AddressType.Validator)

    def _make_address(self, address_type):
        blake2b = hashlib.blake2b(digest_size=32)
        blake2b.update(self.data)
        hash_256 = blake2b.digest()
        hash_160 = ripemd160(hash_256)
        addr = address.Address(address_type, hash_160)
        return addr
