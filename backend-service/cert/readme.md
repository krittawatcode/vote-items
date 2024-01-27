## Certificates

The `certs` folder contains the SSL certificates and keys necessary for secure communication. This includes:

- `server.pem`: The public SSL certificate for the server.
- `server-key.pem`: The private key for the server's SSL certificate.

These files are used to establish secure connections between the server and clients. They should be kept secure and not shared publicly.

**Note:** The `certs` folder and its contents are not included in the version control system for security reasons. Make sure to add this folder's content to your `.gitignore` file to prevent accidentally committing sensitive data.