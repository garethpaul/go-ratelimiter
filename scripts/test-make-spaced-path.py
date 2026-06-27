#!/usr/bin/env python3
import os
from pathlib import Path
import shutil
import subprocess
import tempfile


ROOT = Path(__file__).resolve().parents[1]
CHILD_MARKER = "GO_RATELIMITER_MAKE_SPACE_CHILD"


def make_command(copied_root, *arguments, environment=None):
    return subprocess.run(
        ["make", "--no-print-directory", "-f", str(copied_root / "Makefile"), *arguments],
        cwd=copied_root.parent,
        env=environment,
        capture_output=True,
        text=True,
        timeout=180,
    )


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

        root_result = subprocess.run(
            [
                "make",
                "--no-print-directory",
                "-f",
                str(copied_root / "Makefile"),
                "--eval",
                'print-root:;@printf "%s\\n" "$(ROOT)"',
                "print-root",
            ],
            cwd=caller_root,
            capture_output=True,
            text=True,
            check=True,
            timeout=30,
        )
        if root_result.stdout.strip() != str(copied_root.resolve()):
            raise AssertionError("spaced Makefile invocation resolved the wrong repository root")

        extra_makefile = temporary_root / "extra.mk"
        extra_makefile.write_text("all:\n\t@:\n")
        hostile_environment = os.environ.copy()
        hostile_environment["MAKEFILES"] = str(extra_makefile)
        hostile_result = make_command(copied_root, "check", environment=hostile_environment)
        if hostile_result.returncode == 0 or "MAKEFILES must be empty" not in hostile_result.stderr:
            raise AssertionError("MAKEFILES contamination must fail closed")

        list_override_result = make_command(copied_root, "MAKEFILE_LIST=hostile", "check")
        if list_override_result.returncode == 0 or "MAKEFILE_LIST must not be overridden" not in list_override_result.stderr:
            raise AssertionError("MAKEFILE_LIST override must fail closed")

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
