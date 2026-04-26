# samba4-manager
Web Administration Panel for Samba 4 Active Directory

> **Status:** ✅ Stable (v1.2.0) — Fork ativo com correções e melhorias
> **Stack:** Golang, Echo Framework, GORM, Tailwind CSS (Neo-Brutalism), jQuery, DataTables, Lucide Icons, LDAP, Kerberos.

> 💡 Fork do projeto original [go-samba4](https://github.com/jniltinho/go-samba4) por [@jniltinho](https://github.com/jniltinho), com correções críticas e melhorias de usabilidade.

## Overview
The **Samba4 AD Web Admin Panel** is a modern web administration tool to manage Samba 4 Active Directory environments. Built to replace legacy interfaces (such as SWAT or manual `samba-tool` commands), it offers a fast, secure, and extensible solution under modern, high-performance technologies.

The interface focuses on functional clarity (Neo-Brutalist style) aimed at IT teams and systems administrators managing Samba domains without Microsoft's infrastructure and tools (RSAT). The panel operates as a modular monolith, where the backend and frontend (templates and static assets) are unified into a single self-sufficient binary via `go:embed`.

## Novidades nesta versão (v1.2.0) ✨

- **🔧 RBAC corrigido:** Autenticação agora funciona com grupos em formato DN completo (ex: `CN=Domain Admins,CN=Users,DC=domain,DC=local`)
- **👤 Botão "+ NOVO USUÁRIO" visível:** Removida verificação incorreta de `{{if .IsAdmin}}` no template
- **🔒 StartTLS habilitado por padrão:** Configuração segura de LDAP já vem com `use_tls = true` no `config.toml.example`
- **🏷️ Projeto renomeado:** De `go-samba4` para `samba4-manager` em todos os imports, Dockerfile, Makefile e workflows

## Key Features 🚀

- **Comprehensive User and Group Management:** Full CRUD operations in AD (LDAP), disable accounts, reset passwords.
- **OU Tree Navigation:** Hierarchical view and movement of AD objects.
- **Secure Authentication:** Login via LDAP bind, support for Kerberos (SSO), and Two-Factor Authentication (TOTP).
- **Access Control (RBAC):** Conditional permissions (Admin, Operator, Helpdesk) based on AD groups.
- **Auditing and Monitoring:** Detailed tracking with local change logs and searchable DataTables views.
- **Advanced Search:** Find objects with advanced customizable LDAP filters.
- **Independent Local Database:** Embedded SQLite by default (or configurable MySQL/MariaDB via GORM) for audit logs, sessions, and settings.
- **Internationalization:** UI fully translated in English, Portuguese (pt_BR), and Spanish, including table controls.
- **Active Navigation:** Sidebar highlights the current section, preserving context across nested routes.

## Technology & Architecture 🛠️

The architecture guarantees zero dependency on external assets in production.

- **Backend:** Go (`1.26+`), `Echo` (HTTP), `GORM` (Data Modeling)
- **AD/Samba Integration:** `go-ldap/ldap` and `gokrb5` (Kerberos)
- **Frontend:** Server-Side Rendering (SSR) via `html/template`, `TailwindCSS 4.2+`, `jQuery 4`, `DataTables 2.3.7`, `Lucide Icons`
- **i18n:** Portuguese (pt_BR), English (en), Spanish (es) — including DataTables localization
- **CLI Tooling:** `Cobra` & `Viper` (`config.toml` Configurations)

## Build ⚙️

```bash
# Compile Tailwind CSS
make css

# Build the binary (runs css + go build + UPX compression)
make build

# Run the server after building
make run
```

## CLI Usage ⚙️

Because the application is structured around a powerful `cobra` CLI, it provides helper commands alongside starting the server:

```bash
# Start the Web Admin Server (Default Port 8080)
./samba4-manager serve --port 8080

# Run local application database migrations (Sessions and Logs)
./samba4-manager migrate

# Use the emergency CLI for skeleton tasks (local bypass users)
./samba4-manager user

# Print the current version and build date
./samba4-manager version
```

> **Note:** The core application configuration, including LDAP and TLS communication, is managed externally via a `config.toml` file pointed by the `--config` flag if needed (defaults to `./config.toml`).

## Security Requirements 🔒

The environment should implement the following mandatory production practices:
- **Transport Security:** Active Server HTTPS/TLS certificates along with encrypted LDAPS communication (port `636`) to the native Domain Controller.
- **Immutable Auditing:** Appended local files recording the access and delegations performed on the interface.

## Reference
For further architectural and deep roadmap details, please refer to the [Application PRD](DOCUMENTS/PRD-Samba4-AD-WebAdmin.md).
