"""
app/models.py
-------------
Shared data structures for both the offensive scanner and defensive
detector. Keeping these in one place means report.py can format output
from either module without caring which one produced it.
"""

from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum

# ---------------------------------------------------------------------
# Offensive: port scan
# ---------------------------------------------------------------------

class PortState(str, Enum):
    OPEN = "open"
    CLOSED = "closed"
    FILTERED = "filtered" # timed out - likely firewalled, not necessarily closed

@dataclass
class PortResult:
    port: int
    state: PortState
    banner: str | None = None # raw service banner, if grabbed

@dataclass
class ScanResult:
    target: str
    started_at: datetime
    finished_at: datetime
    ports: list[PortResult] = field(default_factory=list)

    @property
    def open_ports(self) -> list[PortResult]:
        return [p for p in self.ports if p.state == PortState.OPEN]

    @property
    def duration_seconds(self) -> float:
        return (self.finished_at - self.started_at).total_seconds()

# ------------------------------------------------------------------
# Defensive: Auth Log Analysis
# -----------------------------------------------------------------

class AuthhEventType(str, Enum):
    FAILED_PASSWORD = "failed_password"
    ACCEPTED_PASSWORD = "accepted_password"
    INVALID_USER = "invalid_user"

@dataclass
class AuthEvent:
    timestamp: datetime
    event_type: AuthEventType
    source_ip: str
    username: str | None = None
    raw_line: str = ""

@dataclass
class SuspectedBruteForce:
    source_ip: str
    attempted_count: int
    window_start: datetime
    window_end: datetime
    usernames_tried: list[str] = field(default_factory=list)

    @property
    def window_minutes(self) -> float:
        return (self.window_end - self.window_start).total_seconds() / 60

@dataclass
class DetectionResults:
    log_source: str
    analyzed_at: datetime
    total_events: int
    suspects: list[SuspectedBruteForce] = field(default_factory=list)
