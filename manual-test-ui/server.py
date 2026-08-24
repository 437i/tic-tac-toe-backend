#!/usr/bin/env python3
import http.client
import http.server
import json
import os
import socketserver
import urllib.parse

BACKEND_HOST = os.environ.get("TICTACTOE_BACKEND_HOST", "127.0.0.1")
BACKEND_PORT = int(os.environ.get("TICTACTOE_BACKEND_PORT", "8080"))
HOST = os.environ.get("TICTACTOE_UI_HOST", "127.0.0.1")
PORT = int(os.environ.get("TICTACTOE_UI_PORT", "3000"))
ROOT = os.path.join(os.path.dirname(__file__), "static")

class Handler(http.server.SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory=ROOT, **kwargs)

    def _proxy(self):
        parsed = urllib.parse.urlsplit(self.path)
        path = parsed.path.removeprefix("/api") or "/"
        if parsed.query:
            path += "?" + parsed.query
        length = int(self.headers.get("Content-Length", "0"))
        body = self.rfile.read(length) if length else None
        headers = {}
        for key in ("Authorization", "Content-Type", "Accept"):
            value = self.headers.get(key)
            if value:
                headers[key] = value
        conn = http.client.HTTPConnection(BACKEND_HOST, BACKEND_PORT, timeout=30)
        try:
            conn.request(self.command, path, body=body, headers=headers)
            response = conn.getresponse()
            data = response.read()
            self.send_response(response.status, response.reason)
            content_type = response.getheader("Content-Type") or "application/json"
            self.send_header("Content-Type", content_type)
            self.send_header("Content-Length", str(len(data)))
            self.end_headers()
            self.wfile.write(data)
        except Exception as exc:
            payload = json.dumps({"error": f"proxy error: {exc}"}).encode()
            self.send_response(502)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(payload)))
            self.end_headers()
            self.wfile.write(payload)
        finally:
            conn.close()

    def do_GET(self):
        if self.path.startswith("/api/") or self.path == "/api":
            self._proxy(); return
        super().do_GET()

    def do_POST(self):
        if self.path.startswith("/api/") or self.path == "/api":
            self._proxy(); return
        self.send_error(405)

    def log_message(self, fmt, *args):
        print(f"[ui] {self.address_string()} - {fmt % args}")

class ThreadedHTTPServer(socketserver.ThreadingMixIn, http.server.HTTPServer):
    daemon_threads = True

if __name__ == "__main__":
    print(f"UI:      http://{HOST}:{PORT}")
    print(f"Backend: http://{BACKEND_HOST}:{BACKEND_PORT}")
    ThreadedHTTPServer((HOST, PORT), Handler).serve_forever()
