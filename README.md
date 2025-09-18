# Go RA GitLab Keys - JWKS Serverless Function

This project provides a serverless function deployed on Netlify that serves JSON Web Key Set (JWKS) data for OAuth discovery. The function exposes a GET endpoint at `/oauth/discovery/keys` that returns RSA public keys used for JWT signature verification.

## Features

- **Serverless Function**: Deployed as a Netlify Go function
- **JWKS Endpoint**: Serves JSON Web Key Set at `/oauth/discovery/keys`
- **CORS Support**: Properly configured CORS headers for cross-origin requests
- **Caching**: HTTP caching headers for optimal performance
- **Security Headers**: Includes security headers for protection

## Project Structure

```
.
├── netlify.toml                    # Netlify configuration and routing
├── go.mod                          # Go module definition
├── netlify/
│   └── functions/
│       └── keys/
│           └── main.go            # Main serverless function handler
└── README.md                      # This file
```

## API Endpoint

### GET /oauth/discovery/keys

Returns a JSON Web Key Set (JWKS) containing RSA public keys for JWT signature verification.

**Response Format:**
```json
{
  "keys": [
    {
      "kty": "RSA",
      "kid": "Q93ngDAyTNgdjaPblfooyeT00CFLJV5pWulhwAEg8Sw",
      "e": "AQAB",
      "n": "xTcVMXDqrYep...",
      "use": "sig",
      "alg": "RS256"
    },
    {
      "kty": "RSA", 
      "kid": "CjZ2bP54fm1lEkeHGo_E4UVyc4MbN4fye1i6DrxFaqQ",
      "e": "AQAB",
      "n": "nN8D8DxzqoQJ...",
      "use": "sig",
      "alg": "RS256"
    }
  ]
}
```

**Response Headers:**
- `Content-Type: application/json`
- `Cache-Control: public, max-age=3600` (1 hour cache)
- `Access-Control-Allow-Origin: *` (CORS enabled)

## Deployment

### Prerequisites

1. **Netlify Account**: Sign up at [netlify.com](https://netlify.com)
2. **Git Repository**: Code must be in a Git repository (GitHub, GitLab, etc.)
3. **Go 1.21+**: Ensure Go version compatibility

### Deploy to Netlify

#### Option 1: Netlify Web Interface

1. **Connect Repository**:
   - Log into Netlify
   - Click "Add new site" → "Import an existing project"
   - Connect your Git provider and select this repository

2. **Configure Build Settings**:
   - **Build command**: Leave empty (auto-detected)
   - **Publish directory**: Leave empty
   - **Functions directory**: `netlify/functions` (auto-detected)

3. **Deploy**:
   - Click "Deploy site"
   - Netlify will automatically build and deploy your Go function

#### Option 2: Netlify CLI

```bash
# Install Netlify CLI
npm install -g netlify-cli

# Login to Netlify
netlify login

# Initialize the site
netlify init

# Deploy
netlify deploy

# Deploy to production
netlify deploy --prod
```

### Environment Variables

No environment variables are required for this deployment as the RSA keys are embedded in the source code.

## Configuration Details

### netlify.toml

The `netlify.toml` file configures:

- **Build Environment**: Go 1.21
- **URL Routing**: Maps `/oauth/discovery/keys` to the serverless function
- **CORS Headers**: Enables cross-origin requests
- **Security Headers**: Adds security-focused HTTP headers
- **Caching**: Sets appropriate cache headers

### Go Function

The main function (`netlify/functions/keys/main.go`) includes:

- **HTTP Method Validation**: Only accepts GET and OPTIONS requests
- **CORS Handling**: Responds to preflight OPTIONS requests
- **Error Handling**: Proper error responses with appropriate status codes
- **JSON Serialization**: Converts Go structs to JSON response
- **Logging**: Request logging for debugging

## Testing Locally

### Using Netlify CLI

```bash
# Install dependencies
go mod tidy

# Start local development server
netlify dev

# The function will be available at:
# http://localhost:8888/oauth/discovery/keys
```

### Manual Testing

```bash
# Test the endpoint
curl http://localhost:8888/oauth/discovery/keys

# Test CORS preflight
curl -X OPTIONS http://localhost:8888/oauth/discovery/keys \
  -H "Access-Control-Request-Method: GET" \
  -H "Access-Control-Request-Headers: Content-Type"
```

## Security Considerations

1. **Public Keys Only**: This endpoint serves public keys only - no private key material
2. **CORS Policy**: Currently allows all origins (`*`) - consider restricting in production
3. **Rate Limiting**: Consider implementing rate limiting for production use
4. **HTTPS Only**: Netlify provides HTTPS by default
5. **Security Headers**: Includes standard security headers

## Key Management

The RSA keys are currently hardcoded in the source code. For production use, consider:

1. **Environment Variables**: Store keys in Netlify environment variables
2. **External Key Management**: Integrate with AWS KMS, HashiCorp Vault, etc.
3. **Key Rotation**: Implement automated key rotation mechanisms
4. **Multiple Environments**: Separate keys for development, staging, and production

## Monitoring and Maintenance

### Netlify Functions Dashboard

Monitor function performance in the Netlify dashboard:
- Function invocations
- Error rates
- Response times
- Bandwidth usage

### Logging

Function logs are available in:
- Netlify Functions dashboard
- Real-time logs via `netlify logs`

### Updates and Deployment

1. Make changes to the code
2. Commit and push to your Git repository
3. Netlify will automatically rebuild and deploy

## Troubleshooting

### Common Issues

1. **Function Not Found (404)**:
   - Verify `netlify.toml` redirects are correct
   - Check function is in `netlify/functions/keys/` directory
   - Ensure `main.go` is present and compiles

2. **CORS Errors**:
   - Verify CORS headers in both function and `netlify.toml`
   - Check preflight OPTIONS handling

3. **Build Failures**:
   - Verify Go version compatibility
   - Check `go.mod` dependencies
   - Review build logs in Netlify dashboard

### Debug Commands

```bash
# Check function locally
netlify functions:list

# View function logs
netlify logs

# Test function build
go build ./netlify/functions/keys/

# Validate JSON output
go run ./netlify/functions/keys/ | jq .
```

## License

This project is provided as-is for internal use within Rockwell Automation.