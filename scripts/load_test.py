#!/usr/bin/env python3
"""Async load tester for Python (Moscow time) and Go (Wordle) apps.
Usage:
  # Python app - run for 30 seconds with 50 clients
  python3 scripts/load_test.py --url http://localhost:8000 --mode python --clients 50 --duration 30

  # Go app - run for 60 seconds with 50 clients
  python3 scripts/load_test.py --url http://localhost:8080 --mode go --clients 50 --duration 60

  # Or use fixed count (legacy)
  python3 scripts/load_test.py --url http://localhost:8000 --mode python --clients 50 --requests 100
"""
import argparse
import asyncio
import random
import re
import time
from typing import Optional

import httpx

WORDS = ["AUDIO", "ADIEU", "ARISE", "RAISE", "SLATE", "HOUSE", "MOUSE", "CRANE", "POUND", "CRATE",
         "TRAIN", "BRAIN", "LEMON", "PEACH", "BREAD", "PLANT", "SMART", "BEACH", "CLOUD", "SPORT",
         "STONE", "LIGHT", "NIGHT", "FIGHT", "TIGHT"]

async def python_worker(client: httpx.AsyncClient, url: str, duration: Optional[float] = None, requests: Optional[int] = None):
    """Load test for Python Moscow time app."""
    endpoints = ["/", "/health", "/metrics"]
    request_count = 0

    if duration is not None:
        # Duration-based: run requests for N seconds
        start = time.time()
        while time.time() - start < duration:
            path = random.choice(endpoints)
            try:
                await client.get(url + path, timeout=10.0)
                request_count += 1
            except Exception:
                pass
    else:
        # Count-based: run fixed number of requests
        for _ in range(requests or 1):
            path = random.choice(endpoints)
            try:
                await client.get(url + path, timeout=10.0)
                request_count += 1
            except Exception:
                pass

    return request_count

async def go_worker(client: httpx.AsyncClient, url: str, duration: Optional[float] = None, games: Optional[int] = None):
    """Load test for Go Wordle app - create games and submit guesses."""
    game_count = 0

    if duration is not None:
        # Duration-based: run games for N seconds
        start = time.time()
        while time.time() - start < duration:
            try:
                # Create a new game
                r = await client.get(url + "/", follow_redirects=True, timeout=10.0)
                # Extract game ID from URL or response
                game_id = None
                if r.history:
                    # Check redirect location
                    location = r.history[-1].headers.get("location", "")
                    match = re.search(r"/game/([^/]+)", location)
                    if match:
                        game_id = match.group(1)

                if not game_id:
                    continue

                game_count += 1

                # Make 2-4 guesses per game
                num_guesses = random.randint(2, 4)
                for _ in range(num_guesses):
                    guess = random.choice(WORDS)
                    try:
                        await client.post(
                            url + "/guess",
                            data={"game_id": game_id, "guess": guess},
                            timeout=10.0,
                            follow_redirects=False
                        )
                    except Exception:
                        pass
            except Exception:
                pass
    else:
        # Count-based: run fixed number of games
        for _ in range(games or 1):
            try:
                # Create a new game
                r = await client.get(url + "/", follow_redirects=True, timeout=10.0)
                # Extract game ID from URL or response
                game_id = None
                if r.history:
                    # Check redirect location
                    location = r.history[-1].headers.get("location", "")
                    match = re.search(r"/game/([^/]+)", location)
                    if match:
                        game_id = match.group(1)

                if not game_id:
                    continue

                game_count += 1

                # Make 2-4 guesses per game
                num_guesses = random.randint(2, 4)
                for _ in range(num_guesses):
                    guess = random.choice(WORDS)
                    try:
                        await client.post(
                            url + "/guess",
                            data={"game_id": game_id, "guess": guess},
                            timeout=10.0,
                            follow_redirects=False
                        )
                    except Exception:
                        pass
            except Exception:
                pass

    return game_count

async def run_python(url: str, clients: int, duration: Optional[float] = None, total_requests: Optional[int] = None):
    """Run load test for Python app."""
    async with httpx.AsyncClient(follow_redirects=True) as client:
        tasks = [
            asyncio.create_task(python_worker(client, url, duration=duration, requests=total_requests))
            for _ in range(clients)
        ]
        start = time.time()
        results = await asyncio.gather(*tasks)
        elapsed = time.time() - start

    total_requests = sum(results)
    print(f"Python app: Completed {total_requests} requests in {elapsed:.2f}s ({total_requests/elapsed:.0f} req/s)")

async def run_go(url: str, clients: int, duration: Optional[float] = None, total_games: Optional[int] = None):
    """Run load test for Go Wordle app."""
    async with httpx.AsyncClient(follow_redirects=False) as client:
        tasks = [
            asyncio.create_task(go_worker(client, url, duration=duration, games=total_games))
            for _ in range(clients)
        ]
        start = time.time()
        results = await asyncio.gather(*tasks)
        elapsed = time.time() - start

    total_games = sum(results)
    print(f"Go Wordle app: Completed {total_games} games in {elapsed:.2f}s ({total_games/elapsed:.0f} games/s)")

if __name__ == '__main__':
    p = argparse.ArgumentParser(description="Async load tester for Moscow time (Python) and Wordle (Go) apps")
    p.add_argument("--url", required=True, help="Base URL, e.g. http://localhost:8080")
    p.add_argument("--mode", choices=["python", "go"], default="python", help="Which app to test")
    p.add_argument("--clients", type=int, default=10, help="Number of parallel clients")
    p.add_argument("--duration", type=int, help="Duration in seconds to run load test (recommended)")
    p.add_argument("--requests", type=int, help="Total requests for Python app mode (legacy, use --duration instead)")
    p.add_argument("--games", type=int, help="Total games for Go app mode (legacy, use --duration instead)")
    args = p.parse_args()

    url = args.url.rstrip('/')

    # Prefer duration, fall back to count-based
    if args.mode == "python":
        if args.duration:
            asyncio.run(run_python(url, args.clients, duration=args.duration))
        else:
            asyncio.run(run_python(url, args.clients, total_requests=args.requests or 100))
    else:  # go
        if args.duration:
            asyncio.run(run_go(url, args.clients, duration=args.duration))
        else:
            asyncio.run(run_go(url, args.clients, total_games=args.games or 100))
