"""
app/main.py
--------------
CLI entrypoint. Wires the offensive scanner and defensive detector together
as subcommands of one tool:

    sentinel scan <target> [--ports 22,80,443] [--full]
    sentinel detect <logfile> [--threshold 5] [--window 10]

Same pattern as CloudSentry: one entrypoint, one subcommand per capability,
shared config and reporting underneath.
"""

import argparse
import asyncio
import sys

from app.config import load_config
from app.detector import analyze_log_file
from app.report import format_detection_report, format_scan_report, to_json
from app.scanner import resolve_target, scan_target


def _parse_ports(raw: str | None) -> list[int] | None:
    """Parses '22,80,443' or '1-1024' style port specs into a list of ints."""
    if not raw:
        return None
    ports: list[int] = []
    for chunk in raw.split(","):
        chunk = chunk.strip()
        if "-" in chunk:
            start, end = chunk.split("-")
            ports.extend(range(int(start), int(end) + 1))
        else:
            ports.append(int(chunk))
    return ports


def _run_scan(args: argparse.Namespace) -> int:
    config = load_config().scanner

    try:
        resolved_ip = resolve_target(args.target)
        print(f"Target {args.target} resolved to {resolved_ip}")
    except Exception:
        print(f"Could not resolve {args.target}, attempting scan anyway...")

    if args.full:
        ports = list(range(1, 65536))
        print("⚠️  Full port sweep (1-65535) — this is loud and slow. "
              "Only run this against systems you own or are authorized to test.")
    else:
        ports = _parse_ports(args.ports)

    result = asyncio.run(scan_target(args.target, config, ports=ports))

    if args.json:
        print(to_json(result))
    else:
        print(format_scan_report(result))

    return 0


def _run_detect(args: argparse.Namespace) -> int:
    config = load_config().detector
    config.failed_attempt_threshold = args.threshold
    config.time_window_minutes = args.window

    result = analyze_log_file(args.logfile, config)

    if args.json:
        print(to_json(result))
    else:
        print(format_detection_report(result))

    return 1 if result.suspects else 0


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="sentinel",
        description="Network Sentinel — offensive port scanning + defensive "
                     "brute-force detection in one tool.",
    )
    subparsers = parser.add_subparsers(dest="command", required=True)

    scan_parser = subparsers.add_parser("scan", help="Scan a target for open ports (offensive)")
    scan_parser.add_argument("target", help="Hostname or IP to scan")
    scan_parser.add_argument("--ports", help="Comma/range port spec, e.g. '22,80,443' or '1-1024'")
    scan_parser.add_argument("--full", action="store_true", help="Scan all 65535 ports (slow, loud)")
    scan_parser.add_argument("--json", action="store_true", help="Output JSON instead of a table")
    scan_parser.set_defaults(handler=_run_scan)

    detect_parser = subparsers.add_parser("detect", help="Analyze an auth log for brute-force attempts (defensive)")
    detect_parser.add_argument("logfile", help="Path to an SSH-style auth log file")
    detect_parser.add_argument("--threshold", type=int, default=5, help="Failed attempts to trigger a flag")
    detect_parser.add_argument("--window", type=int, default=10, help="Time window in minutes")
    detect_parser.add_argument("--json", action="store_true", help="Output JSON instead of a table")
    detect_parser.set_defaults(handler=_run_detect)

    return parser


def main(argv=None) -> int:
    parser = build_parser()
    args = parser.parse_args(argv)
    return args.handler(args)


if __name__ == "__main__":
    sys.exit(main())
