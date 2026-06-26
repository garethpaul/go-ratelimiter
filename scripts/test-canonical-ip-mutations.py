#!/usr/bin/env python3

import subprocess
from pathlib import Path


root = Path(__file__).resolve().parent.parent
source_path = root / "libstring" / "libstring.go"
original = source_path.read_text()
mutations = (
    (
        "RemoteAddr canonicalization",
        "\treturn ip.String()\n}\n\nfunc ipAddrFromHeaderValue",
        "\treturn host\n}\n\nfunc ipAddrFromHeaderValue",
        "TestLimitHandlerSharesBucketAcrossEquivalentIPv6Forms|TestRemoteIPCanonicalizesEquivalentIPv6Forms",
    ),
    (
        "proxy header canonicalization",
        "\treturn ip.String()\n}\n\nfunc ipAddrFromForwardedFor",
        "\treturn value\n}\n\nfunc ipAddrFromForwardedFor",
        "TestRemoteIPCanonicalizesEquivalentIPv6Forms",
    ),
)

try:
    for name, target, replacement, test_pattern in mutations:
        if original.count(target) != 1:
            raise SystemExit(f"missing mutation target: {name}")
        source_path.write_text(original.replace(target, replacement, 1))
        result = subprocess.run(
            ["go", "test", "./...", "-run", test_pattern],
            cwd=root,
            check=False,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
        )
        if result.returncode == 0:
            raise SystemExit(f"mutation survived: {name}\n{result.stdout}")
        source_path.write_text(original)
finally:
    source_path.write_text(original)

print(f"{len(mutations)} canonical IP hostile mutations rejected")
