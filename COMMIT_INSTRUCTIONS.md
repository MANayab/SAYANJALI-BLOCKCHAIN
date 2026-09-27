# Phase 9.3 Manual Commit Instructions

The Phase 9.3 build was created locally only. No GitHub branch, commit, push, or PR was created by this build.

## 1. Fetch origin

```bash
git fetch origin --prune
```

## 2. Verify authoritative main

```bash
git checkout main
git pull --ff-only origin main
git rev-parse HEAD
```

Expected:

```text
b821d158e2f5fb231025ada0d1417890d3fbdb0fcf
```

If the printed SHA differs, STOP and reconcile the baseline before applying Phase 9.3.

> Correction: the authoritative baseline supplied for this phase is `b821d158e2f5fb231025ada0d1417890d3fb0fcf`. Use that exact SHA.

## 3. Create the branch

```bash
git checkout -b phase-9-3-security-remediation b821d158e2f5fb231025ada0d1417890d3fb0fcf
```

## 4. Copy the ZIP contents

Extract `SYJ-BLOCKCHAIN-phase-9.3-security-remediation.zip`, then copy the contents of its `phase-9.3-security-remediation/` directory into the repository root.

Example:

```bash
unzip SYJ-BLOCKCHAIN-phase-9.3-security-remediation.zip -d /tmp/syj-phase93
cp -a /tmp/syj-phase93/phase-9.3-security-remediation/. .
```

## 5. Inspect the diff

```bash
git status --short
git diff --stat
git diff --check
git diff -- . ':!database' ':!logs'
```

Do not commit generated databases, logs, caches, secrets, private keys, or binaries.

## 6. Run validation

```bash
go version
python3 --version

gofmt -l cmd internal pkg tests

go test ./...
go vet ./...
go build ./...
go test -race ./...

govulncheck ./...
staticcheck ./...
gosec ./...

python3 -m compileall -q blockchain api cli config tests
pytest -q
pip-audit -r requirements.lock.txt
bandit -r blockchain api cli config -q
```

Do not treat blocked scanners as PASS.

## 7. Commit

Only after manual review and successful/explicitly documented validation:

```bash
git add .
git diff --cached --check
git commit -m "security: remediate Phase 9.3 consensus and P2P findings"
```

## 8. Push

```bash
git push -u origin phase-9-3-security-remediation
```

## 9. Open the next PR

Open a pull request from:

```text
phase-9-3-security-remediation
```

to:

```text
main
```

Suggested title:

```text
security: remediate Phase 9.3 consensus and P2P findings
```

No PR is claimed to exist until the user creates it.
