#!/usr/bin/env python3
"""Replace {{PLACEHOLDER}} in MDL templates from environment variables."""
import os
import re
import sys


def render(text: str) -> str:
    missing = set()

    def repl(match: re.Match[str]) -> str:
        key = match.group(1)
        val = os.environ.get(key)
        if val is None or val == "":
            missing.add(key)
            return match.group(0)
        return val

    out = re.sub(r"\{\{([A-Z0-9_]+)\}\}", repl, text)
    if missing:
        keys = ", ".join(sorted(missing))
        sys.exit(f"missing env for placeholders: {keys}")
    return out


def main() -> None:
    if len(sys.argv) != 3:
        sys.exit("usage: render-mdl.py <template.mdl> <output.mdl>")
    src, dst = sys.argv[1], sys.argv[2]
    with open(src, encoding="utf-8") as f:
        rendered = render(f.read())
    with open(dst, "w", encoding="utf-8") as f:
        f.write(rendered)


if __name__ == "__main__":
    main()
