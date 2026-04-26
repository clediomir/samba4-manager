# samba4-manager

[![Release](https://img.shields.io/github/v/release/clediomir/samba4-manager?sort=semver)](https://github.com/clediomir/samba4-manager/releases)
[![Docker](https://img.shields.io/badge/Docker-available-2496ED?logo=docker)](https://github.com/clediomir/samba4-manager/pkgs/container/samba4-manager)
[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)](https://golang.org)
[![License](https://img.shields.io/github/license/clediomir/samba4-manager)](LICENSE)

Web Administration Panel for Samba 4 Active Directory

> 💡 Fork do projeto original [go-samba4](https://github.com/jniltinho/go-samba4) por [@jniltinho](https://github.com/jniltinho), com correções críticas e melhorias de usabilidade.

## Overview

`samba4-manager` é um painel web moderno para gerenciar ambientes Samba 4 Active Directory. Construído para substituir interfaces legadas (como SWAT ou comandos manuais `samba-tool`), oferece uma solução rápida, segura e extensível.

A interface foca em clareza funcional (estilo Neo-Brutalist) para equipes de TI e administradores de sistema que gerenciam domínios Samba sem a infraestrutura e ferramentas da Microsoft (RSAT).

**Destaques desta versão (v1.2.0):**

- 🔧 **RBAC corrigido:** Autenticação funciona com grupos em formato DN completo
- 👤 **Botão "+ NOVO USUÁRIO"** visível para admins
- 🔒 **StartTLS habilitado** por padrão no `config.toml.example`
- 🐳 **Imagem Docker** publicada no GitHub Container Registry
- 🏷️ **Projeto renomeado** para `samba4-manager`

## Key Features

- **CRUD Completo de Usuários e Grupos:** Operações full no AD (LDAP), desabilitar contas, reset de senhas
- **Navegação por OUs:** Visualização hierárquica e movimentação de objetos AD
- **Autenticação Segura:** Login via LDAP bind, suporte a Kerberos (SSO) e 2FA (TOTP)
- **Controle de Acesso (RBAC):** Permissões condicionais (Admin, Operator, Helpdesk) baseadas em grupos AD
- **Auditoria:** Tracking detalhado com logs de alterações locais e visualização pesquisável
- **Busca Avançada:** Encontre objetos com filtros LDAP personalizáveis
- **Banco Local Independente:** SQLite embarcado por padrão (ou MySQL/MariaDB via GORM) para logs, sessões e configurações
- **Internacionalização:** UI traduzida em Inglês, Português (pt_BR) e Espanhol

## Quick Start

### Docker (Recomendado)

```bash
# 1. Crie o arquivo de configuração
cp config.toml.example config.toml
# Edite config.toml com os dados do seu Samba 4 AD

# 2. Suba o container
docker run -d \
  --name samba4-manager \
  -p 8080:8080 \
  -v $(pwd)/config.toml:/etc/samba4-manager/config.toml:ro \
  -e SAMBA4_LDAP_PASS=suasenha \
  ghcr.io/clediomir/samba4-manager:latest
```

### Docker Compose

```bash
# Clone o repositório
git clone https://github.com/clediomir/samba4-manager.git
cd samba4-manager

# Configure
cp config.toml.example config.toml
# Edite config.toml

# Suba
docker compose up -d
```

### Build Manual

```bash
# Dependências
sudo apt install -y golang nodejs upx

# Build
make build

# Execute
./samba4-manager serve
```

## Configuração

Crie `config.toml` baseado no `config.toml.example`:

```toml
[ldap]
host            = "dc1.empresa.local"
port            = 636
use_tls         = true
skip_tls_verify = true
base_dn         = "DC=empresa,DC=local"
bind_user       = "CN=admin,CN=Users,DC=empresa,DC=local"

[rbac]
admin_group    = "Domain Admins"
operator_group = "SambaWebOperators"
```

> A senha do bind LDAP deve ser fornecida via variável de ambiente `SAMBA4_LDAP_PASS`.

## Acesso

Após iniciar, acesse: **http://localhost:8080**

Faça login com uma conta do Active Directory que pertença ao grupo `Domain Admins` (ou o grupo configurado como `admin_group`).

## RBAC — Níveis de Acesso

| Perfil | Permissões |
|--------|-----------|
| **Admin** | Acesso completo (CRUD, configurações) |
| **Operator** | CRUD de usuários e grupos |
| **Helpdesk** | Reset de senha, desabilitar contas |
| **ReadOnly** | Apenas visualização |

## Architecture

| Camada | Tecnologia |
|--------|-----------|
| Backend | Go 1.26+, Echo Framework, GORM |
| AD/Samba | go-ldap/ldap, gokrb5 (Kerberos) |
| Frontend | Server-Side Rendering, TailwindCSS 4, DataTables, Lucide Icons |
| i18n | Português (pt_BR), Inglês, Espanhol |
| CLI | Cobra & Viper |
| Database | SQLite (default) ou MySQL/MariaDB |

## Contribuindo

1. Fork o repositório
2. Crie uma branch feature (`git checkout -b feature/nova-funcionalidade`)
3. Commit suas mudanças (`git commit -m 'feat: adiciona nova funcionalidade'`)
4. Push para a branch (`git push origin feature/nova-funcionalidade`)
5. Abra um Pull Request

Veja as [issues](https://github.com/clediomir/samba4-manager/issues) para saber o que está sendo trabalhado.

## License

[MIT License](LICENSE) — sinta-se à vontade para usar, modificar e distribuir.

## Agradecimentos

- [@jniltinho](https://github.com/jniltinho) — Criador do projeto original [go-samba4](https://github.com/jniltinho/go-samba4)
- Comunidade Samba — Por manter o Samba 4 AD em constante evolução
