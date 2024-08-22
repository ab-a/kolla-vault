# kolla-vault

Push kolla-ansible passwords in Hashicorp Vault and replace plaintext password by Vault lookups.

## Blog posts
[Hashicorp Vault and Kolla Ansible, Part I: Integrate Vault secrets in your playbook](https://abayard.com/hashicorp-vault-and-kolla-ansible-part-i-integrate-secrets-in-playbook/)

[Hashicorp Vault and Kolla Ansible, Part II: integration with Gitlab CI](https://abayard.com/hashicorp-vault-and-kolla-ansible-part-ii-integration-with-gitlab-ci/)

## What are the scripts?
- `store_kolla_passwords.go`: push the passwords from `passwords.yml` into Hashicorp Vault. Equivalent of `kolla-writepwd`.
- `replace_passwords.go`: replace the plaintext passwords by lookups.

## Initialization
```bash
kolla-genpwd
export VAULT_TOKEN=$(vault print token)
go mod init kolla-vault
go get github.com/hashicorp/vault/api
go get gopkg.in/yaml.v2
```

## Run the scripts
```bash
go run store_kolla_passwords.go
go run replace_passwords.go
```

## Compile
```bash
go build -o export_kolla_passwords store_kolla_passwords.go
go build -o replace_passwords replace_passwords.go
```

## Snippet of `passwords.yml` lookups
```yml
nova_database_password: '{{ lookup('community.general.hashi_vault', 'secret/data/kolla/default/nova_database_password', 'url={{ vault_url }}', token=lookup('env', 'VAULT_TOKEN')) }}'
nova_keystone_password: '{{ lookup('community.general.hashi_vault', 'secret/data/kolla/default/nova_keystone_password', 'url={{ vault_url }}', token=lookup('env', 'VAULT_TOKEN')) }}'
nova_ssh_key: 'map[private_key:{{ lookup('community.general.hashi_vault', 'secret/data/kolla/default/nova_ssh_key/private_key', 'url={{ vault_url }}', token=lookup('env', 'VAULT_TOKEN')) }} public_key:{{ lookup('community.general.hashi_vault', 'secret/data/kolla/default/nova_ssh_key/public_key', 'url={{ vault_url }}', token=lookup('env', 'VAULT_TOKEN')) }}]'
nova_api_database_password: '{{ lookup('community.general.hashi_vault', 'secret/data/kolla/default/nova_api_database_password', 'url={{ vault_url }}', token=lookup('env', 'VAULT_TOKEN')) }}'
```

## Gitlab CI Pipeline
![pipeline dependencies](https://i.imgur.com/NyG296r.png)
