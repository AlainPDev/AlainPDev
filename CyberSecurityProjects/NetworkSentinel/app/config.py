"""
app/config.py
-------------
Central configuration for Network Sentinel. Both the offensive scanner and the defensive
detector pull their settings from here, so thesholds and timeouts are tuned in one place instead of
scattered across modules.
"""

import os
from dataclasses import dataclass, field


@dataclass
class ScannerConfig:
    """Settings for thhe offensive port scanner."""
    # Keep the default port list small and common - this is a learning tool,
    # not a stealth tool. Full sweeps are opt-in via CLI flag.
    default_ports: list[int] = field(default_factory=lambda: [
        21, 22, 23, 25, 53, 80, 110, 143, 443, 445,
        993, 995, 3306, 3389, 5432, 6379, 8080, 8443,
    ])
    connect_timeout_seconds: float = 1.5
    max_concurrent_connections: int = 200
    banner_read_bytes: int = 1024
    banner_read_timeout_seconds: float = 2.0

@dataclass
class DetectorConfig:
    """Settings for the defensive auth-log brute-force detector."""
    # If the same source IP has this many failed auth attempts within the
    # time window below, it's flagged as a suspected brute-force attempt.
    failed_attempt_threshold: int = 5
    time_window_minutes: int = 10
    # Common auth log line patterns (SSH-style syslog format by default).
    log_timestamp_format: str = "%b %d %H:%M:%S"

@dataclass
class AppConfig:
    scanner: ScannerConfig = field(default_factory=ScannerConfig)
    detector: DetectorConfig = field(default_factory=DetectorConfig)
    data_dir: str = os.environ.get("SENTINEL_DATA_DIR", "data")

def load_config() -> AppConfig:
    """Single entrypoint for getting config - makes it easy to later swap in env-var overrides or a config file without
    touching call sites."""
    return AppConfig()
