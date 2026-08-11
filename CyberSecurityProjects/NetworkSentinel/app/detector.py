"""
app/detector.py
------------------
DEFENSIVE MODULE: auth-log parser + brute-force detector.

Parses SSH-style auth log lines (the format found in /var/log/auth.log on
Debian/Ubuntu or /var/log/secure on RHEL/CentOS), extracts failed-login
events, and flags source IPs that cross a failed-attempt threshold within
a sliding time window — the same core logic behind real SOC detection
rules and SIEM correlation searches.

Expected log line formats (standard OpenSSH via syslog):
    Mar 15 02:14:22 web01 sshd[1234]: Failed password for root from 203.0.113.7 port 51422 ssh2
    Mar 15 02:14:25 web01 sshd[1234]: Failed password for invalid user admin from 203.0.113.7 port 51430 ssh2
    Mar 15 02:15:01 web01 sshd[1234]: Accepted password for deploy from 198.51.100.4 port 40112 ssh2
"""

import re
from datetime import datetime, timedelta
from collections import defaultdict

from app.config import DetectorConfig
from app.models import AuthEvent, AuthEventType, DetectionResult, SuspectedBruteForce

# Matches both "Failed password for USER" and "Failed password for invalid user USER"
_FAILED_PATTERN = re.compile(
    r"^(?P<timestamp>\w{3}\s+\d+\s+\d{2}:\d{2}:\d{2})\s+\S+\s+sshd\[\d+\]:\s+"
    r"Failed password for (?:invalid user )?(?P<username>\S+) from (?P<ip>\d{1,3}(?:\.\d{1,3}){3})"
)

_ACCEPTED_PATTERN = re.compile(
    r"^(?P<timestamp>\w{3}\s+\d+\s+\d{2}:\d{2}:\d{2})\s+\S+\s+sshd\[\d+\]:\s+"
    r"Accepted password for (?P<username>\S+) from (?P<ip>\d{1,3}(?:\.\d{1,3}){3})"
)

_INVALID_USER_MARKER = "invalid user"


def _parse_timestamp(raw: str, config: DetectorConfig, year: int) -> datetime:
    """Syslog timestamps omit the year, so we assume the current year.
    Good enough for a portfolio/lab tool; a production system would need
    to handle year rollover (Dec 31 -> Jan 1) explicitly."""
    parsed = datetime.strptime(raw, config.log_timestamp_format)
    return parsed.replace(year=year)


def parse_log_line(line: str, config: DetectorConfig, year: int) -> AuthEvent | None:
    """Parse a single auth log line into an AuthEvent, or None if the line
    doesn't match a recognized pattern (most lines in a real log won't be
    auth events at all, and that's expected)."""
    line = line.strip()
    if not line:
        return None

    failed_match = _FAILED_PATTERN.match(line)
    if failed_match:
        event_type = (
            AuthEventType.INVALID_USER
            if _INVALID_USER_MARKER in line
            else AuthEventType.FAILED_PASSWORD
        )
        return AuthEvent(
            timestamp=_parse_timestamp(failed_match.group("timestamp"), config, year),
            event_type=event_type,
            source_ip=failed_match.group("ip"),
            username=failed_match.group("username"),
            raw_line=line,
        )

    accepted_match = _ACCEPTED_PATTERN.match(line)
    if accepted_match:
        return AuthEvent(
            timestamp=_parse_timestamp(accepted_match.group("timestamp"), config, year),
            event_type=AuthEventType.ACCEPTED_PASSWORD,
            source_ip=accepted_match.group("ip"),
            username=accepted_match.group("username"),
            raw_line=line,
        )

    return None


def parse_log_file(path: str, config: DetectorConfig, year: int | None = None) -> list[AuthEvent]:
    """Parse an entire log file into a list of AuthEvents, skipping any
    lines that don't match a known pattern."""
    target_year = year or datetime.now().year
    events: list[AuthEvent] = []

    with open(path, "r", errors="replace") as f:
        for line in f:
            event = parse_log_line(line, config, target_year)
            if event:
                events.append(event)

    return events


def detect_brute_force(
    events: list[AuthEvent],
    config: DetectorConfig,
) -> list[SuspectedBruteForce]:
    """Group failed-login events by source IP, then flag any IP whose
    failed attempts within a sliding time window meet or exceed the
    configured threshold.

    This uses a simple sliding-window scan per IP rather than fixed
    buckets, so a burst that straddles a bucket boundary isn't missed —
    a real gap that fixed-interval bucketing tools often have."""
    failed_events = [
        e for e in events
        if e.event_type in (AuthEventType.FAILED_PASSWORD, AuthEventType.INVALID_USER)
    ]

    by_ip: dict[str, list[AuthEvent]] = defaultdict(list)
    for event in failed_events:
        by_ip[event.source_ip].append(event)

    window = timedelta(minutes=config.time_window_minutes)
    suspects: list[SuspectedBruteForce] = []

    for ip, ip_events in by_ip.items():
        ip_events.sort(key=lambda e: e.timestamp)

        # Sliding window: for each event, count how many failures fall
        # within `window` of it (looking forward).
        for i, start_event in enumerate(ip_events):
            window_events = [
                e for e in ip_events[i:]
                if e.timestamp - start_event.timestamp <= window
            ]
            if len(window_events) >= config.failed_attempt_threshold:
                suspects.append(SuspectedBruteForce(
                    source_ip=ip,
                    attempt_count=len(window_events),
                    window_start=window_events[0].timestamp,
                    window_end=window_events[-1].timestamp,
                    usernames_tried=sorted({e.username for e in window_events if e.username}),
                ))
                break  # one flag per IP is enough; avoid duplicate overlapping windows

    return suspects


def analyze_log_file(path: str, config: DetectorConfig) -> DetectionResult:
    """End-to-end: parse a log file and run brute-force detection on it."""
    events = parse_log_file(path, config)
    suspects = detect_brute_force(events, config)

    return DetectionResult(
        log_source=path,
        analyzed_at=datetime.now(),
        total_events=len(events),
        suspects=suspects,
    )
