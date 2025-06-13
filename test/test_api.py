import pytest
import requests

BASE_URL = "http://localhost:8080"

def test_get_posts_empty():
    resp = requests.get(f"{BASE_URL}/posts")
    assert resp.status_code == 200
    data = resp.json()
    assert isinstance(data, list)
    assert data == []

def test_create_and_get_post():
    payload = {"title": "Test", "content": "Content"}
    # Create post
    resp = requests.post(f"{BASE_URL}/posts", json=payload)
    assert resp.status_code == 200
    created = resp.json()
    assert "ID" in created or "id" in created
    post_id = created.get("ID") or created.get("id")
    assert isinstance(post_id, int)
    assert created["title"] == payload["title"]
    assert created["content"] == payload["content"]
    # Get the post
    resp = requests.get(f"{BASE_URL}/posts/{post_id}")
    assert resp.status_code == 200
    fetched = resp.json()
    assert fetched["title"] == payload["title"]
    assert fetched["content"] == payload["content"]

def test_get_post_not_found():
    resp = requests.get(f"{BASE_URL}/posts/999999")
    assert resp.status_code == 404
    err = resp.json()
    assert err.get("error") == "post not found"

def test_create_post_bad_request():
    # Send invalid JSON
    headers = {"Content-Type": "application/json"}
    resp = requests.post(f"{BASE_URL}/posts", data="{invalid_json", headers=headers)
    assert resp.status_code == 400

def test_get_post_invalid_id():
    resp = requests.get(f"{BASE_URL}/posts/abc")
    assert resp.status_code == 404
    err = resp.json()
    assert err.get("error") == "post not found"

def test_posts_order():
    # Create two posts
    one = {"title": "First", "content": "1"}
    two = {"title": "Second", "content": "2"}
    requests.post(f"{BASE_URL}/posts", json=one)
    requests.post(f"{BASE_URL}/posts", json=two)
    # Fetch list
    resp = requests.get(f"{BASE_URL}/posts")
    assert resp.status_code == 200
    posts = resp.json()
    titles = [p.get("title") for p in posts]
    assert titles[0] == "Second"
    assert titles[1] == "First"