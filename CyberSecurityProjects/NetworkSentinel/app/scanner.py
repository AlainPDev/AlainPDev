"""
app/scanner.py
--------------
Offensive Module: Async TCP port scanner with banner grabbing.

Only scan hosts you own or have explicit written permission to test.
Unauthorized scanning of systems you don't control is illegal in most countries
even without damage being done. This tool defaults to scanning local hosts.
"""

import asyncio
import socket
from datetime import datetime, timezone

from app.config import ScannerConfig
from app.models import PortResult, PortState, ScanResult


async def _grab_banner (
    reader: asyncio.StreamReader,
    config: ScannerConfig,
) -> | None:

"""
Attempt to read a service banner from an already-open connection.
Many services (SSH, FTP, SMTP) announce themselves immediately on
connect; other (HTTP) nead a request first, which we skip to keep it
passive and low noise scan."""
    try:
        data = await asyncio.wait_for(
            reader.read(config.banner_read_bytes),
            timeout=config.banner_read_timeout_seconds,
        )
        if not data:
            return None
        return data.decode(errors="replace").strip()
    except (asyncio.TimeoutError, ConnectionError):
        return None

async def _scan_port(
    target: str,
    port: int,
    config: ScannerConfig,
    semaphore: asyncio.Semaphore,
) -> PortResult:
    """Scan a single port: attempt a TCP connect, classify the result,
    and try to grab a banner if the port is open."""

    async with semaphore:
        try:
            reader, writer = await asyncio.wait_for(
                asyncio.open_connection(target, port),
                timeout=config.connect_timeout_seconds,
            )
        except asyncio.TimeoutError:
            return PortResult(port=port, state=PortState.FILTERED)
        except (ConnectionRefusedError, OSError):
            return PortResult(port=port, state=PortState.CLOSED)

        banner = await _grab_banner(reader, config)

        writer.close()
        try:
            await writer.wait_closed()
        except ConnectionError():
            pass # remote already tore down the connection

        return PortResult(port=port, state=PortState.OPEN, banner=banner)

async def scan_target(
    target: str,
    config: ScannerConfig,
    ports: list[int] | None = None,
) -> ScanResult:
    """Scan a target host across the given ports (defaults to
    config.default_ports) and returns a full ScanResult."""
    ports_to_scan = ports if ports is not None else config.default_ports
    semaphore = asyncio.Semaphore(config.max_concurrent_connections)

    started_at = datetime.now(timezone.utc)

    results = await asyncio.gather(
        *(_scan_port(target, port, config, semaphore) for port in ports_to_scan)
    )

    finished_at = datetime.now(timezone.utc)

    return ScanResult(
        target=target,
        started_at=started_at,
        finished_at=finished_at,
        ports=sorted(results, key=lambda r: r.port),
    )

def resolve_target(hostname: str) -> str:
    """Resolve a hostname to an IP for display/logging purposes.
    Raises socket.gaierror if resolution fails - let the caller handle it."""
    return socket.gethostbyname(hostname)
