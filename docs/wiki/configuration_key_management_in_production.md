# Configuration and Key Management in Production

This guide outlines strategies for securely managing configuration and sensitive data for the XPay application in production environments.

## Current Setup (Local Development)

```yaml
# config.yaml
app:
  env: dev
  gin_mode: release
  server_address: ":8080"

db:
  url: "postgres://ash:lol@127.0.0.1:5432/xpay?sslmode=disable&timezone=UTC"
  max_open_conns: 18
  max_idle_conns: 18
  conn_max_lifetime: "1h"
  conn_max_idle_time: "30m"

jwt:
  private_key: "LS0tLS1CRUdJTiBFQyBQUklWQVRFIEtFWS0t..."
  public_key: "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZrd..."

card:
  aes_key: "CWcKy/Jl/FOwCevQfkWDSGU5QZt0WMZCh/kC68k1LmM="
```

## Production Configuration Management

### 1. AWS Specific Solutions (KMS, Parameter Store, Secrets Manager)

#### AWS Key Management Service (KMS)
- Purpose: Manage encryption keys for sensitive data.
- Usage: Encrypt secrets stored in Parameter Store and Secrets Manager.

#### AWS Systems Manager Parameter Store
- Purpose: Store regular configuration and encrypted sensitive data.
- Usage:
  - Regular config: Store as String type.
  - Sensitive data: Store as SecureString type, encrypted with KMS.

#### AWS Secrets Manager
- Purpose: Store and automatically rotate database credentials.
- Usage: Store DB credentials with automatic rotation.

### 2. HashiCorp Vault

- Purpose: Centralized secrets management, suitable for multi-cloud environments.
- Usage:
  - Store secrets and configuration in key-value store.
  - Leverage dynamic secrets for database credentials.

### 3. Ansible Vault

- Purpose: Encrypt sensitive data within Ansible playbooks or YAML files.
- Usage:
  - Encrypt entire config.yaml file.
  - Deploy encrypted file to servers using Ansible playbooks.

## Cost Comparison (As of 2024) & Choosing a Solution

1. AWS Solutions:
   - KMS: $1/month per CMK, $0.03 per 10,000 API calls
   - Parameter Store: Free for standard parameters, $0.05 per advanced parameter per month
   - Secrets Manager: $0.40 per secret per month, $0.05 per 10,000 API calls

2. HashiCorp Vault:
   - Open Source: Free (self-managed)
   - Enterprise: From $500/node/year
   - Cloud: Starting at ~$0.04/hour

3. Ansible Vault:
   - Free (part of Ansible)

Choosing a Solution:
- AWS-native applications: Use AWS services for tight integration and simplicity.
- Multi-cloud or hybrid environments: Consider HashiCorp Vault for flexibility.
- Simple deployments or existing Ansible users: Ansible Vault can be a cost-effective choice.

Consider factors like scalability needs, existing infrastructure, compliance requirements, and team expertise when making your decision.

## Viper Integration

Sample functions for loading configurations:

```go
// AWS Parameter Store and Secrets Manager
func loadConfigFromAWS(v *viper.Viper) error {
    if err := loadFromParameterStore(v); err != nil {
        return err
    }
    return loadFromSecretsManager(v)
}

// HashiCorp Vault
func loadConfigFromVault(v *viper.Viper) error {
    client, err := vault.NewClient(vault.DefaultConfig())
    if err != nil {
        return err
    }
    secret, err := client.Logical().Read("secret/xpay/config")
    if err != nil {
        return err
    }
    for key, value := range secret.Data {
        v.Set(key, value)
    }
    return nil
}

// Ansible Vault
func loadConfigFromAnsibleVault(v *viper.Viper) error {
    v.SetConfigName("config")
    v.SetConfigType("yaml")
    v.AddConfigPath(".")
    v.AutomaticEnv()
    return v.ReadInConfig()
}
```

## GitHub Actions

For AWS:

```yaml
- name: Configure AWS Credentials
  uses: aws-actions/configure-aws-credentials@v1
  with:
    aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
    aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
    aws-region: us-west-2

- name: Load Configuration
  run: |
    # Script to load configuration from Parameter Store/Secrets Manager
```

For HashiCorp Vault:

```yaml
- name: Load Vault Configuration
  env:
    VAULT_ADDR: ${{ secrets.VAULT_ADDR }}
    VAULT_TOKEN: ${{ secrets.VAULT_TOKEN }}
  run: |
    # Script to load configuration from Vault
```

For Ansible Vault:

```yaml
- name: Decrypt Configuration
  run: ansible-vault decrypt config.yaml --vault-password-file vault_pass.txt
  env:
    ANSIBLE_VAULT_PASSWORD_FILE: ${{ secrets.ANSIBLE_VAULT_PASSWORD }}
```

## Further Reading

- AWS Services:
  - [AWS KMS Developer Guide](https://docs.aws.amazon.com/kms/latest/developerguide/overview.html)
  - [AWS Systems Manager Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html)
  - [AWS Secrets Manager User Guide](https://docs.aws.amazon.com/secretsmanager/latest/userguide/intro.html)

- HashiCorp Vault:
  - [Vault Documentation](https://www.vaultproject.io/docs)
  - [Vault vs. Other Software](https://www.vaultproject.io/intro/vs)

- Ansible Vault:
  - [Ansible Vault Guide](https://docs.ansible.com/ansible/latest/user_guide/vault.html)
  - [Using encrypted variables and files](https://docs.ansible.com/ansible/latest/user_guide/playbooks_best_practices.html#keep-vaulted-variables-safely-visible)

- General Security Best Practices:
  - [OWASP Secrets Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Secrets_Management_Cheat_Sheet.html)
  - [NIST Cryptographic Standards and Guidelines](https://csrc.nist.gov/projects/cryptographic-standards-and-guidelines)

- Configuration Management:
  - [The Twelve-Factor App: Config](https://12factor.net/config)
  - [Viper Documentation](https://github.com/spf13/viper)
