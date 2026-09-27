"""Persistent P2P identity using cryptography's maintained SECP256k1 implementation."""
from __future__ import annotations
from cryptography.exceptions import InvalidSignature
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives.asymmetric.utils import decode_dss_signature, encode_dss_signature
from blockchain.utils import get_logger

logger = get_logger("blockchain.network.identity")
ORDER = int("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)

class P2PIdentity:
    def __init__(self, node_id: str, private_key_hex: str, public_key_hex: str) -> None:
        self.node_id, self.private_key_hex, self.public_key_hex = node_id, private_key_hex, public_key_hex

    @classmethod
    def generate(cls, node_id: str) -> "P2PIdentity":
        key = ec.generate_private_key(ec.SECP256K1()); pub = key.public_key().public_numbers()
        return cls(node_id, key.private_numbers().private_value.to_bytes(32, "big").hex(), pub.x.to_bytes(32, "big").hex()+pub.y.to_bytes(32, "big").hex())

    def sign(self, message: str) -> str:
        key = ec.derive_private_key(int.from_bytes(bytes.fromhex(self.private_key_hex), "big"), ec.SECP256K1())
        r, s = decode_dss_signature(key.sign(message.encode(), ec.ECDSA(hashes.SHA256())))
        if s > ORDER // 2: s = ORDER - s
        return r.to_bytes(32, "big").hex()+s.to_bytes(32, "big").hex()

def verify_identity_signature(public_key_hex: str, message: str, signature_hex: str) -> bool:
    try:
        pub, sig = bytes.fromhex(public_key_hex), bytes.fromhex(signature_hex)
        if len(pub) != 64 or len(sig) != 64: return False
        s = int.from_bytes(sig[32:], "big")
        if s == 0 or s > ORDER // 2: return False
        key = ec.EllipticCurvePublicNumbers(int.from_bytes(pub[:32], "big"), int.from_bytes(pub[32:], "big"), ec.SECP256K1()).public_key()
        key.verify(encode_dss_signature(int.from_bytes(sig[:32], "big"), s), message.encode(), ec.ECDSA(hashes.SHA256()))
        return True
    except (InvalidSignature, ValueError, TypeError): return False

def is_valid_public_key(public_key_hex: str) -> bool:
    try:
        raw = bytes.fromhex(public_key_hex)
        if len(raw) != 64: return False
        ec.EllipticCurvePublicNumbers(int.from_bytes(raw[:32], "big"), int.from_bytes(raw[32:], "big"), ec.SECP256K1()).public_key()
        return True
    except Exception: return False
