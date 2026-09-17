"""Shared auth/API helpers for sdk example verify scripts."""
import json
import os
import subprocess
import sys
import urllib.error
import urllib.request


def die(msg: str) -> None:
    print(f"FAIL: {msg}", file=sys.stderr)
    sys.exit(1)


def resolve_auth() -> tuple[str, str, str]:
    endpoint = os.environ.get("VIVARCUS_ENDPOINT", "").rstrip("/")
    token = os.environ.get("VIVARCUS_TOKEN", "")
    vault = os.environ.get("VIVARCUS_VAULT", "")
    vivarcus = os.environ.get("VIVARCUS_BIN", "vivarcus")
    if endpoint and token and vault:
        return endpoint, token, vault
    try:
        out = subprocess.run(
            [vivarcus, "auth", "status", "--json"],
            capture_output=True,
            text=True,
            check=True,
        )
        status = json.loads(out.stdout or "{}")
    except (subprocess.CalledProcessError, json.JSONDecodeError) as e:
        die(f"not authenticated — run vivarcus auth login ({e})")
    if not token:
        token = status.get("session_token") or status.get("token") or ""
    if not endpoint:
        endpoint = (status.get("endpoint") or os.environ.get("VIVARCUS_ENDPOINT") or "").rstrip("/")
    if not vault:
        vault = status.get("default_vault") or os.environ.get("VIVARCUS_VAULT") or ""
    if not endpoint:
        cfg = subprocess.run([vivarcus, "config", "get", "endpoint"], capture_output=True, text=True)
        if cfg.returncode == 0:
            endpoint = cfg.stdout.strip().rstrip("/")
    if not vault:
        cfg = subprocess.run([vivarcus, "config", "get", "default_vault"], capture_output=True, text=True)
        if cfg.returncode == 0:
            vault = cfg.stdout.strip()
    if not token:
        cfg = subprocess.run([vivarcus, "config", "get", "token"], capture_output=True, text=True)
        if cfg.returncode == 0:
            token = cfg.stdout.strip()
    if not endpoint or not token or not vault:
        die("need endpoint, session token, and vault (CLI config or env)")
    return endpoint, token, vault


class DemoClient:
    def __init__(self) -> None:
        self.endpoint, self.token, self.vault = resolve_auth()
        self.vivarcus_bin = os.environ.get("VIVARCUS_BIN", "vivarcus")

    def api(self, method: str, path: str, body=None) -> dict:
        headers = {
            "Authorization": "Bearer " + self.token,
            "X-Vault-Id": self.vault,
            "Content-Type": "application/json",
        }
        data = None if body is None else json.dumps(body).encode()
        req = urllib.request.Request(self.endpoint + path, data=data, headers=headers, method=method)
        try:
            with urllib.request.urlopen(req) as resp:
                raw = resp.read().decode()
        except urllib.error.HTTPError as e:
            die(f"{method} {path} ({e.code}): {e.read().decode()[:1200]}")
        return json.loads(raw) if raw.strip() else {}

    def vivarcus(self, *args: str) -> dict:
        env = os.environ.copy()
        env["VIVARCUS_ENDPOINT"] = self.endpoint
        env["VIVARCUS_TOKEN"] = self.token
        env["VIVARCUS_VAULT"] = self.vault
        cmd = [self.vivarcus_bin, *args, "--vault", self.vault, "--json"]
        proc = subprocess.run(cmd, env=env, capture_output=True, text=True)
        if proc.returncode != 0:
            die(f"vivarcus {' '.join(args)}: {(proc.stdout or proc.stderr)[:1200]}")
        return json.loads(proc.stdout or "{}")
