#!/usr/bin/env python3

import subprocess
from pathlib import Path


root = Path(__file__).resolve().parent.parent
source_path = root / "libstring" / "libstring.go"
original = source_path.read_text()
mutations = (
    (
        "canonical address text",
        "\treturn addr.Unmap().String()",
        "\treturn s",
        "TestLimitHandlerSharesBucketAcrossEquivalentIPv6Forms|TestRemoteIPCanonicalizesEquivalentIPv6Forms",
    ),
    (
        "scoped IPv6 zone preservation",
        "\treturn addr.Unmap().String()",
        "\treturn addr.WithZone(\"\").Unmap().String()",
        "ScopedIPv6",
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
