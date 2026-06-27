#!/usr/bin/env python3
import os
from pathlib import Path
import shutil
import subprocess
import tempfile


ROOT = Path(__file__).resolve().parents[1]
CHILD_MARKER = "GO_RATELIMITER_MAKE_SPACE_CHILD"


def main():
    if os.environ.get(CHILD_MARKER) == "1":
        return

    with tempfile.TemporaryDirectory(prefix="go-ratelimiter-make-space-") as temporary_directory:
        temporary_root = Path(temporary_directory)
        copied_root = temporary_root / "repository with spaces"
        caller_root = temporary_root / "external caller"
        shutil.copytree(
            ROOT,
            copied_root,
            ignore=shutil.ignore_patterns(".git", "__pycache__", "*.pyc", "*.pyo"),
        )
        caller_root.mkdir()

        environment = os.environ.copy()
        environment[CHILD_MARKER] = "1"
        subprocess.run(
            ["make", "-f", str(copied_root / "Makefile"), "check"],
            cwd=caller_root,
            env=environment,
            check=True,
            timeout=180,
        )


if __name__ == "__main__":
    main()
