from http.server import SimpleHTTPRequestHandler, HTTPServer
from urllib.parse import urlparse
import signal
import sys

# Define the custom MIME types
custom_mime_types = {
    ".tgz": "application/x-gzip",
}

class CustomMimeTypesHTTPRequestHandler(SimpleHTTPRequestHandler):
    def guess_type(self, path):
        url = urlparse(path)
        file_ext = url.path
        pos = file_ext.rfind('.')
        if pos != -1:
            file_ext = file_ext[pos:]
        else:
            file_ext = ""

        # Check if the file extension has a custom MIME type
        if file_ext in custom_mime_types:
            return custom_mime_types[file_ext]

        # Fallback to the default MIME type guessing
        return super().guess_type(path)

# Set the handler to use the custom class
handler = CustomMimeTypesHTTPRequestHandler

port = 8666

try:
    print(f'Serving on port {port} directory {sys.argv[1]}')
    with HTTPServer(('', port), lambda *_: CustomMimeTypesHTTPRequestHandler(*_, directory=sys.argv[1])) as server:
        def stop():
            print('SIGTERM, shutting down.')
            server.shutdown()
            sys.exit(0)

        signal.signal(signal.SIGTERM, stop)
        server.serve_forever()
except KeyboardInterrupt:
    print('KeyboardInterrupt, shutting down.')