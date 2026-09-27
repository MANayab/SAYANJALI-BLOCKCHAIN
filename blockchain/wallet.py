"""Wallet primitives using cryptography's maintained SECP256k1 implementation."""
from __future__ import annotations
from dataclasses import dataclass
from cryptography.exceptions import InvalidSignature
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives.asymmetric.utils import decode_dss_signature, encode_dss_signature
from blockchain.utils import get_logger, sha256, WalletError

logger = get_logger("blockchain.wallet")
ADDRESS_PREFIX = "SYJ"
ADDRESS_HASH_LENGTH = 40
SECP256K1_ORDER = int("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)

@dataclass
class Wallet:
    private_key_hex: str
    public_key_hex: str
    address: str

    @classmethod
    def create(cls) -> "Wallet":
        key = ec.generate_private_key(ec.SECP256K1())
        priv = key.private_numbers().private_value.to_bytes(32, "big").hex()
        pub = key.public_key().public_numbers()
        pub_hex = pub.x.to_bytes(32, "big").hex() + pub.y.to_bytes(32, "big").hex()
        return cls(priv, pub_hex, derive_address(pub_hex))

    @classmethod
    def from_private_key(cls, private_key_hex: str) -> "Wallet":
        try:
            raw = bytes.fromhex(private_key_hex)
            if len(raw) != 32 or int.from_bytes(raw, "big") == 0 or int.from_bytes(raw, "big") >= SECP256K1_ORDER:
                raise ValueError("private key scalar out of range")
            key = ec.derive_private_key(int.from_bytes(raw, "big"), ec.SECP256K1())
        except Exception as exc:
            raise WalletError(f"Invalid private key: {exc}") from exc
        pub = key.public_key().public_numbers()
        pub_hex = pub.x.to_bytes(32, "big").hex() + pub.y.to_bytes(32, "big").hex()
        return cls(private_key_hex.lower(), pub_hex, derive_address(pub_hex))

    def export_keys(self) -> dict[str, str]:
        return {"public_key": self.public_key_hex, "address": self.address}

    def sign(self, message: str) -> str:
        key = ec.derive_private_key(int.from_bytes(bytes.fromhex(self.private_key_hex), "big"), ec.SECP256K1())
        r, s = decode_dss_signature(key.sign(message.encode(), ec.ECDSA(hashes.SHA256())))
        if s > SECP256K1_ORDER // 2:
            s = SECP256K1_ORDER - s
        return r.to_bytes(32, "big").hex() + s.to_bytes(32, "big").hex()

def derive_address(public_key_hex: str) -> str:
    return f"{ADDRESS_PREFIX}{sha256(public_key_hex)[:ADDRESS_HASH_LENGTH]}"

def is_valid_address(address: str) -> bool:
    if not isinstance(address, str) or not address.startswith(ADDRESS_PREFIX): return False
    remainder = address[len(ADDRESS_PREFIX):]
    if len(remainder) != ADDRESS_HASH_LENGTH or remainder != remainder.lower(): return False
    try: int(remainder, 16); return True
    except ValueError: return False

def verify_signature(public_key_hex: str, message: str, signature_hex: str) -> bool:
    try:
        pub = bytes.fromhex(public_key_hex); sig = bytes.fromhex(signature_hex)
        if len(pub) != 64 or len(sig) != 64: return False
        s = int.from_bytes(sig[32:], "big")
        if s == 0 or s > SECP256K1_ORDER // 2: return False
        key = ec.EllipticCurvePublicNumbers(int.from_bytes(pub[:32], "big"), int.from_bytes(pub[32:], "big"), ec.SECP256K1()).public_key()
        key.verify(encode_dss_signature(int.from_bytes(sig[:32], "big"), s), message.encode(), ec.ECDSA(hashes.SHA256()))
        return True
    except (InvalidSignature, ValueError, TypeError): return False
