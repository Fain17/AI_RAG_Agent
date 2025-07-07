import pytest
from fastapi.testclient import TestClient
from app.main import app

client = TestClient(app)


def test_read_main():
    """Test that the main app is accessible"""
    response = client.get("/docs")
    assert response.status_code == 200


def test_app_title():
    """Test that the app has the correct title"""
    assert app.title == "Business Logic API"
    assert app.version == "1.0"


@pytest.mark.asyncio
async def test_health_check():
    """Basic health check test"""
    # This is a placeholder test - you can add actual health endpoints later
    assert True
