"""
Unit tests for the Moscow Time Display FastAPI application.
"""
import pytest
from fastapi.testclient import TestClient
from datetime import datetime
from zoneinfo import ZoneInfo

from app import app

client = TestClient(app)


class TestMoscowTimeEndpoint:
    """Test suite for the main Moscow time display endpoint."""

    def test_root_endpoint_returns_200(self):
        """Test that the root endpoint returns a successful response."""
        response = client.get("/")
        assert response.status_code == 200

    def test_root_endpoint_returns_html(self):
        """Test that the root endpoint returns HTML content."""
        response = client.get("/")
        assert response.headers["content-type"] == "text/html; charset=utf-8"

    def test_html_contains_moscow_time(self):
        """Test that the HTML response contains Moscow time elements."""
        response = client.get("/")
        html_content = response.text

        # Check for key elements in the page
        assert "Moscow" in html_content
        assert "MSK" in html_content
        assert "Current Time" in html_content

    def test_displayed_time_is_current(self):
        """Test that the displayed time is approximately current Moscow time."""
        response = client.get("/")
        html_content = response.text

        # Get current Moscow time
        moscow_tz = ZoneInfo("Europe/Moscow")
        current_moscow_time = datetime.now(moscow_tz)

        # Check if the current hour is in the response
        current_hour = current_moscow_time.strftime("%H")
        assert current_hour in html_content

    def test_page_has_refresh_button(self):
        """Test that the page includes a refresh button."""
        response = client.get("/")
        html_content = response.text
        assert "Refresh" in html_content or "refresh" in html_content


class TestHealthCheckEndpoint:
    """Test suite for the health check endpoint."""

    def test_health_endpoint_returns_200(self):
        """Test that the health endpoint returns a successful response."""
        response = client.get("/health")
        assert response.status_code == 200

    def test_health_endpoint_returns_json(self):
        """Test that the health endpoint returns JSON content."""
        response = client.get("/health")
        assert response.headers["content-type"] == "application/json"

    def test_health_endpoint_structure(self):
        """Test that the health endpoint returns the correct JSON structure."""
        response = client.get("/health")
        data = response.json()

        assert "status" in data
        assert "service" in data

    def test_health_endpoint_values(self):
        """Test that the health endpoint returns the correct values."""
        response = client.get("/health")
        data = response.json()

        assert data["status"] == "healthy"
        assert data["service"] == "moscow-time-app"


class TestApplicationBehavior:
    """Test suite for general application behavior."""

    def test_invalid_endpoint_returns_404(self):
        """Test that invalid endpoints return 404."""
        response = client.get("/invalid-endpoint")
        assert response.status_code == 404

    def test_post_to_root_not_allowed(self):
        """Test that POST requests to root are not allowed."""
        response = client.post("/")
        assert response.status_code == 405

    def test_multiple_requests_succeed(self):
        """Test that multiple consecutive requests succeed."""
        for _ in range(5):
            response = client.get("/")
            assert response.status_code == 200


class TestTimezoneHandling:
    """Test suite for timezone handling logic."""

    def test_moscow_timezone_offset(self):
        """Test that Moscow timezone is correctly configured."""
        moscow_tz = ZoneInfo("Europe/Moscow")
        moscow_time = datetime.now(moscow_tz)

        # Moscow is UTC+3
        utc_offset_hours = moscow_time.utcoffset().total_seconds() / 3600
        assert utc_offset_hours == 3.0

    def test_time_format_in_response(self):
        """Test that time is formatted correctly (HH:MM:SS)."""
        response = client.get("/")
        html_content = response.text

        # Look for time pattern HH:MM:SS using regex pattern in HTML
        import re
        time_pattern = r'\d{2}:\d{2}:\d{2}'
        assert re.search(time_pattern, html_content) is not None


if __name__ == "__main__":
    pytest.main([__file__, "-v", "--cov=app", "--cov-report=term-missing"])
