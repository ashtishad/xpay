# Generating Secrets

This guide covers best practices for generating cryptographic keys and secrets for the XPay application.

## 1. JWT ES256 ECDSA Keys

Generate ECDSA key pair using OpenSSL:

```bash
# Generate private key
openssl ecparam -name prime256v1 -genkey -noout -out private.pem

# Extract public key
openssl ec -in private.pem -pubout -out public.pem

# Base64 encode keys
base64 -w 0 private.pem > private_base64.txt
base64 -w 0 public.pem > public_base64.txt
```

Store base64 encoded keys in `config.yaml`:

```yaml
jwt:
  private_key: "<contents of private_base64.txt>"
  public_key: "<contents of public_base64.txt>"
```

### Further Reading:
- [NIST SP 800-56A: Recommendation for Pair-Wise Key Establishment Schemes Using Discrete Logarithm Cryptography](https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-56Ar3.pdf)
- [RFC 7518 - JSON Web Algorithms (JWA)](https://tools.ietf.org/html/rfc7518#section-3.4)
- [OpenSSL Command-Line HOWTO](https://www.madboa.com/geek/openssl/)

## 2. Card AES Key

Generate a secure AES key:

```bash
openssl rand -base64 32 > aes_key.txt
```

Store in `config.yaml`:

```yaml
card:
  aes_key: "<contents of aes_key.txt>"
```

### Further Reading:
- [NIST SP 800-38A: Recommendation for Block Cipher Modes of Operation](https://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-38a.pdf)
- [NIST SP 800-57 Part 1 Rev. 5: Recommendation for Key Management](https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-57pt1r5.pdf)
- [OpenSSL Cookbook](https://www.feistyduck.com/library/openssl-cookbook/)

## 3. GitHub Actions Secrets

For CI/CD pipelines, use GitHub Actions secrets:

1. Go to repository Settings > Secrets and variables > Actions
2. Click "New repository secret"
3. Add secrets (e.g., `JWT_PRIVATE_KEY`, `CARD_AES_KEY`, `DB_URL`)

Usage in workflow:

```yaml
steps:
  - name: Use secrets
    env:
      JWT_PRIVATE_KEY: ${{ secrets.JWT_PRIVATE_KEY }}
      CARD_AES_KEY: ${{ secrets.CARD_AES_KEY }}
    run: |
      # Use environment variables in your scripts
```

### Further Reading:
- [GitHub Actions: Encrypted secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
- [Security hardening for GitHub Actions](https://docs.github.com/en/actions/security-guides/security-hardening-for-github-actions)
- [GitHub Actions: Creating and storing encrypted secrets](https://docs.github.com/en/actions/reference/encrypted-secrets)

## Best Practices

1. Never commit sensitive keys to version control.
2. Use different keys for each environment (dev, staging, prod).
3. Implement key rotation policies.
4. Apply least privilege principle for key access.
5. Audit and monitor key usage.
6. Use strong, randomly generated keys.
7. Encrypt keys at rest and in transit.
8. Use managed identities or IAM roles instead of long-lived access keys.
9. Enable logging for all key operations.
10. Regularly audit access and usage patterns.

### Further Reading on Best Practices:
- [OWASP Cryptographic Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cryptographic_Storage_Cheat_Sheet.html)
- [Cloud Security Alliance: Cloud Controls Matrix](https://cloudsecurityalliance.org/research/cloud-controls-matrix/)
- [CIS Controls](https://www.cisecurity.org/controls/cis-controls-list)

## For Key Management in Production

For detailed information on managing these secrets in a production environment, please refer to the [Configuration and Key Management in Production](https://github.com/ashtishad/xpay/blob/main/guides/configuration_key_management_in_production.md) guide.
