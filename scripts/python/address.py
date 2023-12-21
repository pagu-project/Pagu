import io
from enum import Enum
import utils

# Address format: hrp + `1` + type + data + checksum


class AddressType(Enum):
    Treasury = 0
    Validator = 1
    BLSAccount = 2


AddressSize = 21
TreasuryAddressString = "000000000000000000000000000000000000000000"
AddressHRP = "tpc"


class Address:
    def __init__(self, address_type, data):
        if len(data) != AddressSize - 1:
            raise ValueError("Data must be 21 bytes long")

        self.data = bytearray()
        self.data.append(address_type.value)
        self.data.extend(data)

    @classmethod
    def from_string(cls, text):
        if text == TreasuryAddressString:
            return bytes([0])

        hrp, typ, data = utils.decode_to_base256_with_type(text)
        if hrp != AddressHRP:
            raise ValueError(f"Invalid HRP: {hrp}")

        typ = AddressType(typ)
        if typ in (AddressType.Validator, AddressType.BLSAccount):
            if len(data) != 20:
                raise ValueError(f"Invalid length: {len(data) + 1}")
        else:
            raise ValueError(f"Invalid address type: {typ}")

        return cls(typ, data)

    def bytes(self):
        return bytes(self.data)

    def string(self):
        if self.data == bytes([0]):
            return TreasuryAddressString

        return utils.encode_from_base256_with_type(AddressHRP, self.data[0], self.data[1:])

    def address_type(self):
        return AddressType(self.data[0])

    def is_treasury_address(self):
        return self.address_type() == AddressType.Treasury

    def is_account_address(self):
        t = self.address_type()
        return t in (AddressType.Treasury, AddressType.BLSAccount)

    def is_validator_address(self):
        return self.address_type() == AddressType.Validator
