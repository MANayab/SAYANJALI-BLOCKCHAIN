from blockchain.state_root import state_root


def test_state_root_deterministic_and_order_independent():
    a = state_root(balances={"SYJa": 10, "SYJb": 20}, nonces={"SYJa": 2}, genesis_supply=30, mining_issued=5, supply=35)
    b = state_root(balances={"SYJb": 20, "SYJa": 10}, nonces={"SYJa": 2}, genesis_supply=30, mining_issued=5, supply=35)
    assert a == b


def test_state_root_changes_on_balance_nonce_or_supply():
    base = dict(balances={"SYJa": 10}, nonces={"SYJa": 1}, genesis_supply=10, mining_issued=0, supply=10)
    roots = {state_root(**base), state_root(**{**base, "balances": {"SYJa": 11}}), state_root(**{**base, "nonces": {"SYJa": 2}}), state_root(**{**base, "supply": 11})}
    assert len(roots) == 4
