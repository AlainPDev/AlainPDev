"""
app/report.py
----------------
Formats output from either module (scanner.py or detector.py) into
consistent, readable reports. Keeping formatting separate from logic
means scanner.py and detector.py never need to know or care how their
results get displayed — you could swap this for JSON/HTML output later
without touching either core module.
"""

import json
from dataclasses import asdict

from app.models import DetectionResult, PortState, ScanResult


# ---------------------------------------------------------------------------
# Offensive: scan report formatting
# ---------------------------------------------------------------------------

def format_scan_report(result: ScanResult) -> str:
    lines = [
        f"\n=== Scan Report: {result.target} ===",
        f"Duration: {result.duration_seconds:.2f}s | "
        f"Ports scanned: {len(result.ports)} | "
        f"Open: {len(result.open_ports)}",
        "",
    ]

    if not result.open_ports:
        lines.append("No open ports found.")
        return "\n".join(lines)

    lines.append(f"{'PORT':<8}{'STATE':<10}BANNER")
    lines.append("-" * 60)
    for port_result in result.open_ports:
        banner = (port_result.banner or "").replace("\n", " ")[:60]
        lines.append(f"{port_result.port:<8}{port_result.state.value:<10}{banner}")

    return "\n".join(lines)


# ---------------------------------------------------------------------------
# Defensive: detection report formatting
# ---------------------------------------------------------------------------

def format_detection_report(result: DetectionResult) -> str:
    lines = [
        f"\n=== Auth Log Analysis: {result.log_source} ===",
        f"Total auth events parsed: {result.total_events}",
        f"Suspected brute-force sources: {len(result.suspects)}",
        "",
    ]

    if not result.suspects:
        lines.append("No brute-force patterns detected.")
        return "\n".join(lines)

    for suspect in result.suspects:
        lines.append(f"⚠️  {suspect.source_ip}")
        lines.append(
            f"   {suspect.attempt_count} failed attempts in "
            f"{suspect.window_minutes:.1f} min "
            f"({suspect.window_start.strftime('%H:%M:%S')} - "
            f"{suspect.window_end.strftime('%H:%M:%S')})"
        )
        usernames = ", ".join(suspect.usernames_tried) or "unknown"
        attack_type = "targeted" if len(suspect.usernames_tried) <= 1 else "credential spray"
        lines.append(f"   Usernames tried: {usernames} ({attack_type})")
        lines.append("")

    return "\n".join(lines)


# ---------------------------------------------------------------------------
# Shared: JSON export (useful for feeding into other tools later)
# ---------------------------------------------------------------------------

def _json_default(obj):
    """Handle enums and datetimes when serializing dataclasses to JSON."""
    if hasattr(obj, "value"):  # Enum
        return obj.value
    if hasattr(obj, "isoformat"):  # datetime
        return obj.isoformat()
    raise TypeError(f"Object of type {type(obj)} is not JSON serializable")


def to_json(result: ScanResult | DetectionResult) -> str:
    """Serialize either result type to JSON — handy for piping into a
    future compliance/reporting module or storing in data/."""
    return json.dumps(asdict(result), default=_json_default, indent=2)
