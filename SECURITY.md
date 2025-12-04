# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 2.x.x   | true |
| 1.x.x   | false                |

## Reporting a Vulnerability

We take the security of `go-infrastructure` seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### Where to Report

**Please do NOT report security vulnerabilities through public GitHub issues.**

Instead, please report them via one of the following methods:

1. **GitHub Security Advisories** (Preferred)
   - Navigate to the [Security tab](https://github.com/janmbaco/go-infrastructure/security/advisories/new)
   - Click "Report a vulnerability"
   - Fill in the details

2. **Email**
   - Send an email to the repository owner
   - Include the word "SECURITY" in the subject line
   - Provide detailed information about the vulnerability

### What to Include

Please include the following information in your report:

- **Type of vulnerability** (e.g., injection, authentication bypass, etc.)
- **Full paths of source file(s)** related to the vulnerability
- **Location of the affected source code** (tag/branch/commit or direct URL)
- **Step-by-step instructions** to reproduce the issue
- **Proof-of-concept or exploit code** (if possible)
- **Impact of the vulnerability** and potential attack scenarios
- **Any mitigations** you have identified

### What to Expect

- **Acknowledgment**: We will acknowledge receipt of your vulnerability report within 48 hours
- **Initial Assessment**: We will provide an initial assessment within 5 business days
- **Status Updates**: We will keep you informed about our progress
- **Disclosure**: We will work with you to understand and resolve the issue before any public disclosure
- **Credit**: We will credit you in the security advisory (unless you prefer to remain anonymous)

## Security Best Practices for Users

### Dependency Injection Context Usage

When using context-aware dependency injection:

1. **Sensitive Data in Context**
   - Be cautious when storing sensitive data in context values
   - Never store passwords, API keys, or tokens directly in context
   - Use secure context key types (not strings) to avoid collisions

   ```go
   // ❌ Bad - string keys are public
   ctx := context.WithValue(ctx, "apiKey", "secret123")
   
   // ✅ Good - unexported type prevents external access
   type contextKey string
   const apiKeyCtx contextKey = "apiKey"
   ctx := context.WithValue(ctx, apiKeyCtx, apiKey)
   ```

2. **Timeout and Cancellation**
   - Always set appropriate timeouts to prevent resource exhaustion
   - Use cancellation to stop long-running operations
   
   ```go
   ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
   defer cancel()
   ```

3. **Context Propagation**
   - Never reuse contexts across different requests
   - Each HTTP request should have its own context
   - Be aware that context values propagate through the entire dependency chain

### Error Handling

1. **ErrorCatcher Usage**
   - Always handle errors properly; don't silence them
   - Log errors with appropriate context
   - Never expose internal error details to end users

2. **Validation**
   - Use `validation.ValidateNotNil()` for parameter validation
   - Validate user inputs before processing
   - Don't trust external data sources

### Configuration Management

1. **File Configuration**
   - Store sensitive configuration in environment variables or secure vaults
   - Never commit secrets to version control
   - Use appropriate file permissions for configuration files (e.g., 0600)

2. **Encryption**
   - Use strong encryption keys (minimum 32 bytes for AES-256)
   - Rotate encryption keys periodically
   - Never hardcode encryption keys in source code

   ```go
   // ❌ Bad
   key := []byte("my-secret-key")
   
   // ✅ Good
   key := []byte(os.Getenv("ENCRYPTION_KEY"))
   if len(key) < 32 {
       return errors.New("encryption key must be at least 32 bytes")
   }
   ```

### Database Security

1. **Connection Strings**
   - Never log or expose database connection strings
   - Use environment variables for credentials
   - Enable SSL/TLS for database connections when possible

2. **SQL Injection Prevention**
   - The `persistence/dataaccess` package uses GORM parameterization
   - Always use the provided query builders
   - Never concatenate user input into queries

### Server Security

1. **HTTPS/TLS**
   - Always use HTTPS in production
   - Keep TLS certificates up to date
   - Use strong cipher suites

2. **CORS Configuration**
   - Configure CORS restrictively
   - Don't use `*` for allowed origins in production
   - Validate Origin headers

3. **Rate Limiting**
   - Implement rate limiting to prevent abuse
   - Monitor for unusual traffic patterns

## Security Updates

Security updates will be released as patch versions (e.g., 2.0.x) and documented in:
- GitHub Security Advisories
- CHANGELOG.md
- Release notes

Subscribe to repository notifications to stay informed about security updates.

## Vulnerability Disclosure Policy

We follow a **coordinated disclosure** process:

1. Security vulnerabilities are reported privately
2. We work with the reporter to understand and fix the issue
3. A fix is developed and tested
4. A security advisory is prepared
5. The fix is released
6. The advisory is published 24-48 hours after release

We aim to release security fixes within:
- **Critical vulnerabilities**: 7 days
- **High severity**: 30 days
- **Medium/Low severity**: 90 days

## Contact

For security-related questions or concerns that are not vulnerabilities, you can:
- Open a discussion in the GitHub Discussions tab
- Contact the maintainers through GitHub

---

Thank you for helping keep `go-infrastructure` and its users safe!
